package marketplace

import (
	"math/rand"
	"time"

	"github.com/gorilla/sessions"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/settings"
)

var (
	appSettings  = settings.GetSettings()
	sessionStore *sessions.CookieStore
	rs           = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func init() {

	if !appSettings.Debug {
		sessionStore = sessions.NewCookieStore([]byte(MARKETPLACE_SETTINGS.CookieEncryptionSalt))
	} else {
		sessionStore = sessions.NewCookieStore([]byte("debug"))
	}

	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24,
		HttpOnly: true,
	}
}
