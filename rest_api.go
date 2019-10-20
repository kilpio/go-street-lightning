	package main

	//imports
	import (
		"fmt"
		"log"
		"net/http"
		"github.com/gorilla/mux"
		bolt "go.etcd.io/bbolt"
	//	"bytes"
		"encoding/json"
	//	"io/ioutil"
	)

	//Lightning type for mosdata records
	type Lightning struct {
		GlobalID           string
	DurationOfLighting string
	Date               string
	OnTime             string
	OffTime            string
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func returnLightning(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
//header('Access-Control-Allow-Origin: *');
//header('Content-Type: application/json');

	params := mux.Vars(r)
	date := params["date"]
	//fmt.Fprintln(string(date))	
	
	selected_date := string(date)

	db, err := bolt.Open(`C:\Users\kilpio\go\zip\my.db`, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("daylightning"))
		v := b.Get([]byte(selected_date))
		var dada Lightning
		json.Unmarshal(v, &dada)
	//	fmt.Printf("%s", dada.DurationOfLighting)

	response := `{"Lightning Durtion":"`+ string(dada.DurationOfLighting) + `",`
	response  += `"OnTime":"`+ string(dada.OnTime) + `",`
	response  += `"OffTime":"`+ string(dada.OffTime) + `"}`

	fmt.Fprintf(w,response)
		return nil
	})
	
	
	
	
} 

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/Lightning/{date}", returnLightning)
	log.Fatal(http.ListenAndServe(":8080", router))
}

