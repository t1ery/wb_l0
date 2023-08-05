package server

import (
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/t1ery/wb_l0/internal/cache"
	"github.com/t1ery/wb_l0/internal/database"
	"github.com/t1ery/wb_l0/internal/service"
	"github.com/t1ery/wb_l0/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPHandler_GetOrder(t *testing.T) {
	// Подготовка к тестированию HTTP-сервера
	db := test.InitTestDatabase(t)
	defer db.Close()
	pgDB := database.NewPostgresDB(db)
	cache := cache.NewCache()
	s := service.NewService(pgDB, cache)

	// Создаем экземпляр HTTP-сервера
	handler := NewHTTPServer(s)

	// Создаем тестовый заказ и сохраняем его в базу данных и кеш
	err := s.SaveData(test.TestOrder)
	assert.NoError(t, err)

	// Создаем тестовый HTTP-запрос
	url := "/get_data?id=" + test.TestOrder.OrderUID
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	// Создаем HTTP Response Recorder
	rr := httptest.NewRecorder()

	// Обрабатываем HTTP-запрос с помощью HTTP-обработчика
	handler.ServeHTTP(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем содержимое ответа
	expectedResponse := `{"order_uid":"b563feb7b2b84b6test","track_number":"WBILMTESTTRACK","entry":"WBIL","delivery":{"name":"Test Testov","phone":"+9720000000","zip":"2639809","city":"Kiryat Mozkin","address":"Ploshad Mira 15","region":"Kraiot","email":"test@gmail.com"},"payment":{"transaction":"b563feb7b2b84b6test","request_id":"","currency":"USD","provider":"wbpay","amount":1817,"payment_dt":1637907727,"bank":"alpha","delivery_cost":1500,"goods_total":317,"custom_fee":0},"items":[{"chrt_id":9934930,"track_number":"WBILMTESTTRACK","price":453,"rid":"ab4219087a764ae0btest","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,"brand":"Vivienne Sabo","status":202}],"locale":"en","internal_signature":"","customer_id":"test","delivery_service":"meest","shardkey":"9","sm_id":99,"date_created":"2021-11-26T06:22:19Z","oof_shard":"1"}
`

	assert.Equal(t, expectedResponse, rr.Body.String())
}
