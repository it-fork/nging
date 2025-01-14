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
	"os"

	"github.com/admpub/nging/v3/application/handler"
	"github.com/admpub/nging/v3/application/library/config"
	"github.com/admpub/nging/v3/application/library/notice"
	"github.com/webx-top/echo"
)

func Restart(ctx echo.Context) error {
	wOut, wErr, err := handler.NoticeWriter(ctx, ctx.T(`FTP服务`))
	if err != nil {
		return ctx.String(err.Error())
	}
	if err := config.DefaultCLIConfig.FTPRestart(wOut, wErr); err != nil {
		return ctx.String(err.Error())
	}
	return ctx.String(ctx.T(`已经重启FTP服务`))
}

func Stop(ctx echo.Context) error {
	if err := config.DefaultCLIConfig.FTPStop(); err != nil {
		return ctx.String(err.Error())
	}
	return ctx.String(ctx.T(`已经关闭FTP服务`))
}

func Log(ctx echo.Context) error {
	on := ctx.Formx(`on`).Bool()
	if on {
		wOut, wErr, err := handler.NoticeWriter(ctx, ctx.T(`FTP服务`))
		if err != nil {
			return ctx.String(err.Error())
		}
		err = config.DefaultCLIConfig.SetLogWriter(`ftp`, wOut, wErr)
		if err != nil {
			return ctx.String(err.Error())
		}
		return ctx.String(ctx.T(`已经开始直播FTP服务状态`))
	}
	err := config.DefaultCLIConfig.SetLogWriter(`ftp`, os.Stdout, os.Stderr)
	if err != nil {
		return ctx.String(err.Error())
	}
	user := handler.User(ctx)
	if user == nil {
		return ctx.String(ctx.T(`请先登录`))
	}
	typ := `service:` + ctx.T(`FTP服务`)
	notice.CloseMessage(user.Username, typ)
	return ctx.String(ctx.T(`已经停止直播FTP服务状态`))
}
