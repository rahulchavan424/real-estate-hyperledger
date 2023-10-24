package model

// Selling Sale Proposal
// Need to determine if ObjectOfSale belongs to Seller
// Buyer is initially empty
// Seller and ObjectOfSale together serve as a composite key to ensure that all sales initiated by a seller can be queried through the seller's AccountID
type Selling struct {
	ObjectOfSale  string  `json:"objectOfSale"`  // Object for sale (RealEstateID of the real estate currently for sale)
	Seller        string  `json:"seller"`        // Initiator of the sale, seller (Seller's AccountID)
	Buyer         string  `json:"buyer"`         // Participant in the sale, buyer (Buyer's AccountID)
	Price         float64 `json:"price"`         // Price
	CreateTime    string  `json:"createTime"`    // Creation time
	SalePeriod    int     `json:"salePeriod"`    // Validity period of the smart contract (in days)
	SellingStatus string  `json:"sellingStatus"` // Sales status
}

// SellingStatusConstant Sales Status
var SellingStatusConstant = func() map[string]string {
	return map[string]string{
		"saleStart": "In Sale",     // In the process of sale, waiting for buyers to visit
		"cancelled": "Canceled",    // Sale canceled by the seller or canceled due to a buyer refund
		"expired":   "Expired",     // Sale period has expired
		"delivery":  "In Delivery", // Buyer has made the purchase and payment, waiting for seller to confirm receipt; if the seller fails to confirm receipt, the buyer can cancel and get a refund
		"done":      "Completed",   // Seller confirms receipt of funds, transaction completed
	}
}

// Donating Donation Proposal
// Need to determine if ObjectOfDonating belongs to Donor
// Specify the Grantee and wait for their agreement to receive
type Donating struct {
	ObjectOfDonating string `json:"objectOfDonating"` // Object for donation (RealEstateID of the real estate currently being donated)
	Donor            string `json:"donor"`            // Donor (Donor's AccountID)
	Grantee          string `json:"grantee"`          // Grantee (Grantee's AccountID)
	CreateTime       string `json:"createTime"`       // Creation time
	DonatingStatus   string `json:"donatingStatus"`   // Donation status
}

// DonatingStatusConstant Donation Status
var DonatingStatusConstant = func() map[string]string {
	return map[string]string{
		"donatingStart": "In Donation", // Donor initiates the donation contract, waiting for the Grantee to confirm acceptance
		"cancelled":     "Canceled",   // Donor cancels the donation before the Grantee confirms acceptance or the Grantee cancels the acceptance of the donation
		"done":          "Completed",  // Grantee confirms acceptance, transaction completed
	}
}
