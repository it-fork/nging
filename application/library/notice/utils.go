package notice

import (
	"context"
	"sync/atomic"
	"time"
)

type HTTPNoticerConfig struct {
	User     string
	Type     string
	ClientID string
	ID       interface{}
	IsExited IsExited
	Timeout  time.Duration
	Mode     string
}

func NewHTTPNoticerConfig() *HTTPNoticerConfig {
	return &HTTPNoticerConfig{}
}

func (c *HTTPNoticerConfig) SetUser(user string) *HTTPNoticerConfig {
	c.User = user
	return c
}

func (c *HTTPNoticerConfig) SetType(typ string) *HTTPNoticerConfig {
	c.Type = typ
	return c
}

func (c *HTTPNoticerConfig) SetClientID(clientID string) *HTTPNoticerConfig {
	c.ClientID = clientID
	return c
}

func (c *HTTPNoticerConfig) SetID(id interface{}) *HTTPNoticerConfig {
	c.ID = id
	return c
}

func (c *HTTPNoticerConfig) SetTimeout(t time.Duration) *HTTPNoticerConfig {
	c.Timeout = t
	return c
}

func (c *HTTPNoticerConfig) SetIsExited(isExited IsExited) *HTTPNoticerConfig {
	c.IsExited = isExited
	return c
}

func (c *HTTPNoticerConfig) SetMode(mode string) *HTTPNoticerConfig {
	c.Mode = mode
	return c
}

func (c *HTTPNoticerConfig) Noticer(ctx context.Context) Noticer {
	return NewNoticer(ctx, c)
}

func NewControlWithContext(ctx context.Context, timeout time.Duration) IsExited {
	defaultCtrl := &Control{}
	defaultCtrl.ListenContextAndTimeout(ctx, timeout)
	return defaultCtrl
}

func NewWithProgress(noticer Noticer, progresses ...*Progress) *NoticeAndProgress {
	var progress *Progress
	if len(progresses) > 0 {
		progress = progresses[0]
	}
	if progress == nil {
		progress = NewProgress()
	}
	return &NoticeAndProgress{
		send: noticer,
		prog: progress,
	}
}

type NProgressor interface {
	Send(message interface{}, statusCode int) error
	Success(message interface{}) error
	Failure(message interface{}) error
	Add(n int64) NProgressor
	Done(n int64) NProgressor
	AutoComplete(on bool) NProgressor
	Complete() NProgressor
}

type NoticeAndProgress struct {
	send         Noticer
	prog         *Progress
	autoComplete bool
}

// - Noticer -

func (a *NoticeAndProgress) Send(message interface{}, statusCode int) error {
	return a.send(message, statusCode, a.prog)
}

func (a *NoticeAndProgress) Success(message interface{}) error {
	return a.send(message, 1, a.prog)
}

func (a *NoticeAndProgress) Failure(message interface{}) error {
	return a.send(message, 0, a.prog)
}

// - Progress -

func (a *NoticeAndProgress) Add(n int64) NProgressor {
	if a.prog.Finish > 0 {
		atomic.StoreInt64(&a.prog.Finish, 0)
	}
	atomic.AddInt64(&a.prog.Total, n)
	return a
}

func (a *NoticeAndProgress) Done(n int64) NProgressor {
	atomic.AddInt64(&a.prog.Finish, n)
	if a.autoComplete && a.prog.Finish == a.prog.Total {
		a.prog.Complete = true
	}
	return a
}

func (a *NoticeAndProgress) AutoComplete(on bool) NProgressor {
	a.autoComplete = on
	return a
}

func (a *NoticeAndProgress) Complete() NProgressor {
	a.prog.Complete = true
	return a
}

func NewNoticer(ctx context.Context, config *HTTPNoticerConfig) Noticer {
	var noticeSender Noticer
	if config.IsExited == nil && config.Timeout != 0 {
		config.IsExited = NewControlWithContext(ctx, config.Timeout)
	}
	if len(config.Mode) == 0 {
		if config.ID != nil {
			config.Mode = `element`
		} else {
			config.Mode = `notify`
		}
	}
	progress := NewProgress().SetControl(config.IsExited)
	if len(config.User) > 0 {
		OpenMessage(config.User, config.Type)
		//defer CloseMessage(config.User, config.Type)
		noticeSender = func(message interface{}, statusCode int, progs ...*Progress) error {
			msg := NewMessageWithValue(
				config.Type,
				``,
				message,
				statusCode,
			).SetMode(config.Mode).SetID(config.ID)
			var prog *Progress
			if len(progs) > 0 {
				prog = progs[0]
			}
			if prog == nil {
				prog = progress
			}
			msg.SetProgress(prog).CalcPercent().SetClientID(config.ClientID)
			sendErr := Send(config.User, msg)
			return sendErr
		}
	} else {
		noticeSender = DefaultNoticer
	}
	return noticeSender
}
