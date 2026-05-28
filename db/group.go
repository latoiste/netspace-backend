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
