package cachemanager

import (
	"./cacher"
)

func NewCache() *cacher.Cacher {
	return &cacher.Cacher{}
}