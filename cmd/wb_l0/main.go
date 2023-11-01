package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"sync"

	"github.com/t1ery/wb_l0/internal/broker"
	"github.com/t1ery/wb_l0/internal/cache"
	"github.com/t1ery/wb_l0/internal/database"
	"github.com/t1ery/wb_l0/internal/server"
	"github.com/t1ery/wb_l0/internal/service"
)

func main() {

	var wg sync.WaitGroup

	// Инициализируем базу данных PostgreSQL.
	db := database.InitDB()
	defer db.Close()

	// Создаем таблицы в базе данных, если они еще не созданы.
	if err := database.CreateTables(db); err != nil {
		log.Fatal("Ошибка при создании таблиц в базе данных:", err)
	}

	// Инициализируем кеш для хранения данных.
	cache := cache.NewCache()

	// Создаем сервис для работы с данными.
	dataService := service.NewService(database.NewPostgresDB(db), cache)

	// При старте сервиса восстанавливаем данные из базы данных в кеш.
	data, err := dataService.GetAllDataFromDB()
	if err != nil {
		log.Fatal("Ошибка при чтении данных из базы данных:", err)
	}
	for _, orderData := range data {
		cache.SetByID(orderData.OrderUID, orderData)
	}

	// Инициализируем NATS
	broker.InitNATS()

	// Подписываемся на канал NATS для получения данных
	broker.SubscribeToNATS(dataService, &wg)

	// Регистрируем обработчик HTTP для корневого пути.
	http.Handle("/", server.NewHTTPServer(dataService))

	// Запускаем HTTP-сервер.
	addr := ":8080"
	fmt.Println("Сервер запущен и слушает на", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
