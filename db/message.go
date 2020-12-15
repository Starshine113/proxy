package db

import (
	"context"

	"github.com/google/uuid"
)

// Message holds info about proxied messages
type Message struct {
	ID        string
	ChannelID string
	Member    uuid.UUID
	Sender    string
	Original  string
}

// SaveMessage saves a message to the database
func (db *Db) SaveMessage(msg, channel, member, sender, original string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "insert into public.messages (id, channel, member, sender, original_id) values ($1, $2, $3, $4, $5)", msg, channel, member, sender, original)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}

	return
}

// MessageExists checks if a message exists in the database
func (db *Db) MessageExists(msgID string) (b bool) {
	db.Pool.QueryRow(context.Background(), "select exists (select id from public.messages where id = $1 or original_id = $1)", msgID).Scan(&b)
	return b
}

// GetMessage gets info about a message
func (db *Db) GetMessage(msgID string) (m *Message, err error) {
	m = &Message{}

	err = db.Pool.QueryRow(context.Background(), "select id, channel, member, sender, original_id from public.messages where id = $1 or original_id = $1", msgID).Scan(&m.ID, &m.ChannelID, &m.Member, &m.Sender, &m.Original)
	return
}
