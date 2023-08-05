package database

import (
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/t1ery/wb_l0/test"
	"testing"
)

func TestPostgresDB_SaveAndGetDataByID(t *testing.T) {
	// Подготовка к тестированию базы данных
	db := test.InitTestDatabase(t)
	defer db.Close()

	// Создаем экземпляр PostgresDB
	pgDB := NewPostgresDB(db)

	// Сохраняем заказ в базу данных
	err := pgDB.SaveData(test.TestOrder)
	assert.NoError(t, err)

	// Получаем заказ по его идентификатору
	data, err := pgDB.GetDataByID(test.TestOrder.OrderUID)
	assert.NoError(t, err)
	assert.Equal(t, test.TestOrder, data)
}
