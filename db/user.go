package db

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/latoiste/netspace/model"
)

func (e *Env) InsertUser(user model.User, ctx context.Context) error {
	db := e.db

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

	err = e.insertInterests(user.Id, user.Interests)
	if err != nil {
		return err
	}
	log.Println("Inserted to user successfully")

	return nil
}

func (e *Env) insertInterests(userId string, interests []model.Interest) error {
	interestIds := make([]int, 0)
	for _, interest := range interests {
		if !interest.IsCustom {
			interestIds = append(interestIds, interest.Id)
		} else {
			if err := e.insertCustomInterest(userId, interest); err != nil {
				return err
			}
		}
	}
	if err := e.insertUserInterests(userId, interestIds); err != nil {
		return err
	}
	return nil
}

func (e *Env) insertUserInterests(userId string, interestIds []int) error {
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

	_, err := e.db.Exec(query, args...)
	return err
}

func (e *Env) insertCustomInterest(userId string, interest model.Interest) error {
	const query = `
        INSERT INTO UserCustomInterests (userId, emoji, label)
        VALUES ($1, $2, $3)
    `

	_, err := e.db.Exec(
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

func (e *Env) UserInterests(userId string, ctx context.Context) ([]model.Interest, error) {
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

	rows, err := e.db.QueryContext(ctx, query, userId)
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

func (e *Env) UsersInLocation(locationId int, ctx context.Context) ([]model.User, error) {
	const query = `
		SELECT u.id, u.name, u.age, u.gender, u.slug, u.locationId
		FROM Locations l
		JOIN Users u
		ON l.Id=u.locationId
		WHERE locationId=$1;
	`

	users := make([]model.User, 0)

	rows, err := e.db.QueryContext(ctx, query, locationId)
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

		user.Interests, err = e.UserInterests(user.Id, ctx2)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
