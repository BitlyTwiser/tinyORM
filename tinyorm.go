package tinyorm

func Connect(connections ...string) error {
	if len(connections) == 0 {
		// Default to development
	}

	return nil
}
