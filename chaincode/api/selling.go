package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// CreateSelling initiates a sale
func CreateSelling(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Validate parameters
	if len(args) != 4 {
		return shim.Error("Insufficient number of parameters")
	}
	objectOfSale := args[0]
	seller := args[1]
	price := args[2]
	salePeriod := args[3]
	if objectOfSale == "" || seller == "" || price == "" || salePeriod == "" {
		return shim.Error("Parameters contain empty values")
	}
	// Convert parameter data formats
	var formattedPrice float64
	if val, err := strconv.ParseFloat(price, 64); err != nil {
		return shim.Error(fmt.Sprintf("Failed to convert the price parameter: %s", err))
	} else {
		formattedPrice = val
	}
	var formattedSalePeriod int
	if val, err := strconv.Atoi(salePeriod); err != nil {
		return shim.Error(fmt.Sprintf("Failed to convert the salePeriod parameter: %s", err))
	} else {
		formattedSalePeriod = val
	}
	// Check if 'objectOfSale' belongs to 'seller'
	resultsRealEstate, err := utils.GetStateByPartialCompositeKeys2(stub, model.RealEstateKey, []string{seller, objectOfSale})
	if err != nil || len(resultsRealEstate) != 1 {
		return shim.Error(fmt.Sprintf("Validation failed: %s does not belong to %s: %s", objectOfSale, seller, err))
	}
	var realEstate model.RealEstate
	if err = json.Unmarshal(resultsRealEstate[0], &realEstate); err != nil {
		return shim.Error(fmt.Sprintf("CreateSelling - Deserialization error: %s", err))
	}
	// Check if the record already exists; a sale cannot be initiated more than once
	// If Encumbrance is true, it means the real estate is already in a collateralized state
	if realEstate.Encumbrance {
		return shim.Error("This real estate is already in a collateralized state and cannot be initiated for sale again")
	}
	createTime, _ := stub.GetTxTimestamp()
	selling := &model.Selling{
		ObjectOfSale:  objectOfSale,
		Seller:        seller,
		Buyer:         "",
		Price:         formattedPrice,
		CreateTime:    time.Unix(int64(createTime.GetSeconds()), int64(createTime.GetNanos())).Local().Format("2006-01-02 15:04:05"),
		SalePeriod:    formattedSalePeriod,
		SellingStatus: model.SellingStatusConstant()["saleStart"],
	}
	// Write to the ledger
	if err := utils.WriteLedger(selling, stub, model.SellingKey, []string{selling.Seller, selling.ObjectOfSale}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	// Set the real estate status to collateralized
	realEstate.Encumbrance = true
	if err := utils.WriteLedger(realEstate, stub, model.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	// Return information about the successful creation
	sellingByte, err := json.Marshal(selling)
	if err != nil {
		return shim.Error(fmt.Sprintf("Serialization error for the successful creation: %s", err))
	}
	// Success response
	return shim.Success(sellingByte)
}

// CreateSellingByBuy participates in a sale (buyer purchases)
func CreateSellingByBuy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Validate parameters
	if len(args) != 3 {
		return shim.Error("Insufficient number of parameters")
	}
	objectOfSale := args[0]
	seller := args[1]
	buyer := args[2]
	if objectOfSale == "" || seller == "" || buyer == "" {
		return shim.Error("Parameters contain empty values")
	}
	if seller == buyer {
		return shim.Error("The buyer and seller cannot be the same person")
	}
	// Obtain real estate information to be purchased based on 'objectOfSale' and 'seller' and ensure it exists
	resultsRealEstate, err := utils.GetStateByPartialCompositeKeys2(stub, model.RealEstateKey, []string{seller, objectOfSale})
	if err != nil || len(resultsRealEstate) != 1 {
		return shim.Error(fmt.Sprintf("Failed to retrieve real estate information based on %s and %s: %s", objectOfSale, seller, err))
	}
	// Obtain sale information based on 'objectOfSale' and 'seller'
	resultsSelling, err := utils.GetStateByPartialCompositeKeys2(stub, model.SellingKey, []string{seller, objectOfSale})
	if err != nil || len(resultsSelling) != 1 {
		return shim.Error(fmt.Sprintf("Failed to retrieve sale information based on %s and %s: %s", objectOfSale, seller, err))
	}
	var selling model.Selling
	if err = json.Unmarshal(resultsSelling[0], &selling); err != nil {
		return shim.Error(fmt.Sprintf("CreateSellingBuy - Deserialization error: %s", err))
	}
	// Check if the selling status is 'saleStart'
	if selling.SellingStatus != model.SellingStatusConstant()["saleStart"] {
		return shim.Error("This transaction is not in the 'saleStart' status and cannot be purchased")
	}
	// Obtain buyer information based on 'buyer'
	resultsAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{buyer})
	if err != nil || len(resultsAccount) != 1 {
		return shim.Error(fmt.Sprintf("Failed to validate buyer information: %s", err))
	}
	var buyerAccount model.Account
	if err = json.Unmarshal(resultsAccount[0], &buyerAccount); err != nil {
		return shim.Error(fmt.Sprintf("Failed to query buyer information - Deserialization error: %s", err))
	}
	if buyerAccount.UserName == "admin" {
		return shim.Error(fmt.Sprintf("The admin cannot make purchases: %s", err))
	}
	// Check if the balance is sufficient
	if buyerAccount.Balance < selling.Price {
		return shim.Error(fmt.Sprintf("The selling price is %f, and your current balance is %f. The purchase has failed.", selling.Price, buyerAccount.Balance))
	}
	// Write the buyer information into the selling transaction and change the status to 'delivery'
	selling.Buyer = buyer
	selling.SellingStatus = model.SellingStatusConstant()["delivery"]
	if err := utils.WriteLedger(selling, stub, model.SellingKey, []string{selling.Seller, selling.ObjectOfSale}); err != nil {
		return shim.Error(fmt.Sprintf("Failed to write buyer information into the selling transaction and change the status - %s", err))
	}
	createTime, _ := stub.GetTxTimestamp()
	// Write this purchase transaction to the ledger for buyer's reference
	sellingBuy := &model.SellingBuy{
		Buyer:      buyer,
		CreateTime: time.Unix(int64(createTime.GetSeconds()), int64(createTime.GetNanos())).Local().Format("2006-01-02 15:04:05"),
		Selling:    selling,
	}
	if err := utils.WriteLedger(sellingBuy, stub, model.SellingBuyKey, []string{sellingBuy.Buyer, sellingBuy.CreateTime}); err != nil {
		return shim.Error(fmt.Sprintf("Failed to write this purchase transaction to the ledger - %s", err))
	}
	sellingBuyByte, err := json.Marshal(sellingBuy)
	if err != nil {
		return shim.Error(fmt.Sprintf("Serialization error for the successful creation: %s", err))
	}
	// Purchase successful; deduct the balance. Note that the payment will be transferred to the seller's account after the seller confirms receipt. The balance is deducted from the buyer's account at this stage.
	buyerAccount.Balance -= selling.Price
	if err := utils.WriteLedger(buyerAccount, stub, model.AccountKey, []string{buyerAccount.AccountId}); err != nil {
		return shim.Error(fmt.Sprintf("Failed to deduct the buyer's balance - %s", err))
	}
	// Success response
	return shim.Success(sellingBuyByte)
}

