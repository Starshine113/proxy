package router

import "github.com/bwmarrin/discordgo"

// MessageCreate ...
func (r *Router) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var err error

	// if message was sent by a bot return; not only to ignore bots, but also to make sure PluralKit users don't trigger commands twice.
	if m.Author.Bot {
		return
	}

	ctx, err := Context(r.Bot.Config.Bot.Prefixes, m.Content, m, r.Bot)
	if err != nil {
		r.Bot.Sugar.Errorf("Error getting context for %v: %v", m.ID, err)
		return
	}

	// check if the message might be a command
	if ctx.MatchPrefix() {
		r.execute(ctx)
		return
	}

	r.Proxy.MessageCreate(s, m)
}

func (r *Router) execute(ctx *Ctx) {
	err := r.Execute(ctx)
	if err != nil {
		r.Bot.Sugar.Errorf("Error running command %v: %v", ctx.Message.ID, err)
	}
}
