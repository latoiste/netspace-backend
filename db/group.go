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

// UpdateGroupName renames a group. The new name is what the chat list and group
// header show after a re-fetch; the live group_renamed event updates open views.
func (r *Repository) UpdateGroupName(groupId string, name string, ctx context.Context) error {
	const query = `UPDATE groups SET name = $1 WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, name, groupId)
	if err != nil {
		return err
	}

	return nil
}

// RemoveGroupMember drops one membership row. Called when a user explicitly
// leaves a group, so the group disappears from THEIR chat list while staying
// intact for everyone else.
func (r *Repository) RemoveGroupMember(groupId string, userId string, ctx context.Context) error {
	const query = `DELETE FROM groupmembers WHERE groupid = $1 AND userid = $2`

	_, err := r.db.ExecContext(ctx, query, groupId, userId)
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

// GroupMembers returns every user currently recorded as a member of the group.
// Used to seed the member bar (incl. for users who join mid-session and would
// otherwise miss the existing roster).
func (r *Repository) GroupMembers(groupId string, ctx context.Context) ([]model.User, error) {
	const query = `
		SELECT u.id, u.locationId, u.name, u.slug, u.age, u.gender, u.createdAt, u.isActive
		FROM GroupMembers gm
		JOIN Users u ON u.id = gm.userId
		WHERE gm.groupId = $1
	`

	users := make([]model.User, 0)

	rows, err := r.db.QueryContext(ctx, query, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.Id,
			&user.LocationId,
			&user.Name,
			&user.Slug,
			&user.Age,
			&user.Gender,
			&user.CreatedAt,
			&user.IsActive,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
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
