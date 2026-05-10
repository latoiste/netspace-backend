package db

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/latoiste/netspace/api"
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

func (e *Env) GetUserInterests(userId string, ctx context.Context) ([]api.InterestOutput, error) {
	const query = `
		SELECT i.emoji, i.label
		FROM UserInterests ui
		JOIN Interests i ON i.id = ui.interestId
		WHERE ui.userId = $1

		UNION ALL

		SELECT uci.emoji, uci.label
		FROM UserCustomInterests uci
		WHERE uci.userId = $1
	`

	interests := make([]api.InterestOutput, 0)

	rows, err := e.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var it api.InterestOutput
		err = rows.Scan(&it.Emoji, &it.Label)
		if err != nil {
			return nil, err
		}
		interests = append(interests, it)
	}
	return interests, nil
}

func (e *Env) UsersInLocation(locationId int, ctx context.Context) ([]api.UserOutput, error) {
	const query = `
		SELECT
			id,
			slug,
			name
		FROM Users
		WHERE locationId = $1;
	`

	userOutput := make([]api.UserOutput, 0)

	rows, err := e.db.QueryContext(ctx, query, locationId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user api.UserOutput
		err = rows.Scan(
			&user.Id,
			&user.Slug,
			&user.Name,
		)
		if err != nil {
			return nil, err
		}

		ctx2, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		user.Interests, err = e.GetUserInterests(user.Id, ctx2)
		if err != nil {
			return nil, err
		}
		userOutput = append(userOutput, user)
	}

	return userOutput, nil
}
