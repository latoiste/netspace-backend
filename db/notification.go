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
			secondarylabel
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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
	)

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
