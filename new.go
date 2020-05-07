package cachemanager

import (
	"github.com/suremeo/cachemanager/cacher"
)

func NewCache() *cacher.Cacher {
	return &cacher.Cacher{}
}