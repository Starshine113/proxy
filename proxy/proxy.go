package proxy

import (
	"fmt"
	"strings"
	"time"

	"github.com/Starshine113/proxy/bot"
	"github.com/Starshine113/proxy/db"
	"github.com/bwmarrin/discordgo"
)

// Proxy handles proxying messages
type Proxy struct {
	Bot     *bot.Bot
	Session *discordgo.Session
}

// MessageCreate ...
func (p *Proxy) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var err error

	if m.GuildID == "" {
		return
	}

	if m.Author.Bot {
		return
	}

	if !p.Bot.Db.HasSystem(m.Author.ID) {
		return
	}

	var webhook *db.Webhook
	if !p.Bot.Db.HasWebhook(m.ChannelID) {
		webhook, err = p.Bot.Db.CreateWebhook(s, m.ChannelID)
		if err != nil {
			p.Bot.Sugar.Errorf("Error creating webhook in %v: %v", m.ChannelID, err)
			return
		}
	} else {
		webhook, err = p.Bot.Db.GetWebhook(m.ChannelID)
		if err != nil {
			p.Bot.Sugar.Errorf("Error getting webhook for %v: %v", m.ChannelID, err)
			return
		}
	}

	members, err := p.Bot.Db.GetAccountMembers(m.Author.ID)
	if err != nil {
		p.Bot.Sugar.Errorf("Error getting members for %v: %v", m.Author.ID, err)
	}

	for _, member := range members {
		if member.Prefix != "" {
			if strings.HasPrefix(m.Content, member.Prefix) {
				c := strings.TrimPrefix(m.Content, member.Prefix)
				if c != "" {
					system, err := p.Bot.Db.GetUserSystem(m.Author.ID)
					if err != nil {
						p.Bot.Sugar.Errorf("Error getting system for %v: %v", m.Author.ID, err)
						return
					}

					err = p.exec(webhook, member, system, c, m)
					if err != nil {
						p.Bot.Sugar.Errorf("Error executing proxy: %v", err)
					}
					break
				}
			}
		}
		if member.Suffix != "" {
			if strings.HasSuffix(m.Content, member.Suffix) {
				c := strings.TrimSuffix(m.Content, member.Suffix)
				if c != "" {
					system, err := p.Bot.Db.GetUserSystem(m.Author.ID)
					if err != nil {
						p.Bot.Sugar.Errorf("Error getting system for %v: %v", m.Author.ID, err)
						return
					}

					err = p.exec(webhook, member, system, c, m)
					if err != nil {
						p.Bot.Sugar.Errorf("Error executing proxy: %v", err)
					}
					break
				}
			}
		}
	}
}

func (p *Proxy) exec(w *db.Webhook, m *db.Member, s *db.System, content string, msg *discordgo.MessageCreate) (err error) {
	proxy, err := p.Session.WebhookExecute(w.ID, w.Token, true, &discordgo.WebhookParams{
		Content:   content,
		Username:  fmt.Sprintf("%v %v", m.DisplayedName(), s.Tag),
		AvatarURL: m.AvatarURL,
	})
	if err != nil {
		return err
	}

	err = p.Bot.Db.SaveMessage(proxy.ID, proxy.ChannelID, m.ID.String(), msg.Author.ID, msg.ID)
	if err != nil {
		return err
	}

	// wait a second before doing post-proxy tasks
	time.Sleep(time.Second)

	// delete original message
	err = p.Session.ChannelMessageDelete(msg.ChannelID, msg.ID)
	if err != nil {
		return
	}

	// ...

	return nil
}
