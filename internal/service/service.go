package service

import (
	"github.com/t1ery/wb_l0/internal/cache"
	"github.com/t1ery/wb_l0/internal/database"
	"github.com/t1ery/wb_l0/internal/models"
)

// Service представляет сервис для работы с данными заказов.
type Service struct {
	db    *database.PostgresDB
	cache *cache.Cache
}

// NewService создает новый экземпляр сервиса и инициализирует его зависимости.
func NewService(db *database.PostgresDB, cache *cache.Cache) *Service {
	return &Service{
		db:    db,
		cache: cache,
	}
}

// GetAllDataFromDB возвращает все данные заказов из базы данных.
func (s *Service) GetAllDataFromDB() ([]models.OrderJSON, error) {
	return s.db.GetAllData()
}

// GetDataByID получает данные заказа по его идентификатору.
// Сначала ищем данные в кеше, если они есть - возвращаем их.
// Если данных нет в кеше, обращаемся к базе данных, сохраняем их в кеш и возвращаем.
func (s *Service) GetDataByID(id string) (models.OrderJSON, error) {
	data, ok := s.cache.GetByID(id)
	if ok {
		return data, nil
	}

	data, err := s.db.GetDataByID(id)
	if err != nil {
		return models.OrderJSON{}, err
	}

	s.cache.SetByID(id, data)
	return data, nil
}

// SaveData сохраняет данные заказа.
// Проверяем, есть ли данные в кеше. Если есть, не сохраняем их повторно.
// Если данных нет в кеше, сохраняем их в базе данных и обновляем кеш.
func (s *Service) SaveData(data models.OrderJSON) error {
	_, ok := s.cache.GetByID(data.OrderUID)
	if ok {
		// Данные уже есть в кеше, не сохраняем их повторно
		return nil
	}

	if err := s.db.SaveData(data); err != nil {
		return err
	}

	s.cache.SetByID(data.OrderUID, data)
	return nil
}
