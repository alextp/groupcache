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

type HeapItem struct {
	priority float64
	position int
	Value    interface{}
}

type ItemArray []*HeapItem

type Heap struct {
	size     int
	elements ItemArray
}

func NewHeap() *Heap {
	return &Heap{
		size:     0,
		elements: ItemArray{},
	}
}

func (heap *Heap) Swap(i, j int) {
	a := heap.elements[i]
	heap.elements[i] = heap.elements[j]
	heap.elements[j] = a
	heap.elements[i].position = i
	heap.elements[j].position = j
}

func (heap *Heap) Up(index int) {
	for {
		i := (index - 1) / 2 // parent
		if i == index || !(heap.elements[index].priority < heap.elements[i].priority) {
			break
		}
		heap.Swap(i, index)
		index = i
	}
}

func (heap *Heap) Down(index int) {
	for {
		j1 := 2*index + 1
		if j1 >= heap.size || j1 < 0 {
			break
		}
		j := j1
		p1 := heap.elements[j1].priority
		j2 := j1 + 1
		p2 := heap.elements[j2].priority
		if j2 < heap.size && !(p1 < p2) {
			j = j2
		}
		if !(heap.elements[j].priority < heap.elements[index].priority) {
			break
		}
		heap.Swap(index, j)
		index = j
	}
}

func (heap *Heap) Push(element *HeapItem) {
	element.position = heap.size
	heap.size += 1
	heap.elements = append(heap.elements, element)
	heap.Up(element.position)
}

func (heap *Heap) Insert(element interface{}, priority float64) *HeapItem {
	item := &HeapItem{
		Value:    element,
		priority: priority,
		position: -1,
	}
	heap.Push(item)
	return item
}

func (heap *Heap) Remove(index int) *HeapItem {
	n := heap.size - 1
	heap.Swap(index, n)
	heap.Down(index)
	e := heap.elements[n]
	heap.size -= 1
	heap.elements = heap.elements[0:heap.size]
	return e
}

func (heap *Heap) Reinsert(index int, priority float64) {
	item := heap.Remove(index)
	item.priority = priority
	heap.Push(item)
}

func (heap *Heap) Pop() *HeapItem {
	return heap.Remove(0)
}

// Cache is an LRU cache. It is not safe for concurrent access.
type Cache struct {
	// OnEvicted optionally specificies a callback function to be
	// executed when an entry is purged from the cache.
	OnEvicted func(key Key, value interface{})

	operations int
	heap       *Heap
	cache      map[interface{}]*HeapItem
	lastUse    map[interface{}]int
}

// A Key may be any value that is comparable. See http://golang.org/ref/spec#Comparison_operators
type Key interface{}

type entry struct {
	key   Key
	value interface{}
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func New() *Cache {
	return &Cache{
		heap:       NewHeap(),
		operations: 0,
		cache:      make(map[interface{}]*HeapItem),
		lastUse:    make(map[interface{}]int),
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key Key, value interface{}) {
	if c.cache == nil {
		c.cache = make(map[interface{}]*HeapItem)
		c.lastUse = make(map[interface{}]int)
		c.heap = NewHeap()
		c.operations = 0
	}
	if ee, ok := c.cache[key]; ok {
		c.operations += 1
		c.heap.Reinsert(ee.position, float64(-c.operations)) // TODO(apassos): perceptron decision goes here
		ee.Value.(*entry).value = value
		return
	}
	c.operations += 1
	ele := c.heap.Insert(&entry{key, value}, float64(-c.operations)) // TODO(apassos): perceptron decision goes here
	c.cache[key] = ele
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	c.operations += 1
	if ele, hit := c.cache[key]; hit {
		c.heap.Reinsert(ele.position, float64(-c.operations)) // TODO(apassos): perceptron decision goes here
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

func (c *Cache) removeElement(e *HeapItem) {
	// TODO(apassos): perceptron update goes here
	c.heap.Remove(e.position)
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
	return c.heap.size
}
