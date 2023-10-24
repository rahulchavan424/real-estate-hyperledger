package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// CreateDonating initiates a donation.
func CreateDonating(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Validate parameters
	if len(args) != 3 {
		return shim.Error("Incorrect number of parameters")
	}
	objectOfDonating := args[0]
	donor := args[1]
	grantee := args[2]
	if objectOfDonating == "" || donor == "" || grantee == "" {
		return shim.Error("Empty parameters are not allowed")
	}
	if donor == grantee {
		return shim.Error("The donor and grantee cannot be the same person")
	}
	// Check if objectOfDonating belongs to the donor
	resultsRealEstate, err := utils.GetStateByPartialCompositeKeys2(stub, model.RealEstateKey, []string{donor, objectOfDonating})
	if err != nil || len(resultsRealEstate) != 1 {
		return shim.Error(fmt.Sprintf("Verification failed: %s", err))
	}
	var realEstate model.RealEstate
	if err = json.Unmarshal(resultsRealEstate[0], &realEstate); err != nil {
		return shim.Error(fmt.Sprintf("CreateDonating - Deserialization error: %s", err))
	}
	// Get grantee information
	resultsAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{grantee})
	if err != nil || len(resultsAccount) != 1 {
		return shim.Error(fmt.Sprintf("Grantee information verification failed: %s", err))
	}
	var accountGrantee model.Account
	if err = json.Unmarshal(resultsAccount[0], &accountGrantee); err != nil {
		return shim.Error(fmt.Sprintf("Querying operator information - Deserialization error: %s", err))
	}
	if accountGrantee.UserName == "Admin" {
		return shim.Error("Cannot donate to the admin")
	}
	// Check if the record already exists, no duplicate donations allowed
	// If Encumbrance is true, it means the real estate is already under collateral
	if realEstate.Encumbrance {
		return shim.Error("This real estate is already being used as collateral and cannot be donated")
	}
	createTime, _ := stub.GetTxTimestamp()
	donating := &model.Donating{
		ObjectOfDonating: objectOfDonating,
		Donor:            donor,
		Grantee:          grantee,
		CreateTime:       time.Unix(int64(createTime.GetSeconds()), int64(createTime.GetNanos())).Local().Format("2006-01-02 15:04:05"),
		DonatingStatus:   model.DonatingStatusConstant()["donatingStart"],
	}
	// Write to the ledger
	if err := utils.WriteLedger(donating, stub, model.DonatingKey, []string{donating.Donor, donating.ObjectOfDonating, donating.Grantee}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	// Set the real estate as under collateral status
	realEstate.Encumbrance = true
	if err := utils.WriteLedger(realEstate, stub, model.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	// Write the donation transaction for grantee to query
	donatingGrantee := &model.DonatingGrantee{
		Grantee:    grantee,
		CreateTime: time.Unix(int64(createTime.GetSeconds()), int64(createTime.GetNanos())).Local().Format("2006-01-02 15:04:05"),
		Donating:   *donating,
	}
	if err := utils.WriteLedger(donatingGrantee, stub, model.DonatingGranteeKey, []string{donatingGrantee.Grantee, donatingGrantee.CreateTime}); err != nil {
		return shim.Error(fmt.Sprintf("Failed to write this donation transaction: %s", err))
	}
	donatingGranteeByte, err := json.Marshal(donatingGrantee)
	if err != nil {
		return shim.Error(fmt.Sprintf("Serialization of created information failed: %s", err))
	}
	// Success response
	return shim.Success(donatingGranteeByte)
}

// QueryDonatingList queries the list of donations (all or by the donor) for donors to query.
func QueryDonatingList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var donatingList []model.Donating
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.DonatingKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var donating model.Donating
			err := json.Unmarshal(v, &donating)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryDonatingList - Deserialization error: %s", err))
			}
			donatingList = append(donatingList, donating)
		}
	}
	donatingListByte, err := json.Marshal(donatingList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryDonatingList - Serialization error: %s", err))
	}
	return shim.Success(donatingListByte)
}

// QueryDonatingListByGrantee queries the list of donations by grantee (for grantees to query).
func QueryDonatingListByGrantee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Must specify the grantee AccountId to query")
	}
	var donatingGranteeList []model.DonatingGrantee
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.DonatingGranteeKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var donatingGrantee model.DonatingGrantee
			err := json.Unmarshal(v, &donatingGrantee)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryDonatingListByGrantee - Deserialization error: %s", err))
			}
			donatingGranteeList = append(donatingGranteeList, donatingGrantee)
		}
	}
	donatingGranteeListByte, err := json.Marshal(donatingGranteeList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryDonatingListByGrantee - Serialization error: %s", err))
	}
	return shim.Success(donatingGranteeListByte)
}

