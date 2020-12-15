package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) Ready(_ *discordgo.Session, _ *discordgo.Ready) {
	err := b.Session.UpdateStatus(0, fmt.Sprintf("%vhelp | in %v servers", b.Prefix, len(b.Session.State.Guilds)))
	if err != nil {
		b.Sugar.Errorf("Error updating status: %v", err)
	}
}
