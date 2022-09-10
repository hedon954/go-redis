package dict

import "sync"

// SyncDict a concurrently secure dict
type SyncDict struct {
	m sync.Map
}

func MakeSyncDict() *SyncDict {
	return &SyncDict{}
}

func (dict *SyncDict) Get(key string) (val interface{}, exists bool) {
	return dict.m.Load(key)
}

func (dict *SyncDict) Len() int {
	length := 0
	dict.m.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	return length
}

func (dict *SyncDict) Put(key string, val interface{}) (result int) {
	_, ok := dict.m.Load(key)
	if !ok {
		dict.m.Store(key, val)
		return 1
	}
	dict.m.Store(key, val)
	return 0
}

func (dict *SyncDict) PutIfAbsent(key string, val interface{}) (result int) {
	_, ok := dict.m.Load(key)
	if ok {
		return 0
	}
	dict.m.Store(key, val)
	return 1
}

func (dict *SyncDict) PutIfExists(key string, val interface{}) (result int) {
	_, ok := dict.m.Load(key)
	if !ok {
		return 0
	}
	dict.m.Store(key, val)
	return 1
}

func (dict *SyncDict) Remove(key string) (result int) {
	_, ok := dict.m.Load(key)
	dict.m.Delete(key)
	if ok {
		return 1
	}
	return 0
}

func (dict *SyncDict) Foreach(consumer Consumer) {
	dict.m.Range(func(key, value interface{}) bool {
		return consumer(key.(string), value)
	})
}

func (dict *SyncDict) Keys() []string {
	keys := make([]string, dict.Len())
	i := 0
	dict.m.Range(func(key, value interface{}) bool {
		keys[i] = key.(string)
		i++
		return true
	})
	return keys
}

func (dict *SyncDict) RandomKeys(limit int) []string {
	if limit <= 0 {
		return nil
	}
	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		dict.m.Range(func(key, value interface{}) bool {
			result[i] = key.(string)
			return false
		})
	}
	return result
}

func (dict *SyncDict) RandomDistinctKeys(limit int) []string {
	if limit <= 0 {
		return nil
	}
	i := 0
	result := make([]string, limit)
	dict.m.Range(func(key, value interface{}) bool {
		result[i] = key.(string)
		i++
		if i == limit {
			return false
		}
		return true
	})
	return result
}

func (dict *SyncDict) clear() {
	*dict = *MakeSyncDict()
}
