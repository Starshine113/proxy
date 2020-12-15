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
	"time"

	"github.com/google/uuid"
)

// Member is a single member
type Member struct {
	ID          uuid.UUID
	System      uuid.UUID
	AvatarURL   string
	Name        string
	DisplayName string
	Prefix      string
	Suffix      string
	Created     time.Time
}

// DisplayedName returns the displayed name--either display name or normal name
func (m Member) DisplayedName() string {
	if m.DisplayName != "" {
		return m.DisplayName
	}
	return m.Name
}

// Member returns a member by UUID
func (db *Db) Member(uuid string) (m *Member, err error) {
	m = &Member{}

	err = db.Pool.QueryRow(context.Background(), "select id, system, avatar_url, name, display_name, prefix, suffix, created from public.members where id = $1", uuid).Scan(&m.ID, &m.System, &m.AvatarURL, &m.Name, &m.DisplayName, &m.Prefix, &m.Suffix, &m.Created)
	return m, err
}

// GetSystemMembers gets all members for a system by UUID
func (db *Db) GetSystemMembers(uuid string) (members []*Member, err error) {
	members = make([]*Member, 0)

	rows, err := db.Pool.Query(context.Background(), "select id, system, avatar_url, name, display_name, prefix, suffix, created from public.members where system = $1", uuid)
	if err != nil {
		return
	}

	for rows.Next() {
		m := &Member{}

		rows.Scan(&m.ID, &m.System, &m.AvatarURL, &m.Name, &m.DisplayName, &m.Prefix, &m.Suffix, &m.Created)
		members = append(members, m)
	}
	return members, err
}

// GetAccountMembers gets all members for a system by user ID
func (db *Db) GetAccountMembers(id string) (members []*Member, err error) {
	members = make([]*Member, 0)

	rows, err := db.Pool.Query(context.Background(), "select id, system, avatar_url, name, display_name, prefix, suffix, created from public.members where system = (select system from public.accounts where account = $1)", id)
	if err != nil {
		return
	}

	for rows.Next() {
		m := &Member{}

		rows.Scan(&m.ID, &m.System, &m.AvatarURL, &m.Name, &m.DisplayName, &m.Prefix, &m.Suffix, &m.Created)
		members = append(members, m)
	}
	return members, err
}

// NewMember creates a new member for the given system
func (db *Db) NewMember(system, name string) (m *Member, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}

	m = &Member{
		ID:   id,
		Name: name,
	}

	err = db.Pool.QueryRow(context.Background(), "insert into public.members (id, name, system) values ($1, $2, $3) returning created", id, name, system).Scan(&m.Created)
	return
}

// SetProxy sets a member's proxy
func (db *Db) SetProxy(uuid, prefix, suffix string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.members set prefix = $1, suffix = $2 where id = $3", prefix, suffix, uuid)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}

// SetAvatar sets the avatar for a member
func (db *Db) SetAvatar(uuid, url string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.members set avatar_url = $1 where id = $2", url, uuid)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}

// SetDisplayName sets the display name for a member
func (db *Db) SetDisplayName(uuid, name string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.members set display_name = $1 where id = $2", name, uuid)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}

// DeleteMember deletes a member
func (db *Db) DeleteMember(uuid string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "delete from public.members where id = $1", uuid)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}