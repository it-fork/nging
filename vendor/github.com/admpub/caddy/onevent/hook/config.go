package hook

import (
	"github.com/admpub/caddy"
)

// Config describes how Hook should be configured and used.
type Config struct {
	ID      string
	Event   caddy.EventName
	Command string
	Args    []string
}

// SupportedEvents is a map of supported events.
var SupportedEvents = map[string]caddy.EventName{
	"startup":   caddy.InstanceStartupEvent,
	"shutdown":  caddy.ShutdownEvent,
	"certrenew": caddy.CertRenewEvent,
}
