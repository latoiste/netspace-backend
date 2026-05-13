package db

import (
	"context"
	"fmt"

	"github.com/latoiste/netspace/model"
	"github.com/lib/pq"
)

func (e *Env) GetActiveSessions(ctx context.Context) ([]model.ActiveSession, error) {
	rows, err := e.db.QueryContext(ctx, `
		SELECT user_id, name, age, gender, table_number, interest, current_job, location_id, last_active_at
		FROM users
		WHERE status = 'logon'
		ORDER BY last_active_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []model.ActiveSession
	for rows.Next() {
		var s model.ActiveSession
		if err := rows.Scan(
			&s.UserID,
			&s.Name,
			&s.Age,
			&s.Gender,
			&s.TableNumber,
			pq.Array(&s.Interest),
			&s.CurrentJob,
			&s.LocationID,
			&s.LastActiveAt,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	return sessions, nil
}

func (e *Env) ForceLogoutUser(ctx context.Context, userID int) error {
	result, err := e.db.ExecContext(ctx, `
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

func (e *Env) GetAnalytics(ctx context.Context, from, to string) ([]model.AnalyticsData, error) {
	rows, err := e.db.QueryContext(ctx, `
		SELECT 
			DATE(created_at) as date,
			COUNT(DISTINCT user_id) as active_users,
			COUNT(DISTINCT session_token) as total_sessions,
			0 as message_count
		FROM users
		WHERE DATE(created_at) BETWEEN $1 AND $2
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.AnalyticsData
	for rows.Next() {
		var a model.AnalyticsData
		if err := rows.Scan(
			&a.Date,
			&a.ActiveUsers,
			&a.TotalSessions,
			&a.MessageCount,
		); err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	return result, nil
}
