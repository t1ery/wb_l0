package broker

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/t1ery/wb_l0/internal/models"
	"github.com/t1ery/wb_l0/internal/service"
	"log"
	"sync"
)

var Nc *nats.Conn
var _ *sync.WaitGroup

// InitNATS устанавливает соединение с сервером NATS.
func InitNATS() {
	var err error
	Nc, err = nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
}

// SubscribeToNATS подписывается на указанный канал в сервере NATS и обрабатывает полученные сообщения.
// Когда сообщение получено, оно декодируется из JSON в объект OrderJSON.
// Затем данные проходят валидацию, и если они проходят проверку, они сохраняются в кеш и базу данных через сервис s.
// Если происходит ошибка на любом из этапов, она регистрируется в журнале.
func SubscribeToNATS(s *service.Service, wg *sync.WaitGroup) {
	Nc.Subscribe("orders", func(msg *nats.Msg) {
		defer wg.Done() // Указываем, что операция завершена
		var orderJSON models.OrderJSON
		err := json.Unmarshal(msg.Data, &orderJSON)
		if err != nil {
			fmt.Printf("Ошибка при разборе JSON: %v\n", err)
			return
		}

		fmt.Println("Получен заказ из NATS:", orderJSON)

		err = orderJSON.Validate()
		if err != nil {
			fmt.Printf("Ошибка валидации данных: %v\n", err)
			return
		}

		fmt.Println("Данные успешно прошли валидацию")

		err = s.SaveData(orderJSON)
		if err != nil {
			fmt.Printf("Ошибка при сохранении данных: %v\n", err)
		} else {
			fmt.Println("Данные успешно сохранены")
		}
		fmt.Println("Функция успешно отработала")
	})
}
