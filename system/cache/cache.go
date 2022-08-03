// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package cache

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/cache"
	"github.com/e154/smart-home-node/common/logger"
)

var (
	log = logger.MustGetLogger("cache")
)

type Cache struct {
	bm        cache.Cache
	Cachetime int64
	Name      string "maincache"
	Verbose   bool
}

func (c *Cache) init() (*Cache, error) {
	c.log("init")

	var err error
	c.bm, err = cache.NewCache("memory", fmt.Sprintf(`{"interval":%d}`, time.Duration(c.Cachetime)*time.Second))
	if err != nil {
		c.log("error %s", err.Error())
	}
	return c, err
}

func (c *Cache) ClearAll() (*Cache, error) {
	c.log("clear all")

	if c.bm == nil {
		c.init()
	}

	err := c.bm.ClearAll()

	return c, err
}

func (c *Cache) GetKey(key interface{}) string {
	return fmt.Sprintf("%s_%s", c.Name, key.(string))
}

func (c *Cache) Clear(key interface{}) (*Cache, error) {
	cacheKey := c.GetKey(key)
	c.log("clear %s", cacheKey)

	if c.bm == nil {
		c.init()
	}

	err := c.bm.Delete(cacheKey)

	return c, err
}

func (c *Cache) addToGroup(group, key string) (*Cache, error) {

	if c.bm == nil {
		c.init()
	}

	g := []string{}
	w := c.bm.Get(group)
	if w != nil {
		g = w.([]string)
	}

	exist := false
	for _, v := range g {
		if key == v {
			exist = true
		}
	}

	var err error
	if !exist {
		c.log("add to group %s", group)
		g = append(g, key)
		err = c.bm.Put(group, g, time.Duration(c.Cachetime)*time.Second)
	}

	return c, err
}

func (c *Cache) ClearGroup(group string) (*Cache, error) {
	c.log("clear group %s", group)

	if c.bm == nil {
		c.init()
	}

	g := []string{}
	w := c.bm.Get(group)
	if w == nil {
		return c, nil
	}

	g = w.([]string)
	if len(g) == 0 {
		return c, nil
	}

	for _, key := range g {
		c.bm.Delete(key)
	}

	_, err := c.Clear(group)

	return c, err
}

func (c *Cache) Put(group, key string, val interface{}) (*Cache, error) {
	c.log("put key %s", key)

	if c.bm == nil {
		c.init()
	}

	if err := c.bm.Put(key, val, time.Duration(c.Cachetime)*time.Second); err != nil {
		return c, err
	}

	return c.addToGroup(group, key)
}

func (c *Cache) IsExist(key string) bool {
	if c.bm == nil {
		c.init()
	}

	return c.bm.IsExist(key)
}

func (c *Cache) Get(key string) interface{} {
	c.log("get key %s", key)

	if c.bm == nil {
		c.init()
	}

	return c.bm.Get(key)
}

func (c *Cache) Delete(key string) *Cache {
	c.log("delete value by key %s", key)

	if c.bm == nil {
		c.init()
	}

	c.bm.Delete(key)

	return c
}

func (c *Cache) log(format string, a ...interface{}) {
	if c.Verbose {
		log.Debugf("cache: %s", fmt.Sprintf(format, a...))
	}
}
