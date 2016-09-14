/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/hyperledger/fabric/events/consumer"
	pb "github.com/hyperledger/fabric/protos"
)

type adapter struct {
	notfy              chan *pb.Event_Block
	rejected           chan *pb.Event_Rejection
	cEvent             chan *pb.Event_ChaincodeEvent
	listenToRejections bool
	chaincodeID        string
}

//GetInterestedEvents implements consumer.EventAdapter interface for registering interested events
func (a *adapter) GetInterestedEvents() ([]*pb.Interest, error) {
	if a.chaincodeID != "" {

		fmt.Printf("Listening to chain code %s", a.chaincodeID)

		return []*pb.Interest{
			{EventType: pb.EventType_BLOCK},
			{EventType: pb.EventType_REJECTION},
			{EventType: pb.EventType_CHAINCODE,
				RegInfo: &pb.Interest_ChaincodeRegInfo{
					ChaincodeRegInfo: &pb.ChaincodeReg{
						ChaincodeID: a.chaincodeID,
						EventName:   ""}}}}, nil

	}

	return []*pb.Interest{{EventType: pb.EventType_BLOCK}, {EventType: pb.EventType_REJECTION}}, nil
	//return []*pb.Interest{{EventType: pb.EventType_BLOCK}, {EventType: pb.EventType_REJECTION}, {EventType: pb.EventType_CHAINCODE}}, nil
}

//Recv implements consumer.EventAdapter interface for receiving events
func (a *adapter) Recv(msg *pb.Event) (bool, error) {
	if o, e := msg.Event.(*pb.Event_Block); e {
		a.notfy <- o
		return true, nil
	}
	if o, e := msg.Event.(*pb.Event_Rejection); e && a.listenToRejections {
		a.rejected <- o
		return true, nil
	}
	if o, e := msg.Event.(*pb.Event_ChaincodeEvent); e {
		a.cEvent <- o
		return true, nil
	}
	a.notfy <- nil
	return false, nil
}

//Disconnected implements consumer.EventAdapter interface for disconnecting
func (a *adapter) Disconnected(err error) {
	fmt.Printf("Disconnected...exiting\n")
	os.Exit(1)
}

func createEventClient(eventAddress string, listenToRejections bool, cid string) *adapter {
	var obcEHClient *consumer.EventsClient

	done := make(chan *pb.Event_Block)
	reject := make(chan *pb.Event_Rejection)
	adapter := &adapter{notfy: done, rejected: reject, listenToRejections: listenToRejections, chaincodeID: cid, cEvent: make(chan *pb.Event_ChaincodeEvent)}
	obcEHClient, _ = consumer.NewEventsClient(eventAddress, 5, adapter)

	if err := obcEHClient.Start(); err != nil {
		fmt.Printf("could not start chat %s\n", err)
		obcEHClient.Stop()
		return nil
	}

	return adapter
}

func main() {
	var eventAddress string
	var listenToRejections bool
	var chaincodeID string
	flag.StringVar(&eventAddress, "events-address", "0.0.0.0:7053", "address of events server")
	flag.BoolVar(&listenToRejections, "listen-to-rejections", false, "whether to listen to rejection events")
	flag.StringVar(&chaincodeID, "events-from-chaincode", "", "listen to events from given chaincode")
	flag.Parse()

	fmt.Printf("Event Address: %s\n", eventAddress)

	a := createEventClient(eventAddress, listenToRejections, chaincodeID)
	if a == nil {
		fmt.Printf("Error creating event client\n")
		return
	}

	for {
		select {
		case b := <-a.notfy:
			fmt.Printf("\n")
			fmt.Printf("\n")
			fmt.Printf("Received block\n")
			fmt.Printf("--------------\n")
			for _, r := range b.Block.Transactions {
				fmt.Printf("Transaction:\n\t[%v]\n", r)
			}
		case r := <-a.rejected:
			fmt.Printf("\n")
			fmt.Printf("\n")
			fmt.Printf("Received rejected transaction\n")
			fmt.Printf("--------------\n")
			fmt.Printf("Transaction error:\n%s\t%s\n", r.Rejection.Tx.Txid, r.Rejection.ErrorMsg)
		case ce := <-a.cEvent:
			fmt.Printf("\n")
			fmt.Printf("\n")
			fmt.Printf("Received chaincode event\n")
			fmt.Printf("------------------------\n")
			fmt.Printf("Chaincode Event:%v\n", ce)

			processChainCodeEvent(ce, chaincodeID)
		}
	}
}

func processChainCodeEvent(cEvent *pb.Event_ChaincodeEvent, chaincodeID string) {

	eventName := *(&cEvent.ChaincodeEvent.EventName)
	payload := *(&cEvent.ChaincodeEvent.Payload)

	fmt.Printf("Chaincode Event Name :%s\n", eventName)
	fmt.Printf("Chaincode Event Payload :%s\n", string(payload))

	if eventName == "createAsset" {

		url := "http://10.0.2.2:8081/create"

		returnPayload, err := postToBigChainProxy(url, payload)

		if err != nil {
			fmt.Printf("Error Sending create to big chain db %s", err.Error())
			return
		}

		returnFromHlPayload, err := postToHyperledgerProxy("http://localhost:7050/chaincode", returnPayload, chaincodeID)

		if err != nil {
			fmt.Printf("Error Sending createconfirm to hyperledger %s %s", err.Error(), string(returnFromHlPayload))
			return
		}

		return
	}

	if eventName == "transferAsset" {

		url := "http://10.0.2.2:8081/transfer"

		postToBigChainProxy(url, payload)
		return
	}

	fmt.Printf("Not processing event :%s\n", eventName)
}

func postToBigChainProxy(url string, payload []byte) ([]byte, error) {

	fmt.Println("URL:>", url)

	var jsonStr = payload

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	if resp.Status != "200" {
		return body, errors.New("Bad response from bigchain proxy" + resp.Status)
	}

	return body, nil
}

func postToHyperledgerProxy(url string, payload []byte, chaincodeID string) ([]byte, error) {

	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{
            "jsonrpc": "2.0",
            "method": "invoke",
            "params": {
                "type": 1,
                "chaincodeID": { "name": "` + chaincodeID + `" },
                "ctorMsg": { "function": "createconfirm",  "args": ["` + string(payload) + `"]  }
            },
            "id": 10
        }`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	if resp.Status != "200" {
		return body, errors.New("Bad response from hyperledger" + resp.Status)
	}

	return body, nil
}
