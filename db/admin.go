package db

import (
	"context"
	"time"

	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/model"
)

func (r *Repository) AdminByUsername(username string, ctx context.Context) (*model.Admin, error) {
	const query = `
		SELECT 
			"id",
			username,
			"password",
			"role",
			"plan",
			"avatar",
			"name"
		FROM admins
		WHERE username = $1
	`

	row := r.db.QueryRowContext(
		ctx,
		query,
		username,
	)

	var admin model.Admin

	if err := row.Scan(
		&admin.Id,
		&admin.Username,
		&admin.Password,
		&admin.Role,
		&admin.Plan,
		&admin.Avatar,
		&admin.Name,
	); err != nil {
		return nil, err
	}

	return &admin, nil
}

func (r *Repository) TotalCheckInRange(start time.Time, end time.Time, ctx context.Context) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM users
		WHERE createdAt >= $1 AND createdAt < $2
	`

	row := r.db.QueryRowContext(
		ctx,
		query,
		start,
		end,
	)

	var count int

	if err := row.Scan(
		&count,
	); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) TotalActiveUsersRange(start time.Time, end time.Time, ctx context.Context) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM users
		WHERE createdAt >= $1 AND createdAt < $2 AND isactive = true
	`

	row := r.db.QueryRowContext(
		ctx,
		query,
		start,
		end,
	)

	var count int

	if err := row.Scan(
		&count,
	); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) TotalMessagesRange(start time.Time, end time.Time, ctx context.Context) (int, error) {
	const query = `
		SELECT
		(
			SELECT COUNT(*)
			FROM PrivateMessages
			WHERE timestamp >= $1 AND timestamp < $2
		)
		+
		(
			SELECT COUNT(*)
			FROM PublicMessages
			WHERE timestamp >= $1 AND timestamp < $2
		)
		+
		(
			SELECT COUNT(*)
			FROM GroupMessages
			WHERE timestamp >= $1 AND timestamp < $2
		)
    `

	row := r.db.QueryRowContext(
		ctx,
		query,
		start,
		end,
	)

	var count int

	if err := row.Scan(
		&count,
	); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) ActiveUsers(locationId int, ctx context.Context) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM users
		WHERE isactive = true AND locationid = $1
	`

	row := r.db.QueryRowContext(ctx, query, locationId)

	var count int

	if err := row.Scan(
		&count,
	); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) GetTopInterests(ctx context.Context) ([]api.InterestPercentageDTO, error) {
	const query = `
		SELECT
			i.emoji,
			i.label,
			ROUND(
				COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (),
				0
			) AS percentage
		FROM UserInterests ui
		JOIN Interests i
			ON i.id = ui.interestId
		GROUP BY i.emoji, i.label
		ORDER BY percentage DESC;
	`

	topInterests := make([]api.InterestPercentageDTO, 0)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var it api.InterestPercentageDTO
		err := rows.Scan(
			&it.Emoji,
			&it.Label,
			&it.Percentage,
		)
		if err != nil {
			return nil, err
		}
		topInterests = append(topInterests, it)
	}

	return topInterests, nil
}

// func (r *Repository) GetAnalytics(ctx context.Context, from, to string) ([]model.AnalyticsData, error) {
// 	rows, err := r.db.QueryContext(ctx, `
// 		SELECT
// 			DATE(created_at) as date,
// 			COUNT(DISTINCT user_id) as active_users,
// 			COUNT(DISTINCT session_token) as total_sessions,
// 			0 as message_count
// 		FROM users
// 		WHERE DATE(created_at) BETWEEN $1 AND $2
// 		GROUP BY DATE(created_at)
// 		ORDER BY date ASC
// 	`, from, to)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var result []model.AnalyticsData
// 	for rows.Next() {
// 		var a model.AnalyticsData
// 		if err := rows.Scan(
// 			&a.Date,
// 			&a.ActiveUsers,
// 			&a.TotalSessions,
// 			&a.MessageCount,
// 		); err != nil {
// 			return nil, err
// 		}
// 		result = append(result, a)
// 	}

// 	return result, nil
// }
