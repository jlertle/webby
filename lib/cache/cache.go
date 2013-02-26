// Data Cache
package cache

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type cache struct {
	Content interface{}
	Expire  time.Time
}

func init() {
	gob.Register(cache{})
}

var (
	cache_list = struct {
		sync.Mutex
		m map[string]interface{}
	}{m: map[string]interface{}{}}
	cacheExpiryCheckActive = false
)

type CacheHandler interface {
	Set(string, interface{})
	SetAdv(string, interface{}, time.Time)
	Get(string) interface{}
	Delete(string)
	Purge(string)
}

type CacheMemory struct{}

func (c CacheMemory) Set(key string, value interface{}) {
	c.SetAdv(key, value, time.Now().Add(1*time.Hour))
}

func (_ CacheMemory) SetAdv(key string, value interface{}, expire time.Time) {
	cache_list.Lock()
	defer cache_list.Unlock()

	if !cacheExpiryCheckActive {
		cacheExpiryCheckActive = true
		go cacheExpiryCheck()
	}
	cache_list.m[key] = cache{value, expire}
}

func (c CacheMemory) Get(key string) interface{} {
	cache_list.Lock()
	defer cache_list.Unlock()

	switch t := cache_list.m[key].(type) {
	case cache:
		if time.Now().Unix() > t.Expire.Unix() {
			c.delete(key)
			return nil
		}
		return t.Content
	}
	return nil
}

func (_ CacheMemory) delete(key string) {
	delete(cache_list.m, key)
}

func (c CacheMemory) Delete(key string) {
	cache_list.Lock()
	defer cache_list.Unlock()
	c.delete(key)
}

func (c CacheMemory) Purge(beginWith string) {
	cache_list.Lock()
	defer cache_list.Unlock()

	beginWith_len := len(beginWith)
	for key, _ := range cache_list.m {
		if len(key) < beginWith_len {
			continue
		}
		if beginWith == key[:beginWith_len] {
			c.delete(key)
		}
	}
}

const CacheFileExt = ".wbc"

// Note: Remember the filename limit is 255 (minus '.wbc' 251) with the majority of modern file systems!
// Avoid using reserved characters such as / ? \ % * : | " < >
// Hashes such as SHA256 is advisable when possible!
type CacheFile struct {
	Path string
}

func (c CacheFile) Set(key string, value interface{}) {
	c.SetAdv(key, value, time.Now().Add(1*time.Hour))
}

func (c CacheFile) SetAdv(key string, value interface{}, expire time.Time) {
	file, err := os.Create(c.Path + "/" + key + CacheFileExt)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	err = enc.Encode(cache{value, expire})
	if err != nil {
		panic(err)
	}
}

func (c CacheFile) Get(key string) interface{} {
	file, err := os.Open(c.Path + "/" + key + CacheFileExt)
	if err != nil {
		return nil
	}
	defer file.Close()

	cac := cache{}

	dec := gob.NewDecoder(file)
	err = dec.Decode(&cac)
	if err != nil {
		return nil
	}

	if time.Now().Unix() > cac.Expire.Unix() {
		c.Delete(key)
		return nil
	}

	return cac.Content
}

func (c CacheFile) Delete(key string) {
	os.Remove(c.Path + "/" + key + CacheFileExt)
}

func (c CacheFile) Purge(beginWith string) {
	matches, err := filepath.Glob(c.Path + "/" + beginWith + "*" + CacheFileExt)
	if err != nil {
		return
	}

	for _, match := range matches {
		os.Remove(match)
	}
}

var DefaultCacheHandler CacheHandler = CacheMemory{}

func Set(key string, value interface{}) {
	DefaultCacheHandler.Set(key, value)
}

func SetAdv(key string, value interface{}, expire time.Time) {
	DefaultCacheHandler.SetAdv(key, value, expire)
}

func Get(key string) interface{} {
	return DefaultCacheHandler.Get(key)
}

func Delete(key string) {
	DefaultCacheHandler.Delete(key)
}

func Purge(beginWith string) {
	DefaultCacheHandler.Purge(beginWith)
}

func cacheExpiryCheck() {
	for {
		time.Sleep(10 * time.Minute)

		cache_list.Lock()

		if len(cache_list.m) <= 0 {
			cacheExpiryCheckActive = false
			cache_list.Unlock()
			break
		}
		curtime := time.Now()
		for key, value := range cache_list.m {
			switch t := value.(type) {
			case cache:
				if curtime.Unix() > t.Expire.Unix() {
					delete(cache_list.m, key)
				}
			}
		}

		cache_list.Unlock()
	}
}
