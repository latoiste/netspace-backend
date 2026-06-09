package db

import (
	"context"
	"database/sql"
	"log"
	"time"

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

// PrivateMessagesBetween returns the full DM history between two users,
// oldest first, so the chat page can render past messages on load. isread is
// included so the sender's own bubbles can show the right (grey/blue) receipt.
func (r *Repository) PrivateMessagesBetween(userId string, partnerId string, ctx context.Context) ([]model.PrivateMessage, error) {
	const query = `
		SELECT messageid, senderId, "message", "timestamp", isread
		FROM PrivateMessages
		WHERE (senderId = $1 AND recipientId = $2)
		   OR (senderId = $2 AND recipientId = $1)
		ORDER BY "timestamp" ASC;
	`

	messages := make([]model.PrivateMessage, 0)

	rows, err := r.db.QueryContext(ctx, query, userId, partnerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m model.PrivateMessage
		if err := rows.Scan(&m.MessageId, &m.SenderId, &m.Message, &m.Timestamp, &m.IsRead); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}

// MarkPrivateMessagesRead flips every message sent BY senderId TO recipientId to
// read. Called when recipientId opens the chat with senderId; the sender's
// bubbles then show blue double-checks. Returns the number of rows changed so
// the caller can skip notifying the sender when nothing actually changed.
func (r *Repository) MarkPrivateMessagesRead(senderId string, recipientId string, ctx context.Context) (int64, error) {
	const query = `
		UPDATE PrivateMessages
		SET isread = TRUE
		WHERE senderId = $1 AND recipientId = $2 AND isread = FALSE;
	`

	result, err := r.db.ExecContext(ctx, query, senderId, recipientId)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// GroupMessagesByGroupId returns the full message history for a group,
// oldest first.
func (r *Repository) GroupMessagesByGroupId(groupId string, ctx context.Context) ([]model.GroupMessage, error) {
	const query = `
		SELECT messageId, senderId, groupId, "message", "timestamp"
		FROM GroupMessages
		WHERE groupId = $1
		ORDER BY "timestamp" ASC;
	`

	messages := make([]model.GroupMessage, 0)

	rows, err := r.db.QueryContext(ctx, query, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m model.GroupMessage
		if err := rows.Scan(&m.MessageId, &m.SenderId, &m.GroupId, &m.Message, &m.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
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

// RecentPublicMessages returns up to `limit` of the most recent public-room
// messages for a location that are newer than `since`, oldest-first so the chat
// renders top-to-bottom. A LEFT JOIN keeps messages from people who have since
// logged out (their Users row is retained, just marked inactive), so the shared
// timeline stays intact for everyone still in the room.
func (r *Repository) RecentPublicMessages(locationId int, limit int, since time.Time, ctx context.Context) ([]api.PublicChatMessageDTO, error) {
	const query = `
		SELECT pm.messageid, pm.senderId, COALESCE(u.name, ''), pm."message", pm."timestamp"
		FROM PublicMessages pm
		LEFT JOIN Users u ON u.id = pm.senderId
		WHERE pm.locationId = $1 AND pm."timestamp" >= $2
		ORDER BY pm."timestamp" DESC
		LIMIT $3;
	`

	rows, err := r.db.QueryContext(ctx, query, locationId, since, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]api.PublicChatMessageDTO, 0, limit)
	for rows.Next() {
		var id, senderId, name, message string
		var ts time.Time
		if err := rows.Scan(&id, &senderId, &name, &message, &ts); err != nil {
			return nil, err
		}
		messages = append(messages, api.PublicChatMessageDTO{
			Id:          id,
			SenderId:    senderId,
			SenderName:  name,
			SenderEmoji: api.EmojiForUser(senderId),
			Message:     message,
			Timestamp:   ts.Local().Format("15:04"),
		})
	}

	// Query is newest-first (so LIMIT keeps the latest); flip to oldest-first.
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// DeletePublicMessagesBefore is the public-room retention sweep: it removes
// messages older than the cutoff so the shared timeline stays bounded and fresh
// to the venue's current crowd. Returns how many rows were removed.
func (r *Repository) DeletePublicMessagesBefore(cutoff time.Time, ctx context.Context) (int64, error) {
	const query = `DELETE FROM PublicMessages WHERE "timestamp" < $1;`

	result, err := r.db.ExecContext(ctx, query, cutoff)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// GroupsForUser returns every ACTIVE group the user is a member of — including
// groups with no messages yet — each with its latest message (if any). This is
// what keeps a group in the chat list while the user is still a member, even
// after they navigate away (Home); it only drops out when they leave the group
// (membership removed) or it dissolves (isActive = false).
func (r *Repository) GroupsForUser(userId string, ctx context.Context) ([]api.MessageDTO, error) {
	const query = `
		SELECT g.id, g.name, lm."message", lm."timestamp"
		FROM GroupMembers gm
		JOIN Groups g ON g.id = gm.groupId AND g.isActive = true
		LEFT JOIN LATERAL (
			SELECT m."message", m."timestamp"
			FROM GroupMessages m
			WHERE m.groupId = g.id
			ORDER BY m."timestamp" DESC
			LIMIT 1
		) lm ON true
		WHERE gm.userId = $1
		ORDER BY lm."timestamp" DESC NULLS LAST;
	`

	dtos := make([]api.MessageDTO, 0)

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		var lastMsg sql.NullString
		var ts sql.NullTime
		if err := rows.Scan(&id, &name, &lastMsg, &ts); err != nil {
			return nil, err
		}
		dtos = append(dtos, api.ConstructGroupSummaryDTO(
			id, name, lastMsg.String, lastMsg.Valid, ts.Time, ts.Valid,
		))
	}

	return dtos, nil
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
