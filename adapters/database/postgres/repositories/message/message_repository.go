package message_repository

import (
	"context"

	message_entity "github.com/LeonardoSola/wagon-api/internal/entities/message"
	"github.com/LeonardoSola/wagon-api/internal/logger"
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
	dml := `INSERT INTO messages (id, external_id, device_id, text, recipient, recipient_type, created_at) VALUES (@id, @external_id, @device_id, @text, @recipient, @recipient_type, @created_at)`

	logger.Log.Debug("MessageRepository.Create", "message.ID", message.ID)
	logger.Log.Debug("MessageRepository.Create", "message.ExternalID", message.ExternalID)
	logger.Log.Debug("MessageRepository.Create", "message.DeviceID", message.DeviceID)
	logger.Log.Debug("MessageRepository.Create", "message.DeviceJID", message.DeviceJID)
	logger.Log.Debug("MessageRepository.Create", "message.Text", message.Text)
	logger.Log.Debug("MessageRepository.Create", "message.Recipient", message.Recipient)
	logger.Log.Debug("MessageRepository.Create", "message.RecipientType", message.RecipientType)
	logger.Log.Debug("MessageRepository.Create", "message.CreatedAt", message.CreatedAt)

	_, err := r.db.Exec(ctx, dml, pgx.NamedArgs{
		"id":             message.ID,
		"external_id":    message.ExternalID,
		"device_id":      message.DeviceID,
		"text":           message.Text,
		"recipient":      message.Recipient,
		"recipient_type": message.RecipientType,
		"created_at":     message.CreatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *MessageRepository) GetAll(ctx context.Context) ([]message_entity.Message, error) {
	dml := `SELECT * FROM messages`

	rows, err := r.db.Query(ctx, dml)
	if err != nil {
		return []message_entity.Message{}, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[message_entity.Message])
}

func (r *MessageRepository) Get(ctx context.Context, messageID string) (message_entity.Message, error) {
	dml := `SELECT * FROM messages WHERE id = $1 LIMIT 1`

	rows, err := r.db.Query(ctx, dml, messageID)
	if err != nil {
		return message_entity.Message{}, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[message_entity.Message])
}

func (r *MessageRepository) Delete(ctx context.Context, messageID string) error {
	dml := `UPDATE messages SET deleted_at = NOW() WHERE id = $1`

	_, err := r.db.Exec(ctx, dml, messageID)
	if err != nil {
		return err
	}

	return nil
}

func (r *MessageRepository) Update(ctx context.Context, message message_entity.Message) error {
	dml := `UPDATE messages SET text = $1, recipient = $2, recipient_type = $3, updated_at = NOW() WHERE id = $4`

	_, err := r.db.Exec(ctx, dml, message.Text, message.Recipient, message.RecipientType, message.ID)
	if err != nil {
		return err
	}

	return nil
}
