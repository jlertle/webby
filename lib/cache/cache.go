// Data Cache
package cache

import (
	"encoding/gob"
	"os"
	"path/filepath"
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
	cache_list             = map[string]interface{}{}
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
	if !cacheExpiryCheckActive {
		cacheExpiryCheckActive = true
		go cacheExpiryCheck()
	}
	cache_list[key] = cache{value, expire}
}

func (_ CacheMemory) Get(key string) interface{} {
	switch t := cache_list[key].(type) {
	case cache:
		if time.Now().Unix() > t.Expire.Unix() {
			delete(cache_list, key)
			return nil
		}
		return t.Content
	}
	return nil
}

func (_ CacheMemory) Delete(key string) {
	delete(cache_list, key)
}

func (_ CacheMemory) Purge(beginWith string) {
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

const CacheFileExt = ".wbc"

// Note: Remember the filename limit is 255 with the majority of modern file systems!
// Avoid using reserved characters such as / ? \ % * : | " < >
// Hashes such as SHA256 is advisable and recommended!
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
		if len(cache_list) <= 0 {
			cacheExpiryCheckActive = false
			break
		}
		curtime := time.Now()
		for key, value := range cache_list {
			switch t := value.(type) {
			case cache:
				if curtime.Unix() > t.Expire.Unix() {
					delete(cache_list, key)
				}
			}
		}
	}
}
