/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package model

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/model/base"
)

func NewDbSyncLog(ctx echo.Context) *DbSyncLog {
	return &DbSyncLog{
		NgingDbSyncLog: &dbschema.NgingDbSyncLog{},
		Base:           base.New(ctx),
	}
}

type DbSyncLog struct {
	*dbschema.NgingDbSyncLog
	*base.Base
}

func (a *DbSyncLog) Add() (interface{}, error) {
	return a.NgingDbSyncLog.Add()
}

func (a *DbSyncLog) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return a.NgingDbSyncLog.Edit(mw, args...)
}
