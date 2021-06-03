package mem

import (
	"testing"
)

func TestCache(t *testing.T) {
	size := 3
	c := NewCache(size)

	c.Set("foo", []byte("bar"))
	c.Set("key", []byte("value"))
	c.Set("some", []byte("thing"))

	if data, _ := c.Get("foo"); string(data) != "bar" {
		t.Fatalf("ожидалось значение \"bar\", получено - %s", string(data))
	}

	c.Set("year", []byte("2021"))

	if data, _ := c.Get("key"); data != nil {
		t.Fatalf("получено значение %s для ключа, которого не должно быть в кэше", string(data))
	}

	c.Set("some", []byte("new thing"))

	if data, _ := c.Get("some"); string(data) != "new thing" {
		t.Fatalf("ожидалось значение \"new thing\", получено - %s", string(data))
	}

	c.Print()
}
