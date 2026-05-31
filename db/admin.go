package db

import (
	"context"
	"fmt"
	"time"

	"github.com/latoiste/netspace/api"
)

// func (r *Repository) GetActiveSessions(ctx context.Context) (int, error) {

// }

func (r *Repository) ForceLogoutUser(ctx context.Context, userID int) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET status = 'logoff'
		WHERE user_id = $1 AND status = 'logon'
	`, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user %d not found or already logged out", userID)
	}

	return nil
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
