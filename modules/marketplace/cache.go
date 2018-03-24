package marketplace

import (
	"time"

	"github.com/bluele/gcache"
)

var (
	duration, _   = time.ParseDuration("15m")
	gc            = gcache.New(1024).LRU().Expiration(duration).Build()
	icDuration, _ = time.ParseDuration("1m")
	ic            = gcache.New(10240).LRU().Expiration(icDuration).Build()
)

func CacheGetUserUuid(username string) string {
	cUuid, _ := ic.Get(username)
	if cUuid == nil {
		user, err := FindUserByUsername(username)
		if err != nil {
			return ""
		}
		ic.Set(username, user.Uuid)
		return user.Uuid
	}

	return cUuid.(string)
}
