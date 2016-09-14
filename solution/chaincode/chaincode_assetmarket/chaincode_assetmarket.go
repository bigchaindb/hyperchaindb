/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Asset complex type to hold some kind of asset
type Asset struct {
	Name       string  `json:"name"`
	Identifier string  `json:"identifier"`
	Type       string  `json:"type"`
	Owner      string  `json:"owner"`
	Price      float64 `json:"amount"`
	BigChainId string  `json:"bigchainid"`
}

type AssetTransfer struct {
	OwnerFrom  string `json:"ownerfrom"`
	OwnerTo    string `json:"ownerto"`
	Identifier string `json:"identifer"`
	BigChainId string `json:"bigchainid"`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func assetsFromBytes(assetsAsBytes []byte) (map[string]*Asset, error) {

	var decodedMap map[string]*Asset

	buf := bytes.NewBuffer(assetsAsBytes)

	d := gob.NewDecoder(buf)

	// Decoding the serialized data
	err := d.Decode(&decodedMap)
	if err != nil {

		return nil, err
	}

	return decodedMap, nil
}

func assetsToBytes(assets map[string]*Asset) ([]byte, error) {

	gob.Register(Asset{})

	b := new(bytes.Buffer)

	e := gob.NewEncoder(b)

	// Encoding the map
	err := e.Encode(assets)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

//Init - initialise storage
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error

	//Save initial assets
	assets := make(map[string]*Asset)

	assetsAsBytes, err := assetsToBytes(assets)

	if err != nil {
		return nil, errors.New("can't serialise assets, error:" + err.Error())
	}

	for _, owner := range args {

		fmt.Printf("Initialising assets for owner %s", owner)

		err = stub.PutState(owner, assetsAsBytes)

		if err != nil {
			return nil, err
		}
	}

	err = stub.PutState("TransferTemp", assetsAsBytes)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

//Invoke -  Invoke functions that change the state of the ledger
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Println("in Invoke", function, args)

	if function == "create" {
		//Creates a new asset
		return t.create(stub, args)
	}

	if function == "transfer" {
		//transfers an asset between parties
		return t.transfer(stub, args)
	}

	if function == "createconfirmed" {
		//call back from bigchain to confirm asset has been created
		return t.createConfirmed(stub, args)
	}

	if function == "transferconfirmed" {
		//call back from bigchain to confirm asset has been transfered
		return t.transferConfirmed(stub, args)
	}

	return nil, nil
}

// sell - puts an asset on the market
/*
func (t *SimpleChaincode) sell(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("sell - Incorrect number of arguments. Expecting 1")
	}

	assetAsJson := args[0]

	asset := Asset{}

	json.Unmarshal([]byte(assetAsJson), &asset)

	//Get the assets from state
	assetsAsBytes, err := stub.GetState("sales")

	if err != nil {
		return nil, errors.New("can't get sales from state")
	}

	//deserialise
	assets := assetsFromBytes(assetsAsBytes)

	assets[asset.Identifier] = &asset

	updatedassetsBytes := assetsToBytes(assets)

	err = stub.PutState("sales", updatedassetsBytes)

	return nil, nil
}
*/

// create asset and associates it with user
func (t *SimpleChaincode) create(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("sell - Incorrect number of arguments. Expecting 1 which should be the asset")
	}

	var err error

	assetAsJson := args[0]

	asset := Asset{}

	json.Unmarshal([]byte(assetAsJson), &asset)

	//Get the assets from state
	assetsAsBytes, err := stub.GetState(asset.Owner)

	if err != nil {
		return nil, errors.New("can't get assets from state, error:" + err.Error())
	}

	assets, err := assetsFromBytes(assetsAsBytes)

	if err != nil {
		return nil, errors.New("can't deserialise assets, error:" + err.Error())

	}

	assets[asset.Identifier] = &asset

	updatedassetsBytes, err := assetsToBytes(assets)

	if err != nil {
		return nil, errors.New("can't serialise assets error:" + err.Error())
	}

	err = stub.PutState(asset.Owner, updatedassetsBytes)

	createAssetEvent := assetAsJson

	err = stub.SetEvent("createAsset", []byte(createAssetEvent))

	return nil, nil
}

func (t *SimpleChaincode) createConfirmed(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("sell - Incorrect number of arguments. Expecting 1 which should be the asset")
	}

	var err error

	assetAsJson := args[0]

	asset := Asset{}

	json.Unmarshal([]byte(assetAsJson), &asset)

	//Get the assets from state
	assetsAsBytes, err := stub.GetState(asset.Owner)

	if err != nil {
		return nil, errors.New("can't get assets for " + asset.Owner + " from state")
	}

	//deserialise
	assets, err := assetsFromBytes(assetsAsBytes)

	if err != nil {
		return nil, errors.New("can't deserialise assets, error:" + err.Error())
	}

	originalAsset := assets[asset.Identifier]

	//Update the original asset with the big chain Id
	originalAsset.BigChainId = asset.BigChainId

	assets[asset.Identifier] = originalAsset

	updatedAssetsBytes, err := assetsToBytes(assets)

	if err != nil {
		return nil, errors.New("can't serialise assets, error:" + err.Error())
	}

	err = stub.PutState(asset.Owner, updatedAssetsBytes)

	return nil, nil
}

