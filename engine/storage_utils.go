/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package engine

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/utils"
)

// Various helpers to deal with database

func ConfigureDataStorage(db_type, host, port, name, user, pass, marshaler string, cacheCfg config.CacheConfig, loadHistorySize int) (db DataDB, err error) {
	var d DataDB
	switch db_type {
	case utils.REDIS:
		var db_nb int
		db_nb, err = strconv.Atoi(name)
		if err != nil {
			utils.Logger.Crit("Redis db name must be an integer!")
			return nil, err
		}
		if port != "" {
			host += ":" + port
		}
		d, err = NewRedisStorage(host, db_nb, pass, marshaler, utils.REDIS_MAX_CONNS, cacheCfg, loadHistorySize)
	case utils.MONGO:
		d, err = NewMongoStorage(host, port, name, user, pass, utils.DataDB, nil, cacheCfg, loadHistorySize)
		db = d.(DataDB)
	default:
		err = errors.New(fmt.Sprintf("Unknown db '%s' valid options are '%s' or '%s'",
			db_type, utils.REDIS, utils.MONGO))
	}
	if err != nil {
		return nil, err
	}
	return d, nil
}

func ConfigureStorStorage(db_type, host, port, name, user, pass, marshaler string, maxConn, maxIdleConn, connMaxLifetime int, cdrsIndexes []string) (db Storage, err error) {
	var d Storage
	switch db_type {
	case utils.MONGO:
		d, err = NewMongoStorage(host, port, name, user, pass, utils.StorDB, cdrsIndexes, nil, 1)
	case utils.POSTGRES:
		d, err = NewPostgresStorage(host, port, name, user, pass, maxConn, maxIdleConn, connMaxLifetime)
	case utils.MYSQL:
		d, err = NewMySQLStorage(host, port, name, user, pass, maxConn, maxIdleConn, connMaxLifetime)
	default:
		err = errors.New(fmt.Sprintf("Unknown db '%s' valid options are [%s, %s, %s]",
			db_type, utils.MYSQL, utils.MONGO, utils.POSTGRES))
	}
	if err != nil {
		return nil, err
	}
	return d, nil
}

func ConfigureLoadStorage(db_type, host, port, name, user, pass, marshaler string, maxConn, maxIdleConn, connMaxLifetime int, cdrsIndexes []string) (db LoadStorage, err error) {
	var d LoadStorage
	switch db_type {
	case utils.POSTGRES:
		d, err = NewPostgresStorage(host, port, name, user, pass, maxConn, maxIdleConn, connMaxLifetime)
	case utils.MYSQL:
		d, err = NewMySQLStorage(host, port, name, user, pass, maxConn, maxIdleConn, connMaxLifetime)
	case utils.MONGO:
		d, err = NewMongoStorage(host, port, name, user, pass, utils.StorDB, cdrsIndexes, nil, 1)
	default:
		err = errors.New(fmt.Sprintf("Unknown db '%s' valid options are [%s, %s, %s]",
			db_type, utils.MYSQL, utils.MONGO, utils.POSTGRES))
	}
	if err != nil {
		return nil, err
	}
	return d, nil
}

func ConfigureCdrStorage(db_type, host, port, name, user, pass string, maxConn, maxIdleConn, connMaxLifetime int, cdrsIndexes []string) (db CdrStorage, err error) {
	var d CdrStorage
	switch db_type {
	case utils.POSTGRES:
		d, err = NewPostgresStorage(host, port, name, user, pass, maxConn, maxIdleConn, connMaxLifetime)
	case utils.MYSQL:
		d, err = NewMySQLStorage(host, port, name, user, pass, maxConn, maxIdleConn, connMaxLifetime)
	case utils.MONGO:
		d, err = NewMongoStorage(host, port, name, user, pass, utils.StorDB, cdrsIndexes, nil, 1)
	default:
		err = errors.New(fmt.Sprintf("Unknown db '%s' valid options are [%s, %s, %s]",
			db_type, utils.MYSQL, utils.MONGO, utils.POSTGRES))
	}
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Stores one Cost coming from SM
type SMCost struct {
	CGRID       string
	RunID       string
	OriginHost  string
	OriginID    string
	CostSource  string
	Usage       float64
	CostDetails *CallCost
}

type AttrCDRSStoreSMCost struct {
	Cost           *SMCost
	CheckDuplicate bool
}

type ArgsV2CDRSStoreSMCost struct {
	Cost           *V2SMCost
	CheckDuplicate bool
}

type V2SMCost struct {
	CGRID       string
	RunID       string
	OriginHost  string
	OriginID    string
	CostSource  string
	Usage       float64
	CostDetails *EventCost
}
