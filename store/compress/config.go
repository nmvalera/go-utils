package compress

import (
	store "github.com/kkrt-labs/go-utils/store"
	multistore "github.com/kkrt-labs/go-utils/store/multi"
)

type Config struct {
	ContentEncoding  store.ContentEncoding
	MultiStoreConfig multistore.Config
}
