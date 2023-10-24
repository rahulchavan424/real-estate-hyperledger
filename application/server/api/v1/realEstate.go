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

type RealEstateRequestBody struct {
	AccountId   string  `json:"accountId"`   // Operator ID
	Proprietor  string  `json:"proprietor"`  // Proprietor (Owner) (Owner's Account ID)
	TotalArea   float64 `json:"totalArea"`   // Total Area
	LivingSpace float64 `json:"livingSpace"` // Living Space
}

type RealEstateQueryRequestBody struct {
	Proprietor string `json:"proprietor"` // Proprietor (Owner) (Owner's Account ID)
}

func CreateRealEstate(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(RealEstateRequestBody)
	// Parse the body parameters
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "Failure", fmt.Sprintf("Parameter error: %s", err.Error()))
		return
	}
	if body.TotalArea <= 0 || body.LivingSpace <= 0 || body.LivingSpace > body.TotalArea {
		appG.Response(http.StatusBadRequest, "Failure", "TotalArea and LivingSpace must be greater than 0, and LivingSpace must be less than or equal to TotalArea")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.AccountId))
	bodyBytes = append(bodyBytes, []byte(body.Proprietor))
	bodyBytes = append(bodyBytes, []byte(strconv.FormatFloat(body.TotalArea, 'E', -1, 64)))
	bodyBytes = append(bodyBytes, []byte(strconv.FormatFloat(body.LivingSpace, 'E', -1, 64)))
	// Invoke the smart contract
	resp, err := bc.ChannelExecute("createRealEstate", bodyBytes)
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

func QueryRealEstateList(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(RealEstateQueryRequestBody)
	// Parse the body parameters
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "Failure", fmt.Sprintf("Parameter error: %s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.Proprietor != "" {
		bodyBytes = append(bodyBytes, []byte(body.Proprietor))
	}
	// Invoke the smart contract
	resp, err := bc.ChannelQuery("queryRealEstateList", bodyBytes)
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
