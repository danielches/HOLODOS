package service

import (
	"HOLODOS/internal/models"
	"HOLODOS/internal/storage"
	"errors"
	"github.com/google/uuid"
	"time"
)

type FridgeService struct {
	storage *storage.MemoryStorage // Храним указатель
}

func NewFridgeService(storage *storage.MemoryStorage) *FridgeService {
	return &FridgeService{storage: storage}
}

func (s *FridgeService) AddProduct(product models.Product) (models.Product, error) {
	if product.Name == "" || product.Quantity <= 0 {
		return models.Product{}, errors.New("invalid product data")
	}

	product.ID = generateID()
	product.DateAdded = time.Now()

	return s.storage.AddProduct(product)
}

func (s *FridgeService) GetProduct(id string) (models.Product, error) {
	return s.storage.GetProduct(id)
}

func (s *FridgeService) ListProducts() ([]models.Product, error) {
	return s.storage.ListProducts()
}

func (s *FridgeService) RemoveProduct(id string) error {
	return s.storage.RemoveProduct(id)
}

func (s *FridgeService) CheckProductExpiry(id string) (models.Product, bool, int, error) {
	product, err := s.storage.GetProduct(id)
	if err != nil {
		return models.Product{}, false, 0, err
	}

	now := time.Now()
	if product.ExpiryDate.IsZero() {
		return product, false, 0, nil
	}

	daysRemaining := int(product.ExpiryDate.Sub(now).Hours() / 24)
	isExpired := now.After(product.ExpiryDate)

	return product, isExpired, daysRemaining, nil
}

func (s *FridgeService) GetExpiringProducts(daysThreshold int) ([]models.Product, error) {
	products, err := s.storage.ListProducts()
	if err != nil {
		return nil, err
	}

	var expiring []models.Product
	now := time.Now()

	for _, p := range products {
		if p.ExpiryDate.IsZero() {
			continue
		}

		daysRemaining := int(p.ExpiryDate.Sub(now).Hours() / 24)
		if daysRemaining <= daysThreshold && daysRemaining >= 0 {
			expiring = append(expiring, p)
		}
	}

	return expiring, nil
}

func generateID() string {
	return uuid.New().String()
}
