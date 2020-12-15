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
	"errors"
	"time"

	"github.com/google/uuid"
)

// System holds the information for a system
type System struct {
	ID      uuid.UUID
	Name    string
	Tag     string
	Token   string
	Created time.Time
}

// Errors
var (
	ErrorNoRowsAffected = errors.New("no rows affected")
)

// HasSystem returns true if the specified user has a system, otherwise false
func (db *Db) HasSystem(userID string) (b bool) {
	db.Pool.QueryRow(context.Background(), "select exists (select id from public.systems where id = (select system from public.accounts where account = $1))", userID).Scan(&b)
	return b
}

// GetUserSystem gets a system by user ID
func (db *Db) GetUserSystem(userID string) (s *System, err error) {
	s = &System{}

	err = db.Pool.QueryRow(context.Background(), "select id, name, tag, token, created from public.systems where id = (select system from public.accounts where account = $1)", userID).Scan(&s.ID, &s.Name, &s.Tag, &s.Token, &s.Created)
	return
}

// GetSystem gets a system by UUID
func (db *Db) GetSystem(uuid string) (s *System, err error) {
	err = db.Pool.QueryRow(context.Background(), "select id, name, tag, token, created from public.systems where id = $1", uuid).Scan(&s.ID, &s.Name, &s.Tag, &s.Token, &s.Created)
	return
}

// CreateSystem creates a system for the given user ID, with the given name
func (db *Db) CreateSystem(userID, name string) (s *System, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}

	s = &System{}

	err = db.Pool.QueryRow(context.Background(), "insert into public.systems (id, name) values ($1, $2) returning id, name, tag, token, created", id, name).Scan(&s.ID, &s.Name, &s.Tag, &s.Token, &s.Created)
	if err != nil {
		return
	}

	commandTag, err := db.Pool.Exec(context.Background(), "insert into public.accounts (account, system) values ($1, $2)", userID, id)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return s, ErrorNoRowsAffected
	}
	return
}

// GetSystemUsers gets all user IDs associated with a system
func (db *Db) GetSystemUsers(uuid string) (users []string, err error) {
	err = db.Pool.QueryRow(context.Background(), "select array(select account from public.accounts where system = $1)", uuid).Scan(&users)
	return
}

// SetTag sets the tag for a system
func (db *Db) SetTag(uuid, tag string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.systems set tag = $1 where id = $2", tag, uuid)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}

// SetName sets the name for a system
func (db *Db) SetName(uuid, name string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.systems set name = $1 where id = $2", name, uuid)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}

// DeleteSystem deletes a system
func (db *Db) DeleteSystem(uuid string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "delete from public.systems where id = $1", uuid)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}
