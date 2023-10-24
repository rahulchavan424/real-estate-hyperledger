package utils

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// WriteLedger writes to the ledger
func WriteLedger(obj interface{}, stub shim.ChaincodeStubInterface, objectType string, keys []string) error {
	// Create a composite key
	var key string
	if val, err := stub.CreateCompositeKey(objectType, keys); err != nil {
		return errors.New(fmt.Sprintf("%s - Error creating composite key: %s", objectType, err))
	} else {
		key = val
	}
	bytes, err := json.Marshal(obj)
	if err != nil {
		return errors.New(fmt.Sprintf("%s - Error serializing JSON data: %s", objectType, err))
	}
	// Write to the blockchain ledger
	if err := stub.PutState(key, bytes); err != nil {
		return errors.New(fmt.Sprintf("%s - Error writing to the blockchain ledger: %s", objectType, err))
	}
	return nil
}

// DelLedger deletes from the ledger
func DelLedger(stub shim.ChaincodeStubInterface, objectType string, keys []string) error {
	// Create a composite key
	var key string
	if val, err := stub.CreateCompositeKey(objectType, keys); err != nil {
		return errors.New(fmt.Sprintf("%s - Error creating composite key: %s", objectType, err))
	} else {
		key = val
	}
	// Delete from the blockchain ledger
	if err := stub.DelState(key); err != nil {
		return errors.New(fmt.Sprintf("%s - Error deleting from the blockchain ledger: %s", objectType, err))
	}
	return nil
}

// GetStateByPartialCompositeKeys queries data based on composite keys (suitable for getting all, multiple, or single data)
// Splits the keys for querying
func GetStateByPartialCompositeKeys(stub shim.ChaincodeStubInterface, objectType string, keys []string) (results [][]byte, err error) {
	if len(keys) == 0 {
		// If the length of keys passed is 0, it means to retrieve and return all data
		// Retrieve relevant data from the blockchain by primary key, performing a fuzzy search on the primary key
		resultIterator, err := stub.GetStateByPartialCompositeKey(objectType, keys)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("%s - Error getting all data: %s", objectType, err))
		}
		defer resultIterator.Close()

		// Check if returned data is not empty, then iterate through the data; otherwise, return an empty array
		for resultIterator.HasNext() {
			val, err := resultIterator.Next()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("%s - Error with returned data: %s", objectType, err))
			}

			results = append(results, val.GetValue())
		}
	} else {
		// If the length of keys passed is not 0, it means to retrieve the corresponding data
		for _, v := range keys {
			// Create a composite key
			key, err := stub.CreateCompositeKey(objectType, []string{v})
			if err != nil {
				return nil, errors.New(fmt.Sprintf("%s - Error creating composite key: %s", objectType, err))
			}
			// Retrieve data from the ledger
			bytes, err := stub.GetState(key)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("%s - Error getting data: %s", objectType, err))
			}

			if bytes != nil {
				results = append(results, bytes)
			}
		}
	}

	return results, nil
}

// GetStateByPartialCompositeKeys2 queries data based on composite keys (suitable for getting all or specific data)
func GetStateByPartialCompositeKeys2(stub shim.ChaincodeStubInterface, objectType string, keys []string) (results [][]byte, err error) {
	// Retrieve relevant data from the blockchain by primary key, performing a fuzzy search on the primary key
	resultIterator, err := stub.GetStateByPartialCompositeKey(objectType, keys)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s - Error getting all data: %s", objectType, err))
	}
	defer resultIterator.Close()

	// Check if returned data is not empty, then iterate through the data; otherwise, return an empty array
	for resultIterator.HasNext() {
		val, err := resultIterator.Next()
		if err != nil {
			return nil, errors.New(fmt.Sprintf("%s - Error with returned data: %s", objectType, err))
		}

		results = append(results, val.GetValue())
	}
	return results, nil
}
