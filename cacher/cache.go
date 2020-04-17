package cacher

import (
	"errors"
	"io/ioutil"
	"sync"
	"time"
)

type Cacher struct {
	Expire int
	Subsections map[string]*Cacher
	items map[string]*Item
	running bool
	mutex sync.Mutex
}

type Item struct {
	Identifier string
	Data interface{}
	added int
}

func (cache *Cacher) Set(label string, item interface{}) *Cacher {
	if !cache.running {
		cache.Run()
	}
	cache.items[label] = &Item{
		Identifier: label,
		Data: item,
		added: int(time.Now().Unix()),
	}
	return cache
}

func (cache *Cacher) Get(identifier string) (*Item, error) {
	if !cache.running {
		cache.Run()
	}
	dat := cache.items[identifier]
	if dat == nil {
		return dat, errors.New("cache: item not cached")
	}
	return dat, nil
}

func (cache *Cacher) Run() *Cacher {
	if cache.Expire == 0 {
		cache.Expire = 30
	}
	if cache.items == nil {
		cache.items = map[string]*Item{}
	}
	if cache.running {
		return cache
	}
	cache.running = true
	go func(){
		defer func() {
			if i := recover(); i != nil {
				cache.running = false
				cache.Run()
			}
		}()
		for {
			cache.mutex.Lock()
			now := int(time.Now().Unix())
			for key, item := range cache.items {
				if (item.added + cache.Expire) < now {
					delete(cache.items, key)
				}
			}
			cache.mutex.Unlock()
			time.Sleep(1*time.Second)
		}
	}()
	return cache
}

func (cache *Cacher) Remove(identifier string) {
	if !cache.running {
		cache.Run()
	}
	defer func() {
		if i := recover(); i != nil {}
	}()

	delete(cache.items, identifier)
}

func (cache *Cacher) Clear() *Cacher {
	if !cache.running {
		cache.Run()
	}
	cache.items = map[string]*Item{}
	return cache
}

func (cache *Cacher) File(path string) (data []byte, err error, wascached bool) {
	if !cache.running {
		cache.Run()
	}
	id := "FILECACHE:" + path
	item, err := cache.Get(id)
	if err != nil {
		data, err = ioutil.ReadFile(path)
		if err == nil {
			cache.Set(id, data)
		}
	} else {
		byt, ok := item.Data.([]byte)
		if ok {
			return byt, err, true
		} else {
			data, err = ioutil.ReadFile(path)
			if err == nil {
				cache.Set(id, data)
			}
		}
	}
	return data, err, false
}