// UpdateDonating updates the donation status (confirm or cancel).
func UpdateDonating(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Validate parameters
	if len(args) != 4 {
		return shim.Error("Incorrect number of parameters")
	}
	objectOfDonating := args[0]
	donor := args[1]
	grantee := args[2]
	status := args[3]
	if objectOfDonating == "" || donor == "" || grantee == "" || status == "" {
		return shim.Error("Empty parameters are not allowed")
	}
	if donor == grantee {
		return shim.Error("The donor and grantee cannot be the same person")
	}
	// Get the real estate information that the donor wants to donate to, confirm its existence
	resultsRealEstate, err := utils.GetStateByPartialCompositeKeys2(stub, model.RealEstateKey, []string{donor, objectOfDonating})
	if err != nil || len(resultsRealEstate) != 1 {
		return shim.Error(fmt.Sprintf("Failed to get real estate information for %s and %s: %s", objectOfDonating, donor, err))
	}
	var realEstate model.RealEstate
	if err = json.Unmarshal(resultsRealEstate[0], &realEstate); err != nil {
		return shim.Error(fmt.Sprintf("UpdateDonating - Deserialization error: %s", err))
	}
	// Get grantee information
	resultsGranteeAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{grantee})
	if err != nil || len(resultsGranteeAccount) != 1 {
		return shim.Error(fmt.Sprintf("Grantee information verification failed: %s", err))
	}
	var accountGrantee model.Account
	if err = json.Unmarshal(resultsGranteeAccount[0], &accountGrantee); err != nil {
		return shim.Error(fmt.Sprintf("Querying grantee information - Deserialization error: %s", err))
	}
	// Get the donation information for the objectOfDonating, donor, and grantee
	resultsDonating, err := utils.GetStateByPartialCompositeKeys2(stub, model.DonatingKey, []string{donor, objectOfDonating, grantee})
	if err != nil || len(resultsDonating) != 1 {
		return shim.Error(fmt.Sprintf("Failed to get donation information for %s, %s, and %s: %s", objectOfDonating, donor, grantee, err))
	}
	var donating model.Donating
	if err = json.Unmarshal(resultsDonating[0], &donating); err != nil {
		return shim.Error(fmt.Sprintf("UpdateDonating - Deserialization error: %s", err))
	}
	// Regardless of completion or cancellation, ensure the donation is in the "donatingStart" status
	if donating.DonatingStatus != model.DonatingStatusConstant()["donatingStart"] {
		return shim.Error("This transaction is not in the 'donatingStart' status and cannot be confirmed/cancelled")
	}
	// Get the donation transaction for grantee to purchase
	var donatingGrantee model.DonatingGrantee
	resultsDonatingGrantee, err := utils.GetStateByPartialCompositeKeys2(stub, model.DonatingGranteeKey, []string{grantee})
	if err != nil || len(resultsDonatingGrantee) == 0 {
		return shim.Error(fmt.Sprintf("Failed to get grantee information for %s: %s", grantee, err))
	}
	for _, v := range resultsDonatingGrantee {
		if v != nil {
			var s model.DonatingGrantee
			err := json.Unmarshal(v, &s)
			if err != nil {
				return shim.Error(fmt.Sprintf("UpdateDonating - Deserialization error: %s", err))
			}
			if s.Donating.ObjectOfDonating == objectOfDonating && s.Donating.Donor == donor && s.Grantee == grantee {
				// Must also check that the status is "donatingStart" to prevent cases where the real estate is already transacted but got canceled
				if s.Donating.DonatingStatus == model.DonatingStatusConstant()["donatingStart"] {
					donatingGrantee = s
					break
				}
			}
		}
	}
	var data []byte
	// Check the donation status
	switch status {
	case "done":
		// Transfer real estate information to the grantee and reset the collateral status
		realEstate.Proprietor = grantee
		realEstate.Encumbrance = false
		if err := utils.WriteLedger(realEstate, stub, model.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		// Clear the original real estate information
		if err := utils.DelLedger(stub, model.RealEstateKey, []string{donor, objectOfDonating}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		// Set the donation status to "done" and update the real estate ID
		donating.DonatingStatus = model.DonatingStatusConstant()["done"]
		donating.ObjectOfDonating = realEstate.RealEstateID
		if err := utils.WriteLedger(donating, stub, model.DonatingKey, []string{donating.Donor, objectOfDonating, grantee}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		donatingGrantee.Donating = donating
		if err := utils.WriteLedger(donatingGrantee, stub, model.DonatingGranteeKey, []string{donatingGrantee.Grantee, donatingGrantee.CreateTime}); err != nil {
			return shim.Error(fmt.Sprintf("Failed to write this donation transaction: %s", err))
		}
		data, err = json.Marshal(donatingGrantee)
		if err != nil {
			return shim.Error(fmt.Sprintf("Serialization of donation transaction information failed: %s", err))
		}
		break
	case "cancelled":
		// Reset the collateral status of real estate information
		realEstate.Encumbrance = false
		if err := utils.WriteLedger(realEstate, stub, model.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		// Update the donation status to "cancelled"
		donating.DonatingStatus = model.DonatingStatusConstant()["cancelled"]
		if err := utils.WriteLedger(donating, stub, model.DonatingKey, []string{donating.Donor, donating.ObjectOfDonating, donating.Grantee}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		donatingGrantee.Donating = donating
		if err := utils.WriteLedger(donatingGrantee, stub, model.DonatingGranteeKey, []string{donatingGrantee.Grantee, donatingGrantee.CreateTime}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		data, err = json.Marshal(donatingGrantee)
		if err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		break
	default:
		return shim.Error(fmt.Sprintf("Status %s is not supported", status))
	}
	return shim.Success(data)
}
