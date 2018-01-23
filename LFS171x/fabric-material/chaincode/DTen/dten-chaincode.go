// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario

 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
*/

package main

/* Imports
* 4 utility libraries for handling bytes, reading and writing JSON,
formatting, and string manipulation
* 2 specific Hyperledger Fabric specific libraries for Smart Contracts
*/
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

/* Define Tender structure, with 4 properties.
Structure tags are used by encoding/json library
*/
type Tender struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Size      string `json:"size"`
	Lots      string `json:"lots"`
}

/*
 * The Init method *
 called when the Smart Contract "dten-chaincode" is instantiated by the network
 * Best practice is to have any Ledger initialization in separate function
 -- see initLedger()
*/
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "dten-chaincode"
 The app also specifies the specific smart contract function to call with args
*/
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "queryTender" {
		return s.queryTender(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "recordTender" {
		return s.recordTender(APIstub, args)
	} else if function == "queryAllTenders" {
		return s.queryAllTenders(APIstub)
	} //else if function == "changeTenderHolder" {
	// 	return s.changeTenderHolder(APIstub, args)
	// }

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The queryTuna method *
Used to view the records of one particular tuna
It takes one argument -- the key for the tuna in question
*/
func (s *SmartContract) queryTender(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	tenderAsBytes, _ := APIstub.GetState(args[0])
	if tenderAsBytes == nil {
		return shim.Error("Could not locate tender")
	}
	return shim.Success(tenderAsBytes)
}

/*
 * The initLedger method *
Will add test data (10 tenders)to our network
*/
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	tender := []Tender{
		Tender{Type: "Service", Size: "10", Timestamp: "1504054225", Lots: "100"},
		Tender{Type: "Supply", Size: "10", Timestamp: "1504057825", Lots: "100"},
		Tender{Type: "Work", Size: "10", Timestamp: "1493517025", Lots: "100"},
		Tender{Type: "Building", Size: "10", Timestamp: "1496105425", Lots: "100"},
		Tender{Type: "Service", Size: "10", Timestamp: "1493512301", Lots: "100"},

		Tender{Type: "Supply", Size: "10", Timestamp: "1494117101", Lots: "100"},
		Tender{Type: "Work", Size: "10", Timestamp: "1496104301", Lots: "100"},
		Tender{Type: "Service", Size: "10", Timestamp: "1485066691", Lots: "100"},
		Tender{Type: "Building", Size: "10", Timestamp: "1485153091", Lots: "100"},
		Tender{Type: "Service", Size: "10", Timestamp: "1487745091", Lots: "100"},
	}

	i := 0
	for i < len(tender) {
		fmt.Println("i is ", i)
		tenderAsBytes, _ := json.Marshal(tender[i])
		APIstub.PutState(strconv.Itoa(i+1), tenderAsBytes)
		fmt.Println("Added", tender[i])
		i = i + 1
	}

	return shim.Success(nil)
}

/*
 * The recordTender method *
Fisherman like Sarah would use to record each of her tuna catches.
This method takes in five arguments (attributes to be saved in the ledger).
*/
func (s *SmartContract) recordTender(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var tender = Tender{Type: args[1], Size: args[2], Timestamp: args[3], Lots: args[4]}

	tenderAsBytes, _ := json.Marshal(tender)
	err := APIstub.PutState(args[0], tenderAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record Tender: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The queryAllTender method *
allows for assessing all the records added to the ledger(all tender catches)
This method does not take any arguments. Returns JSON string containing results.
*/
func (s *SmartContract) queryAllTenders(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllTenders:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The changeTunaHolder method *
The data in the world state can be updated with who has possession.
This function takes in 2 arguments, tuna id and new holder name.
*/
// func (s *SmartContract) changeTenderHolder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

// 	if len(args) != 2 {
// 		return shim.Error("Incorrect number of arguments. Expecting 2")
// 	}

// 	tenderAsBytes, _ := APIstub.GetState(args[0])
// 	if tenderAsBytes == nil {
// 		return shim.Error("Could not locate tuna")
// 	}
// 	tender := Tender{}

// 	json.Unmarshal(tenderAsBytes, &tender)
// 	// Normally check that the specified argument is a valid holder of tuna
// 	// we are skipping this check for this example
// 	tender.Holder = args[1]

// 	tunaAsBytes, _ = json.Marshal(tuna)
// 	err := APIstub.PutState(args[0], tunaAsBytes)
// 	if err != nil {
// 		return shim.Error(fmt.Sprintf("Failed to change tuna holder: %s", args[0]))
// 	}

// 	return shim.Success(nil)
// }

/*
 * main function *
calls the Start function
The main function starts the chaincode in the container during instantiation.
*/
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
