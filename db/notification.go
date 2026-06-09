package db

import (
	"context"

	"github.com/latoiste/netspace/model"
)

func (r *Repository) InsertNotification(notif model.Notification, userId string, ctx context.Context) error {
	query := `
		INSERT INTO Notifications (
			id,
			userid,
			"type",
			emoji,
			avatargradient,
			title,
			description,
			"timestamp",
			unread,
			primarylabel,
			secondarylabel,
			groupid,
			senderid
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		notif.Id,
		userId,
		notif.Type,
		notif.Emoji,
		notif.AvatarGradient,
		notif.Title,
		notif.Description,
		notif.Timestamp,
		notif.Unread,
		notif.PrimaryLabel,
		notif.SecondaryLabel,
		notif.GroupId,
		notif.SenderId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) NotificationByUserId(userId string, ctx context.Context) ([]model.Notification, error) {
	query := `
		SELECT 
			"id",
			userid,
			"type",
			emoji,
			avatargradient,
			title,
			description,
			"timestamp",
			unread,
			COALESCE(primaryLabel, ''),
			COALESCE(secondaryLabel, ''),
			COALESCE(groupid, ''),
			COALESCE(senderid, '')
		FROM notifications
		WHERE userid = $1
	`

	notifications := make([]model.Notification, 0)

	rows, err := r.db.QueryContext(
		ctx,
		query,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notif model.Notification
		err = rows.Scan(
			&notif.Id,
			&notif.UserId,
			&notif.Type,
			&notif.Emoji,
			&notif.AvatarGradient,
			&notif.Title,
			&notif.Description,
			&notif.Timestamp,
			&notif.Unread,
			&notif.PrimaryLabel,
			&notif.SecondaryLabel,
			&notif.GroupId,
			&notif.SenderId,
		)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notif)
	}
	return notifications, nil
}

// DeleteNotification removes a single notification, but only if it belongs to
// the requester (ownership check) — so tapping a notification dismisses it for
// good instead of it reappearing on the next /api/notifications fetch.
func (r *Repository) DeleteNotification(notificationId string, userId string, ctx context.Context) error {
	const query = `DELETE FROM notifications WHERE id = $1 AND userid = $2`

	_, err := r.db.ExecContext(ctx, query, notificationId, userId)
	if err != nil {
		return err
	}

	return nil
}

// DeleteGroupInviteNotif removes a user's group-invite notification(s) for a
// given group, once they've accepted or rejected it — so a resolved invite
// disappears for good instead of reappearing (with stale Accept/Reject buttons)
// on the next /api/notifications fetch.
func (r *Repository) DeleteGroupInviteNotif(userId string, groupId string, ctx context.Context) error {
	const query = `
		DELETE FROM notifications
		WHERE userid = $1 AND "type" = 'group_invite' AND groupid = $2
	`

	_, err := r.db.ExecContext(ctx, query, userId, groupId)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateNotificationUnread(notificationId string, unread bool, ctx context.Context) error {
	query := `
		UPDATE notifications
		SET unread = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		unread,
		notificationId,
	)

	if err != nil {
		return err
	}

	return nil
}
