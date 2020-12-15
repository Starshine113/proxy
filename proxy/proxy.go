package proxy

// Discord proxy bot
// Copyright (C) 2020  Starshine System

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"codeberg.org/eviedelta/dwhook"
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

	gs, err := p.Bot.Db.GetGuildSystem(m.Author.ID, m.GuildID)
	if err != nil {
		p.Bot.Sugar.Errorf("Error getting guild/system settings: %v", err)
		return
	}

	if !gs.ProxyEnabled {
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
				if c == "" && len(m.Attachments) == 0 {
					continue
				}
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
		if member.Suffix != "" {
			if strings.HasSuffix(m.Content, member.Suffix) {
				c := strings.TrimSuffix(m.Content, member.Suffix)
				if c == "" && len(m.Attachments) == 0 {
					continue
				}
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

	// try autoproxy
	if gs.AutoproxyMode == db.AutoproxyModeOff {
		return
	}

	if gs.AutoproxyMode == db.AutoproxyModeLatch {
		if gs.LastProxiedMember == "" {
			return
		}

		member, err := p.Bot.Db.Member(gs.LastProxiedMember)
		if err != nil {
			p.Bot.Sugar.Errorf("Error fetching latch member %v: %v", gs.LastProxiedMember, err)
			return
		}

		system, err := p.Bot.Db.GetUserSystem(m.Author.ID)
		if err != nil {
			p.Bot.Sugar.Errorf("Error getting system for %v: %v", m.Author.ID, err)
			return
		}

		err = p.exec(webhook, member, system, m.Content, m)
		if err != nil {
			p.Bot.Sugar.Errorf("Error executing proxy: %v", err)
		}
	}
}

func (p *Proxy) exec(w *db.Webhook, m *db.Member, s *db.System, content string, msg *discordgo.MessageCreate) (err error) {
	perms, err := p.Session.State.UserChannelPermissions(msg.Author.ID, msg.ChannelID)
	if err == discordgo.ErrStateNotFound {
		perms, err = p.Session.UserChannelPermissions(msg.Author.ID, msg.ChannelID)
	}
	if err != nil {
		return err
	}

	var mentions *discordgo.MessageAllowedMentions
	// if the user can mention @everyone/@here, mirror that permission
	if perms&discordgo.PermissionMentionEveryone == discordgo.PermissionMentionEveryone {
		mentions = &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeRoles, discordgo.AllowedMentionTypeUsers, discordgo.AllowedMentionTypeEveryone},
		}
	} else {
		mentions = &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
			Roles: make([]string, 0),
		}
		roles, err := p.Session.GuildRoles(msg.GuildID)
		if err == nil {
			for _, r := range roles {
				if r.Mentionable {
					mentions.Roles = append(mentions.Roles, r.ID)
				}
			}
		}
	}
	embeds := make([]*discordgo.MessageEmbed, 0)
	if perms&discordgo.PermissionEmbedLinks == discordgo.PermissionEmbedLinks {
		if len(msg.Embeds) != 0 {
			embeds = msg.Embeds
		}
	}

	var proxy *discordgo.Message
	if len(msg.Attachments) == 0 {
		proxy, err = p.Session.WebhookExecute(w.ID, w.Token, true, &discordgo.WebhookParams{
			Content:         content,
			Username:        fmt.Sprintf("%v %v", m.DisplayedName(), s.Tag),
			AvatarURL:       m.AvatarURL,
			AllowedMentions: mentions,
			Embeds:          embeds,
		})
		if err != nil {
			return err
		}
	} else if len(msg.Attachments) == 1 {
		var mentions dwhook.AllowedMentions
		// if the user can mention @everyone/@here, mirror that permission
		if perms&discordgo.PermissionMentionEveryone == discordgo.PermissionMentionEveryone {
			mentions = dwhook.AllowedMentions{
				Parse: []string{dwhook.MentionRoles, dwhook.MentionUsers, dwhook.MentionEveryone},
			}
		} else {
			mentions = dwhook.AllowedMentions{
				Parse: []string{dwhook.MentionUsers},
				Roles: make([]string, 0),
			}
			roles, err := p.Session.GuildRoles(msg.GuildID)
			if err == nil {
				for _, r := range roles {
					if r.Mentionable {
						mentions.Roles = append(mentions.Roles, r.ID)
					}
				}
			}
		}

		resp, err := http.Get(msg.Attachments[0].URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)

		if len(b) > 8*1000*1000 {
			_, err = p.Bot.Session.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("❌ This message cannot be proxied, its attachment is too large (%v MB > 8 MB).\nUnfortunately webhooks aren't considered to have Discord Nitro :(", int(len(b)/1024/1024)))
			return err
		}

		b, err = dwhook.SendFileToToken(w.ID, w.Token, dwhook.Message{
			Content:         content,
			Username:        fmt.Sprintf("%v %v", m.DisplayedName(), s.Tag),
			AvatarURL:       m.AvatarURL,
			AllowedMentions: mentions,
		}, msg.Attachments[0].Filename, b)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(b, &proxy); err != nil {
			return err
		}
	} else {
		msg, err := p.Session.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("❌ Unfortunately, %v doesn't support proxying multiple attachments in one file, %v :(", p.Bot.Session.State.User.Username, msg.Author.Mention()))
		if err != nil {
			return err
		}
		time.Sleep(10 * time.Second)
		return p.Session.ChannelMessageDelete(msg.ChannelID, msg.ID)
	}

	err = p.Bot.Db.SaveMessage(proxy.ID, proxy.ChannelID, m.ID.String(), msg.Author.ID, msg.ID)
	if err != nil {
		return err
	}

	// update last-proxied system/guild member
	err = p.Bot.Db.SetLastProxiedMember(s.ID.String(), msg.GuildID, m.ID.String())
	if err != nil {
		p.Bot.Sugar.Errorf("Error saving last-proxied member: %v", err)
		return nil
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
