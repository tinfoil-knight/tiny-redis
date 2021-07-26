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
	underlying sync.Map
}

func New() *Store {
	kv := Store{
		underlying: *new(sync.Map),
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
	var tmp map[string][]byte
	if err = gob.NewDecoder(f).Decode(&tmp); err != nil {
		panic(err)
	}
	kv.underlying = *toSyncMap(&tmp)
	return true
}

func (kv *Store) Save() {
	f, err := os.OpenFile(defaultPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b := new(bytes.Buffer)
	tmp := fromSyncMap(&kv.underlying)
	if err = gob.NewEncoder(b).Encode(tmp); err != nil {
		panic(err)
	}
	if _, err = io.Copy(f, b); err != nil {
		panic(err)
	}
}

func fromSyncMap(sm *sync.Map) *map[string][]byte {
	tmp := make(map[string][]byte)
	sm.Range(func(k, v interface{}) bool {
		tmp[k.(string)] = v.([]byte)
		return true
	})
	return &tmp
}

func toSyncMap(m *map[string][]byte) *sync.Map {
	sm := sync.Map{}
	for k, v := range *m {
		sm.Store(k, v)
	}
	return &sm
}

func (kv *Store) Set(key []byte, value []byte) {
	kv.underlying.Store(string(key), value)
}

func (kv *Store) Get(key []byte) (value []byte, ok bool) {
	if v, ok := kv.underlying.Load(string(key)); ok {
		return v.([]byte), true
	}
	return nil, false
}

func (kv *Store) Del(key []byte) {
	kv.underlying.Delete(string(key))
}
