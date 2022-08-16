package main

import (
	"os"

	"github.com/ApogeeNetworking/ciscowireless"
	"github.com/ApogeeNetworking/insights"
	"github.com/subosito/gotenv"
)

var baseURL, apiToken string

func init() {
	gotenv.Load()
	baseURL = os.Getenv("BASEURL")
	apiToken = os.Getenv("API_TOKEN")
}

func main() {
	da := insights.NewService(baseURL, apiToken, true)

	schoolID := "361439dc-7d9f-4536-9706-acdbd7e1c06a"

	cc := ciscowireless.NewService(
		"ipAddress",
		"user",
		"password",
		"",
		true,
	)
	apDb, _ := cc.AccessPoints.Get()
	var daAps []insights.SyncAp
	for _, ap := range apDb {
		daAps = append(daAps, insights.SyncAp{
			Name:    ap.Name,
			MacAddr: []string{ap.MacAddr},
			Serial:  ap.Serial,
			Switch:  "aSwitch",
		})
		// if len(daAps) == 100 {
		// 	syncCount += len(daAps)
		// 	insert, err := da.SyncAps(schoolID, daAps)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	fmt.Println(insert)
		// 	daAps = nil
		// }
	}
	da.BulkSyncAps(schoolID, daAps)
}
