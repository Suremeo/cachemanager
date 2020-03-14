package cacher

import (
	"errors"
	"io/ioutil"
	"sync"
	"time"
)

type Cacher struct {
	Expire int
	items []*Item
	running bool
	mutex sync.Mutex
}

type Item struct {
	Identifier string
	Data interface{}
	added int
}

func (cache *Cacher) Add (label string, item interface{}) *Cacher {
	cache.mutex.Lock()
	cache.items = append(cache.items, &Item{
		Identifier: label,
		Data: item,
		added: int(time.Now().Unix()),
	})
	cache.mutex.Unlock()
	return cache
}

func (cache *Cacher) Get(identifier string) (*Item, error) {
	for _, element := range cache.items {
		if element.Identifier == identifier {
			return element, nil
		}
	}
	return nil, errors.New("cache: item not cached")
}

func (cache *Cacher) Run() *Cacher {
	if cache.Expire == 0 {
		cache.Expire = 30
	}
	cache.running = true
	go func(){
		defer func() {
			if i := recover(); i != nil {
				cache.Run()
			}
		}()
		for {
			cache.mutex.Lock()
			now := int(time.Now().Unix())
			for index, item := range cache.items {
				if (item.added + cache.Expire) < now {
					cache.items = remove(cache.items, index)
				}
			}
			cache.mutex.Unlock()
			time.Sleep(1*time.Second)
		}
	}()
	return cache
}

func (cache *Cacher) Remove(identifier string) error {
	defer func() {
		if i := recover(); i != nil {}
	}()
	cache.mutex.Lock()
	for index, element := range cache.items {
		if element.Identifier == identifier {
			cache.items = remove(cache.items, index)
			cache.mutex.Unlock()
			return nil
		}
	}
	cache.mutex.Unlock()
	return errors.New("cache: Item not found")
}

func (cache *Cacher) Clear() *Cacher {
	cache.items = []*Item{}
	return cache
}

func (cache *Cacher) File(path string) (data []byte, err error, wascached bool) {
	id := "FILECACHE:" + path
	item, err := cache.Get(id)
	if err != nil {
		data, err = ioutil.ReadFile(path)
		if err == nil {
			cache.Add(id, data)
		}
	} else {
		byt, ok := item.Data.([]byte)
		if ok {
			return byt, err, true
		} else {
			data, err = ioutil.ReadFile(path)
			if err == nil {
				cache.Add(id, data)
			}
		}
	}
	return data, err, false
}

func remove(slice []*Item, s int) []*Item {
	defer func() {
		if i := recover(); i != nil {}
	}()
	return append(slice[:s], slice[s+1:]...)
}