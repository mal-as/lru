package client

import (
	"errors"
	"testing"

	"github.com/tarantool/go-tarantool"
)

func TestClient(t *testing.T) {
	addr := "localhost:3301"
	conf := tarantool.Opts{
		User: "go",
		Pass: "passwd",
	}
	size := 3

	c, err := NewCache(addr, conf, size, WithTruncate)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	testValues := map[string]string{
		"foo":  "bar",
		"key":  "value",
		"test": "thing",
	}

	for key, value := range testValues {
		if err := c.Set(key, []byte(value)); err != nil {
			t.Fatal(err)
		}
	}

	data, err := c.Get("foo")
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "bar" {
		t.Fatalf("ожидалось значение \"bar\", получено - %s", string(data))
	}

	if err := c.Set("year", []byte("2021")); err != nil {
		t.Fatal(err)
	}

	if err := c.Set("month", []byte("june")); err != nil {
		t.Fatal(err)
	}

	value, err := c.Get("key")
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			t.Fatal(err)
		}
	} else {
		t.Fatalf("получено значение %s для ключа, которого не должно быть в кэше", string(value))
	}

	if err := c.Set("year", []byte("2022")); err != nil {
		t.Fatal(err)
	}

	data, err = c.Get("year")
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "2022" {
		t.Fatalf("ожидалось значение \"2022\", получено - %s", string(data))
	}
}
