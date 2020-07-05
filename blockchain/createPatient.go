package blockchain

// CreatePatient : create Patient
func (setup *FabricSetup) CreatePatient(firstName string, lastName string,
	contactNo string, gender string, birthday string, address string) (string, error) {

	// Prepare arguments
	funcName := "createPatient"
	args := []string{firstName, lastName, contactNo, gender, birthday, address}

	return setup.Invoke(funcName, args)
}
