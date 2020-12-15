package router

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/Starshine113/proxy/bot"
	"github.com/Starshine113/proxy/proxy"
)

// PermLevel is the permission level of the command
type PermLevel int

const (
	// PermLevelNone can be used by all users
	PermLevelNone PermLevel = iota
	// PermLevelManager requires "Manage Server"
	PermLevelManager
	// PermLevelOwner requires the person to be the bot owner
	PermLevelOwner
)

// String gives the string representation of a permission level
func (p PermLevel) String() string {
	switch p {
	case PermLevelNone:
		return "None"
	case PermLevelManager:
		return "Manage Server"
	case PermLevelOwner:
		return "Bot Owner"
	}
	return fmt.Sprintf("PermLevel(%d)", p)
}

// Router is the command router
type Router struct {
	Commands []*Command
	Groups   []*Group

	BotOwners []string
	Cooldowns *ttlcache.Cache
	Bot       *bot.Bot

	Proxy *proxy.Proxy
}

// Command is a single command
type Command struct {
	Name    string
	Aliases []string
	Regex   *regexp.Regexp

	Description     string
	LongDescription string
	Usage           string

	Command func(*Ctx) error

	Permissions PermLevel
	GuildOnly   bool
	Cooldown    time.Duration

	Router *Router
}

// NewRouter creates a Router object
func NewRouter(b *bot.Bot, p *proxy.Proxy) *Router {
	cache := ttlcache.NewCache()
	cache.SkipTTLExtensionOnHit(true)

	router := &Router{
		BotOwners: b.Config.Bot.BotOwners,
		Bot:       b,
		Cooldowns: cache,
		Proxy:     p,
	}

	router.Bot.Session.AddHandler(router.messageCreate)

	router.AddCommand(&Command{
		Name:        "Commands",
		Description: "Show a list of commands",
		Usage:       "[command]",
		Permissions: PermLevelNone,
		Command:     router.dummy,
	})

	return router
}

// dummy is used when a command isn't handled with the normal process
func (r *Router) dummy(ctx *Ctx) error {
	return nil
}

// AddCommand adds a command to the router
func (r *Router) AddCommand(cmd *Command) {
	cmd.Router = r
	if cmd.Cooldown == 0 {
		cmd.Cooldown = 500 * time.Millisecond
	}
	r.Commands = append(r.Commands, cmd)
}

// GetCommand gets a command by name
func (r *Router) GetCommand(name string) (c *Command) {
	for _, cmd := range r.Commands {
		if strings.ToLower(cmd.Name) == strings.ToLower(name) {
			return cmd
		}
		for _, a := range cmd.Aliases {
			if strings.ToLower(a) == strings.ToLower(name) {
				return cmd
			}
		}
	}
	return nil
}
