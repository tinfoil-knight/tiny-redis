package store

import (
	"encoding/gob"
	"os"
)

type Store struct {
	underlying map[string]([]byte)
}

func New() *Store {
	kv := Store{
		underlying: make(map[string]([]byte)),
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

func (kv *Store) Set(key []byte, value []byte) {
	kv.underlying[string(key)] = value
}

func (kv *Store) Get(key []byte) (value []byte, ok bool) {
	v, ok := kv.underlying[string(key)]
	return v, ok
}

func (kv *Store) Del(key []byte) {
	delete(kv.underlying, string(key))
}

func (kv *Store) GetUnderlying() *map[string]([]byte) {
	return &kv.underlying
}
