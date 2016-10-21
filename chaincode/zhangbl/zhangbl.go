/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

	if len(args) != 4 {
		//if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}

	if function == "createAccNO" {
		accNO := args[0]
		var amount int
		var err error

		if 1 == len(args) {
			amount = 0
		} else {
			amount, err = strconv.Atoi(args[1])
			if err != nil {
				return nil, errors.New("Expecting integer value for asset holding")
			}
		}

		err = stub.PutState(accNO, []byte(strconv.Itoa(amount)))
		if err != nil {
			return nil, err
		}
	}

	if function == "transferAcc" {
		if len(args) != 3 {
			return nil, errors.New("Incorrect number of arguments. Expecting 3")
		}

		srcAcc := args[0]                    // 转出户名
		tarAcc := args[1]                    // 转入户名
		amount, err := strconv.Atoi(args[2]) // 金额

		// 获取A的余额
		srcAccvalbytes, err := stub.GetState(srcAcc)
		if err != nil {
			return nil, errors.New("Failed to get state")
		}
		if srcAccvalbytes == nil {
			return nil, errors.New("Entity not found")
		}
		srcAccVal, _ := strconv.Atoi(string(srcAccvalbytes))

		// 获取B的月
		tarAccvalbytes, err := stub.GetState(tarAcc)
		if err != nil {
			return nil, errors.New("Failed to get state")
		}
		if tarAccvalbytes == nil {
			return nil, errors.New("Entity not found")
		}
		tarAccVal, _ := strconv.Atoi(string(tarAccvalbytes))

		srcAccVal -= amount
		tarAccVal += amount

		// Write the state back to the ledger
		err = stub.PutState(srcAcc, []byte(strconv.Itoa(srcAccVal)))
		if err != nil {
			return nil, err
		}

		err = stub.PutState(tarAcc, []byte(strconv.Itoa(tarAccVal)))
		if err != nil {
			return nil, err
		}
	}

	if function == "ECashPrinter" {
		accNO := args[0]                     // 转出户名
		amount, err := strconv.Atoi(args[1]) // 金额

		// 获取A的余额
		accValbytes, err := stub.GetState(accNO)
		if err != nil {
			return nil, errors.New("Failed to get state")
		}
		if accValbytes == nil {
			return nil, errors.New("Entity not found")
		}
		accVal, _ := strconv.Atoi(string(accValbytes))

		accVal += amount

		// Write the state back to the ledger
		err = stub.PutState(accNO, []byte(strconv.Itoa(accVal)))
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var A string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return Avalbytes, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
