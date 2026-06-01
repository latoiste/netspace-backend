package db

import (
	"context"
	"log"

	"github.com/latoiste/netspace/api"
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

func (r *Repository) AllPrivateMessages(senderId string, ctx context.Context) ([]api.MessageDTO, error) {
	const query = `
		SELECT DISTINCT ON (chat_partner)
			messageid,
			chat_partner,
			"message",
			"timestamp"
		FROM (
			SELECT
				messageid,
				CASE
					WHEN senderId = $1 THEN recipientId
					ELSE senderId
				END AS chat_partner,
				"message",
				"timestamp"
			FROM PrivateMessages
			WHERE senderId = $1 OR recipientId = $1
		) t
		ORDER BY chat_partner, "timestamp" DESC;
	`

	messageDTOs := make([]api.MessageDTO, 0)

	rows, err := r.db.QueryContext(ctx, query, senderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message model.PrivateMessage
		err = rows.Scan(
			&message.MessageId,
			&message.RecipientId,
			&message.Message,
			&message.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		recipient, err := r.UserById(message.RecipientId, ctx)
		if err != nil {
			return nil, err
		}
		messageDTO := api.ConstructPrivateMessageDTO(message, *recipient)

		messageDTOs = append(messageDTOs, messageDTO)
	}

	return messageDTOs, nil
}

func (r *Repository) AllGroupMessages(senderId string, ctx context.Context) ([]api.MessageDTO, error) {
	const query = `
		SELECT DISTINCT ON (gm.groupId)
			gm.groupId,
			gm.messageId,
			gm.message,
			gm.timestamp
		FROM GroupMessages gm
		JOIN GroupMembers m
			ON m.groupId = gm.groupId
		WHERE m.userId = $1
		ORDER BY gm.groupId, gm.timestamp DESC;
	`

	messageDTOs := make([]api.MessageDTO, 0)

	rows, err := r.db.QueryContext(ctx, query, senderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message model.GroupMessage
		err = rows.Scan(
			&message.GroupId,
			&message.MessageId,
			&message.Message,
			&message.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		group, err := r.GroupById(message.GroupId, ctx)
		if err != nil {
			return nil, err
		}

		messageDTO := api.ConstructGroupMessageDTO(message, *group)

		messageDTOs = append(messageDTOs, messageDTO)
	}

	return messageDTOs, nil
}
