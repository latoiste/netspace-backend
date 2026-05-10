package db

import (
	"context"

	"github.com/latoiste/netspace/model"
)

func (e *Env) LocationBySlug(slug string, ctx context.Context) (*model.Location, error) {
	db := e.db

	const query = `
		SELECT slug, name, isActive
		FROM Locations 
		WHERE slug=$1
	`

	row := db.QueryRowContext(ctx, query, slug)

	var location model.Location

	if err := row.Scan(
		&location.Slug,
		&location.Name,
		&location.IsActive,
	); err != nil {
		return nil, err
	}

	return &location, nil
}

func (e *Env) LocationIdBySlug(slug string, ctx context.Context) (int, error) {
	var id int

	err := e.db.QueryRowContext(ctx,
		`SELECT id FROM Locations WHERE slug = $1`,
		slug,
	).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}
