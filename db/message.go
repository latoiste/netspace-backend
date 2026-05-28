package db

import (
	"context"
	"log"

	"github.com/latoiste/netspace/model"
)

func (r *Repository) InsertPrivateMessage(privateMsg model.PrivateMessage, ctx context.Context) error {
	query := `
		INSERT INTO PrivateMessages (
			messageid,
			locationId,
			senderId,
			recipientId,
			"message",
			"timestamp"
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		privateMsg.MessageId,
		privateMsg.LocationId,
		privateMsg.SenderId,
		privateMsg.RecipientId,
		privateMsg.Message,
		privateMsg.Timestamp,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) InsertPublicMessage(publicMsg model.PublicMessage, ctx context.Context) error {
	query := `
		INSERT INTO PublicMessages (
			messageId,
			locationId,
			senderId,
			"message",
			"timestamp"
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		publicMsg.MessageId,
		publicMsg.LocationId,
		publicMsg.SenderId,
		publicMsg.Message,
		publicMsg.Timestamp,
	)

	log.Println(result.RowsAffected())

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) InsertGroupMessage(groupMsg model.GroupMessage, ctx context.Context) error {
	query := `
		INSERT INTO GroupMessages (
			messageid,
			locationid,
			senderid,
			groupid,
			message,
			timestamp
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		groupMsg.MessageId,
		groupMsg.LocationId,
		groupMsg.SenderId,
		groupMsg.GroupId,
		groupMsg.Message,
		groupMsg.Timestamp,
	)

	if err != nil {
		return err
	}

	return nil
}
