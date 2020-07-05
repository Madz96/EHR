package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	layout = "2006-01-02"
)

// Patient : Patient Details
type Patient struct {
	ID string `json:"id"`
	Name string `json:"name"`
	ContactNo string `json:"contactNo"`
}

// EHR : Electronic Health Record
type EHR struct {
	ID                string        `json:"id"`
	Firstname         string        `json:"firstname"`
	Lastname          string        `json:"lastname"`
	SocialSecurityNum string        `json:"socialSecurityNum"`
	Birthday          time.Time     `json:"birthday"`
	Appointments      []Appointment `json:"visits"`
}

// Appointment public for access outside the CC
type Appointment struct {
	DrID    string    `json:"drId"`
	Date    time.Time `json:"date"`
	Comment string    `json:"comment"`
}

// HeroesServiceChaincode implementation of Chaincode
type HeroesServiceChaincode struct {
}

// Init of the chaincode
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (t *HeroesServiceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### HeroesServiceChaincode Init ###########")

	// Get the function and arguments from the request
	function, _ := stub.GetFunctionAndParameters()

	// Check if the request is the init function
	if function != "init" {
		return shim.Error("Unknown function call")
	}

	var ehr EHR
	ehr.ID = "ID"
	ehr.Firstname = "firstname"
	ehr.Lastname = "lastname"
	ehr.SocialSecurityNum = "socialsecuritynumber"
	ehr.Birthday = time.Now()
	ehr.Appointments = nil

	behr, err := json.Marshal(ehr)
	if err != nil {
		return shim.Error("error marshalling EHR to Json")
	}

	//patient stuff
	var pat Patient
	pat.ID = "ID"
	pat.Name = "NAME"
	pat.ContactNo = "1234567"

	bpat, err := json.Marshal(pat)
	if err != nil {
		return shim.Error("error marshalling pat to json in Init")
	}

	// Put in the ledger the key/value hello/world
	// err = stub.PutState("hello", behr)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	// Put in the ledger the key/value hello/world
	err = stub.PutState("hello", bpat)
	if err != nil {
		fmt.Println("something went wrong at PutState()")
		return shim.Error(err.Error())
	}

	// Return a successful message
	return shim.Success(nil)
}




