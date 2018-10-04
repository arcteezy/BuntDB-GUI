package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/buntdb"
)

// Data strucutre : key-value pair for Bunt DB
type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// APIResponse structure : send data to UI
type APIResponse struct {
	Body string `json:"body"`
}

// Array to keep whole data
var dbContent []Data

// Response variable
var response APIResponse

// Database pointer
var db *buntdb.DB

// Error interface
var err error

func main() {
	// Open the data.db file. It will be created if it doesn't exist.
	db, err = buntdb.Open("bunt.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Starting web server
	http.HandleFunc("/writeData", writeData)
	http.HandleFunc("/getAllData", getAllData)
	if err = http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

// Function to write data
func writeData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Page hit : /writeData")
	// Variable to check save success
	writeSuccess := false
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	// Unmarshal
	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data.Key, data.Value)
	if data.Key != "" {
		// Write data
		err = db.Update(func(tx *buntdb.Tx) error {
			_, _, err := tx.Set(data.Key, data.Value, nil)
			return err
		})
		if err != nil {
			fmt.Println(err)
		} else {
			writeSuccess = true
		}
	}
	// Response
	if writeSuccess {
		response.Body = "SUCCESSFUL"
	} else {
		response.Body = "FAILED"
	}
	// Marshal response
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	// Sending response
	w.Write(jsonResponse)
}

// Function to read data
func getAllData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Page hit : /getAllData")
	// Reset data
	dbContent = nil
	// Open a view transaction
	err = db.View(func(tx *buntdb.Tx) error {
		// Iterate through content
		err := tx.Ascend("", func(key, value string) bool {
			dbContent = append(dbContent, Data{Key: key, Value: value})
			return true
		})
		return err
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dbContent)
	// Marshal into string
	jsondata, err := json.Marshal(dbContent)
	if err != nil {
		fmt.Println(err)
	}
	// Data into response structure
	response.Body = string(jsondata)
	// Marshal response
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	// Sending response
	w.Write(jsonResponse)
}
