/*
Copyright 2013 Alexandre Passos

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lfu

import "math"

// Cache is an LFU cache. It is not safe for concurrent access.
type Cache struct {
	// OnEvicted optionally specificies a callback function to be
	// executed when an entry is purged from the cache.
	OnEvicted func(key string, value interface{})

	operations int32
	version    int64
	alpha      float64
	limit      int32
	Heap       *Heap
	cache      map[interface{}]*HeapItem
}

type entry struct {
	key            string
	useAccumulator float64
	version        int64
	weights        float64
	value          interface{}
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func New(alpha float64, limit int32) *Cache {
	return &Cache{
		Heap:       NewHeap(),
		alpha:      alpha,
		limit:      limit,
		version:    0,
		operations: 0,
		cache:      make(map[interface{}]*HeapItem),
	}
}

func (c *Cache) doOperation(ee *entry) {
	c.operations += 1
	if c.operations == c.limit {
		c.operations = 0
		c.limit += 1
	}
	if ee != nil {
		if c.version-ee.version != 0 {
			dif := c.version - ee.version
			ee.useAccumulator *= math.Pow(c.alpha, math.Abs(float64(c.limit*int32(dif))))
			ee.version = c.version
		}
		ee.useAccumulator += (1 - c.alpha) / math.Pow(c.alpha, float64(c.operations))
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value interface{}, cost float64) {
	if ee, ok := c.cache[key]; ok {
		entry := ee.Value.(*entry)
		c.doOperation(entry)
		priority := entry.weights * entry.useAccumulator
		c.Heap.Reinsert(ee.Position, priority, c.version)
		return
	}
	entry := &entry{
		key:            key,
		value:          value,
		weights:        cost,
		version:        0,
		useAccumulator: 0,
	}
	c.doOperation(entry)
	priority := entry.weights * entry.useAccumulator
	ele := c.Heap.Insert(entry, priority, c.version)
	c.cache[key] = ele
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key string) (value interface{}, ok bool) {
	if ele, hit := c.cache[key]; hit {
		ee := ele.Value.(*entry)
		c.doOperation(ee)
		priority := ee.weights * ee.useAccumulator
		c.Heap.Reinsert(ele.Position, priority, c.version)
		return ele.Value.(*entry).value, true
	} else {
		c.doOperation(nil)
	}
	return
}

// Remove removes the provided key from the cache.
func (c *Cache) Remove(key string) {
	if c.cache == nil {
		return
	}
	c.doOperation(nil)
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// RemoveOldest removes the oldest item from the cache.
func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	if c.Heap.Size > 0 {
		ele := c.Heap.Head()
		c.removeElement(ele)
	}
}

func (c *Cache) removeElement(e *HeapItem) {
	c.Heap.Remove(e.Position)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.Heap.Size
}
