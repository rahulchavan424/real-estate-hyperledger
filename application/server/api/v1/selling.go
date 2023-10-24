package v1

import (
	bc "application/blockchain"
	"application/pkg/app"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SellingRequestBody struct {
	ObjectOfSale string  `json:"objectOfSale"` // Sale object (RealEstateID being sold)
	Seller       string  `json:"seller"`       // Initiator of the sale, seller (Seller's Account ID)
	Price        float64 `json:"price"`        // Price
	SalePeriod   int     `json:"salePeriod"`   // Validity period of the smart contract (in days)
}

type SellingByBuyRequestBody struct {
	ObjectOfSale string `json:"objectOfSale"` // Sale object (RealEstateID being sold)
	Seller       string `json:"seller"`       // Initiator of the sale, seller (Seller's Account ID)
	Buyer        string `json:"buyer"`        // Buyer (Buyer's Account ID)
}

type SellingListQueryRequestBody struct {
	Seller string `json:"seller"` // Initiator of the sale, seller (Seller's Account ID)
}

type SellingListQueryByBuyRequestBody struct {
	Buyer string `json:"buyer"` // Buyer (Buyer's Account ID)
}

type UpdateSellingRequestBody struct {
	ObjectOfSale string `json:"objectOfSale"` // Sale object (RealEstateID being sold)
	Seller       string `json:"seller"`       // Initiator of the sale, seller (Seller's Account ID)
	Buyer        string `json:"buyer"`        // Buyer (Buyer's Account ID)
	Status       string `json:"status"`       // Status to be changed
}

func CreateSelling(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SellingRequestBody)
	// Parse the body parameters
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "Failure", fmt.Sprintf("Parameter error: %s", err.Error()))
		return
	}
	if body.ObjectOfSale == "" || body.Seller == "" {
		appG.Response(http.StatusBadRequest, "Failure", "ObjectOfSale and Seller cannot be empty")
		return
	}
	if body.Price <= 0 || body.SalePeriod <= 0 {
		appG.Response(http.StatusBadRequest, "Failure", "Price and SalePeriod (in days) must be greater than 0")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.ObjectOfSale))
	bodyBytes = append(bodyBytes, []byte(body.Seller))
	bodyBytes = append(bodyBytes, []byte(strconv.FormatFloat(body.Price, 'E', -1, 64)))
	bodyBytes = append(bodyBytes, []byte(strconv.Itoa(body.SalePeriod)))
	// Invoke the smart contract
	resp, err := bc.ChannelExecute("createSelling", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "Failure", err.Error())
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "Failure", err.Error())
		return
	}
	appG.Response(http.StatusOK, "Success", data)
}

func CreateSellingByBuy(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SellingByBuyRequestBody)
	// Parse the body parameters
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "Failure", fmt.Sprintf("Parameter error: %s", err.Error()))
		return
	}
	if body.ObjectOfSale == "" || body.Seller == "" || body.Buyer == "" {
		appG.Response(http.StatusBadRequest, "Failure", "Parameters cannot be empty")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.ObjectOfSale))
	bodyBytes = append(bodyBytes, []byte(body.Seller))
	bodyBytes = append(bodyBytes, []byte(body.Buyer))
	// Invoke the smart contract
	resp, err := bc.ChannelExecute("createSellingByBuy", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "Failure", err.Error())
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "Failure", err.Error())
		return
	}
	appG.Response(http.StatusOK, "Success", data)
}

func QuerySellingList(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SellingListQueryRequestBody)
	// Parse the body parameters
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "Failure", fmt.Sprintf("Parameter error: %s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.Seller != "" {
		bodyBytes = append(bodyBytes, []byte(body.Seller))
	}
	// Invoke the smart contract
	resp, err := bc.ChannelQuery("querySellingList", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "Failure", err.Error())
		return
	}
	// Deserialize JSON
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "Failure", err.Error())
		return
	}
	appG.Response(http.StatusOK, "Success", data)
}

func QuerySellingListByBuyer(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SellingListQueryByBuyRequestBody)
	// Parse the body parameters
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "Failure", fmt.Sprintf("Parameter error: %s", err.Error()))
		return
	}
	if body.Buyer == "" {
		appG.Response(http.StatusBadRequest, "Failure", "Buyer Account ID must be specified for the query")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.Buyer))
	// Invoke the smart contract
	resp, err := bc.ChannelQuery("querySellingListByBuyer", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "Failure", err.Error())
		return
	}
	// Deserialize JSON
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "Failure", err.Error())
		return
	}
	appG.Response(http.StatusOK, "Success", data)
}

func UpdateSelling(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(UpdateSellingRequestBody)
	// Parse the body parameters
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "Failure", fmt.Sprintf("Parameter error: %s", err.Error()))
		return
	}
	if body.ObjectOfSale == "" || body.Seller == "" || body.Status == "" {
		appG.Response(http.StatusBadRequest, "Failure", "Parameters cannot be empty")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.ObjectOfSale))
	bodyBytes = append(bodyBytes, []byte(body.Seller))
	bodyBytes = append(bodyBytes, []byte(body.Buyer))
	bodyBytes = append(bodyBytes, []byte(body.Status))
	// Invoke the smart contract
	resp, err := bc.ChannelExecute("updateSelling", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "Failure", err.Error())
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "Failure", err.Error())
		return
	}
	appG.Response(http.StatusOK, "Success", data)
}
