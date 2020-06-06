package cache

import (
	"container/list"
	"fmt"
)

// Entry describes each row of a Cache
type Entry struct {
	Key   interface{}
	Value interface{}
}

// Cache is a LRU cache container
type Cache struct {
	capacity int
	order    *list.List
	data     map[interface{}]*list.Element
}

// New creates a new Cache instance
func New(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		order:    list.New(),
		data:     make(map[interface{}]*list.Element),
	}
}

// Write inserts the key-value pair into the Cache
func (c *Cache) Write(key interface{}, value interface{}) {
	if c.capacity == c.Size() {
		element := c.order.Back()
		if element == nil {
			return
		}
		c.remove(element.Value.(*Entry).Key)
	}

	lruEntry := &Entry{Key: key, Value: value}
	listElement := c.order.PushFront(lruEntry)
	c.data[key] = listElement
}

// Read retrieves the value of the key in the Cache
func (c *Cache) Read(key interface{}) (interface{}, bool) {
	listElement, exists := c.data[key]
	if !exists {
		return nil, false
	}

	c.order.MoveToFront(listElement)
	lruEntry := listElement.Value.(*Entry)
	return lruEntry.Value, true
}

// Remove removes an Entry in the Cache by key
func (c *Cache) Remove(key interface{}) {
	c.remove(key)
}

// Size gets the current size of the Cache
func (c *Cache) Size() int {
	return len(c.data)
}

func (c *Cache) remove(key interface{}) {
	listElement, exists := c.data[key]
	if !exists {
		return
	}

	c.order.Remove(listElement)
	delete(c.data, key)
}

// String defines an string representation of a Cache
func (c Cache) String() string {
	return fmt.Sprintf("<LRU Cache %d/%d>\n", len(c.data), c.capacity)
}
