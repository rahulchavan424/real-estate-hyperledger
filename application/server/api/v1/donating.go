package v1

import (
	bc "application/blockchain"
	"application/pkg/app"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DonatingRequestBody struct {
	ObjectOfDonating string `json:"objectOfDonating"` // Object of Donation
	Donor            string `json:"donor"`            // Donor
	Grantee          string `json:"grantee"`          // Grantee
}

type DonatingListQueryRequestBody struct {
	Donor string `json:"donor"`
}

type DonatingListQueryByGranteeRequestBody struct {
	Grantee string `json:"grantee"`
}

type UpdateDonatingRequestBody struct {
	ObjectOfDonating string `json:"objectOfDonating"` // Object of Donation
	Donor            string `json:"donor"`            // Donor
	Grantee          string `json:"grantee"`          // Grantee
	Status           string `json:"status"`           // Status to be updated
}

func CreateDonating(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(DonatingRequestBody)
	// Parse Body parameters
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "Failed", fmt.Sprintf("Parameter error: %s", err.Error()))
		return
	}
	if body.ObjectOfDonating == "" || body.Donor == "" || body.Grantee == "" {
		appG.Response(http.StatusBadRequest, "Failed", "ObjectOfDonating, Donor, and Grantee cannot be empty")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.ObjectOfDonating))
	bodyBytes = append(bodyBytes, []byte(body.Donor))
	bodyBytes = append(bodyBytes, []byte(body.Grantee))
	// Invoke smart contract
	resp, err := bc.ChannelExecute("createDonating", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "Failed", err.Error())
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "Failed", err.Error())
		return
	}
	appG.Response(http.StatusOK, "Success", data)
}

func QueryDonatingList(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(DonatingListQueryRequestBody)
	// Parse Body parameters
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "Failed", fmt.Sprintf("Parameter error: %s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.Donor != "" {
		bodyBytes = append(bodyBytes, []byte(body.Donor))
	}
	// Invoke smart contract
	resp, err := bc.ChannelQuery("queryDonatingList", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "Failed", err.Error())
		return
	}
	// Deserialize JSON
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "Failed", err.Error())
		return
	}
	appG.Response(http.StatusOK, "Success", data)
}

func QueryDonatingListByGrantee(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(DonatingListQueryByGranteeRequestBody)
	// Parse Body parameters
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "Failed", fmt.Sprintf("Parameter error: %s", err.Error()))
		return
	}
	if body.Grantee == "" {
		appG.Response(http.StatusBadRequest, "Failed", "AccountId must be specified")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.Grantee))
	// Invoke smart contract
	resp, err := bc.ChannelQuery("queryDonatingListByGrantee", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "Failed", err.Error())
		return
	}
	// Deserialize JSON
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "Failed", err.Error())
		return
	}
	appG.Response(http.StatusOK, "Success", data)
}

func UpdateDonating(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(UpdateDonatingRequestBody)
	// Parse Body parameters
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "Failed", fmt.Sprintf("Parameter error: %s", err.Error()))
		return
	}
	if body.ObjectOfDonating == "" || body.Donor == "" || body.Grantee == "" || body.Status == "" {
		appG.Response(http.StatusBadRequest, "Failed", "Parameters cannot be empty")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.ObjectOfDonating))
	bodyBytes = append(bodyBytes, []byte(body.Donor))
	bodyBytes = append(bodyBytes, []byte(body.Grantee))
	bodyBytes = append(bodyBytes, []byte(body.Status))
	// Invoke smart contract
	resp, err := bc.ChannelExecute("updateDonating", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "Failed", err.Error())
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "Failed", err.Error())
		return
	}
	appG.Response(http.StatusOK, "Success", data)
}
