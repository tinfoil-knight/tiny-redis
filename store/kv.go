package store

import (
	"encoding/gob"
	"os"
)

type Store struct {
	underlying map[string]string
}

func New() *Store {
	kv := Store{
		underlying: make(map[string]string),
	}
	defaultPath := "dump.trdb"
	f, err := os.Open(defaultPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &kv
		}
		panic(err)
	}
	defer f.Close()
	if err = gob.NewDecoder(f).Decode(&kv.underlying); err != nil {
		panic(err)
	}
	// fmt.Printf("%#v\n", kv)
	return &kv
}

func (kv *Store) Set(key, value string) {
	kv.underlying[key] = value
}

func (kv *Store) Get(key string) (value string, ok bool) {
	v, ok := kv.underlying[key]
	return v, ok
}

func (kv *Store) Del(key string) {
	delete(kv.underlying, key)
}

func (kv *Store) GetUnderlying() *map[string]string {
	return &kv.underlying
}
