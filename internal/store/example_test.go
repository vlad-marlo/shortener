package store_test

import (
	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/inmemory"
)

func Example() {
	var storage store.Store
	// init storage
	// for example inmemory storage
	storage = inmemory.New()
	defer func() {
		// always close storage after work with it
		if err := storage.Close(); err != nil {
			// ...
		}
	}()
	// ...
}
