package adapter

func Register(resource Resource) error {
	return RegistryInstance().Register(resource)
}

func MustRegister(resource Resource) {
	if err := Register(resource); err != nil {
		panic(err)
	}
}
