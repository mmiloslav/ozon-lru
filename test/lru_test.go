package test

import (
	"context"
	"testing"
	"time"

	"lru"
)

func TestPutAndGet(t *testing.T) {
	cache := lru.NewLRUCache(2)
	ctx := context.TODO()

	putKey := "key1"
	putValue := "value1"

	err := cache.Put(ctx, putKey, putValue, 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to put data with error [%s]", err.Error())
	}

	value, expiresAt, err := cache.Get(ctx, putKey)
	if err != nil {
		t.Fatalf("failed to get data with error [%s]", err.Error())
	}

	if value != putValue {
		t.Fatalf("expected [%s], got [%s]", putValue, value)
	}

	if time.Until(expiresAt) < 59*time.Minute {
		t.Fatalf("expected TTL to be around 1 hour, got [%v]", expiresAt)
	}
}

func TestExpiration(t *testing.T) {
	cache := lru.NewLRUCache(2)
	ctx := context.TODO()

	putKey := "key1"
	putValue := "value1"

	err := cache.Put(ctx, putKey, putValue, 1*time.Second)
	if err != nil {
		t.Fatalf("failed to put data with error [%s]", err.Error())
	}

	time.Sleep(2 * time.Second)

	_, _, err = cache.Get(ctx, putKey)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if err.Error() != lru.ErrKeyExpired {
		t.Fatalf("expected [%s], got [%s]", lru.ErrKeyExpired, err.Error())
	}
}

func TestEvict(t *testing.T) {
	cache := lru.NewLRUCache(2)
	ctx := context.TODO()

	putKey := "key1"
	putValue := "value1"

	err := cache.Put(ctx, putKey, putValue, 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to put data with error [%s]", err.Error())
	}

	value, err := cache.Evict(ctx, putKey)
	if err != nil {
		t.Fatalf("failed to evict data with error [%s]", err.Error())
	}

	if value != putValue {
		t.Fatalf("expected [%s], got [%s]", putValue, value)
	}

	_, _, err = cache.Get(ctx, putKey)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if err.Error() != lru.ErrKeyNotFound {
		t.Fatalf("expected [%s], got [%s]", lru.ErrKeyNotFound, err.Error())
	}
}

func TestEvictAll(t *testing.T) {
	cache := lru.NewLRUCache(2)
	ctx := context.TODO()

	putKey1 := "key1"
	putValue1 := "value1"

	putKey2 := "key2"
	putValue2 := "value2"

	err := cache.Put(ctx, putKey1, putValue1, 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to put data1 with error [%s]", err.Error())
	}

	err = cache.Put(ctx, putKey2, putValue2, 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to put data2 with error [%s]", err.Error())
	}

	err = cache.EvictAll(ctx)
	if err != nil {
		t.Fatalf("failed to evict all with error [%s]", err.Error())
	}

	_, _, err = cache.Get(ctx, putKey1)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if err.Error() != lru.ErrKeyNotFound {
		t.Fatalf("expected [%s], got [%s]", lru.ErrKeyNotFound, err.Error())
	}

	_, _, err = cache.Get(ctx, putKey2)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if err.Error() != lru.ErrKeyNotFound {
		t.Fatalf("expected [%s], got [%s]", lru.ErrKeyNotFound, err.Error())
	}
}

func TestGetAll(t *testing.T) {
	cache := lru.NewLRUCache(2)
	ctx := context.TODO()

	putKey1 := "key1"
	putValue1 := "value1"

	putKey2 := "key2"
	putValue2 := "value2"

	err := cache.Put(ctx, putKey1, putValue1, 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to put data1 with error [%s]", err.Error())
	}

	err = cache.Put(ctx, putKey2, putValue2, 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to put data2 with error [%s]", err.Error())
	}

	keys, values, err := cache.GetAll(ctx)
	if err != nil {
		t.Fatalf("failed to get all with error [%s]", err.Error())
	}

	if len(keys) != 2 || len(values) != 2 {
		t.Fatalf("expected 2 keys and 2 values, got [%d] keys and [%d] values", len(keys), len(values))
	}

	expectedKeys := []string{putKey2, putKey1}
	expectedValues := []interface{}{putValue2, putValue1}

	for i, key := range keys {
		if key != expectedKeys[i] {
			t.Fatalf("expected key [%s], got [%s]", expectedKeys[i], key)
		}
	}

	for i, value := range values {
		if value != expectedValues[i] {
			t.Fatalf("expected value [%s], got [%s]", expectedValues[i], value)
		}
	}
}

func TestLRUCapacity(t *testing.T) {
	cache := lru.NewLRUCache(2)
	ctx := context.TODO()

	putKey1 := "key1"
	putValue1 := "value1"

	putKey2 := "key2"
	putValue2 := "value2"

	putKey3 := "key3"
	putValue3 := "value3"

	err := cache.Put(ctx, putKey1, putValue1, 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to put data1 with error [%s]", err.Error())
	}

	err = cache.Put(ctx, putKey2, putValue2, 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to put data2 with error [%s]", err.Error())
	}

	err = cache.Put(ctx, putKey3, putValue3, 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to put data3 with error [%s]", err.Error())
	}

	_, _, err = cache.Get(ctx, putKey1)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if err.Error() != lru.ErrKeyNotFound {
		t.Fatalf("expected %v, got %v", lru.ErrKeyNotFound, err)
	}

	_, _, err = cache.Get(ctx, putKey2)
	if err != nil {
		t.Fatalf("failed to get data with error [%s]", err.Error())
	}
	_, _, err = cache.Get(ctx, putKey3)
	if err != nil {
		t.Fatalf("failed to get data with error [%s]", err.Error())
	}
}
