package db

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/latoiste/netspace/model"
)

func (r *Repository) UserById(userId string, ctx context.Context) (*model.User, error) {
	const query = `
		SELECT
			id,
			locationId,
			name,
			slug,
			age,
			gender,
			COALESCE(occupation, '') AS occupation,
			createdAt,
			isActive
		FROM Users
		WHERE id=$1
	`

	row := r.db.QueryRowContext(ctx, query, userId)

	var user model.User

	if err := row.Scan(
		&user.Id,
		&user.LocationId,
		&user.Name,
		&user.Slug,
		&user.Age,
		&user.Gender,
		&user.Occupation,
		&user.CreatedAt,
		&user.IsActive,
	); err != nil {
		return nil, err
	}

	ctx2, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	interests, err := r.UserInterests(userId, ctx2)
	if err != nil {
		return nil, err
	}

	user.Interests = interests

	return &user, nil
}

func (r *Repository) InsertUser(user model.User, ctx context.Context) error {
	db := r.db

	const query = `
		INSERT INTO Users (id, locationId, name, slug, age, gender, occupation)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := db.ExecContext(ctx, query,
		user.Id,
		user.LocationId,
		user.Name,
		user.Slug,
		user.Age,
		user.Gender,
		user.Occupation,
	)
	if err != nil {
		return err
	}

	err = r.insertInterests(user.Id, user.Interests)
	if err != nil {
		return err
	}
	log.Println("Inserted to user successfully")

	return nil
}

func (r *Repository) insertInterests(userId string, interests []model.Interest) error {
	interestIds := make([]int, 0)
	for _, interest := range interests {
		if !interest.IsCustom {
			interestIds = append(interestIds, interest.Id)
		} else {
			if err := r.insertCustomInterest(userId, interest); err != nil {
				return err
			}
		}
	}
	if err := r.insertUserInterests(userId, interestIds); err != nil {
		return err
	}
	return nil
}

func (r *Repository) insertUserInterests(userId string, interestIds []int) error {
	if len(interestIds) == 0 {
		return nil
	}

	var (
		values []string
		args   []any
	)

	argPos := 1

	for _, interestId := range interestIds {
		values = append(
			values,
			fmt.Sprintf("($%d, $%d)", argPos, argPos+1),
		)

		args = append(args, userId, interestId)

		argPos += 2
	}

	query := fmt.Sprintf(`
        INSERT INTO UserInterests (userId, interestId)
        VALUES %s
    `, strings.Join(values, ","))

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *Repository) insertCustomInterest(userId string, interest model.Interest) error {
	const query = `
        INSERT INTO UserCustomInterests (userId, emoji, label)
        VALUES ($1, $2, $3)
    `

	_, err := r.db.Exec(
		query,
		userId,
		interest.Emoji,
		interest.Label,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UserInterests(userId string, ctx context.Context) ([]model.Interest, error) {
	const query = `
		SELECT i.id, i.emoji, i.label, false as isCustom
		FROM UserInterests ui
		JOIN Interests i ON i.id = ui.interestId
		WHERE ui.userId = $1

		UNION ALL

		SELECT -1 as id, uci.emoji, uci.label, true as isCustom
		FROM UserCustomInterests uci
		WHERE uci.userId = $1
	`

	interests := make([]model.Interest, 0)

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var it model.Interest
		err = rows.Scan(
			&it.Id,
			&it.Emoji,
			&it.Label,
			&it.IsCustom,
		)
		if err != nil {
			return nil, err
		}
		interests = append(interests, it)
	}
	return interests, nil
}

func (r *Repository) UsersInLocation(locationId int, ctx context.Context) ([]model.User, error) {
	const query = `
		SELECT
			u.id,
			u.locationId,
			u.name,
			u.slug,
			u.age,
			u.gender,
			COALESCE(u.occupation, '') AS occupation,
			u.createdAt,
			u.isActive
		FROM Locations l
		JOIN Users u
		ON l.Id=u.locationId
		WHERE locationId=$1 AND u.isActive = true;
	`

	users := make([]model.User, 0)

	rows, err := r.db.QueryContext(ctx, query, locationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err = rows.Scan(
			&user.Id,
			&user.LocationId,
			&user.Name,
			&user.Slug,
			&user.Age,
			&user.Gender,
			&user.Occupation,
			&user.CreatedAt,
			&user.IsActive,
		)
		if err != nil {
			return nil, err
		}

		ctx2, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		user.Interests, err = r.UserInterests(user.Id, ctx2)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *Repository) UpdateUserIsActive(userId string, isActive bool, ctx context.Context) error {
	query := `
		UPDATE users
		SET isactive = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, isActive, userId)
	if err != nil {
		return err
	}

	return nil
}

// DeactivateAllUsers flips every still-active user to inactive. Called once on
// backend startup to clear "ghosts" — users left marked active because the
// previous process died before its per-socket disconnect handlers could flip
// them off. After a restart nobody has a live socket yet, so the only correct
// state is everyone inactive; real users re-arm isactive when they reconnect.
func (r *Repository) DeactivateAllUsers(ctx context.Context) error {
	query := `
		UPDATE users
		SET isactive = false
		WHERE isactive = true
	`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserChatData purges the private things a user generated — their DMs (on
// both sides of the conversation), group messages, group memberships and
// notifications — so a session leaves no *personal* chat history behind on
// logout. The Users row itself is kept (marked inactive) so analytics still add
// up. Runs in a single transaction; child tables reference Users, so deleting
// the child rows never trips a foreign-key constraint.
//
// Public room messages are deliberately NOT deleted here: the public room is a
// shared, venue-scoped timeline, so removing one person's lines when they leave
// would tear holes in the conversation everyone else can see. Public messages
// instead age out via a time-based retention sweep (DeletePublicMessagesBefore).
func (r *Repository) DeleteUserChatData(userId string, ctx context.Context) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmts := []string{
		`DELETE FROM PrivateMessages WHERE senderId = $1 OR recipientId = $1`,
		`DELETE FROM GroupMessages WHERE senderId = $1`,
		`DELETE FROM GroupMembers WHERE userId = $1`,
		`DELETE FROM Notifications WHERE userId = $1`,
	}

	for _, q := range stmts {
		if _, err := tx.ExecContext(ctx, q, userId); err != nil {
			return err
		}
	}

	return tx.Commit()
}
