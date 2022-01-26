package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ricablanca-gem-acosta/valhalla-merchant/api"
)

func main() {
	err := api.InitDb(true)
	defer api.CloseDb(false)
	if err != nil {
		log.Fatal("Failed to initialize db")
	}
	fmt.Println("Started Merchant API at port 3000")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", api.GetRouter()))
}
