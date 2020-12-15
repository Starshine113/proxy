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

// DBVersion is the current database version
const DBVersion = 1

// DBVersions is a slice of schemas for every database version
var DBVersions []string = []string{""}

// initDBSql is the initial SQL database schema
var initDBSql = `create type autoproxy_setting as enum ('off', 'latch', 'front', 'member');

create table if not exists systems
(
    id      uuid        not null primary key,
    name    text        not null default '',
    tag     text        not null default '',
    token   text        not null default '',
    created timestamp   not null default (current_timestamp at time zone 'utc')
);

create table if not exists system_guilds
(
    system  uuid        references systems (id) on delete cascade,
    guild   text        not null,

    proxy_enabled		bool				not null default true,
	autoproxy_mode		autoproxy_setting	not null default 'off';
	last_proxied_member text				not null default ''; 

    primary key (system, guild)
);

create table if not exists members
(
    id           uuid           not null primary key,
    system       uuid           not null references systems (id) on delete cascade,
    avatar_url   text           not null default '',
    name         text           not null,
    display_name text           not null default '',
    prefix       text           not null default '';
    suffix       text           not null default '';
    created      timestamp      not null default (current_timestamp at time zone 'utc')
);

create table if not exists accounts
(
    account text    not null primary key,
    system  uuid    not null references systems (id) on delete cascade
);

create table if not exists messages
(
    id          text    primary key,
    channel     text    not null,
    member      uuid    not null references members (id) on delete cascade,
    sender      text    not null,
    original_id text    not null
);

create table if not exists webhooks
(
    channel text    primary key,
    webhook text    not null,
    token   text    not null
);

create table if not exists servers
(
    id            text      primary key,
    log_channel   text,
    log_blacklist text[]    not null default array[]::text[],
    blacklist     text[]    not null default array[]::text[] 
);

create table if not exists info
(
    id                      int primary key not null default 1, -- enforced only equal to 1
    schema_version          int,
    constraint singleton    check (id = 1) -- enforce singleton table/row
);

insert into info (schema_version) values (1);`
