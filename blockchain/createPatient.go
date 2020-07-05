package blockchain

// CreatePatient : create Patient
func (setup *FabricSetup) CreatePatient(name string, contactNo string) (string, error) {

	// Prepare arguments
	funcName := "createPatient"
	args := []string{name, contactNo}

	return setup.Invoke(funcName, args)
}
