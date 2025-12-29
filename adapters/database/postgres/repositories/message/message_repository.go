package message_repository

import (
	"context"

	message_entity "github.com/LeonardoSola/wagon-api/internal/entities/message"
	"github.com/LeonardoSola/wagon-api/internal/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) ports.MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, message message_entity.Message) error {
	dml := `INSERT INTO messages (id, external_id, device_id, text, recipient, recipient_type, timestamp) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(ctx, dml, message.ID, message.ExternalID, message.DeviceID, message.Text, message.Recipient, message.RecipientType, message.Timestamp)
	if err != nil {
		return err
	}

	return nil
}

func (r *MessageRepository) Get(ctx context.Context, messageID string) (message_entity.Message, error) {
	dml := `SELECT * FROM messages WHERE id = $1`

	rows, err := r.db.Query(ctx, dml, messageID)
	if err != nil {
		return message_entity.Message{}, err
	}

	message, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[message_entity.Message])
	if err != nil {
		return message_entity.Message{}, err
	}

	return message, nil
}
