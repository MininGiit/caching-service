package http

import "time"

type ResponseItem struct {
	Key	string			    `json:"key"`
	Value interface{}	    `json:"value"`
	ExpiresAt time.Duration `json:"expires_at"`
}

type ResponseItems struct {
	Keys	[]string	 `json:"keys"`
	Values []interface{} `json:"values"`
}

type RequestItem struct {
	Key			string			`json:"key"`
	Value 		interface{}		`json:"value"`
	TtlSeconds	int				`json:"ttl_seconds"`
}

func NewResponseItem(key string, value interface{}, expiresTime time.Time) *ResponseItem {
	expiresAt := expiresTime.Sub(time.Now())
	return &ResponseItem{
		Key: 	key,
		Value:	value,
		ExpiresAt: expiresAt,
	}
}

func NewResponseItems(keys []string, values []interface{}) *ResponseItems {
	return &ResponseItems{
		Keys: 	keys,
		Values:	values,
	}
}