package mem

// Cache - структура для кэша в памяти
type Cache struct {
	length int
	list   *list
	data   map[string]*node
}

// NewCache - конструктор для Cache
func NewCache(n int) *Cache {
	return &Cache{
		length: n,
		list:   newList(),
		data:   make(map[string]*node, n),
	}
}

// Set устанавливает значение в кэше
func (c *Cache) Set(key string, value []byte) error {
	if n, ok := c.data[key]; ok {
		n.value = pair{key: key, data: value}
		c.list.moveToFront(n)
		return nil
	}

	newNode := &node{value: pair{key: key, data: value}}
	c.data[key] = newNode
	c.list.insesrtFront(newNode)

	if c.list.length > c.length {
		last := c.list.retriveLastElement()

		delete(c.data, last.value.key)
		c.list.pop(c.length)
	}

	return nil
}

// Get получает значение из кэша
func (c *Cache) Get(key string) ([]byte, error) {
	n, ok := c.data[key]
	if ok {
		c.list.moveToFront(n)
		return n.value.data, nil
	}

	return nil, nil
}

// Print печатает связный список в основе кэша
func (c *Cache) Print() {
	c.list.print()
}
