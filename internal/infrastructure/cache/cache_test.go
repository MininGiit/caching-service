package cache

import (
	"context"
	"testing"
	"time"
)

func compare(resKeys []string, resValues []interface{}, keys []string, values []interface{}) bool{
	resMap := make(map[string] interface{}, len(resKeys))
	for i, resKey := range resKeys {
		resMap[resKey] = resValues[i]
	}
	for i, key := range keys {
		resVal, ok := resMap[key]
		if !ok {
			return false
		}
		if resVal != values[i] {
			return false
		} 
	}
	return true
}

func TestPut(t *testing.T) {
	maxSize := 10
	len := 4
	keys := []string {"qweer", "qwe", "123", "sdf43"}
	values := []interface{} {1, "34", 5.5, 3}
	var ttl time.Duration = time.Second * 10
	cache := New(maxSize, ttl)
	ctx := context.Background()

	expectedErrors := []error {nil, nil, nil, nil, ErrDataNotValid}

	for i := 0; i < len; i++ {
		err := cache.Put(ctx, keys[i], values[i], ttl)
		if err != expectedErrors[i] {
			t.Errorf("exected error: %v, recived error: %v", expectedErrors[i], err)
		}
	}
	if cache.size != len {
		t.Errorf("the sixe do not match: exected %v, recived %v", len, cache.size)
	} 
	for i := 0; i < len; i++ {
		recived := cache.data[keys[i]].value 
		exected := values[i]
		if recived != exected {
			t.Errorf("the values do not match: exected %v, recived %v", exected, recived)
		} 
	}
}

func TestUpdate(t *testing.T) {
	key := "qwe"
	value := "123"
	newValue := "124"
	var ttl time.Duration = time.Second * 10
	maxSize := 10
	cache := New(maxSize, ttl)
	cache.Put(context.Background(), key, value, ttl)
	cache.Put(context.Background(), key, newValue, ttl)

	resValue, _, _ := cache.Get(context.Background(), key)
	if resValue != newValue {
			t.Errorf("exected %v, recived %v", newValue, resValue)
	}
}

func TestGet(t *testing.T) {
	maxSize := 10
	cashSize := 4
	keys := []string {"qweer", "qwe", "123", "sdf43"}
	values := []interface{} {1, "34", 5.5, 3}
	var ttl time.Duration = time.Second * 10
	cache := New(maxSize, ttl)
	ctx := context.Background()
	for i := 0; i < cashSize; i++ {
		cache.Put(ctx, keys[i], values[i], ttl)
	}

	inputKeys := []string {"qweer", "qwe", "bv", "123", "sdf43", "sd"}
	expectedValues := []interface{} {1, "34", nil, 5.5, 3, nil}
	expectedErrors := []error {nil, nil, ErrKeyNotFound, nil, nil, ErrKeyNotFound}

	for i := 0; i < len(inputKeys); i++ {
		recived, _, err := cache.Get(ctx, inputKeys[i])
		if err != expectedErrors[i] {
			t.Errorf("exected error: %v, recived error: %v", expectedErrors[i], err)
		}
		exected := expectedValues[i]
		if recived != exected {
			t.Errorf("exected %v, recived %v", exected, recived)
		} 
	}
}

func TestGetAll(t *testing.T) {
	maxSize := 10
	cashSize := 4
	keys := []string {"qweer", "qwe", "123", "sdf43"}
	values := []interface{} {1, "34", 5.5, 3}
	var ttl time.Duration = time.Second * 10
	cache := New(maxSize, ttl)
	ctx := context.Background()

	for i := 0; i < cashSize; i++ {
		cache.Put(ctx, keys[i], values[i], ttl)
	}

	resKeys, resValues, resErr := cache.GetAll(ctx)
	if resErr != nil {
		t.Errorf("exected error = nil, recived error: %v", resErr)
	}
	if len(keys) != cashSize || len(values) != cashSize {
		t.Errorf("the length didin't match")
	}
	if !compare(resKeys, resValues, keys, values) {
		t.Errorf("slices do not match %v, %v", resKeys, resValues)
	}
}

