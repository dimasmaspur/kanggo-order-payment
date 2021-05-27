package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

type Order struct {
	ID        int    `form:"id" json:"id"`
	UserId    int    `form:"user_id" json:"user_id"`
	ProductId int    `form:"product_id" json:"product_id"`
	Name      string `form:"name" json:"name"`
	Price     int    `form:"price" json:"price"`
	Quantity  int    `form:"quantity" json:"quantity"`
	Amount    int    `form:"amount" json:"amount"`
	Status    string `form:"status" json:"status"`
}

type Result struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func main() {
	db, err = gorm.Open("mysql", "root:@/payment_kanggo?charset=utf8&parseTime=True")
	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection established")
	}

	db.AutoMigrate(&Order{})
	handleRequests()
}

func handleRequests() {

	log.Println("Start the development server at http://127.0.0.1:8090")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/order", createOrder).Methods("POST")
	myRouter.HandleFunc("/order", getOrders).Methods("GET")
	myRouter.HandleFunc("/order/{id}", getOrder).Methods("GET")
	myRouter.HandleFunc("/order/{id}", updateOrder).Methods("PUT")
	myRouter.HandleFunc("/order/pay/{id}", createPayment).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8090", myRouter))

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var order Order
	json.Unmarshal(payloads, &order)

	db.Create(&order)

	res := Result{Status: "success", Data: order}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: get products")

	orders := []Order{}
	db.Find(&orders)

	res := Result{Status: "success", Data: orders}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	var order Order

	db.First(&order, orderID)

	res := Result{Status: "success", Data: order}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)

	var orderUpdates Order
	json.Unmarshal(payloads, &orderUpdates)

	var order Order
	db.First(&order, orderID)
	db.Model(&order).Updates(orderUpdates)

	res := Result{Status: "success", Data: order}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func createPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)

	var orderUpdates Order
	json.Unmarshal(payloads, &orderUpdates)

	var order Order
	db.First(&order, orderID)
	db.Model(&order).Updates(orderUpdates)

	res := Result{Status: "success", Data: order}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
