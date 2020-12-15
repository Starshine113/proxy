package proxy

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// ReactionAdd ...
func (p *Proxy) ReactionAdd(_ *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if !p.Bot.Db.MessageExists(r.MessageID) {
		return
	}

	m, err := p.Bot.Db.GetMessage(r.MessageID)
	if err != nil {
		p.Bot.Sugar.Errorf("Error getting message %v: %v", r.MessageID, err)
	}

	switch r.MessageReaction.Emoji.APIName() {
	case "❓", "❔":
		p.Session.MessageReactionRemove(m.ChannelID, m.ID, r.MessageReaction.Emoji.APIName(), r.UserID)
		c, err := p.Session.UserChannelCreate(r.UserID)
		if err != nil {
			p.Bot.Sugar.Errorf("Error creating user channel for %v: %v", r.UserID, err)
			return
		}

		member, err := p.Bot.Db.Member(m.Member.String())
		if err != nil {
			p.Bot.Sugar.Errorf("Error fetching member %s: %v", m.Member, err)
		}

		s, err := p.Bot.Db.GetUserSystem(m.Sender)
		if err != nil {
			p.Bot.Sugar.Errorf("Error getting system for %v: %v", m.Sender, err)
		}

		name := member.Name
		if s.Name != "" {
			name += " (" + s.Name + ")"
		}

		sys := fmt.Sprintf("`%s`", s.ID)
		if s.Name != "" {
			sys = fmt.Sprintf("%v (`%s`)", s.Name, s.ID)
		}
		tag := "N/A"
		if s.Tag != "" {
			tag = s.Tag
		}
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: member.AvatarURL,
				Name:    name,
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: member.AvatarURL,
			},
			Title:       "Message",
			Description: "Original account: " + m.Sender,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Display name",
					Value:  member.DisplayedName(),
					Inline: true,
				},
				{
					Name:   "System",
					Value:  sys,
					Inline: true,
				},
				{
					Name:   "Tag",
					Value:  tag,
					Inline: true,
				},
				{
					Name:   "Member",
					Value:  fmt.Sprintf("%v (`%s`)", member.Name, m.Member),
					Inline: false,
				},
			},
		}

		p.Session.ChannelMessageSendEmbed(c.ID, embed)
	case "❌":
		if r.UserID == m.Sender {
			err = p.Session.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				p.Bot.Sugar.Errorf("Error deleting message %v: %v", m.ID, err)
			}
		}
	case "🔔", "🛎️", "🏓", "❗", "❕":
		p.Session.MessageReactionRemove(m.ChannelID, m.ID, r.MessageReaction.Emoji.APIName(), r.UserID)
		_, err = p.Session.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: fmt.Sprintf("Psst, <@!%v>, you have been pinged by <@!%v>.", m.Sender, r.UserID),
			Embed: &discordgo.MessageEmbed{
				Description: fmt.Sprintf("[Jump to pinged message](https://discord.com/channels/%v/%v/%v)", r.GuildID, r.ChannelID, r.MessageID),
			},
		})
		if err != nil {
			p.Bot.Sugar.Errorf("Error sending ping message for %v: %v", r.MessageID, err)
		}
	}
}
