package db

// DBVersion is the current database version
const DBVersion = 1

// DBVersions is a slice of schemas for every database version
var DBVersions []string = []string{""}

// initDBSql is the initial SQL database schema
var initDBSql = `create table if not exists systems
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

    proxy_enabled bool not null default true,

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
