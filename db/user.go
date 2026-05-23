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
		SELECT * FROM Users
		WHERE id=$1
	`

	row := r.db.QueryRowContext(ctx, query, userId)

	var user model.User

	if err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Age,
		&user.Gender,
		&user.Slug,
		&user.LocationId,
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
		INSERT INTO Users 
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := db.ExecContext(ctx, query,
		user.Id,
		user.Name,
		user.Age,
		user.Gender,
		user.Slug,
		user.LocationId,
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
		SELECT u.id, u.name, u.age, u.gender, u.slug, u.locationId
		FROM Locations l
		JOIN Users u
		ON l.Id=u.locationId
		WHERE locationId=$1;
	`

	users := make([]model.User, 0)

	rows, err := r.db.QueryContext(ctx, query, locationId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user model.User
		err = rows.Scan(
			&user.Id,
			&user.Name,
			&user.Age,
			&user.Gender,
			&user.Slug,
			&user.LocationId,
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
