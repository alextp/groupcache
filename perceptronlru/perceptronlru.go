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

// Package lru implements a cache that is LRU-like but tries to guess
// next access time with a perceptron.
package perceptronlru

import 	"github.com/golang/groupcache/perceptronlru/heap"


// Cache is an LRU cache. It is not safe for concurrent access.
type Cache struct {
	// OnEvicted optionally specificies a callback function to be
	// executed when an entry is purged from the cache.
	OnEvicted func(key Key, value interface{})

	operations int
	heap       *heap.Heap
	cache      map[interface{}]*heap.HeapItem
}

// A Key may be any value that is comparable. See http://golang.org/ref/spec#Comparison_operators
type Key interface{}

type entry struct {
	key     Key
	lastUse int
	value   interface{}
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func New() *Cache {
	return &Cache{
		heap:       heap.NewHeap(),
		operations: 0,
		cache:      make(map[interface{}]*heap.HeapItem),
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key Key, value interface{}) {
	if c.cache == nil {
		c.cache = make(map[interface{}]*heap.HeapItem)
		c.heap = heap.NewHeap()
		c.operations = 0
	}
	if ee, ok := c.cache[key]; ok {
		c.operations += 1
		c.heap.Reinsert(ee.Position, float64(-c.operations)) // TODO(apassos): perceptron decision goes here
		ee.Value.(*entry).lastUse = c.operations
		ee.Value.(*entry).value = value
		return
	}
	c.operations += 1
	ele := c.heap.Insert(&entry{key, c.operations, value}, float64(-c.operations)) // TODO(apassos): perceptron decision goes here
	c.cache[key] = ele
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	c.operations += 1
	if ele, hit := c.cache[key]; hit {
		c.heap.Reinsert(ele.Position, float64(-c.operations)) // TODO(apassos): perceptron decision goes here
		ele.Value.(*entry).lastUse = c.operations
		return ele.Value.(*entry).value, true
	}
	return
}

// Remove removes the provided key from the cache.
func (c *Cache) Remove(key Key) {
	if c.cache == nil {
		return
	}
	c.operations += 1
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// RemoveOldest removes the oldest item from the cache.
func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.heap.Pop()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *Cache) removeElement(e *heap.HeapItem) {
	// TODO(apassos): perceptron update goes here
	c.heap.Remove(e.Position)
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
	return c.heap.Size
}
