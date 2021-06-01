package store

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sync"
)

// should be part of the Store struct
var defaultPath = "dump.trdb"

type Store struct {
	underlying map[string]([]byte)
	mu         sync.RWMutex
}

func New() *Store {
	kv := Store{
		underlying: make(map[string]([]byte)),
	}
	ok := kv.Load(defaultPath)
	if ok {
		fmt.Printf("DB loaded from disk: %s\n", defaultPath)
	}
	// fmt.Printf("%#v\n", kv)
	return &kv
}

func (kv *Store) Load(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	}
	defer f.Close()
	if err = gob.NewDecoder(f).Decode(&kv.underlying); err != nil {
		panic(err)
	}
	return true
}

func (kv *Store) Save() {
	f, err := os.OpenFile(defaultPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b := new(bytes.Buffer)
	if err = gob.NewEncoder(b).Encode(kv.underlying); err != nil {
		panic(err)
	}
	if _, err = io.Copy(f, b); err != nil {
		panic(err)
	}
}

func (kv *Store) Set(key []byte, value []byte) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.underlying[string(key)] = value
}

func (kv *Store) Get(key []byte) (value []byte, ok bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	v, ok := kv.underlying[string(key)]
	return v, ok
}

func (kv *Store) Del(key []byte) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	delete(kv.underlying, string(key))
}
