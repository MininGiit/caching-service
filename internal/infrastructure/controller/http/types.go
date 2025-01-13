package http

import "time"

type responseItem struct {
	Key       string        `json:"key"`
	Value     interface{}   `json:"value"`
	ExpiresAt time.Duration `json:"expires_at"`
}

type responseItems struct {
	Keys   []string      `json:"keys"`
	Values []interface{} `json:"values"`
}

type requestItem struct {
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	TtlSeconds int         `json:"ttl_seconds"`
}

func newResponseItem(key string, value interface{}, expiresTime time.Time) *responseItem {
	expiresAt := expiresTime.Sub(time.Now())
	return &responseItem{
		Key:       key,
		Value:     value,
		ExpiresAt: expiresAt,
	}
}

func newResponseItems(keys []string, values []interface{}) *responseItems {
	return &responseItems{
		Keys:   keys,
		Values: values,
	}
}
