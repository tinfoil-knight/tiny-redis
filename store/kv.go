package store

type Store struct {
	underlying map[string]string
}

func New() *Store {
	kv := Store{
		underlying: make(map[string]string),
	}
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
