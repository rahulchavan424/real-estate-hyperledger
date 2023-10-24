package model

// Account represents an account, including virtual administrators and several owner accounts.
type Account struct {
	AccountId string  `json:"accountId"` // Account ID
	UserName  string  `json:"userName"`  // Account name
	Balance   float64 `json:"balance"`   // Balance
}

// RealEstate is used as collateral for sale, donation, or pledge with Encumbrance set to true by default.
// Initiating sale, donation, or pledge is only possible when Encumbrance is false.
// Proprietor and RealEstateID together form a composite key, ensuring that all real estate information can be queried by Proprietor.
type RealEstate struct {
	RealEstateID string  `json:"realEstateId"` // Real estate ID
	Proprietor   string  `json:"proprietor"`   // Owner (proprietor) (Owner's AccountId)
	Encumbrance  bool    `json:"encumbrance"`  // Whether it is used as collateral
	TotalArea    float64 `json:"totalArea"`    // Total area
	LivingSpace  float64 `json:"livingSpace"`  // Living space
}

// Selling represents a sales offer.
// It's necessary to confirm if ObjectOfSale belongs to Seller.
// The buyer is initially empty.
// Seller and ObjectOfSale together form a composite key, ensuring that all sales initiated by the seller can be queried.
type Selling struct {
	ObjectOfSale  string  `json:"objectOfSale"`  // Object being sold (RealEstateID currently for sale)
	Seller        string  `json:"seller"`        // Initiator of the sale, seller (Seller's AccountId)
	Buyer         string  `json:"buyer"`         // Participant in the sale, buyer (Buyer's AccountId)
	Price         float64 `json:"price"`         // Price
	CreateTime    string  `json:"createTime"`    // Creation time
	SalePeriod    int     `json:"salePeriod"`    // Validity period of the smart contract (in days)
	SellingStatus string  `json:"sellingStatus"` // Sale status
}

// SellingStatusConstant defines constants for selling status.
var SellingStatusConstant = func() map[string]string {
	return map[string]string{
		"saleStart": "Selling",   // Currently in the sale state, waiting for the buyer's visit
		"cancelled": "Cancelled", // Canceled by the seller or buyer's refund operation
		"expired":   "Expired",   // Sale period has expired
		"delivery":  "In Progress", // Buyer has bought and paid, waiting for the seller to confirm the payment; if the seller fails to confirm the payment, the buyer can cancel and get a refund
		"done":      "Completed",   // Seller confirms the receipt of funds, completing the transaction
	}
}

// SellingBuy represents a buyer's participation in the sale.
// The Object of Sale cannot be initiated by the buyer.
// Buyer and CreateTime form a composite key, ensuring that all sales participated by the buyer can be queried.
type SellingBuy struct {
	Buyer      string  `json:"buyer"`      // Participant in the sale, buyer (Buyer's AccountId)
	CreateTime string  `json:"createTime"` // Creation time
	Selling    Selling `json:"selling"`    // Sale object
}

// Donating represents a donation offer.
// It's necessary to confirm if ObjectOfDonating belongs to Donor.
// Specify the Grantee and wait for the Grantee's agreement to accept.
type Donating struct {
	ObjectOfDonating string `json:"objectOfDonating"` // Object being donated (RealEstateID currently being donated)
	Donor            string `json:"donor"`            // Donor (Donor's AccountId)
	Grantee          string `json:"grantee"`          // Grantee (Grantee's AccountId)
	CreateTime       string `json:"createTime"`       // Creation time
	DonatingStatus   string `json:"donatingStatus"`   // Donation status
}

// DonatingStatusConstant defines constants for donation status.
var DonatingStatusConstant = func() map[string]string {
	return map[string]string{
		"donatingStart": "In Progress", // Donor initiates a donation contract, waiting for the Grantee to confirm the donation
		"cancelled":     "Cancelled",  // Donor cancels the donation before the Grantee confirms the donation, or the Grantee cancels the acceptance of the donation
		"done":          "Completed",  // Grantee confirms receipt, completing the transaction
	}
}

// DonatingGrantee is used for Grantee to query donations.
type DonatingGrantee struct {
	Grantee    string   `json:"grantee"`    // Grantee (Grantee's AccountId)
	CreateTime string   `json:"createTime"` // Creation time
	Donating   Donating `json:"donating"`   // Donation object
}

const (
	AccountKey         = "account-key"
	RealEstateKey      = "real-estate-key"
	SellingKey         = "selling-key"
	SellingBuyKey      = "selling-buy-key"
	DonatingKey        = "donating-key"
	DonatingGranteeKey = "donating-grantee-key"
)
