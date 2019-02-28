//+build !release

package oauth_storage

import (
	"encoding/json"
	"expvar"

	lru "github.com/hashicorp/golang-lru"
)

type expvarStorage struct {
	*Storage
}

func init() {
	Initializer(func(store *Storage) {
		// ticker := time.NewTicker(2 * time.Second)
		// go func() {
		// 	for t := range ticker.C {
		// 		fmt.Println("Tick at", t)
		// 		expvar.Publish("store", &expvarStorage{store})
		// 	}
		// }()
		//
		expvar.Publish("store", &expvarStorage{store})
	})
}

var NoHandler = func(m interface{}) {}

func (store *expvarStorage) String() string {
	var s = make(map[string]interface{})
	clients := store.Storage.clients.cache
	tokens := store.Storage.access.cache
	authorizes := store.Storage.authorize.cache
	s["clients"] = getItems(clients, NoHandler)
	s["tokens"] = getItems(tokens, NoHandler)
	s["authorizes"] = getItems(authorizes, NoHandler)

	ret, _ := json.Marshal(s)
	// s["clients"] = string(ret)
	return string(ret)
}

func getItems(cache *lru.Cache, handle func(interface{})) []interface{} {
	var a = make([]interface{}, 0)
	for _, key := range cache.Keys() {
		m, _ := cache.Get(key)
		handle(m)
		a = append(a, m)
	}
	return a
}
