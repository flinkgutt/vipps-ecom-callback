package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type OrderCallback struct {
	MerchantSerialNumber string    `json:"merchantSerialNumber"`
	OrderID              string `json:"orderId"`
	ShippingDetails      struct {
		Address struct {
			AddressLine1 string `json:"addressLine1"`
			AddressLine2 string `json:"addressLine2"`
			City         string `json:"city"`
			Country      string `json:"country"`
			PostCode     string `json:"postCode"`
		} `json:"address"`
		ShippingCost   float64 `json:"shippingCost"`
		ShippingMethod string  `json:"shippingMethod"`
	} `json:"shippingDetails"`
	TransactionInfo struct {
		Amount        int    `json:"amount"`
		Status        string `json:"status"`
		TimeStamp     string `json:"timeStamp"`
		TransactionID string `json:"transactionId"`
	} `json:"transactionInfo"`
	UserDetails struct {
		BankIDVerified string `json:"bankIdVerified"`
		DateOfBirth    string `json:"dateOfBirth"`
		Email          string `json:"email"`
		FirstName      string `json:"firstName"`
		LastName       string `json:"lastName"`
		MobileNumber   string `json:"mobileNumber"`
		Ssn            string `json:"ssn"`
		UserID         string `json:"userId"`
	} `json:"userDetails"`
}

type ShippingRequest struct {
	AddressID    int    `json:"addressId"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressType  string `json:"addressType"`
	City         string `json:"city"`
	Country      string `json:"country"`
	PostCode     string `json:"postCode"`
}

func payments(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("Callback on order:", params.ByName("orderId"))
	orderId := params.ByName("orderId")
	decoder := json.NewDecoder(r.Body)
	var order OrderCallback
	err := decoder.Decode(&order)
	if err != nil {
		fmt.Println("Failed to decode JSON for order:", orderId, "with error message: '", err, "'")
		w.WriteHeader(500) // Not strictly speaking true... buut whatever..
		return
	}
	fmt.Printf("%+v\n", order)
	fmt.Printf("OrderID: %v, MSN: %v \n", order.OrderID, order.MerchantSerialNumber)
	fmt.Println("----------------------------------------------------------------------")

}

func shipping(_ http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Printf("Vipps asking for Shipping Details for order %s - EXPRESS CHECKOUT\n", params.ByName("orderId"))

	decoder := json.NewDecoder(r.Body)
	var sr ShippingRequest
	err := decoder.Decode(&sr)
	if err != nil {
		fmt.Println("Error decoding shippingDetails request:", err)
		return
	}
	fmt.Println("ID\tType\tLine 1\tPostCode\tCity")
	fmt.Printf("%v\t%v\t%v\t%v\t%v\n", sr.AddressID, sr.AddressType, sr.AddressLine1, sr.PostCode, sr.City)
	fmt.Println("----------------------------------------------------------------------")
}

func main() {
	r := httprouter.New()
	r.GET("/", index)
	r.POST("/vipps/v2/payments/:orderId", payments)
	r.POST("/vipps/v2/payments/:orderId/shippingDetails", shipping)
	err := http.ListenAndServe(":80", r)
	if err != nil {
		fmt.Println("Error on startup:", err)
	}
}

func index(_ http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("Incoming GET request hitting index function")
	fmt.Println("Remote Address:", r.RemoteAddr)
	fmt.Println("----------------------------------------------------------------------")
}
