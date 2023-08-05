package cache

import (
	"github.com/stretchr/testify/assert"
	"github.com/t1ery/wb_l0/test"
	"testing"
)

func TestCache_GetAndSetByID(t *testing.T) {
	// Создаем экземпляр Cache
	cache := NewCache()

	// Сохраняем заказ в кеш
	cache.SetByID(test.TestOrder.OrderUID, test.TestOrder)

	// Получаем заказ по его идентификатору из кеша
	data, ok := cache.GetByID(test.TestOrder.OrderUID)
	assert.True(t, ok)
	assert.Equal(t, test.TestOrder, data)
}
