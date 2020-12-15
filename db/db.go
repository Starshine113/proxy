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
	"fmt"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/Starshine113/proxy/structs"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

// Db ...
type Db struct {
	Pool       *pgxpool.Pool
	Sugar      *zap.SugaredLogger
	GuildCache *ttlcache.Cache
}

// Init ...
func Init(c *structs.BotConfig, sugar *zap.SugaredLogger) (db *Db, err error) {
	guildCache := ttlcache.NewCache()
	guildCache.SetCacheSizeLimit(100)
	guildCache.SetTTL(10 * time.Minute)

	pool, err := initDB(c)
	if err != nil {
		return nil, err
	}

	db = &Db{
		Pool:       pool,
		Sugar:      sugar,
		GuildCache: guildCache,
	}

	return
}

func initDB(config *structs.BotConfig) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), config.Auth.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}
	if err := initDBIfNotInitialised(db); err != nil {
		return nil, err
	}
	err = updateDB(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initDBIfNotInitialised(db *pgxpool.Pool) error {
	var exists bool
	err := db.QueryRow(context.Background(), "select exists (select from information_schema.tables where table_schema = 'public' and table_name = 'info')").Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return nil // the database has been initialised so we're done
	}

	// ...it's not initialised and we have to do that
	_, err = db.Exec(context.Background(), initDBSql)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully initialised the database.\n")
	return nil
}

func updateDB(db *pgxpool.Pool) (err error) {
	var dbVersion int
	err = db.QueryRow(context.Background(), "select schema_version from public.info").Scan(&dbVersion)
	if err != nil {
		return err
	}
	initialDBVersion := dbVersion
	for dbVersion < DBVersion {
		_, err = db.Exec(context.Background(), DBVersions[dbVersion-1])
		if err != nil {
			return err
		}
		dbVersion++
		fmt.Printf("Updated database to version %v\n", dbVersion)
	}
	if initialDBVersion < DBVersion {
		fmt.Printf("Successfully updated database to target version %v\n", DBVersion)
	}
	return nil
}
