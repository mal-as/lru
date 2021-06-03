package client

import (
	"errors"
	"fmt"

	"github.com/tarantool/go-tarantool"
	"gopkg.in/vmihailenco/msgpack.v2"
)

// ErrNotFound - значение для данного ключа не содержится в кэше
var ErrNotFound = errors.New("value was not found")

// данные для отправки в тарантул (установка значения для ключа)
type requestData struct {
	Key   string
	Value []byte
}

// EncodeMsgpack - преобразование в msgPack
func (d *requestData) EncodeMsgpack(e *msgpack.Encoder) error {
	if err := e.EncodeArrayLen(2); err != nil {
		return err
	}

	if err := e.EncodeString(d.Key); err != nil {
		return err
	}

	if err := e.EncodeBytes(d.Value); err != nil {
		return err
	}

	return nil
}

// данные полученные из тарантула (спейса кэша)
type responseData struct {
	ID    int
	Key   string
	Value string
}

// DecodeMsgpack декодирует данные из msgPack
func (rd *responseData) DecodeMsgpack(d *msgpack.Decoder) error {
	var err error
	var l int
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}

	switch l {
	case 3:
		break
	case 1:
		if err := d.DecodeNil(); err != nil {
			return err
		}
		return ErrNotFound
	default:
		return fmt.Errorf("array len doesn't match: %d", l)
	}

	if rd.ID, err = d.DecodeInt(); err != nil {
		return err
	}
	if rd.Key, err = d.DecodeString(); err != nil {
		return err
	}
	if rd.Value, err = d.DecodeString(); err != nil {
		return err
	}
	return nil
}

// функция, параметризирующая инстанс кэша
type option func(c *Cache) error

// WithTruncate очищает кэш перед созданием нового инстанса
var WithTruncate = func(c *Cache) error {
	return c.Truncate()
}

// Cache - структура для установки значений и получения их из кэша
type Cache struct {
	conn *tarantool.Connection
}

// NewCache - конструктор для Cache
func NewCache(addr string, config tarantool.Opts, size int, params ...option) (*Cache, error) {
	conn, err := tarantool.Connect(addr, config)
	if err != nil {
		return nil, err
	}

	c := &Cache{conn: conn}

	if err = c.setSize(size); err != nil {
		return nil, err
	}

	if len(params) > 0 {
		for _, param := range params {
			if err := param(c); err != nil {
				return nil, err
			}
		}
	}

	return c, nil
}

// Set устанавливает значение в кэше
func (c *Cache) Set(key string, value []byte) error {
	tuple := &requestData{Key: key, Value: value}
	_, err := c.conn.Call("set", tuple)
	return err
}

// Get получает значение из кэша
func (c *Cache) Get(key string) ([]byte, error) {
	var tuples []responseData

	if err := c.conn.CallTyped("get", []interface{}{key}, &tuples); err != nil {
		return nil, err
	}

	if len(tuples) == 0 {
		return nil, ErrNotFound
	}

	return []byte(tuples[0].Value), nil
}

// Truncate очищает кэш
func (c *Cache) Truncate() error {
	_, err := c.conn.Call("truncate", []interface{}{})
	return err
}

// Close закрывает соединение с тарантулом
func (c *Cache) Close() error {
	return c.conn.Close()
}

// устанавлиевает размер кэша
func (c *Cache) setSize(size int) error {
	_, err := c.conn.Call("setCacheSize", []interface{}{size})
	return err
}
