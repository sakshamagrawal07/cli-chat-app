package db

import (
	"log"

	"github.com/sakshamagrawal07/cli-chat-app.git/shared/models"
)

func InsertMessage(msg models.Message) error {
	query := `
		INSERT INTO messages (sender_username, recipient_username, message, channel)
		VALUES ($1, $2, $3, $4)
	`

	_, err := DB.ExecContext(Ctx, query,
		msg.SenderUsername,
		msg.RecipientUsername,
		msg.Message,
		msg.Channel,
	)

	return err
}

func ReadMessages(channel string) ([]models.Message, error) {
	query := `SELECT sender_username, recipient_username, message, channel
		FROM messages
	`

	rows, err := DB.QueryContext(Ctx, query)
	log.Println("[DEBUG][ReadMessages] Executing query:", query)
	log.Println("[DEBUG][ReadMessages] Rows : ", rows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.SenderUsername, &msg.RecipientUsername, &msg.Message, &msg.Channel); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
