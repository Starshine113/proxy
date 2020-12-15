package db

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
func (db *Db) GetMessage(msgID string) (m *Message, member *Member, s *System, err error) {
	m = &Message{}
	member = &Member{}
	s = &System{}

	err = db.Pool.QueryRow(context.Background(), "select m.id, m.channel, m.member, m.sender, m.original_id, s.id, s.name, s.tag, mem.id, mem.name, mem.display_name, mem.avatar_url from public.messages as m, public.members as mem, public.systems as s where (m.id = $1 or m.original_id = $1) and mem.id = m.member and s.id = mem.system", msgID).Scan(&m.ID, &m.ChannelID, &m.Member, &m.Sender, &m.Original, &s.ID, &s.Name, &s.Tag, &member.ID, &member.Name, &member.DisplayName, &member.AvatarURL)
	return
}
