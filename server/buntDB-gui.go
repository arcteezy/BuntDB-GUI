package main

import (
	"encoding/json"
	"fmt"
	"log"
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
		log.Println(err)
	}
	defer db.Close()

	// Starting web server
	http.HandleFunc("/writeData", writeData)
	http.HandleFunc("/getAllData", getAllData)

	fmt.Printf("BuntDB-GUI started at %s\n", "http://localhost:8080")

	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Function to write data
func writeData(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Page hit : /writeData")
	var key, value string

	// open the specified DB
	if dbPath := r.URL.Query()["db"]; len(dbPath) > 1 {
		db, err = buntdb.Open(dbPath[0])
		// log error
		if err != nil {
			log.Println(err)
		}
		defer db.Close()
	}

	// get key
	if keys := r.URL.Query()["key"]; len(keys) > 1 {
		key = keys[0]
	}

	// get value
	if values := r.URL.Query()["value"]; len(values) > 1 {
		value = values[0]
	}

	// Write data
	err = db.Update(func(tx *buntdb.Tx) error {
		_, done, err := tx.Set(key, value, nil)

		if err != nil {
			log.Println(err)
			return err
		}

		if done {

			if _, err = w.Write([]byte("Successfully set new value to DB")); err != nil {
				log.Println(err)
			}
		} else {

			if _, err = w.Write([]byte("write failed")); err != nil {
				log.Println(err)
			}
		}
		return err
	})
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
		log.Println(err)
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

	// open the specified DB
	if dbPath := r.URL.Query()["db"]; len(dbPath) > 1 {
		db, err = buntdb.Open(dbPath[0])
		// log error
		if err != nil {
			log.Println(err)
		}
		defer db.Close()
	}

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
		log.Println(err)
	}
	fmt.Println(dbContent)
	// Marshal into string
	jsondata, err := json.Marshal(dbContent)
	if err != nil {
		log.Println(err)
	}
	// Data into response structure
	response.Body = string(jsondata)
	// Marshal response
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}
	// Sending response
	if _, err = w.Write(jsonResponse); err != nil {
		log.Println(err)
	}
}

// Function to delete data
func deleteData(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Page hit : /deleteData")
	var key string

	// open the specified DB
	if dbPath := r.URL.Query()["db"]; len(dbPath) > 1 {
		db, err = buntdb.Open(dbPath[0])
		// log error
		if err != nil {
			log.Println(err)
		}
		defer db.Close()
	}

	// get key
	if keys := r.URL.Query()["key"]; len(keys) > 1 {
		key = keys[0]
	}

	// Update db
	err = db.Update(func(tx *buntdb.Tx) error {
		_, done, err := tx.Delete(key)

		if err != nil {
			log.Println(err)
			return err
		}

		if done {

			if _, err = w.Write([]byte("Successfully removed data from DB")); err != nil {
				log.Println(err)
			}
		} else {

			if _, err = w.Write([]byte("delete failed")); err != nil {
				log.Println(err)
			}
		}
		return err
	})
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
		log.Println(err)
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

