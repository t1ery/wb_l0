package cache

import (
	"github.com/t1ery/wb_l0/internal/models"
	"sync"
)

// Cache представляет кеш для хранения данных заказов в оперативной памяти.
type Cache struct {
	mu   sync.Mutex
	data map[string]models.OrderJSON
}

// NewCache создает новый экземпляр кеша и возвращает указатель на него.
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]models.OrderJSON),
	}
}

// GetByID возвращает данные заказа по его идентификатору.
// Если данные для указанного идентификатора отсутствуют в кеше, возвращается второй аргумент со значением false.
func (c *Cache) GetByID(id string) (models.OrderJSON, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, ok := c.data[id]
	return data, ok
}

// SetByID сохраняет данные заказа в кеше по его идентификатору.
// Если данные для указанного идентификатора уже присутствуют в кеше, они будут заменены новыми данными.
func (c *Cache) SetByID(id string, data models.OrderJSON) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[id] = data
}
