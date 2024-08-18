package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

func WriteHandler(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Bad Request: unable to decode request body", http.StatusBadRequest)
		return
	}

	data, ok := requestData["data"].(map[string]interface{})
	if !ok {
		http.Error(w, "Bad Request: 'data' field is missing or not a map", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var hasError bool

	for key, value := range data {
		wg.Add(1)
		go func(key string, value interface{}) {
			defer wg.Done()
			if err := tarantoolSet(key, value); err != nil {
				mu.Lock()
				hasError = true
				mu.Unlock()
			}
		}(key, value)
	}
	wg.Wait()

	if hasError {
		http.Error(w, "Internal Server Error: failed to write some keys", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func ReadHandler(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Bad Request: unable to decode request body", http.StatusBadRequest)
		return
	}

	keys, ok := requestData["keys"].([]interface{})
	if !ok {
		http.Error(w, "Bad Request: 'keys' field is missing or not a list", http.StatusBadRequest)
		return
	}

	data := make(map[string]interface{})
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, key := range keys {
		keyStr, ok := key.(string)
		if !ok {
			http.Error(w, "Bad Request: one of the keys is not a string", http.StatusBadRequest)
			return
		}

		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			value, err := tarantoolGet(key)
			if err == nil && value != nil {
				mu.Lock()
				data[key] = value
				mu.Unlock()
			}
		}(keyStr)
	}
	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": data})
}
