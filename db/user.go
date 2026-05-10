package db

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/latoiste/netspace/model"
)

func (e *Env) InsertUser(user model.User, ctx context.Context) error {
	db := e.db

	const query = `
		INSERT INTO Users 
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := db.ExecContext(ctx, query,
		user.Id,
		user.Name,
		user.Age,
		user.Gender,
		user.Slug,
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

// ('☕', 'Kopi'),
//   ('🎮', 'Gaming'),
//   ('📚', 'Buku'),
//   ('🎵', 'Musik'),
//   ('🍜', 'Kuliner'),
//   ('✈️', 'Travel'),
//   ('💻', 'Tech'),
//   ('🎨', 'Seni'),
//   ('🏋️', 'Olahraga'),
//   ('🎬', 'Film'),
//   ('📷', 'Fotografi'),
//   ('🌱', 'Tanaman')
