package db

import (
	"context"
	"github.com/google/uuid"
)

// GuildSystem holds info on a system in a specific server
type GuildSystem struct {
	GuildID           string
	SystemID          uuid.UUID
	ProxyEnabled      bool
	AutoproxyMode     string
	LastProxiedMember string
}

// Constants for autoproxy modes
const (
	AutoproxyModeOff    = "off"
	AutoproxyModeLatch  = "latch"
	AutoproxyModeFront  = "front"
	AutoproxyModeMember = "member"
)

// GetGuildSystem ...
func (db *Db) GetGuildSystem(userID, guildID string) (g *GuildSystem, err error) {
	g = &GuildSystem{}
	var exists bool
	err = db.Pool.QueryRow(context.Background(), "select exists (select guild from public.system_guilds where guild = $1 and system = (select system from accounts where account = $2))", guildID, userID).Scan(&exists)
	if err != nil {
		return
	}

	if !exists {
		commandTag, err := db.Pool.Exec(context.Background(), "insert into public.system_guilds (guild, system) values ($1, (select system from accounts where account = $2))", guildID, userID)
		if err != nil {
			return g, err
		}
		if commandTag.RowsAffected() != 1 {
			return g, ErrorNoRowsAffected
		}
	}

	err = db.Pool.QueryRow(context.Background(), "select guild, system, proxy_enabled, autoproxy_mode, last_proxied_member from public.system_guilds where guild = $1 and system = (select system from accounts where account = $2)", guildID, userID).Scan(&g.GuildID, &g.SystemID, &g.ProxyEnabled, &g.AutoproxyMode, &g.LastProxiedMember)
	return
}

// SetGuildProxy sets the proxying status for a system
func (db *Db) SetGuildProxy(systemID, guildID string, enable bool) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.system_guilds set proxy_enabled = $1 where system = $2 and guild = $3", enable, systemID, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}