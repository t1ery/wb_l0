package database

import (
	"database/sql"
	"fmt"
	"github.com/t1ery/wb_l0/internal/models"
	"log"
)

// PostgresDB представляет базу данных PostgreSQL.
type PostgresDB struct {
	db *sql.DB
}

// NewPostgresDB создает новый экземпляр базы данных PostgreSQL и возвращает указатель на него.
func NewPostgresDB(db *sql.DB) *PostgresDB {
	return &PostgresDB{
		db: db,
	}
}

// SaveData сохраняет данные заказа в базе данных. Реализация в файле service.go
func (pg *PostgresDB) SaveData(data models.OrderJSON) error {
	return SaveToDB(pg.db, data)
}

// InitDB выполняет инициализацию и подключение к базе данных PostgreSQL и возвращает указатель на созданное подключение.
func InitDB() *sql.DB {
	host := "my_postgres"
	port := "5432"
	user := "postgres"
	password := "481516"
	dbname := "wb_l0"

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// SaveToDB сохраняет данные заказа в базе данных PostgreSQL, используя транзакции.
func SaveToDB(db *sql.DB, data models.OrderJSON) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		data.OrderUID, data.TrackNumber, data.Entry, data.Locale, data.InternalSignature, data.CustomerID, data.DeliveryService, data.ShardKey, data.SmID, data.DateCreated, data.OofShard)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		data.OrderUID, data.Delivery.Name, data.Delivery.Phone, data.Delivery.Zip, data.Delivery.City, data.Delivery.Address, data.Delivery.Region, data.Delivery.Email)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		data.OrderUID, data.Payment.Transaction, data.Payment.RequestID, data.Payment.Currency, data.Payment.Provider, data.Payment.Amount, data.Payment.PaymentDt, data.Payment.Bank, data.Payment.DeliveryCost, data.Payment.GoodsTotal, data.Payment.CustomFee)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range data.Items {
		_, err = tx.Exec("INSERT INTO item (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			data.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// GetDataByID получает данные заказа по его идентификатору из базы данных PostgreSQL.
func (pg *PostgresDB) GetDataByID(id string) (models.OrderJSON, error) {
	row := pg.db.QueryRow("SELECT track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders WHERE order_uid=$1", id)
	var orderData models.OrderJSON
	var trackNumber, entry, locale, internalSignature, customerID, deliveryService, shardKey, dateCreated, oofShard string
	var smID int
	err := row.Scan(&trackNumber, &entry, &locale, &internalSignature, &customerID, &deliveryService, &shardKey, &smID, &dateCreated, &oofShard)
	if err != nil {
		return models.OrderJSON{}, err
	}

	orderData.OrderUID = id
	orderData.TrackNumber = trackNumber
	orderData.Entry = entry
	orderData.Locale = locale
	orderData.InternalSignature = internalSignature
	orderData.CustomerID = customerID
	orderData.DeliveryService = deliveryService
	orderData.ShardKey = shardKey
	orderData.SmID = smID
	orderData.DateCreated = dateCreated
	orderData.OofShard = oofShard

	row = pg.db.QueryRow("SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid=$1", id)
	var deliveryData models.Delivery
	err = row.Scan(&deliveryData.Name, &deliveryData.Phone, &deliveryData.Zip, &deliveryData.City, &deliveryData.Address, &deliveryData.Region, &deliveryData.Email)
	if err != nil {
		return models.OrderJSON{}, err
	}

	orderData.Delivery = deliveryData

	row = pg.db.QueryRow("SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid=$1", id)
	var paymentData models.Payment
	err = row.Scan(&paymentData.Transaction, &paymentData.RequestID, &paymentData.Currency, &paymentData.Provider, &paymentData.Amount, &paymentData.PaymentDt, &paymentData.Bank, &paymentData.DeliveryCost, &paymentData.GoodsTotal, &paymentData.CustomFee)
	if err != nil {
		return models.OrderJSON{}, err
	}

	orderData.Payment = paymentData

	rows, err := pg.db.Query("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM item WHERE order_uid=$1", id)
	if err != nil {
		return models.OrderJSON{}, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var itemData models.Item
		err := rows.Scan(&itemData.ChrtID, &itemData.TrackNumber, &itemData.Price, &itemData.RID, &itemData.Name, &itemData.Sale, &itemData.Size, &itemData.TotalPrice, &itemData.NmID, &itemData.Brand, &itemData.Status)
		if err != nil {
			return models.OrderJSON{}, err
		}
		items = append(items, itemData)
	}

	orderData.Items = items

	return orderData, nil
}

// GetAllData получает все данные заказов из базы данных PostgreSQL.
func (pg *PostgresDB) GetAllData() ([]models.OrderJSON, error) {
	rows, err := pg.db.Query("SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []models.OrderJSON
	for rows.Next() {
		var orderData models.OrderJSON
		if err := rows.Scan(
			&orderData.OrderUID,
			&orderData.TrackNumber,
			&orderData.Entry,
			&orderData.Locale,
			&orderData.InternalSignature,
			&orderData.CustomerID,
			&orderData.DeliveryService,
			&orderData.ShardKey,
			&orderData.SmID,
			&orderData.DateCreated,
			&orderData.OofShard,
		); err != nil {
			log.Printf("Ошибка при сканировании данных из базы данных: %v\n", err)
			return nil, err
		}

		deliveryRow := pg.db.QueryRow("SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1", orderData.OrderUID)
		var deliveryData models.Delivery
		if err := deliveryRow.Scan(
			&deliveryData.Name,
			&deliveryData.Phone,
			&deliveryData.Zip,
			&deliveryData.City,
			&deliveryData.Address,
			&deliveryData.Region,
			&deliveryData.Email,
		); err != nil {
			log.Printf("Ошибка при сканировании данных из таблицы delivery: %v\n", err)
			return nil, err
		}
		orderData.Delivery = deliveryData

		paymentRow := pg.db.QueryRow("SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = $1", orderData.OrderUID)
		var paymentData models.Payment
		if err := paymentRow.Scan(
			&paymentData.Transaction,
			&paymentData.RequestID,
			&paymentData.Currency,
			&paymentData.Provider,
			&paymentData.Amount,
			&paymentData.PaymentDt,
			&paymentData.Bank,
			&paymentData.DeliveryCost,
			&paymentData.GoodsTotal,
			&paymentData.CustomFee,
		); err != nil {
			log.Printf("Ошибка при сканировании данных из таблицы payment: %v\n", err)
			return nil, err
		}
		orderData.Payment = paymentData

		itemsRows, err := pg.db.Query("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM item WHERE order_uid = $1", orderData.OrderUID)
		if err != nil {
			log.Printf("Ошибка при выполнении запроса к таблице item: %v\n", err)
			return nil, err
		}
		defer itemsRows.Close()

		var items []models.Item
		for itemsRows.Next() {
			var itemData models.Item
			if err := itemsRows.Scan(
				&itemData.ChrtID,
				&itemData.TrackNumber,
				&itemData.Price,
				&itemData.RID,
				&itemData.Name,
				&itemData.Sale,
				&itemData.Size,
				&itemData.TotalPrice,
				&itemData.NmID,
				&itemData.Brand,
				&itemData.Status,
			); err != nil {
				log.Printf("Ошибка при сканировании данных из таблицы item: %v\n", err)
				return nil, err
			}
			items = append(items, itemData)
		}
		orderData.Items = items

		data = append(data, orderData)
	}

	return data, nil
}
