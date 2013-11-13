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

package greedydual

// Cache is a cache. It is not safe for concurrent access.
type Cache struct {
	// OnEvicted optionally specificies a callback function to be
	// executed when an entry is purged from the cache.
	OnEvicted func(key string, value interface{})

	operations float64
	Heap       *Heap
	cache      map[interface{}]*HeapItem
}

type entry struct {
	key            string
	weights        float64
	value          interface{}
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func New() *Cache {
	return &Cache{
		Heap:       NewHeap(),
		operations: 0,
		cache:      make(map[interface{}]*HeapItem),
	}
}


// Add adds a value to the cache.
func (c *Cache) Add(key string, value interface{}, cost float64) {
	if _, ok := c.cache[key]; ok {
		println("ERROR")
		return
	}
	entry := &entry{
		key:            key,
		value:          value,
		weights:        cost,
	}
	priority := c.operations + cost
	ele := c.Heap.Insert(entry, priority)
	c.cache[key] = ele
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key string) (value interface{}, ok bool) {
	if ele, hit := c.cache[key]; hit {
		ee := ele.Value.(*entry)
		if ee.weights < 0 {
			println("so wrong")
		}
		priority := c.operations + ee.weights
		c.Heap.Reinsert(ele.Position, priority)
		return ele.Value.(*entry).value, true
	}
	return
}

// Remove removes the provided key from the cache.
func (c *Cache) Remove(key string) {
	if c.cache == nil {
		return
	}
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
		if c.operations > ele.Priority {
			println("really weird error", c.operations, ele.Priority)
			//var a *entry
			//a = nil
			//println(a.value)
		}
		c.operations = ele.Priority
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
