package db

import (
	"context"

	"github.com/latoiste/netspace/model"
)

func (e *Env) LocationBySlug(slug string, ctx context.Context) (*model.Location, error) {
	db := e.db

	row := db.QueryRowContext(ctx, "SELECT slug, name, isActive FROM Location WHERE slug=$1", slug)

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
