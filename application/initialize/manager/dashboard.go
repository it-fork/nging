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

package manager

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v3/application/handler"
	"github.com/admpub/nging/v3/application/library/system"
	"github.com/admpub/nging/v3/application/model"
	"github.com/admpub/nging/v3/application/registry/dashboard"
)

func init() {
	dashboard.CardRegister(
		(&dashboard.Card{
			IconName:  `fa-sitemap`,
			IconColor: `primary`,
			Short:     `SITES`,
			Name:      `网站数量`,
			Summary:   ``,
		}).SetContentGenerator(func(ctx echo.Context) interface{} {
			//网站统计
			vhostMdl := model.NewVhost(ctx)
			vhostCount, _ := vhostMdl.Count(nil)
			return vhostCount
		}),
		(&dashboard.Card{
			IconName:  `fa-tasks`,
			IconColor: `danger`,
			Short:     `TASKS`,
			Name:      `计划任务数量`,
			Summary:   ``,
		}).SetContentGenerator(func(ctx echo.Context) interface{} {
			//计划任务统计
			taskMdl := model.NewTask(ctx)
			taskCount, _ := taskMdl.Count(nil)
			return taskCount
		}),
	)

	dashboard.BlockRegister((&dashboard.Block{
		Tmpl:   `server/chart/cpu`,
		Footer: `server/chart/cpu.js`,
	}).SetContentGenerator(func(ctx echo.Context) error {
		ctx.Set(`systemRealtimeStatusIsListening`, system.RealTimeStatusIsListening())
		return nil
	}))
	dashboard.BlockRegister((&dashboard.Block{
		Tmpl: `server/dashbord/cmd_list`,
	}).SetContentGenerator(func(ctx echo.Context) error {
		user := handler.User(ctx)
		//指令集
		cmdMdl := model.NewCommand(ctx)
		if user.Id == 1 {
			cmdMdl.ListByOffset(nil, nil, 0, -1)
		} else {
			roleList := handler.UserRoles(ctx)
			cmdIds := []string{}
			for _, role := range roleList {
				if len(role.PermCmd) > 0 {
					cmdIds = append(cmdIds, strings.Split(role.PermCmd, `,`)...)
				}
			}
			if len(cmdIds) > 0 {
				cmdMdl.ListByOffset(nil, nil, 0, -1, db.Cond{`id`: db.In(cmdIds)})
			}
		}
		ctx.Set(`cmdList`, cmdMdl.Objects())
		return nil
	}))
}