func TestEvictAllForEmpty(t *testing.T) {
	maxSize := 10
	cache := New(maxSize, 0)
	ctx := context.Background()

	_, _, err := cache.GetAll(ctx)
	if err != ErrCacheEmpty {
		t.Errorf("exected %v, recived %v", ErrCacheEmpty, err)
	}
}


func TestEvict(t *testing.T) {
	maxSize := 10
	cashSize := 4
	keys := []string {"qweer", "qwe", "123", "sdf43"}
	values := []interface{} {1, "34", 5.5, 3}
	var ttl time.Duration = time.Second * 10
	cache := New(maxSize, ttl)
	ctx := context.Background()
	for i := 0; i < cashSize; i++ {
		cache.Put(ctx, keys[i], values[i], ttl)
	}
	
	inputKeys := []string {"qweer", "qwe", "123", "sdf43"}
	expectedValues := []interface{} {nil, "34", nil, 3}
	expectedErrors := []error {ErrKeyNotFound, nil, ErrKeyNotFound, nil}

	_, err := cache.Evict(ctx, "qweer")
	if err != nil {
		t.Errorf("exected error = nil, recived error: %v", err)
	}
	_, err = cache.Evict(ctx, "123")
	if err != nil {
		t.Errorf("exected error = nil, recived error: %v", err)
	}

	for i := 0; i < len(inputKeys); i++ {
		recived, _, err := cache.Get(ctx, inputKeys[i])
		if err != expectedErrors[i] {
			t.Errorf("exected error: %v, recived error: %v", expectedErrors[i], err)
		}
		exected := expectedValues[i]
		if recived != exected {
			t.Errorf("exected %v, recived %v", exected, recived)
		} 
	}
}

func TestEvictAll(t *testing.T) {
	maxSize := 10
	cashSize := 4
	keys := []string {"qweer", "qwe", "123", "sdf43"}
	values := []interface{} {1, "34", 5.5, 3}
	var ttl time.Duration = time.Second * 10
	cache := New(maxSize, ttl)
	ctx := context.Background()
	for i := 0; i < cashSize; i++ {
		cache.Put(ctx, keys[i], values[i], ttl)
	}
	
	inputKeys := []string {"qweer", "qwe", "123", "sdf43"}
	expectedValues := []interface{} {nil, nil, nil, nil}
	expectedErrors := []error {ErrKeyNotFound, ErrKeyNotFound, ErrKeyNotFound, ErrKeyNotFound}

	err := cache.EvictAll(ctx)
	if err != nil {
		t.Errorf("exected error = nil, recived error: %v", err)
	}
	for i := 0; i < len(inputKeys); i++ {
		recived, _, err := cache.Get(ctx, inputKeys[i])
		if err != expectedErrors[i] {
			t.Errorf("exected error: %v, recived error: %v", expectedErrors[i], err)
		}
		exected := expectedValues[i]
		if recived != exected {
			t.Errorf("exected %v, recived %v", exected, recived)
		} 
	}
}

func TestLRU(t *testing.T) {
	maxSize := 4
	cashSize := 4
	keys := []string {"qweer", "qwe", "123", "sdf43"}
	values := []interface{} {1, "34", 5.5, 3}
	var ttl time.Duration = time.Second * 10
	cache := New(maxSize, ttl)
	ctx := context.Background()
	for i := 0; i < cashSize; i++ {
		cache.Put(ctx, keys[i], values[i], ttl)
	}

	newItemKey :=  "newKey"
	newItemValue := ":)0)))"
	inputKeys := []string {"qweer", "qwe", "123", "sdf43", newItemKey}
	expectedValues := []interface{} {nil, "34", 5.5, 3, newItemValue}
	expectedErrors := []error {ErrKeyNotFound, nil, nil, nil, nil}
	
	cache.Put(ctx, newItemKey, newItemValue, ttl)

	for i := 0; i < len(inputKeys); i++ {
		recived, _, err := cache.Get(ctx, inputKeys[i])
		if err != expectedErrors[i] {
			t.Errorf("exected error: %v, recived error: %v", expectedErrors[i], err)
		}
		exected := expectedValues[i]
		if recived != exected {
			t.Errorf("exected %v, recived %v", exected, recived)
		} 
	}
}
