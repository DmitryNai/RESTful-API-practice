package main

import (
	"log"

	"github.com/tarantool/go-tarantool"
)

var conn *tarantool.Connection

func init() {
	var err error
	conn, err = tarantool.Connect("tarantool:3301", tarantool.Opts{
		User: "guest",
	})
	if err != nil {
		log.Fatalf("Connection refused: %v", err)
	}

	err = ensureSpaceExists()
	if err != nil {
		log.Fatalf("Failed to ensure space exists: %v", err)
	}
}

func ensureSpaceExists() error {
	_, err := conn.Call("box.space.kv", []interface{}{})
	if err == nil {
		return nil
	}

	_, err = conn.Eval(`
		box.schema.space.create('kv', {
			if_not_exists = true,
			format = {
				{ name = 'key', type = 'string' },
				{ name = 'value', type = 'any' }
			}
		})
		box.space.kv:create_index('primary', {
			parts = {'key'},
			if_not_exists = true
		})
	`, []interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func tarantoolSet(key string, value interface{}) error {
	_, err := conn.Replace("kv", []interface{}{key, value})
	if err != nil {
		log.Printf("Error replacing key %s: %v", key, err)
	}
	return err
}

func tarantoolGet(key string) (interface{}, error) {
	resp, err := conn.Select("kv", "primary", 0, 1, tarantool.IterEq, []interface{}{key})
	if err != nil {
		log.Printf("Error selecting key %s: %v", key, err)
		return nil, err
	}
	if len(resp.Tuples()) == 0 {
		log.Printf("Key %s not found", key)
		return nil, nil
	}
	return resp.Tuples()[0][1], nil
}
