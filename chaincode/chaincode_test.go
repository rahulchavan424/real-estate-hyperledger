package main

import (
	"bytes"
	"chaincode/model"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func initTest(t *testing.T) *shim.MockStub {
	scc := new(BlockChainRealEstate)
	stub := shim.NewMockStub("ex01", scc)
	checkInit(t, stub, [][]byte{[]byte("init")})
	return stub
}

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Initialization failed", string(res.Message))
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) pb.Response {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
	return res
}

// Test chaincode initialization
func TestBlockChainRealEstate_Init(t *testing.T) {
	initTest(t)
}

// Test querying account information
func Test_QueryAccountList(t *testing.T) {
	stub := initTest(t)
	fmt.Println(fmt.Sprintf("1. Test getting all data\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryAccountList"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("2. Test getting multiple data\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryAccountList"),
			[]byte("5feceb66ffc8"),
			[]byte("6b86b273ff34"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("3. Test getting a single data\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryAccountList"),
			[]byte("4e07408562be"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("4. Test getting invalid data\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryAccountList"),
			[]byte("0"),
		}).Payload)))
}

// Test creating real estate
func Test_CreateRealEstate(t *testing.T) {
	stub := initTest(t)
	// Successful
	checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), // Operator
		[]byte("6b86b273ff34"), // Owner
		[]byte("50"),           // Total area
		[]byte("30"),           // Living space
	})
	// Insufficient operator permissions
	checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("6b86b273ff34"), // Operator
		[]byte("4e07408562be"), // Owner
		[]byte("50"),           // Total area
		[]byte("30"),           // Living space
	})
	// Operator should be an administrator and different from the owner
	checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), // Operator
		[]byte("5feceb66ffc8"), // Owner
		[]byte("50"),           // Total area
		[]byte("30"),           // Living space
	})
	// Owner (proprietor) information validation failed
	checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"),    // Operator
		[]byte("6b86b273ff34555"), // Owner
		[]byte("50"),              // Total area
		[]byte("30"),              // Living space
	})
	// Incorrect number of parameters
	checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), // Operator
		[]byte("6b86b273ff34"), // Owner
		[]byte("50"),           // Total area
	})
	// Parameter format conversion error
	checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), // Operator
		[]byte("6b86b273ff34"), // Owner
		[]byte("50f"),          // Total area
		[]byte("30"),           // Living space
	})
}

// Manually create some real estate properties
func checkCreateRealEstate(stub *shim.MockStub, t *testing.T) []model.RealEstate {
	var realEstateList []model.RealEstate
	var realEstate model.RealEstate
	// Successful
	resp1 := checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), // Operator
		[]byte("6b86b273ff34"), // Owner
		[]byte("50"),           // Total area
		[]byte("30"),           // Living space
	})
	resp2 := checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), // Operator
		[]byte("6b86b273ff34"), // Owner
		[]byte("80"),           // Total area
		[]byte("60.8"),         // Living space
	})
	resp3 := checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), // Operator
		[]byte("4e07408562be"), // Owner
		[]byte("60"),           // Total area
		[]byte("40"),           // Living space
	})
	resp4 := checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), // Operator
		[]byte("ef2d127de37b"), // Owner
		[]byte("80"),           // Total area
		[]byte("60"),           // Living space
	})
	json.Unmarshal(bytes.NewBuffer(resp1.Payload).Bytes(), &realEstate)
	realEstateList = append(realEstateList, realEstate)
	json.Unmarshal(bytes.NewBuffer(resp2.Payload).Bytes(), &realEstate)
	realEstateList = append(realEstateList, realEstate)
	json.Unmarshal(bytes.NewBuffer(resp3.Payload).Bytes(), &realEstate)
	realEstateList = append(realEstateList, realEstate)
	json.Unmarshal(bytes.NewBuffer(resp4.Payload).Bytes(), &realEstate)
	realEstateList = append(realEstateList, realEstate)
	return realEstateList
}

// Test querying real estate information
func Test_QueryRealEstateList(t *testing.T) {
	stub := initTest(t)
	realEstateList := checkCreateRealEstate(stub, t)

	fmt.Println(fmt.Sprintf("1. Test getting all data\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryRealEstateList"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("2. Test getting specific data\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryRealEstateList"),
			[]byte(realEstateList[0].Proprietor),
			[]byte(realEstateList[0].RealEstateID),
		}).Payload)))
	fmt.Println(fmt.Sprintf("3. Test getting invalid data\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryRealEstateList"),
			[]byte("0"),
		}).Payload)))
}

