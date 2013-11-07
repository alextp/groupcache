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

import "github.com/golang/groupcache/perceptronlru/heap"
import "github.com/golang/groupcache/perceptronlru/perceptron"
import "hash/fnv"

// Cache is an LRU cache. It is not safe for concurrent access.
type Cache struct {
	// OnEvicted optionally specificies a callback function to be
	// executed when an entry is purged from the cache.
	OnEvicted func(key string, value interface{})

	model      *perceptron.Perceptron
	operations int
	Heap       *heap.Heap
	cache      map[interface{}]*heap.HeapItem
}

// A Key may be any value that is comparable. See http://golang.org/ref/spec#Comparison_operators

func features(str string) []int32 {
	// picks all character 3-grams, 5-grams, and 7-grams from key
	key := []byte(str)
	features := make([]int32, 0, 3*len(key))
	lengths := &[...]int{3, 5, 7}
	for i := 0; i < len(lengths); i++ {
		length := lengths[i]
		for j := 0; j < len(key)-length; j++ {
			hash := fnv.New32a()
			hash.Write(key[j : j+length])
			features = append(features, int32(hash.Sum32()))
		}
	}
	return features
}

type entry struct {
	key     string
	lastUse int
	value   interface{}
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func New(modelSize int32) *Cache {
	return &Cache{
		Heap:       heap.NewHeap(),
		model:      perceptron.New(modelSize),
		operations: 0,
		cache:      make(map[interface{}]*heap.HeapItem),
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value interface{}) {
	if c.cache == nil {
		c.cache = make(map[interface{}]*heap.HeapItem)
		c.Heap = heap.NewHeap()
		c.operations = 0
	}
	if ee, ok := c.cache[key]; ok {
		c.operations += 1
		priority := (float64(c.operations) + c.model.Update(features(key), float64(c.operations-ee.Value.(*entry).lastUse)))
		c.Heap.Reinsert(ee.Position, priority) // TODO(apassos): perceptron decision goes here
		ee.Value.(*entry).lastUse = c.operations
		ee.Value.(*entry).value = value
		return
	}
	c.operations += 1
	priority := (float64(c.operations) + c.model.Score(features(key)))
	ele := c.Heap.Insert(&entry{key, c.operations, value}, priority) // TODO(apassos): perceptron decision goes here
	c.cache[key] = ele
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key string) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	c.operations += 1
	if ele, hit := c.cache[key]; hit {
		priority := (float64(c.operations) + c.model.Update(features(key), float64(c.operations-ele.Value.(*entry).lastUse)))
		c.Heap.Reinsert(ele.Position, priority) // TODO(apassos): perceptron decision goes here
		ele.Value.(*entry).lastUse = c.operations
		return ele.Value.(*entry).value, true
	}
	return
}

// Remove removes the provided key from the cache.
func (c *Cache) Remove(key string) {
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
	if c.Heap.Size > 0 {
		ele := c.Heap.Head()
		// TODO(apassos): perceptron update goes here
		feats := features(ele.Value.(*entry).key)
		prediction := float64(ele.Value.(*entry).lastUse) + c.model.Score(feats)
		if prediction > float64(c.operations) {
			c.model.Update(feats, float64(c.operations-ele.Value.(*entry).lastUse))
		}
		c.removeElement(ele)
	}
}

func (c *Cache) removeElement(e *heap.HeapItem) {
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
