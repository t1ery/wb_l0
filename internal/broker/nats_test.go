package broker

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/t1ery/wb_l0/internal/cache"
	"github.com/t1ery/wb_l0/internal/database"
	"github.com/t1ery/wb_l0/internal/models"
	"github.com/t1ery/wb_l0/internal/service"
	"github.com/t1ery/wb_l0/test"
	"testing"
)

func TestSubscribeToNATS(t *testing.T) {
	// Подготовка к тестированию брокера
	// Создаем экземпляр сервиса и кеша
	db := test.InitTestDatabase(t)
	defer db.Close()
	pgDB := database.NewPostgresDB(db)
	cache := cache.NewCache()
	s := service.NewService(pgDB, cache)

	// Инициализируем брокер
	InitNATS()
	SubscribeToNATS(s)

	// Публикуем тестовый заказ в брокер
	PublishOrderToNATS(test.TestOrder)

	// Проверяем, что данные заказа были сохранены в базе данных и кеше
	data, err := s.GetDataByID(test.TestOrder.OrderUID)
	assert.NoError(t, err)
	assert.Equal(t, test.TestOrder, data)
}

// Вспомогательная функция для публикации заказа в брокер
func PublishOrderToNATS(order models.OrderJSON) {
	payload, err := json.Marshal(order)
	if err != nil {
		fmt.Println("Ошибка при преобразовании в JSON:", err)
		return
	}

	// Выводим информацию о том, что мы собираемся отправить
	fmt.Printf("Отправляем заказ в NATS: %s\n", payload)

	err = Nc.Publish("orders", payload)
	if err != nil {
		fmt.Println("Ошибка при отправке данных в канал:", err)
	} else {
		fmt.Println("Данные успешно отправлены в канал")
	}
}