// Test initiating sales
func Test_CreateSelling(t *testing.T) {
	stub := initTest(t)
	realEstateList := checkCreateRealEstate(stub, t)
	// Successful
	checkInvoke(t, stub, [][]byte{
		[]byte("createSelling"),
		[]byte(realEstateList[0].RealEstateID), // Object for sale (RealEstateID being sold)
		[]byte(realEstateList[0].Proprietor),   // Seller (Seller's AccountId)
		[]byte("50"),                           // Price
		[]byte("30"),                           // Smart contract validity period (in days)
	})
	// Validation fails as object for sale does not belong to the seller
	checkInvoke(t, stub, [][]byte{
		[]byte("createSelling"),
		[]byte(realEstateList[0].RealEstateID), // Object for sale (RealEstateID being sold)
		[]byte(realEstateList[2].Proprietor),   // Seller (Seller's AccountId)
		[]byte("50"),                           // Price
		[]byte("30"),                           // Smart contract validity period (in days)
	})
	checkInvoke(t, stub, [][]byte{
		[]byte("createSelling"),
		[]byte("123"),                        // Object for sale (RealEstateID being sold)
		[]byte(realEstateList[0].Proprietor), // Seller (Seller's AccountId)
		[]byte("50"),                        // Price
		[]byte("30"),                        // Smart contract validity period (in days)
	})
	// Parameter errors
	checkInvoke(t, stub, [][]byte{
		[]byte("createSelling"),
		[]byte(realEstateList[0].RealEstateID), // Object for sale (RealEstateID being sold)
		[]byte(realEstateList[0].Proprietor),   // Seller (Seller's AccountId)
		[]byte("50"),                           // Price
	})
	checkInvoke(t, stub, [][]byte{
		[]byte("createSelling"),
		[]byte(""),                           // Object for sale (RealEstateID being sold)
		[]byte(realEstateList[0].Proprietor), // Seller (Seller's AccountId)
		[]byte("50"),                        // Price
		[]byte("30"),                        // Smart contract validity period (in days)
	})
}

