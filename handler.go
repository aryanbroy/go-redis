package main

import (
	"fmt"
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetAll,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR! Invalid number of arguments for the SET command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR! Invalid number of arguments for the GET command"}
	}

	key := args[0].bulk

	SETsMu.Lock()
	value, ok := SETs[key]
	SETsMu.Unlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR! invalid number of arguments for th HSET command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	fmt.Println("hsets: ", HSETs)

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR! invalid number of arguments for the HGET command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.Lock()
	value, ok := HSETs[hash][key]
	HSETsMu.Unlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hgetAll(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR! invalid number of arguments for the HGETALL command"}
	}

	hash := args[0].bulk

	HSETsMu.Lock()
	key, ok := HSETs[hash]
	if !ok {
		return Value{typ: "array", array: []Value{}}
	}

	result := make([]Value, 0, len(key)*2)
	for k, val := range key {
		result = append(result,
			Value{typ: "bulk", bulk: k},
			Value{typ: "bulk", bulk: val},
		)
	}

	HSETsMu.Unlock()

	return Value{typ: "array", array: result}
}
