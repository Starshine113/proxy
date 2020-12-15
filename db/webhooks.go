package db

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

// Webhook holds info for a webhook
type Webhook struct {
	ChannelID string
	ID        string
	Token     string
}

// HasWebhook checks if the channel has a webhook
func (db *Db) HasWebhook(channel string) (b bool) {
	db.Pool.QueryRow(context.Background(), "select exists (select webhook from public.webhooks where channel = $1)", channel).Scan(&b)
	return b
}

// CreateWebhook creates a webhook and adds it to the database
func (db *Db) CreateWebhook(s *discordgo.Session, channel string) (w *Webhook, err error) {
	wh, err := s.WebhookCreate(channel, "Proxy Webhook", "")
	if err != nil {
		return
	}

	commandTag, err := db.Pool.Exec(context.Background(), "insert into public.webhooks (channel, webhook, token) values ($1, $2, $3)", channel, wh.ID, wh.Token)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return w, ErrorNoRowsAffected
	}
	return &Webhook{
		ChannelID: channel,
		ID:        wh.ID,
		Token:     wh.Token,
	}, nil
}

// GetWebhook gets the webhook for a channel
func (db *Db) GetWebhook(channel string) (w *Webhook, err error) {
	w = &Webhook{}

	err = db.Pool.QueryRow(context.Background(), "select channel, webhook, token from public.webhooks where channel = $1", channel).Scan(&w.ChannelID, &w.ID, &w.Token)
	return
}
