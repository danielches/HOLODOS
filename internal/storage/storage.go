package storage

import (
	"HOLODOS/internal/models"
	"errors"
	"sync"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type MemoryStorage struct {
	mu       sync.RWMutex
	products map[string]models.Product
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		products: make(map[string]models.Product),
	}
}

func (s *MemoryStorage) AddProduct(product models.Product) (models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.products[product.ID] = product
	return product, nil
}

func (s *MemoryStorage) GetProduct(id string) (models.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	product, exists := s.products[id]
	if !exists {
		return models.Product{}, ErrProductNotFound
	}

	return product, nil
}

func (s *MemoryStorage) ListProducts() ([]models.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	products := make([]models.Product, 0, len(s.products))
	for _, p := range s.products {
		products = append(products, p)
	}

	return products, nil
}

func (s *MemoryStorage) RemoveProduct(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.products[id]; !exists {
		return ErrProductNotFound
	}

	delete(s.products, id)
	return nil
}
