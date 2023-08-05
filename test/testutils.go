package test

import (
	"database/sql"
	"github.com/t1ery/wb_l0/internal/models"
	"testing"
)

// TestOrder Объявим как экземпляр данных для тестов
var TestOrder = models.OrderJSON{
	OrderUID:    "b563feb7b2b84b6test",
	TrackNumber: "WBILMTESTTRACK",
	Entry:       "WBIL",
	Delivery: models.Delivery{
		Name:    "Test Testov",
		Phone:   "+9720000000",
		Zip:     "2639809",
		City:    "Kiryat Mozkin",
		Address: "Ploshad Mira 15",
		Region:  "Kraiot",
		Email:   "test@gmail.com",
	},
	Payment: models.Payment{
		Transaction:  "b563feb7b2b84b6test",
		RequestID:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       1817,
		PaymentDt:    1637907727,
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   317,
		CustomFee:    0,
	},
	Items: []models.Item{
		{
			ChrtID:      9934930,
			TrackNumber: "WBILMTESTTRACK",
			Price:       453,
			RID:         "ab4219087a764ae0btest",
			Name:        "Mascaras",
			Sale:        30,
			Size:        "0",
			TotalPrice:  317,
			NmID:        2389212,
			Brand:       "Vivienne Sabo",
			Status:      202,
		},
	},
	Locale:            "en",
	InternalSignature: "",
	CustomerID:        "test",
	DeliveryService:   "meest",
	ShardKey:          "9",
	SmID:              99,
	DateCreated:       "2021-11-26T06:22:19Z",
	OofShard:          "1",
}

// Вспомогательная функция для создания временной базы данных для тестирования
func InitTestDatabase(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=481516 dbname=wb_l0_test sslmode=disable")
	if err != nil {
		t.Fatalf("Ошибка при создании временной базы данных: %v", err)
	}
	// Очистка и создание таблиц перед каждым тестом
	_, err = db.Exec("DROP TABLE IF EXISTS orders, delivery, payment, item")
	if err != nil {
		t.Fatalf("Ошибка при очистке таблиц: %v", err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS orders (\n\torder_uid VARCHAR(50) PRIMARY KEY,\n\ttrack_number VARCHAR(50) NOT NULL,\n\tentry VARCHAR(50) NOT NULL,\n\tlocale VARCHAR(10) NOT NULL,\n\tinternal_signature VARCHAR(100),\n\tcustomer_id VARCHAR(50) NOT NULL,\n\tdelivery_service VARCHAR(50) NOT NULL,\n\tshardkey VARCHAR(50) NOT NULL,\n\tsm_id INT NOT NULL,\n\tdate_created VARCHAR(50) NOT NULL,\n\toof_shard VARCHAR(50) NOT NULL\n);\n\nCREATE TABLE IF NOT EXISTS delivery (\n\torder_uid VARCHAR(50) PRIMARY KEY REFERENCES orders(order_uid),\n\tname VARCHAR(100) NOT NULL,\n\tphone VARCHAR(20) NOT NULL,\n\tzip VARCHAR(10) NOT NULL,\n\tcity VARCHAR(100) NOT NULL,\n\taddress VARCHAR(100) NOT NULL,\n\tregion VARCHAR(100) NOT NULL,\n\temail VARCHAR(100) NOT NULL\n);\n\nCREATE TABLE IF NOT EXISTS payment (\n\torder_uid VARCHAR(50) PRIMARY KEY REFERENCES orders(order_uid),\n\ttransaction VARCHAR(100) NOT NULL,\n\trequest_id VARCHAR(100),\n\tcurrency VARCHAR(10) NOT NULL,\n\tprovider VARCHAR(50) NOT NULL,\n\tamount INT NOT NULL,\n\tpayment_dt INT NOT NULL,\n\tbank VARCHAR(50) NOT NULL,\n\tdelivery_cost INT NOT NULL,\n\tgoods_total INT NOT NULL,\n\tcustom_fee INT\n);\n\nCREATE TABLE IF NOT EXISTS item (\n    order_uid VARCHAR(50) REFERENCES orders(order_uid),\n\tchrt_id INT NOT NULL,\n\ttrack_number VARCHAR(50) NOT NULL,\n\tprice INT NOT NULL,\n\trid VARCHAR(100) NOT NULL,\n\tname VARCHAR(100) NOT NULL,\n\tsale INT NOT NULL,\n\tsize VARCHAR(50) NOT NULL,\n\ttotal_price INT NOT NULL,\n\tnm_id INT NOT NULL,\n\tbrand VARCHAR(100) NOT NULL,\n\tstatus INT NOT NULL,\n    PRIMARY KEY (order_uid, chrt_id)\n);")
	if err != nil {
		t.Fatalf("Ошибка при создании таблицы orders: %v", err)
	}
	// Создайте другие таблицы, если необходимо

	return db
}
