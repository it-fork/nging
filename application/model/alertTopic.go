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
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/model/base"
)

func NewAlertTopic(ctx echo.Context) *AlertTopic {
	m := &AlertTopic{
		NgingAlertTopic: &dbschema.NgingAlertTopic{},
		base:            base.New(ctx),
	}
	m.SetContext(ctx)
	return m
}

type AlertTopic struct {
	*dbschema.NgingAlertTopic
	base *base.Base
}

func (s *AlertTopic) check(row *dbschema.NgingAlertTopic) error {
	row.Topic = strings.TrimSpace(row.Topic)
	if len(row.Topic) == 0 {
		return s.base.E(`topic不能为空`)
	}
	if row.RecipientId <= 0 {
		return s.base.E(`收信账号ID不能为空`)
	}
	var (
		exists bool
		err    error
	)
	if row.Id > 0 {
		exists, err = s.ExistsOther(row.Topic, row.RecipientId, row.Id)
	} else {
		exists, err = s.Exists(row.Topic, row.RecipientId)
	}
	if err != nil {
		return err
	}
	if exists {
		err = s.base.E(`数据已经存在`)
	}
	return err
}

func (s *AlertTopic) Add(rows ...*dbschema.NgingAlertTopic) (pk interface{}, err error) {
	var bean *dbschema.NgingAlertTopic
	if len(rows) > 0 {
		bean = rows[0]
	} else {
		bean = s.NgingAlertTopic
	}
	if err = s.check(bean); err != nil {
		return nil, err
	}
	return bean.Add()
}

func (s *AlertTopic) Exists(topic string, recipientId uint) (bool, error) {
	return s.NgingAlertTopic.Exists(nil, db.And(
		db.Cond{`topic`: topic},
		db.Cond{`recipient_id`: recipientId},
	))
}

func (s *AlertTopic) ExistsOther(topic string, recipientId uint, excludeID uint) (bool, error) {
	return s.NgingAlertTopic.Exists(nil, db.And(
		db.Cond{`topic`: topic},
		db.Cond{`recipient_id`: recipientId},
		db.Cond{`id`: db.NotEq(excludeID)},
	))
}

func (s *AlertTopic) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	if err = s.check(s.NgingAlertTopic); err != nil {
		return err
	}
	return s.NgingAlertTopic.Edit(mw, args...)
}

func (s *AlertTopic) Send(topic string, params param.Store) (err error) {
	skey := `NgingAlertTopics.` + topic
	rows, ok := s.base.Context.Internal().Get(skey).([]*AlertTopicExt)
	if !ok {
		rows = []*AlertTopicExt{}
		_, err = s.ListByOffset(&rows, nil, 0, -1, db.And(
			db.Cond{`topic`: topic},
			db.Cond{`disabled`: `N`},
		))
		if err != nil {
			return
		}
		s.base.Context.Internal().Set(skey, rows)
	}
	for _, row := range rows {
		err = row.Send(params)
	}
	return
}
