package common

import (
	"strings"

	"github.com/webx-top/echo"

	"github.com/admpub/nging/v3/application/library/config/subconfig/scookie"
)

// APIKeyGetter API Key
type APIKeyGetter interface {
	APIKey() string
}

type CookieConfigGetter interface {
	CookieConfig() scookie.Config
}

func CookieConfig() scookie.Config {
	return echo.Get(`DefaultConfig`).(CookieConfigGetter).CookieConfig()
}

func Setting(group ...string) echo.H {
	return echo.GetStoreByKeys(`NgingConfig`, group...)
}

func BackendURL(ctx echo.Context) string {
	backendURL := Setting(`base`).String(`backendURL`)
	if len(backendURL) == 0 {
		if ctx == nil {
			return backendURL
		}
		backendURL = ctx.Site()
	}
	backendURL = strings.TrimSuffix(backendURL, `/`)
	return backendURL
}

func SystemAPIKey() string {
	apiKey := Setting(`base`).String(`apiKey`)
	return apiKey
}
