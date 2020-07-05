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

// EHR : Electronic Health Record
type EHR struct {
	ID                string        `json:"id"`
	Firstname         string        `json:"firstname"`
	Lastname          string        `json:"lastname"`
	ContactNo	  string        `json:"contactNum"`
	Gender		  string	`json:"gender"`
	Birthday          time.Time     `json:"birthday"`
	Address		  string	`json:"address"`
	FileUploads       []f_uploads   `json:"fuploads"`
	EhrUploads        []ehr_uploads `json:"ehruploads"`
}

// Appointment public for access outside the CC
type f_uploads struct {
	IPFS_fHash        string        `json:"ifhash"`
	uploadedDate      time.Time     `json:"udate"`
	fileInfo 	  string        `json:"finfo"`
}

type ehr_uploads struct {
	ehrTID		  string	`json:"ehrtid"`
	DiagnosisTime     time.Time	`json:"dtime"`
	prescriptionInfo  string	`json:"pinfo"`

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
	ehr.ContactNo = "contactNum"
	ehr.Gender = "gender"
	ehr.Birthday = time.Now()
	ehr.Address = "address"
	ehr.FileUploads = nil
	ehr.EhrUploads = nil

	behr, err := json.Marshal(ehr)
	if err != nil {
		return shim.Error("error marshalling EHR to Json")
	}

	// Put in the ledger the key/value hello/world
	err = stub.PutState("hello", behr)
	if err != nil {
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
	case "getPatientDetails":
		return getPatientDetails(stub, args)
	case "viewPatientDetails":
		return viewPatientDetails(stub, args)
	case "updateFileUploads":
		return updateFileUploads(stub, args)
	case "createEHR":
		return createEHR(stub, args)
	case "getEHR":
		return getEHR(stub, args)
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

	ehrID := args[1] //stub.GetTxID
	value := args[2]
	// Check if the ledger key is "hello" and process if it is the case. Otherwise it returns an error.
	if ehrID == "hello" && len(args) == 3 {

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
//	getPatientDetails - Getting Patient Details
// ==========================================================================================
func getPatientDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("running the function getPatientDetais()")

	if len(args) != 6 {
		return shim.Error("Wrong input")
	}

	var ehr EHR
	ehrID := stub.GetTxID()
	ehr.ID = ehrID
	ehr.Firstname = args[0]
	ehr.Lastname = args[1]
	ehr.ContactNo = args[2]
	ehr.Gender = args[3]
	ehr.Birthday, err = time.Parse(layout, args[4])
	ehr.Address = args[5]
	ehr.FileUploads = nil
	ehr.EhrUploads = nil

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
// ViewPatientDetails : query to get patient details by its key
// ==========================================================================================
func viewPatientDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
		fmt.Println("Patient details do not exist")
		return shim.Error("Patient details do not exist")
	}

	return shim.Success(valAsbytes)
}

// ==========================================================================================
// updateFileUploads : 
// ==========================================================================================
func updateFileUploads(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var ehrID string
	var err error
	var ehr *EHR

	if len(args) != 3 { 
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
	err = ehr.addfileuploads(args[1], args[2])
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
// createEHR: 
// ==========================================================================================
func createEHR(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var ehrID string
	var err error
	var ehr *EHR

	if len(args) != 3 {
		return shim.Error("wrong inout")
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
	err = ehr.addehruploads(args[1], args[2])
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
// addFileUploads :
// ==========================================================================================
func (ehr *EHR) addfileuploads( IPFS_fHash string, fileInfo string) error {

	_now := time.Now()
	yyyy := _now.Year()
	MM := _now.Month()
	dd := _now.Day()
	hh := _now.Hour()
	mm := _now.Minute()
	_now = time.Date(yyyy, MM, dd, hh, mm, 0, 0, time.UTC)
	ehr.FileUploads = append(ehr.FileUploads, f_uploads{IPFS_fHash, _now, fileInfo})
	return nil
}

// ==========================================================================================
// addFileUploads :
// ==========================================================================================

func (ehr *EHR) addehruploads(ehrTID string, prescriptionInfo string) error {

	_now := time.Now()
	yyyy := _now.Year()
	MM := _now.Month()
	dd := _now.Day()
	hh := _now.Hour()
	mm := _now.Minute()
	_now = time.Date(yyyy, MM, dd, hh, mm, 0, 0, time.UTC)
	ehr.EhrUploads = append(ehr.EhrUploads, ehr_uploads{ehrTID, _now, prescriptionInfo})
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
