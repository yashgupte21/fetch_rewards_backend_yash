package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"math"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

// receipt json structure
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item
}

// receipt item json structure
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type CacheItf interface {
	Set(key string, data interface{}, expiration time.Duration) error
	Get(key string) ([]byte, error)
}

type AppCache struct {
	client *cache.Cache
}

var myCache CacheItf

func main() {

	InitCache()

	r := mux.NewRouter()

	r.HandleFunc("/receipts/{id}/points", GetFunc).Methods("GET")
	r.HandleFunc("/receipts/process", PostFunc)

	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
	}

	log.Fatal(srv.ListenAndServe())
}

func InitCache() {
	myCache = &AppCache{
		client: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func isIntegral(val float64) bool {
	return val == float64(int(val))
}

// mehtod to calculate receipt points
func CalculatePoints(res Receipt) int { // takes json and returns points

	points := 0

	// One point for every alphanumeric character in the retailer name.
	r, _ := regexp.Compile(`[a-zA-Z0-9]`)
	for _, s := range res.Retailer {
		if r.MatchString(string(s)) {
			points = points + 1
		}
	}

	// 50 points if the total is a round dollar amount with no cents.
	num, _ := strconv.ParseFloat(res.Total, 64)
	if isIntegral(num) {
		points = points + 50
	}

	// 25 points if the total is a multiple of 0.25.
	if int(num*100)%25 == 0 {
		points = points + 25
	}

	// 5 points for every two items on the receipt.
	for i, _ := range res.Items {
		if i%2 != 0 {
			points = points + 5
		}
	}

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, item := range res.Items {
		floatPrice, _ := strconv.ParseFloat(item.Price, 64)
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			points = points + int(math.Ceil(floatPrice*0.2))
		}
	}

	// 6 points if the day in the purchase date is odd.
	intDay, _ := strconv.Atoi(strings.Split(res.PurchaseDate, "-")[2])
	if intDay%2 != 0 {
		points = points + 6
	}
	splitTime := strings.Split(res.PurchaseTime, ":")

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	st1, _ := strconv.Atoi(splitTime[0])
	st2, _ := strconv.Atoi(splitTime[1])
	if st1 == 14 {
		if st2 > 0 {
			points = points + 10
		}
	}
	if st1 == 15 {
		points = points + 10
	}
	return points
}

func PostFunc(w http.ResponseWriter, r *http.Request) {

	// update receipt_json variable to test the API with another receipt
	receipt_json := `{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}`
	res := Receipt{}
	json.Unmarshal([]byte(receipt_json), &res)

	//uuid implementation
	id := uuid.New().String()

	err := myCache.Set(id, res, 1*time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(map[string]interface{}{
		"id": id,
	})

	if err != nil {
		log.Fatal(err)
	}

	w.Write(b)
}

func GetFunc(w http.ResponseWriter, r *http.Request) {

	var result Receipt
	id := strings.Split(r.URL.String(), "/")[2]
	b, err := myCache.Get(id)
	if err != nil {
		log.Fatal(err)
	}

	if b != nil {
		// cache exist
		err := json.Unmarshal(b, &result)

		if err != nil {
			log.Fatal(err)
		}
		points := CalculatePoints(result)
		b, _ := json.Marshal(map[string]interface{}{
			"points": points,
		})
		w.Write([]byte(b))
		return
	}

	w.Write(b)
}

func (r *AppCache) Set(key string, data interface{}, expiration time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	r.client.Set(key, b, expiration)
	return nil
}

func (r *AppCache) Get(key string) ([]byte, error) {
	res, exist := r.client.Get(key)
	if !exist {
		return nil, nil
	}

	resByte, ok := res.([]byte)
	if !ok {
		return nil, errors.New("Format is not array of bytes")
	}

	return resByte, nil
}
