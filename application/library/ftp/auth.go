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

package ftp

import (
	"github.com/admpub/nging/v3/application/model"
	ftpserver "goftp.io/server/v2"
)

func NewAuth() *Auth {
	return &Auth{
		FtpUser: model.NewFtpUser(nil),
	}
}

type Auth struct {
	*model.FtpUser
}

func (a *Auth) CheckPasswd(ftpCtx *ftpserver.Context, username string, password string) (bool, error) {
	return a.FtpUser.CheckPasswd(username, password)
}
