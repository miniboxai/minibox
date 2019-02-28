package oauth_storage

var initializerFuncs = make([]StorageHandler, 0)

func Initializer(handle StorageHandler) {
	initializerFuncs = append(initializerFuncs, handle)
}

func EachInitializer(store *Storage) {
	for _, initfn := range initializerFuncs {
		initfn(store)
	}
}