// Invoke of the chaincode
// All future requests named invoke will arrive here.
func (t *HeroesServiceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### HeroesServiceChaincode Invoke ###########")

	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()

	// Handle different functions
	switch function {
	case "createPatient":
		return createPatient(stub, args)
	case "createEHR":
		return createEHR(stub, args)
	case "getEHR":
		return getEHR(stub, args)
	case "updateEHR":
		return updateEHR(stub, args)
	}

	// Check whether it is an invoke request
	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// In order to manage multiple type of request, we will check the first argument.
	// Here we have one possible argument: query (every query request will read in the ledger without modification)
	if args[0] == "query" {
		return t.query(stub, args)
	}

	// The update argument will manage all update in the ledger
	if args[0] == "invoke" {
		return t.invoke(stub, args)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown action, check the first argument")
}

// query
// Every readonly functions in the ledger will be here
func (t *HeroesServiceChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### HeroesServiceChaincode query ###########")

	// Check whether the number of arguments is sufficient
	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Like the Invoke function, we manage multiple type of query requests with the second argument.
	// We also have only one possible argument: hello
	if args[1] == "hello" {

		// Get the state of the value matching the key hello in the ledger
		state, err := stub.GetState("hello")
		if err != nil {
			return shim.Error("Failed to get state of hello")
		}

		// Return this value in response
		return shim.Success(state)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown query action, check the second argument.")
}

// invoke
// Every functions that read and write in the ledger will be here
func (t *HeroesServiceChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### HeroesServiceChaincode invoke ###########")

	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// ehrID := args[1] //stub.GetTxID

	patID := args[1]
	value := args[2]
	// Check if the ledger key is "hello" and process if it is the case. Otherwise it returns an error.
	if patID == "hello" && len(args) == 3 {

		// Add random suffix to the value
		value = value + strconv.Itoa(time.Now().Nanosecond())
		// Write the new value in the ledger
		err := stub.PutState(ehrID, []byte(value))
		if err != nil {
			return shim.Error("Failed to update state of hello")
		}

		// Notify listeners that an event "eventInvoke" have been executed (check line 19 in the file invoke.go)
		err = stub.SetEvent("eventInvoke", []byte{})
		if err != nil {
			return shim.Error(err.Error())
		}

		// Return this value in response
		return shim.Success(nil)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown invoke action, check the second argument.")
}

// ==========================================================================================
//	createPatient- create a patient
// ==========================================================================================
func createPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("i am now in chaincode main -> creatPatient()")

	if len(args) != 2 {
		return shim.Error("Incorrect input length", len(args))
	}

	var pat Patient
	patID := stub.GetTxID()
	pat.Name = args[0]
	pat.ContactNo = args[1]

	jsonPat, err := json.Marshal(pat)
	if err != nil {|
		return shim.Error("Error marshalling pat to JSON in createPatient()")
	}

	err = stub.PutState(patID, jsonPat)
	if err != nil {
		return shim.Error("createPatient() : Error writing to state")
	}

	// Notify listeners that an event "eventInvoke" has been executed
	err = stub.SetEvent("eventInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(ehrID))
}






// ==========================================================================================
//	createEHR- create a donor-recipient pair of health records
// ==========================================================================================
func createEHR(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("running the function createPair()")

	if len(args) != 4 {
		return shim.Error("Wrong input")
	}

	var ehr EHR
	ehrID := stub.GetTxID()
	ehr.ID = ehrID
	ehr.Firstname = args[0]
	ehr.Lastname = args[1]
	ehr.SocialSecurityNum = args[2]
	ehr.Birthday, err = time.Parse(layout, args[3])
	ehr.Appointments = nil

	if err != nil {
		return shim.Error("Error parsing birthday")
	}

	jsonEHR, err := json.Marshal(ehr)
	if err != nil {
		fmt.Println(err.Error())
		return shim.Error("Error marshalling to JSON")
	}

	err = stub.PutState(ehrID, jsonEHR)
	if err != nil {
		return shim.Error("createEHR() : Error writing to state")
	}

	// Notify listeners that an event "eventInvoke" has been executed
	err = stub.SetEvent("eventInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(ehrID))
}

// ==========================================================================================
// getEHR : query to get a EHR by its key
// ==========================================================================================
func getEHR(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var ehrID string
	var err error

	if len(args) != 1 {
		return shim.Error("Wrong input")
	}
	ehrID = args[0]
	valAsbytes, err := stub.GetState(ehrID)
	if err != nil {
		fmt.Println(err.Error())
		return shim.Error(err.Error())
	} else if valAsbytes == nil {
		fmt.Println("EHR does not exist")
		return shim.Error("EHR does not exist")
	}

	return shim.Success(valAsbytes)
}

// ==========================================================================================
// updateEHR : get a EHR by its key and add a Appointment (drId + date + comment)
// ==========================================================================================
func updateEHR(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var ehrID string
	var err error
	var ehr *EHR

	if len(args) != 3 { // ehrID, DrID, comment
		return shim.Error("Wrong input")
	}
	ehrID = args[0]
	ehr, err = getEHRbyID(stub, ehrID)
	if err != nil {
		fmt.Println(err.Error())
		return shim.Error(err.Error())
	}
	if ehr == nil {
		fmt.Println("Error reading state : EHR is nil")
		return shim.Error("nil ehr")
	}
	err = ehr.addAppointment(args[1], args[2])
	if err != nil {
		fmt.Println(err.Error())
		return shim.Error(err.Error())
	}

	jsonEHR, err := json.Marshal(ehr)
	if err != nil {
		fmt.Println(err.Error())
		return shim.Error("error marshalling json" + err.Error())
	}

	err = stub.PutState(ehrID, jsonEHR)
	if err != nil {
		return shim.Error("updateEHR() : Error put state")
	}

	// Notify listeners that an event "eventInvoke" has been executed
	err = stub.SetEvent("eventInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(jsonEHR)
}

// ==========================================================================================
// add Appointment (date + comment) to EHR
// ==========================================================================================
func (ehr *EHR) addAppointment(DrID string, comment string) error {

	_now := time.Now()
	yyyy := _now.Year()
	MM := _now.Month()
	dd := _now.Day()
	hh := _now.Hour()
	mm := _now.Minute()
	_now = time.Date(yyyy, MM, dd, hh, mm, 0, 0, time.UTC)
	ehr.Appointments = append(ehr.Appointments, Appointment{DrID, _now, comment})
	return nil
}

// ==========================================================================================
// getEHRbyID : get the EHR object by ID - Auxiliary function
// ==========================================================================================
func getEHRbyID(stub shim.ChaincodeStubInterface, ID string) (*EHR, error) {
	valAsbytes, err := stub.GetState(ID)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	} else if valAsbytes == nil {
		return nil, errors.New("EHR does not exist")
	}

	var ehr EHR
	err = json.Unmarshal(valAsbytes, &ehr)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("Error unmarshalling JSON")
	}

	return &ehr, nil
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(HeroesServiceChaincode))
	if err != nil {
		fmt.Printf("Error starting Heroes Service chaincode: %s", err)
	}
}
