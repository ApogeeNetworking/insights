package main

import (
	"fmt"
	"os"

	"github.com/ApogeeNetworking/insights"
	"github.com/subosito/gotenv"
)

var baseURL, apiToken string

func init() {
	gotenv.Load()
	baseURL = os.Getenv("BASEURL")
	apiToken = os.Getenv("API_TOKEN")
}

type JFile struct {
	Data []insights.AccessPoint `json:"data"`
}

func main() {
	da := insights.NewService(baseURL, apiToken, true)
	swuAps, _ := da.GetAps("dad7d5ed-b3c5-48b3-b676-9a53394abc17")
	for _, swuAp := range swuAps.Data {
		fmt.Println(swuAp.Name)
	}
}
