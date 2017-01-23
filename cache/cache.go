package cache

import (
	"github.com/e154/smart-home-node/models"
)

// Singleton
var instantiated *models.Cache = nil

func CachePtr() *models.Cache {
	return instantiated
}

func Init(t int64) {
	instantiated = &models.Cache{
		Cachetime: t,
		Name: "node",
	}
}