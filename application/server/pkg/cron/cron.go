package cron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	bc "application/blockchain"
	"application/model"

	"github.com/robfig/cron/v3"
)

const spec = "0 0 0 * * ?" // Execute every day at 0:00
//const spec = "*/10 * * * * ?" // Execute every 10 seconds, used for testing

func Init() {
	c := cron.New(cron.WithSeconds()) // Supports second-level precision
	_, err := c.AddFunc(spec, GoRun)
	if err != nil {
		log.Printf("Scheduled task failed to start %s", err)
	}
	c.Start()
	log.Printf("Scheduled task has started")
	select {}
}

func GoRun() {
	log.Printf("Scheduled task has started")
	// First, retrieve all sales
	resp, err := bc.ChannelQuery("querySellingList", [][]byte{}) // Call the smart contract
	if err != nil {
		log.Printf("Scheduled task - querySellingList failed: %s", err.Error())
		return
	}
	// Deserialize JSON
	var data []model.Selling
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		log.Printf("Scheduled task - JSON deserialization failed: %s", err.Error())
		return
	}
	for _, v := range data {
		// Select those with statuses "In Sale" and "In Delivery"
		if v.SellingStatus == model.SellingStatusConstant()["saleStart"] ||
			v.SellingStatus == model.SellingStatusConstant()["delivery"] {
			// Calculate the validity period in days
			day, _ := time.ParseDuration(fmt.Sprintf("%dh", v.SalePeriod*24))
			local, _ := time.LoadLocation("Local")
			t, _ := time.ParseInLocation("2006-01-02 15:04:05", v.CreateTime, local)
			vTime := t.Add(day)
			// If time.Now() is greater than vTime, it means it has expired
			if time.Now().Local().After(vTime) {
				// Change the status to "Expired"
				var bodyBytes [][]byte
				bodyBytes = append(bodyBytes, []byte(v.ObjectOfSale))
				bodyBytes = append(bodyBytes, []byte(v.Seller))
				bodyBytes = append(bodyBytes, []byte(v.Buyer))
				bodyBytes = append(bodyBytes, []byte("expired"))
				// Call the smart contract
				resp, err := bc.ChannelExecute("updateSelling", bodyBytes)
				if err != nil {
					return
				}
				var data map[string]interface{}
				if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
					return
				}
				fmt.Println(data)
			}
		}
	}
}
