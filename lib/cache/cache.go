// Data Cache
package cache

import (
	"time"
)

type cache struct {
	content interface{}
	expire  time.Time
}

var cache_list = map[string]interface{}{}

func Set(key string, value interface{}) {
	SetAdv(key, value, time.Now().Add(1*time.Hour))
}

func SetAdv(key string, value interface{}, expire time.Time) {
	cache_list[key] = cache{value, expire}
}

func Get(key string) interface{} {
	switch t := cache_list[key].(type) {
	case cache:
		if time.Now().Unix() > t.expire.Unix() {
			delete(cache_list, key)
			return nil
		}
		return t.content
	}
	return nil
}

func Delete(key string) {
	delete(cache_list, key)
}

func Purge(beginWith string) {
	beginWith_len := len(beginWith)
	for key, _ := range cache_list {
		if len(key) < beginWith_len {
			continue
		}
		if beginWith == key[:beginWith_len] {
			delete(cache_list, key)
		}
	}
}

func cacheCheck() {
	for {
		time.Sleep(10 * time.Minute)
		curtime := time.Now()
		for key, value := range cache_list {
			switch t := value.(type) {
			case cache:
				if curtime.Unix() > t.expire.Unix() {
					delete(cache_list, key)
				}
			}
		}
	}
}

func init() {
	go cacheCheck()
}
