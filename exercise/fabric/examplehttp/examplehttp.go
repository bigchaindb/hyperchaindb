package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getContent() {

	url := "http://localhost:9984/api/v1/transactions/eae4fda7d86667294bdcef3768dc2ef77cb34c5ba9cf25f6a4cc7f1ea5bdf9ae"

	res, err := http.Get(url)

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Results: %v\n", body)

}

func retrieveFromBigChainDb() {

	id := "eae4fda7d86667294bdcef3768dc2ef77cb34c5ba9cf25f6a4cc7f1ea5bdf9ae"

	url := "http://localhost:9984/api/v1/transactions/"

	urlTran := url + id

	resp, err := http.Get(urlTran)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}

func queryFabric(account string) {

	url := "http://localhost:7050/chaincode"
	fmt.Println("URL:>", url)

	//var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)

	var jsonStr = []byte(`{
            "jsonrpc": "2.0",
            "method": "query",
            "params": {
                "type": 1,
                "chaincodeID": { "name": "d2d3fe9e0f52b60d2762a8d990e7ae1365dcbf0ee481a0e61e73dbe63531917f72bcd23380b99f28c2e526305170c7314342dda45b4877889e87e8152a341f4f"  },
                "ctorMsg": { "function": "query",  "args": ["` + account + `"]  }
            },
            "id": 3
        }`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}

func postToBigChain() {

	url := "http://localhost:9984/api/v1/transactions"
	fmt.Println("URL:>", url)

	//var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)

	var jsonStr = []byte(`{
  "id":"7ad5a4b83bc8c70c4fd7420ff3c60693ab8e6d0e3124378ca69ed5acd2578792",
  "transaction":{
      "conditions":[
          {
              "cid":0,
              "condition":{
                  "details":{
                      "bitmask":32,
                      "public_key":"CwA8s2QYQBfNz4WvjEwmJi83zYr7JhxRhidx6uZ5KBVd",
                      "signature":null,
                      "type":"fulfillment",
                      "type_id":4
                  },
                  "uri":"cc:4:20:sVA_3p8gvl8yRFNTomqm6MaavKewka6dGYcFAuPrRXQ:96"
              },
              "owners_after":[
                  "CwA8s2QYQBfNz4WvjEwmJi83zYr7JhxRhidx6uZ5KBVd"
              ]
          }
      ],
      "data":{
          "payload":null,
          "uuid":"a9999d69-6cde-4b80-819d-ed57f6abe257"
      },
      "fulfillments":[
          {
              "owners_before":[
                  "JEAkEJqLbbgDRAtMm8YAjGp759Aq2qTn9eaEHUj2XePE"
              ],
              "fid":0,
              "fulfillment":"cf:4:__Y_Um6H73iwPe6ejWXEw930SQhqVGjtAHTXilPp0P01vE_Cx6zs3GJVoO1jhPL18C94PIVkLTGMUB2aKC9qsbIb3w8ejpOf0_I3OCuTbPdkd6r2lKMeVftMyMxkeWoM",
              "input":{
                  "cid":0,
                  "txid":"598ce4e9a29837a1c6fc337ee4a41b61c20ad25d01646754c825b1116abd8761"
              }
          }
      ],
      "operation":"CREATE",
      "timestamp":"1471423869",
      "version":1
   }
}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}

func main() {

	var account string

	flag.StringVar(&account, "account", "a", "listen to events from given chaincode")
	flag.Parse()

	//queryFabric(account)

	postToBigChain()

	//getContent()
	//retrieveFromBigChainDb()

	os.Exit(0)
}
