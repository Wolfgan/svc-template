package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData = map[string]string{
	"key1": "value1",
	"key2": "value3",
	"key3": "value3",
}

func TestCache_Load(t *testing.T) {
	cache := NewCache[string, string]()
	count := cache.Load(testData)
	assert.Equal(t, 3, count, "The two numbers should be the same")
}

func TestCache_Clear(t *testing.T) {
	cache := NewCache[string, string]()
	cache.Load(testData)
	cache.Clear()
	assert.Equal(t, 0, cache.Count(), "values must be empty")
}

func TestCache_Get(t *testing.T) {
	cache := NewCache[string, string]()
	cache.Load(testData)

	value, found := cache.Get("key3")
	assert.Equal(t, true, found, "key not found")
	assert.Equal(t, testData["key3"], value, "values not equal")
}

func TestCache_Set(t *testing.T) {
	cache := NewCache[string, string]()
	cache.Set("key_test", "value_test")
	value, found := cache.Get("key_test")
	assert.Equal(t, true, found, "Value not added")
	assert.Equal(t, "value_test", value, "Error added value")
}

func TestCache_Exists(t *testing.T) {
	cache := NewCache[string, string]()
	cache.Set("key_test", "value_test")
	found := cache.Exists("key_test")
	assert.Equal(t, true, found, "Value not exists")
}

func TestCache_Delete(t *testing.T) {
	key := "key_test"
	cache := NewCache[string, string]()
	cache.Set(key, "value_test")
	cache.Delete(key)
	_, found := cache.Get(key)
	assert.Equal(t, false, found, "The value has not been deleted")
}

func TestCache_gorutine(t *testing.T) {
	cache := NewCache[string, string]()

	go func() {
		for i := 0; i < 10; i++ {
			cache.Set("key_test", "value_test")
		}
	}()

	var value string
	var found bool
	for !found {
		value, found = cache.Get("key_test")
	}

	assert.Equal(t, true, found, "Value not added")
	assert.Equal(t, "value_test", value, "Error added value")
}
