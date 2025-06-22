package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func ValidateAndGetTimeFromMap(key string, obj map[string]interface{}) (time.Time, error) {
	if obj == nil {
		return time.Time{}, fmt.Errorf("object is nil")
	}

	if !ValidateTimeFromMap(key, obj) {
		return time.Time{}, fmt.Errorf("key %s is missing or not a valid time", key)
	}

	return GetTimeFromMap(key, obj), nil
}

func ValidateTimeFromMap(key string, obj map[string]interface{}) bool {
	if obj == nil {
		return false
	}

	value, exists := obj[key]
	if !exists {
		log.Warnf("Key %s does not exist in the object", key)
		return false
	}

	switch v := value.(type) {
	case time.Time:
		return true
	case string:
		_, err := time.Parse(time.RFC3339, v)
		if err != nil {
			log.Warnf("Key %s is not a valid time string: %v", key, err)
			return false
		}
		return true
	case int64, int:
		return true
	default:
		log.Warnf("Key %s is not a valid time type: %T", key, v)
	}

	return false
}

func GetTimeFromMap(key string, obj map[string]interface{}) time.Time {
	if obj == nil {
		return time.Time{}
	}

	value, exists := obj[key]
	if !exists {
		return time.Time{}
	}

	switch v := value.(type) {
	case time.Time:
		return v
	case string:
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t
		}

		if ValidateIntFromMap(key, obj) {
			value = GetIntFromMap(key, obj)
			if value != 0 {
				return time.Unix(value.(int64), 0)
			}
		}
	case int:
	case int64:
	case float32:
	case float64:
		return time.Unix(int64(v), 0)
	default:
		log.Warnf("Key %s is not a valid time type: %T", key, v)
	}

	return time.Unix(0, 0)
}

func ValidateAndGetStringFromMap(key string, obj map[string]interface{}) (string, error) {
	if obj == nil {
		return "", fmt.Errorf("object is nil")
	}

	if !ValidateStringFromMap(key, obj) {
		return "", fmt.Errorf("key %s is missing or not a string", key)
	}

	return GetStringFromMap(key, obj), nil
}

func ValidateStringFromMap(key string, obj map[string]interface{}) bool {
	if obj == nil {
		return false
	}

	value, exists := obj[key]
	if !exists {
		if MapContainsKey(key, obj) {
			return true
		}
		return false
	}

	if _, ok := value.(string); ok {
		return true
	}

	return true
}

func GetStringFromMap(key string, obj map[string]interface{}) string {
	if obj == nil {
		return ""
	}

	value, exists := obj[key]
	if !exists {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case int64:
		return strconv.FormatInt(v, 10)
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		log.Warnf("Key %s is not a string, int64, int, or float64: %T", key, v)
	}

	return ""
}

func ValidateAndGetIntFromMap(key string, obj map[string]interface{}) (int64, error) {
	if obj == nil {
		return 0, fmt.Errorf("object is nil")
	}

	if !ValidateIntFromMap(key, obj) {
		return 0, fmt.Errorf("key %s is missing or not an int64", key)
	}

	return GetIntFromMap(key, obj), nil
}

func ValidateIntFromMap(key string, obj map[string]interface{}) bool {
	if obj == nil {
		return false
	}

	value, exists := obj[key]
	if !exists {
		log.Warnf("Key %s does not exist in the object", key)
		return false
	}

	switch v := value.(type) {
	case int64:
		return true
	case int:
		return true
	case string:
		if _, err := strconv.ParseInt(v, 10, 64); err == nil {
			return true
		}
	case float64:
		return v == float64(int64(v))
	default:
		log.Warnf("Key %s is not an int64, int, string, or float64: %T", key, v)
	}

	return false
}

func GetIntFromMap(key string, obj map[string]interface{}) int64 {
	if obj == nil {
		return 0
	}

	value, exists := obj[key]
	if !exists {
		return 0
	}

	switch v := value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case string:
		if intValue, err := strconv.ParseInt(v, 10, 64); err == nil {
			return intValue
		}
		return 0
	default:
		return 0
	}
}

func ValidateAndGetBoolFromMap(key string, obj map[string]interface{}) (bool, error) {
	if obj == nil {
		return false, fmt.Errorf("object is nil")
	}

	if !ValidateBoolFromMap(key, obj) {
		return false, fmt.Errorf("key %s is missing or not a bool", key)
	}

	return GetBoolFromMap(key, obj), nil
}

func ValidateBoolFromMap(key string, obj map[string]interface{}) bool {
	if obj == nil {
		return false
	}

	value, exists := obj[key]
	if !exists {
		return false
	}

	_, ok := value.(bool)
	if !ok {
		return false
	}

	return true
}

func GetBoolFromMap(key string, obj map[string]interface{}) bool {
	if obj == nil {
		return false
	}

	value, exists := obj[key]
	if !exists {
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case string:
		if boolValue, err := strconv.ParseBool(v); err == nil {
			return boolValue
		}
	case int64:
		return v != 0
	case int:
		return v != 0
	default:
		log.Warnf("Key %s is not a bool, string, int64, or int: %T", key, v)
	}
	return false
}

func ValidateAndGetFloatFromMap(key string, obj map[string]interface{}) (float64, error) {
	if obj == nil {
		return 0, fmt.Errorf("object is nil")
	}

	if !ValidateFloatFromMap(key, obj) {
		return 0, fmt.Errorf("key %s is missing or not a float64", key)
	}

	return obj[key].(float64), nil
}

func ValidateFloatFromMap(key string, obj map[string]interface{}) bool {
	if obj == nil {
		return false
	}

	value, exists := obj[key]
	if !exists {
		return false
	}

	_, ok := value.(float64)
	if !ok {
		return false
	}

	return true
}

func GetFloatFromMap(key string, obj map[string]interface{}) float64 {
	if obj == nil {
		return 0
	}

	value, exists := obj[key]
	if !exists {
		return 0
	}

	switch v := value.(type) {
	case float64:
		return v
	case int64:
		return float64(v)
	case int:
		return float64(v)
	case string:
		if floatValue, err := strconv.ParseFloat(v, 64); err == nil {
			return floatValue
		}
	default:
		log.Warnf("Key %s is not a float64, int64, int, or string: %T", key, v)
	}

	return 0
}

func MapContainsKey(key string, obj map[string]interface{}) bool {
	if obj == nil {
		return false
	}

	_, exists := obj[key]
	if !exists {
		lowerKey := strings.ToLower(key)
		for k := range obj {
			if strings.ToLower(k) == lowerKey {
				return true
			}
		}
		return false
	}

	return true
}
