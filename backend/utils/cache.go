package utils

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"math"
	"os"
	"path/filepath"
	"time"
)

type CacheValue struct {
	Value      interface{} `json:"value"`
	Expiration int         `json:"expiration"`
}

type Cache struct {
	value map[string]CacheValue
}

func NewCache(filename string) *Cache {
	cache := &Cache{
		value: make(map[string]CacheValue),
	}

	if filename == "" {
		log.Warn("Cache filename is empty, using in-memory cache only.")
		return cache
	}

	if _, err := os.Stat(filepath.Dir(filename)); os.IsNotExist(err) {
		log.Warnf("Cache directory %s does not exist.", filepath.Dir(filename))
		return cache
	}

	if FileExists(filename) == false {
		err := FileTouch(filename, "{}")
		if err != nil {
			log.Errorf("Failed to create cache file %s: %v", filename, err)
			return cache
		}
	}

	if err := cache.Load(filename); err != nil {
		log.Errorf("Failed to load cache from file %s: %v", filename, err)
		return cache
	}

	return cache
}

func (c *Cache) HasKey(key string) (bool, error) {
	if key == "" {
		return false, nil // No error, but key is empty
	}
	_, exists := c.value[key]
	return exists, nil
}

func (c *Cache) Get(key string) (interface{}, error) {
	if key == "" {
		return nil, nil
	}

	cacheValue, exists := c.value[key]
	if !exists {
		return nil, nil
	}

	if cacheValue.Expiration > 0 && cacheValue.Expiration < int(time.Now().Unix()) {
		delete(c.value, key)
		return nil, nil
	}

	return cacheValue.Value, nil
}

func (c *Cache) Set(key string, value interface{}, expiration int) error {
	if expiration >= math.MaxInt {
		log.Warnf("Expiration time %d exceeds maximum allowed value, setting to 0", expiration)
		expiration = 0
	}
	expireTime := 0
	if expiration > 0 {
		expireTime = int(time.Now().Unix()) + expiration
	}

	c.value[key] = CacheValue{
		Value:      value,
		Expiration: expireTime,
	}

	return nil
}

func (c *Cache) Delete(key string) error {
	if _, exists := c.value[key]; exists {
		delete(c.value, key)
		return nil
	}
	return nil
}

func (c *Cache) Clear() error {
	c.value = make(map[string]CacheValue)
	return nil
}

func (c *Cache) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("error closing cache file: %v", err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(c.value); err != nil {
		return err
	}
	return nil
}

func (c *Cache) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("error closing cache file: %v", err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&c.value); err != nil {
		return err
	}
	return nil
}
