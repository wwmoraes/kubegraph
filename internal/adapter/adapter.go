package adapter

func Register(resourceTransformer ResourceTransformer) error {
	return RegistryInstance().Register(resourceTransformer)
}

func MustRegister(resourceTransformer ResourceTransformer) {
	if err := Register(resourceTransformer); err != nil {
		panic(err)
	}
}