// QuerySellingList retrieves sales (can be queried by all or by the initiating seller) - for sellers to query initiated sales
func QuerySellingList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var sellingList []model.Selling
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.SellingKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var selling model.Selling
			err := json.Unmarshal(v, &selling)
			if err != nil {
				return shim.Error(fmt.Sprintf("QuerySellingList - Deserialization error: %s", err))
			}
			sellingList = append(sellingList, selling)
		}
	}
	sellingListByte, err := json.Marshal(sellingList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QuerySellingList - Serialization error: %s", err))
	}
	return shim.Success(sellingListByte)
}

// QuerySellingListByBuyer retrieves sales based on the buyer's Account ID - for buyers to query their participated sales
func QuerySellingListByBuyer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(fmt.Sprintf("Buyer Account ID must be specified for the query"))
	}
	var sellingBuyList []model.SellingBuy
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.SellingBuyKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var sellingBuy model.SellingBuy
			err := json.Unmarshal(v, &sellingBuy)
			if err != nil {
				return shim.Error(fmt.Sprintf("QuerySellingListByBuyer - Deserialization error: %s", err))
			}
			sellingBuyList = append(sellingBuyList, sellingBuy)
		}
	}
	sellingBuyListByte, err := json.Marshal(sellingBuyList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QuerySellingListByBuyer - Serialization error: %s", err))
	}
	return shim.Success(sellingBuyListByte)
}

