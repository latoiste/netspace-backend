package db

import (
	"context"

	"github.com/latoiste/netspace/model"
)

func (r *Repository) InsertGroup(group model.Group, ctx context.Context) error {
	query := `
		INSERT INTO groups (
			id,
			name,
			isActive
		)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		group.Id,
		group.Name,
		group.IsActive,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) InsertGroupMember(groupId string, userId string, ctx context.Context) error {
	const query = `
		INSERT INTO groupmembers (
			groupid,
			userid
		)
		VALUES ($1, $2)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		groupId,
		userId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateGroupIsActive(groupId string, isActive bool, ctx context.Context) error {
	query := `
		UPDATE groups
		SET isActive = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		isActive,
		groupId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GroupById(groupId string, ctx context.Context) (*model.Group, error) {
	const query = `
		SELECT
			id,
			name,
			isactive
		FROM groups
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, groupId)

	var group model.Group

	err := row.Scan(
		&group.Id,
		&group.Name,
		&group.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return &group, nil
}
