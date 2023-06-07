package cache

import "sync"

// Cache хранилище данных в памяти.
type Cache[k comparable, v any] struct {
	sync.RWMutex
	cache map[k]v
}

// NewCache инициализирует новый кеш в памяти.
func NewCache[k comparable, v any]() *Cache[k, v] {
	return &Cache[k, v]{
		cache: make(map[k]v),
	}
}

// Set устанавливает значение.
func (s *Cache[k, v]) Set(key k, value v) {
	s.Lock()
	defer s.Unlock()
	s.cache[key] = value
}

// Get получает значение.
func (s *Cache[k, v]) Get(key k) (v, bool) {
	s.RLock()
	defer s.RUnlock()
	item, found := s.cache[key]

	return item, found
}

// Exists проверяет, существует ли ключ.
func (s *Cache[k, v]) Exists(key k) bool {
	s.RLock()
	defer s.RUnlock()
	_, found := s.cache[key]

	return found
}

// Delete удаляет значение.
func (s *Cache[k, v]) Delete(key k) {
	s.Lock()
	defer s.Unlock()
	delete(s.cache, key)
}

// Clear очищает кеш.
func (s *Cache[k, v]) Clear() {
	s.Lock()
	defer s.Unlock()
	s.cache = make(map[k]v)
}

// Load загружает кеш из карты.
func (s *Cache[k, v]) Load(list map[k]v) int {
	s.Lock()
	defer s.Unlock()
	s.cache = list

	return len(s.cache)
}

// Append добавляет карту к кешу.
func (s *Cache[k, v]) Append(list map[k]v) int {
	s.Lock()
	defer s.Unlock()
	for key, val := range list {
		s.cache[key] = val
	}

	return len(s.cache)
}

// Count количество данных в карте.
func (s *Cache[k, v]) Count() int {
	s.RLock()
	defer s.RUnlock()

	return len(s.cache)
}

// ToList возвращает значения в виде массива.
func (s *Cache[k, v]) ToList() []v {
	s.RLock()
	defer s.RUnlock()

	result := make([]v, 0, len(s.cache))
	for _, val := range s.cache {
		result = append(result, val)
	}

	return result
}
