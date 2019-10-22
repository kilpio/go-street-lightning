package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"log"
	"net/http"

	//https://github.com/etcd-io/bbolt	
	bolt "go.etcd.io/bbolt"
)

//Lightning type for mosdata records
type Lightning struct {
	GlobalID           string
	DurationOfLighting string
	Date               string
	OnTime             string
	OffTime            string
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func main() {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	resp, err := http.Get("https://op.mos.ru/EHDWSREST/catalog/export/get?id=407891")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Fatal(err)
	}

	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {
		//        fmt.Println("Reading file:", zipFile.Name)
		unzippedFileBytes, err := readZipFile(zipFile)
		if err != nil {
			log.Println(err)
			continue
		}

		

		MyJSON := unzippedFileBytes // this is unzipped file bytes
		// MyJSON now contains the uzipped file bytes
		// it's an array of json records for each day described by the Lightning structure

		//days is the slice of an array containing Lightning structures as elements
		days := make([]Lightning, 0)

		
		//unmarshalling th� days slice elements
		json.Unmarshal(MyJSON, &days)

		

		for _, day := range days {

			db.Update(func(tx *bolt.Tx) error {
				b, err := tx.CreateBucketIfNotExists([]byte("daylightning"))
				if err != nil {
					return err
				}


				out, err := json.Marshal(day)
				if err != nil {
					panic(err)
				}
				return b.Put([]byte(day.Date), out)
			})
			
		} //now all timings are in the posts //end for
	}

	// selected_date := "13.10.2019"

	// db.View(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte("daylightning"))
	// 	v := b.Get([]byte(selected_date))
	// 	var dada Lightning
	// 	json.Unmarshal(v, &dada)
	// 	fmt.Printf("%s", dada.DurationOfLighting)
	// 	return nil
	// })

}