// UpdateSelling updates the selling status (buyer confirmation, seller or buyer cancellation)
func UpdateSelling(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Validate parameters
	if len(args) != 4 {
		return shim.Error("Insufficient number of parameters")
	}
	objectOfSale := args[0]
	seller := args[1]
	buyer := args[2]
	status := args[3]
	if objectOfSale == "" || seller == "" || status == "" {
		return shim.Error("Parameters contain empty values")
	}
	if buyer == seller {
		return shim.Error("The buyer and seller cannot be the same person")
	}
	// Obtain real estate information based on 'objectOfSale' and 'seller' and confirm its existence
	resultsRealEstate, err := utils.GetStateByPartialCompositeKeys2(stub, model.RealEstateKey, []string{seller, objectOfSale})
	if err != nil || len(resultsRealEstate) != 1 {
		return shim.Error(fmt.Sprintf("Failed to retrieve real estate information based on %s and %s: %s", objectOfSale, seller, err))
	}
	var realEstate model.RealEstate
	if err = json.Unmarshal(resultsRealEstate[0], &realEstate); err != nil {
		return shim.Error(fmt.Sprintf("UpdateSelling - Deserialization error: %s", err))
	}
	// Obtain selling information based on 'objectOfSale' and 'seller'
	resultsSelling, err := utils.GetStateByPartialCompositeKeys2(stub, model.SellingKey, []string{seller, objectOfSale})
	if err != nil || len(resultsSelling) != 1 {
		return shim.Error(fmt.Sprintf("Failed to retrieve selling information based on %s and %s: %s", objectOfSale, seller, err))
	}
	var selling model.Selling
	if err = json.Unmarshal(resultsSelling[0], &selling); err != nil {
		return shim.Error(fmt.Sprintf("UpdateSelling - Deserialization error: %s", err))
	}
	// Obtain the buying information ('sellingBuy') based on 'buyer'
	var sellingBuy model.SellingBuy
	// If the current status is 'saleStart', there is no buyer
	if selling.SellingStatus != model.SellingStatusConstant()["saleStart"] {
		resultsSellingByBuyer, err := utils.GetStateByPartialCompositeKeys2(stub, model.SellingBuyKey, []string{buyer})
		if err != nil || len(resultsSellingByBuyer) == 0 {
			return shim.Error(fmt.Sprintf("Failed to retrieve buyer's buying information based on %s: %s", buyer, err))
		}
		for _, v := range resultsSellingByBuyer {
			if v != nil {
				var s model.SellingBuy
				err := json.Unmarshal(v, &s)
				if err != nil {
					return shim.Error(fmt.Sprintf("UpdateSelling - Deserialization error: %s", err))
				}
				if s.Selling.ObjectOfSale == objectOfSale && s.Selling.Seller == seller && s.Buyer == buyer {
					// Ensure that the status is 'delivery' to prevent the case where the property has already been sold but was canceled
					if s.Selling.SellingStatus == model.SellingStatusConstant()["delivery"] {
						sellingBuy = s
						break
					}
				}
			}
		}
	}
	var data []byte
	// Determine the selling status
	switch status {
	case "done":
		// If it is a buyer's confirmation operation, ensure that the transaction is in the 'delivery' status
		if selling.SellingStatus != model.SellingStatusConstant()["delivery"] {
			return shim.Error("This transaction is not in 'delivery' status; confirmation of receipt failed")
		}
		// Obtain seller information based on 'seller'
		resultsSellerAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{seller})
		if err != nil || len(resultsSellerAccount) != 1 {
			return shim.Error(fmt.Sprintf("Failed to verify seller information: %s", err))
		}
		var accountSeller model.Account
		if err = json.Unmarshal(resultsSellerAccount[0], &accountSeller); err != nil {
			return shim.Error(fmt.Sprintf("Failed to retrieve seller information - Deserialization error: %s", err))
		}
		// Confirm receipt, transfer the payment to the seller's account
		accountSeller.Balance += selling.Price
		if err := utils.WriteLedger(accountSeller, stub, model.AccountKey, []string{accountSeller.AccountId}); err != nil {
			return shim.Error(fmt.Sprintf("Seller failed to confirm receipt of funds: %s", err))
		}
		// Transfer the property information to the buyer and reset the encumbrance status
		realEstate.Proprietor = buyer
		realEstate.Encumbrance = false
		//realEstate.RealEstateID = stub.GetTxID() // Update the real estate ID
		if err := utils.WriteLedger(realEstate, stub, model.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		// Clear the original property information
		if err := utils.DelLedger(stub, model.RealEstateKey, []string{seller, objectOfSale}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		// Set the order status to 'done' and write to the ledger
		selling.SellingStatus = model.SellingStatusConstant()["done"]
		selling.ObjectOfSale = realEstate.RealEstateID // Update the real estate ID
		if err := utils.WriteLedger(selling, stub, model.SellingKey, []string{selling.Seller, objectOfSale}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		sellingBuy.Selling = selling
		if err := utils.WriteLedger(sellingBuy, stub, model.SellingBuyKey, []string{sellingBuy.Buyer, sellingBuy.CreateTime}); err != nil {
			return shim.Error(fmt.Sprintf("Failed to write this purchase transaction to the ledger: %s", err))
		}
		data, err = json.Marshal(sellingBuy)
		if err != nil {
			return shim.Error(fmt.Sprintf("Serialization error for the purchase transaction: %s", err))
		}
		break
	case "cancelled":
		data, err = closeSelling("cancelled", selling, realEstate, sellingBuy, buyer, stub)
		if err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		break
	case "expired":
		data, err = closeSelling("expired", selling, realEstate, sellingBuy, buyer, stub)
		if err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		break
	default:
		return shim.Error(fmt.Sprintf("Status %s is not supported", status))
	}
	return shim.Success(data)
}

// closeSelling handles both cancellation and expiration cases in two scenarios:
// 1. The transaction is in 'saleStart' status
// 2. The transaction is in 'delivery' status
func closeSelling(closeStart string, selling model.Selling, realEstate model.RealEstate, sellingBuy model.SellingBuy, buyer string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	switch selling.SellingStatus {
	case model.SellingStatusConstant()["saleStart"]:
		selling.SellingStatus = model.SellingStatusConstant()[closeStart]
		// Reset the encumbrance status of the property information
		realEstate.Encumbrance = false
		if err := utils.WriteLedger(realEstate, stub, model.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
			return nil, err
		}
		if err := utils.WriteLedger(selling, stub, model.SellingKey, []string{selling.Seller, selling.ObjectOfSale}); err != nil {
			return nil, err
		}
		data, err := json.Marshal(selling)
		if err != nil {
			return nil, err
		}
		return data, nil
	case model.SellingStatusConstant()["delivery"]:
		selling.SellingStatus = model.SellingStatusConstant()[closeStart]
		// Return the balance to the buyer's account
		resultsAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{buyer})
		if err != nil || len(resultsAccount) != 1 {
			return nil, fmt.Errorf("Failed to verify buyer's information: %s", err)
		}
		var buyerAccount model.Account
		if err = json.Unmarshal(resultsAccount[0], &buyerAccount); err != nil {
			return nil, fmt.Errorf("Failed to retrieve buyer information - Deserialization error: %s", err)
		}
		buyerAccount.Balance += selling.Price
		if err := utils.WriteLedger(buyerAccount, stub, model.AccountKey, []string{buyerAccount.AccountId}); err != nil {
			return nil, fmt.Errorf("Failed to refund the buyer's account: %s", err)
		}
		if err := utils.WriteLedger(selling, stub, model.SellingKey, []string{selling.Seller, selling.ObjectOfSale}); err != nil {
			return nil, err
		}
		data, err := json.Marshal(selling)
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		return nil, fmt.Errorf("The transaction cannot be closed because it is not in 'saleStart' or 'delivery' status")
	}
}