// Test sales initiation, purchase, and related operations
func Test_QuerySellingList(t *testing.T) {
	stub := initTest(t)
	realEstateList := checkCreateRealEstate(stub, t)
	// Initiate sales first
	fmt.Println(fmt.Sprintf("Initiate sales\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("createSelling"),
		[]byte(realEstateList[0].RealEstateID), // Object for sale (RealEstateID being sold)
		[]byte(realEstateList[0].Proprietor),   // Seller (Seller's AccountId)
		[]byte("500000"),                       // Price
		[]byte("30"),                           // Smart contract validity period (in days)
	}).Payload)))
	fmt.Println(fmt.Sprintf("Initiate sales\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("createSelling"),
		[]byte(realEstateList[2].RealEstateID), // Object for sale (RealEstateID being sold)
		[]byte(realEstateList[2].Proprietor),   // Seller (Seller's AccountId)
		[]byte("600000"),                       // Price
		[]byte("40"),                           // Smart contract validity period (in days)
	}).Payload)))
	// Query successfully
	fmt.Println(fmt.Sprintf("1. Query all\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("querySellingList"),
	}).Payload)))
	fmt.Println(fmt.Sprintf("2. Query by %s\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("querySellingList"),
		[]byte(realEstateList[0].Proprietor),
	}).Payload)))
	// Purchase
	fmt.Println(fmt.Sprintf("3. Query the account balance of %s before purchase\n%s", realEstateList[2].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryAccountList"),
		[]byte(realEstateList[2].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("4. Start purchase\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("createSellingByBuy"),
		[]byte(realEstateList[0].RealEstateID), // Object for sale (RealEstateID being sold)
		[]byte(realEstateList[0].Proprietor),   // Seller (Seller's AccountId)
		[]byte(realEstateList[2].Proprietor),   // Buyer (Buyer's AccountId)
	}).Payload)))
	fmt.Println(fmt.Sprintf("5. Query the account balance of %s after purchase\n%s", realEstateList[2].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryAccountList"),
		[]byte(realEstateList[2].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("6. Query purchase success information of the seller\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("querySellingList"),
		[]byte(realEstateList[0].Proprietor), // Buyer (Buyer's AccountId)
	}).Payload)))
	fmt.Println(fmt.Sprintf("7. Query purchase success information of the buyer\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("querySellingListByBuyer"),
		[]byte(realEstateList[2].Proprietor), // Buyer (Buyer's AccountId)
	}).Payload)))
	fmt.Println(fmt.Sprintf("8. Account balance of the seller %s before confirming receipt\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryAccountList"),
		[]byte(realEstateList[0].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("9. Account balance of the buyer %s before confirming receipt\n%s", realEstateList[2].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryAccountList"),
		[]byte(realEstateList[2].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("10. Real estate information of the seller %s before confirming receipt\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryRealEstateList"),
		[]byte(realEstateList[0].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("11. Real estate information of the buyer %s before confirming receipt\n%s", realEstateList[2].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryRealEstateList"),
		[]byte(realEstateList[2].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("12. Seller %s confirms receipt\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("updateSelling"),
		[]byte(realEstateList[0].RealEstateID), // Object for sale (RealEstateID being sold)
		[]byte(realEstateList[0].Proprietor),   // Seller (Seller's AccountId)
		[]byte(realEstateList[2].Proprietor),   // Buyer (Buyer's AccountId)
		[]byte("done"),                         // Confirm receipt
	}).Payload)))
	//fmt.Println(fmt.Sprintf("13. Seller %s cancels receipt\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
	//	[]byte("updateSelling"),
	//	[]byte(realEstateList[0].RealEstateID), // Object for sale (RealEstateID being sold)
	//	[]byte(realEstateList[0].Proprietor),   // Seller (Seller's AccountId)
	//	[]byte(realEstateList[2].Proprietor),   // Buyer (Buyer's AccountId)
	//	[]byte("cancelled"),                    // Cancel receipt
	//}).Payload)))
	fmt.Println(fmt.Sprintf("13. Account balance of the seller %s after confirming receipt\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryAccountList"),
		[]byte(realEstateList[0].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("14. Account balance of the buyer %s after confirming receipt\n%s", realEstateList[2].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryAccountList"),
		[]byte(realEstateList[2].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("15. Real estate information of the seller %s after confirming receipt\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryRealEstateList"),
		[]byte(realEstateList[0].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("16. Real estate information of the buyer %s after confirming receipt\n%s", realEstateList[2].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryRealEstateList"),
		[]byte(realEstateList[2].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("17. Seller %s queries purchase success information\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("querySellingList"),
		[]byte(realEstateList[0].Proprietor), // Buyer (Buyer's AccountId)
	}).Payload)))
	fmt.Println(fmt.Sprintf("18. Buyer %s queries purchase success information\n%s", realEstateList[2].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("querySellingListByBuyer"),
		[]byte(realEstateList[2].Proprietor), // Buyer (Buyer's AccountId)
	}).Payload)))
}

// Test the purchase of properties for sale
func Test_Purchase(t *testing.T) {
	stub := initTest(t)
	realEstateList := checkCreateRealEstate(stub, t)
	// Initiate sales first
	fmt.Println(fmt.Sprintf("Initiate sales\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("createSelling"),
		[]byte(realEstateList[0].RealEstateID), // Object for sale (RealEstateID being sold)
		[]byte(realEstateList[0].Proprietor),   // Seller (Seller's AccountId)
		[]byte("500000"),                       // Price
		[]byte("30"),                           // Smart contract validity period (in days)
	}).Payload)))
	fmt.Println(fmt.Sprintf("Initiate sales\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("createSelling"),
		[]byte(realEstateList[2].RealEstateID), // Object for sale (RealEstateID being sold)
		[]byte(realEstateList[2].Proprietor),   // Seller (Seller's AccountId)
		[]byte("600000"),                       // Price
		[]byte("40"),                           // Smart contract validity period (in days)
	}).Payload)))
	// Purchase
	fmt.Println(fmt.Sprintf("1. Start purchase\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("createSellingByBuy"),
		[]byte(realEstateList[0].RealEstateID), // Object for sale (RealEstateID being sold)
		[]byte(realEstateList[0].Proprietor),   // Seller (Seller's AccountId)
		[]byte(realEstateList[2].Proprietor),   // Buyer (Buyer's AccountId)
	}).Payload)))
	fmt.Println(fmt.Sprintf("2. Seller %s confirms receipt\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("updateSelling"),
		[]byte(realEstateList[0].RealEstateID), // Object for sale (RealEstateID being sold)
		[]byte(realEstateList[0].Proprietor),   // Seller (Seller's AccountId)
		[]byte(realEstateList[2].Proprietor),   // Buyer (Buyer's AccountId)
		[]byte("done"),                         // Confirm receipt
	}).Payload)))
}

func main() {
	test := &BlockChainRealEstate{}
	t := &testing.T{}
	TestBlockChainRealEstate_Init(t)
	Test_QueryAccountList(t)
	Test_CreateRealEstate(t)
	Test_QueryRealEstateList(t)
	Test_CreateSelling(t)
	Test_QuerySellingList(t)
	Test_Purchase(t)
}

