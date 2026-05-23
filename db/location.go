package db

import (
	"context"

	"github.com/latoiste/netspace/model"
)

func (r *Repository) LocationBySlug(slug string, ctx context.Context) (*model.Location, error) {
	db := r.db

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

func (r *Repository) LocationIdBySlug(slug string, ctx context.Context) (int, error) {
	var id int

	err := r.db.QueryRowContext(ctx,
		`SELECT id FROM Locations WHERE slug = $1`,
		slug,
	).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}
