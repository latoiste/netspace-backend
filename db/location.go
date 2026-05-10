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
