package db

import (
	"context"

	"github.com/latoiste/netspace/model"
)

func (r *Repository) LocationBySlug(slug string, ctx context.Context) (*model.Location, error) {
	db := r.db

	const query = `
		SELECT 
		id,
		slug, 
		name, 
		address,
		partnerid,
		joindate,
		capacity,
		timezone,
		isactive,
		qrtoken,
		qrlabel,
		COALESCE(latitude, 0) AS latitude,
		COALESCE(longitude, 0) AS longitude,
		COALESCE(geofenceRadius, 100) AS geofenceRadius
		FROM Locations
		WHERE slug=$1
	`

	row := db.QueryRowContext(ctx, query, slug)

	var location model.Location

	if err := row.Scan(
		&location.Id,
		&location.Slug,
		&location.Name,
		&location.Address,
		&location.PartnerId,
		&location.JoinDate,
		&location.Capacity,
		&location.Timezone,
		&location.IsActive,
		&location.QrToken,
		&location.QrLabel,
		&location.Latitude,
		&location.Longitude,
		&location.GeofenceRadius,
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

func (r *Repository) UpdateLocationIsActive(slug string, isActive bool, ctx context.Context) error {
	const query = `
		UPDATE locations
		SET isactive = $1
		WHERE slug = $2
	`

	_, err := r.db.ExecContext(ctx, query, isActive, slug)
	if err != nil {
		return err
	}
	return nil
}