func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("transfer - Incorrect number of arguments. Expecting 3")
	}

	ownerFrom := args[0]
	ownerTo := args[1]
	assetIdentifier := args[2]

	//Get the 'from' assets

	//Get the assets from state
	assetsAsBytes, err := stub.GetState(ownerFrom)

	if err != nil {
		return nil, errors.New("can't get assets for " + ownerFrom)
	}

	//deserialise
	assets, err := assetsFromBytes(assetsAsBytes)

	if err != nil {
		return nil, errors.New("can't deserialise assets, error:" + err.Error())
	}

	//Find the right one
	asset := assets[assetIdentifier]

	//Put the asset into pending transfer

	transferTempAssetsAsBytes, err := stub.GetState("TransferTemp")

	if err != nil {
		return nil, errors.New("can't get temp transfer assets")
	}

	//deserialise - temp store
	transferTempAssets, err := assetsFromBytes(transferTempAssetsAsBytes)

	if err != nil {
		return nil, errors.New("can't deserialise assets, error:" + err.Error())
	}

	transferTempAssets[asset.Identifier] = asset

	updatedAssetsBytes, err := assetsToBytes(assets)

	if err != nil {
		return nil, errors.New("can't serialise assets, error:" + err.Error())
	}

	err = stub.PutState("TransferTemp", updatedAssetsBytes)

	//Remove from owner map
	delete(assets, assetIdentifier)

	updatedOwnerAssetsBytes, err := assetsToBytes(assets)

	if err != nil {
		return nil, errors.New("can't serialise owner assets, error:" + err.Error())
	}

	err = stub.PutState(ownerFrom, updatedOwnerAssetsBytes)

	//Send out the eventStuff

	transferEvent := &AssetTransfer{OwnerFrom: ownerFrom,
		OwnerTo:    ownerTo,
		Identifier: assetIdentifier,
		BigChainId: asset.BigChainId}

	transferEventAsBytes, err := json.Marshal(transferEvent)

	err = stub.SetEvent("transferAsset", transferEventAsBytes)

	return nil, nil
}

func (t *SimpleChaincode) transferConfirmed(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("transferConfirmed - Incorrect number of arguments. Expecting 1 which is transfer Payload")
	}

	transferAsJson := args[0]

	transfer := AssetTransfer{}

	json.Unmarshal([]byte(transferAsJson), &transfer)

	//Get all Assets in Temp Store
	transferTempAssetsAsBytes, err := stub.GetState("TransferTemp")

	if err != nil {
		return nil, errors.New("can't get assets from TransferTemp")
	}

	transferTempAssets, err := assetsFromBytes(transferTempAssetsAsBytes)

	if err != nil {
		return nil, errors.New("can't deserialise temp assets, error:" + err.Error())
	}

	//Get the asset to transfer from the temp store
	asset := transferTempAssets[transfer.Identifier]

	//Update the owner
	asset.Owner = transfer.OwnerTo

	destinationAssetsAsBytes, err := stub.GetState(transfer.OwnerTo)

	if err != nil {
		return nil, errors.New("can't get the destination assets from ownerTo" + transfer.OwnerTo)
	}

	//deserialise
	destinationAssets, err := assetsFromBytes(destinationAssetsAsBytes)

	if err != nil {
		return nil, errors.New("can't deserialise assets, error:" + err.Error())
	}

	//Save the asset in the new destination
	destinationAssets[asset.Identifier] = asset

	updatedDestinationAssetsBytes, err := assetsToBytes(destinationAssets)

	if err != nil {
		return nil, errors.New("can't serialise assets, error:" + err.Error())
	}

	err = stub.PutState(transfer.OwnerTo, updatedDestinationAssetsBytes)

	//Remove from Transfer Temp

	delete(transferTempAssets, asset.Identifier)

	transferTempUpdatedAsBytes, err := assetsToBytes(transferTempAssets)

	if err != nil {
		return nil, errors.New("can't serialise updated temp assets, error:" + err.Error())
	}

	err = stub.PutState("TransferTemp", transferTempUpdatedAsBytes)

	return nil, nil
}

// Query state of the ledger
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "assets" {
		return t.listAssets(stub, args)
	}

	return nil, errors.New("Invalid query function name. Expecting \"assets\" or \"query\" ")

}

func (t *SimpleChaincode) listAssets(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("list assets - Incorrect number of arguments. Expecting owne to be specified")
	}

	owner := args[0]

	assetsAsBytes, err := stub.GetState(owner)

	if err != nil {
		return nil, errors.New("can't get assets for owner")
	}

	//deserialise
	assets, err := assetsFromBytes(assetsAsBytes)

	if err != nil {
		return nil, errors.New("can't deserialise assets, error:" + err.Error())
	}

	response := "{\"assets\":["

	for _, asset := range assets {

		myassetAsJson, _ := json.Marshal(asset)

		response += string(myassetAsJson)

		response += ","
	}

	response += "{}]}"

	jsonResp := response

	fmt.Printf("Query Response:%s\n", jsonResp)
	return []byte(jsonResp), nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
