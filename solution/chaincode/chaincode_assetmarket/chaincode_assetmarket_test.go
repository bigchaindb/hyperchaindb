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

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args []string) {
	_, err := stub.MockInit("1", "init", args)
	if err != nil {
		fmt.Println("Init failed", err)
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkStateKey(t *testing.T, stub *shim.MockStub, name string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}

}

func checkQuery(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes, err := stub.MockQuery("query", []string{name})
	if err != nil {
		fmt.Println("Query", name, "failed", err)
		t.FailNow()
	}
	if bytes == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("Query value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkQueryWithFunc(t *testing.T, stub *shim.MockStub, function string, name string, value string) {
	bytes, err := stub.MockQuery(function, []string{name})
	if err != nil {
		fmt.Println("Query", name, "failed", err)
		t.FailNow()
	}
	if bytes == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("Query value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkQueryWithFuncAndResult(t *testing.T, stub *shim.MockStub, function string, name string) []byte {
	bytes, err := stub.MockQuery(function, []string{name})
	if err != nil {
		fmt.Println("Query", name, "failed", err)
		t.FailNow()
	}
	if bytes == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	return bytes
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args []string) {
	_, err := stub.MockInvoke("1", "invoke", args)
	if err != nil {
		fmt.Println("Invoke", args, "failed", err)
		t.FailNow()
	}
}

func checkInvokeWithFunc(t *testing.T, stub *shim.MockStub, function string, args []string) {
	_, err := stub.MockInvoke("1", function, args)
	if err != nil {
		fmt.Println("Invoke", args, "failed", err)
		t.FailNow()
	}
}

func TestAssetMarket_Init(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("am01", scc)

	checkInit(t, stub, []string{"aa", "bb", "cc"})

	//emptyAssetMap := make(map[string]*Asset)

	checkStateKey(t, stub, "aa")
	checkStateKey(t, stub, "bb")
	checkStateKey(t, stub, "cc")

}

func TestExample02_Serialise(t *testing.T) {

	myasset := &Asset{Name: "BartCoin",
		Identifier: "A1",
		Type:       "Crypto",
		Owner:      "a",
		Price:      3.14}

	myasset_marshalled, _ := json.Marshal(myasset)

	fmt.Println(myasset_marshalled)

	fmt.Println(string(myasset_marshalled))

	myasset_marshalled_as_string := string(myasset_marshalled)

	myasset_copy := Asset{}

	json.Unmarshal([]byte(myasset_marshalled_as_string), &myasset_copy)

	if myasset_copy.Name != "BartCoin" {
		fmt.Println("cant reserialise")
		t.FailNow()
	}

	fmt.Print("Reserialised:", "Name: ", myasset_copy.Name, " Identifier:", myasset_copy.Identifier, "\n")

}

func TestExample02_SerialiseMap(t *testing.T) {

	sales := make(map[string]*Asset)

	A1 := &Asset{Name: "BartCoin",
		Identifier: "A1",
		Type:       "Crypto",
		Owner:      "a",
		Price:      3.14}

	B1 := &Asset{Name: "LizardCoin",
		Identifier: "B1",
		Type:       "Crypto",
		Owner:      "a",
		Price:      42.0}

	sales[A1.Identifier] = A1
	sales[B1.Identifier] = B1

	for _, p := range sales {
		fmt.Println(p.Name, "likes cheese.")
	}

	var bin_buf bytes.Buffer

	binary.Write(&bin_buf, binary.BigEndian, sales)

	gob.Register(Asset{})

	b := new(bytes.Buffer)

	e := gob.NewEncoder(b)

	// Encoding the map
	err := e.Encode(sales)
	if err != nil {
		panic(err)
	}

	var decodedMap map[string]*Asset
	d := gob.NewDecoder(b)

	// Decoding the serialized data
	err = d.Decode(&decodedMap)
	if err != nil {
		panic(err)
	}

}

func TestExample02_Serious(t *testing.T) {

	sales := make(map[string]*Asset)

	A1 := &Asset{Name: "BartCoin",
		Identifier: "A1",
		Type:       "Crypto",
		Owner:      "aa",
		Price:      3.14}

	B1 := &Asset{Name: "LizardCoin",
		Identifier: "B1",
		Type:       "Crypto",
		Owner:      "bb",
		Price:      42.0}

	sales[A1.Identifier] = A1
	sales[B1.Identifier] = B1

	gob.Register(Asset{})

	b := new(bytes.Buffer)

	e := gob.NewEncoder(b)

	// Encoding the map
	err := e.Encode(sales)
	if err != nil {
		panic(err)
	}

	var decodedMap map[string]*Asset
	d := gob.NewDecoder(b)

	// Decoding the serialized data
	err = d.Decode(&decodedMap)
	if err != nil {
		panic(err)
	}

	bytes, err := assetsToBytes(sales)

	if err != nil {
		panic(err)
	}

	newSales, err := assetsFromBytes(bytes)

	if err != nil {
		panic(err)
	}

	if newSales[A1.Identifier].Price != 3.14 {
		panic(err)
	}
	if newSales[B1.Identifier].Price != 42.00 {
		panic(err)
	}

}

func Test_create(t *testing.T) {

	fmt.Print("\nrunning create test\n")

	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkInit(t, stub, []string{"aa", "bb", "cc"})

	A1 := &Asset{Name: "BartCoin",
		Identifier: "A1",
		Type:       "Crypto",
		Owner:      "aa",
		Price:      3.14}

	myassetMarshalled, _ := json.Marshal(A1)

	fmt.Println(myassetMarshalled)

	fmt.Println(string(myassetMarshalled))

	myasset_marshalled_as_string := string(myassetMarshalled)

	parameters := []string{myasset_marshalled_as_string}

	checkInvokeWithFunc(t, stub, "create", parameters)
	/*
		expected := "{\"assets\":[{\"name\":\"BartCoin\",\"identifier\":\"A1\",\"type\":\"Crypto\",\"owner\":\"aa\",\"amount\":3.14},{}]}"

		checkQueryWithFunc(t, stub, "aa", "", expected)
	*/
}

func Test_createAndConfirm(t *testing.T) {

	fmt.Print("\nrunning CreateAndConfirm test\n")

	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkInit(t, stub, []string{"aa", "bb", "cc"})

	A1 := &Asset{Name: "BartCoin",
		Identifier: "A1",
		Type:       "Crypto",
		Owner:      "aa",
		Price:      3.14}

	myassetMarshalled, _ := json.Marshal(A1)

	fmt.Println(myassetMarshalled)

	fmt.Println(string(myassetMarshalled))

	myasset_marshalled_as_string := string(myassetMarshalled)

	parameters := []string{myasset_marshalled_as_string}

	checkInvokeWithFunc(t, stub, "create", parameters)

	ConfirmedAsset := &Asset{Name: "BartCoin",
		Identifier: "A1",
		Type:       "Crypto",
		Owner:      "aa",
		Price:      3.14,
		BigChainId: "BIG"}

	myConfirmedAasetMarshalled, _ := json.Marshal(ConfirmedAsset)

	fmt.Println(string(myConfirmedAasetMarshalled))

	myconfirmedasset_marshalled_as_string := string(myConfirmedAasetMarshalled)

	confirmedAssetParameters := []string{myconfirmedasset_marshalled_as_string}

	checkInvokeWithFunc(t, stub, "createconfirmed", confirmedAssetParameters)

	/*
		expected := "{\"assets\":[{\"name\":\"BartCoin\",\"identifier\":\"A1\",\"type\":\"Crypto\",\"owner\":\"aa\",\"amount\":3.14},{}]}"

		checkQueryWithFunc(t, stub, "aa", "", expected)
	*/
}

func Test_transfer(t *testing.T) {

	fmt.Print("\nrunning CreateAndConfirm test\n")

	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkInit(t, stub, []string{"aa", "bb", "cc"})

	A1 := &Asset{Name: "BartCoin",
		Identifier: "A1",
		Type:       "Crypto",
		Owner:      "aa",
		Price:      3.14}

	myassetMarshalled, _ := json.Marshal(A1)

	fmt.Println(myassetMarshalled)

	fmt.Println(string(myassetMarshalled))

	myasset_marshalled_as_string := string(myassetMarshalled)

	parameters := []string{myasset_marshalled_as_string}

	checkInvokeWithFunc(t, stub, "create", parameters)

	ConfirmedAsset := &Asset{Name: "BartCoin",
		Identifier: "A1",
		Type:       "Crypto",
		Owner:      "aa",
		Price:      3.14,
		BigChainId: "BIG"}

	myConfirmedAasetMarshalled, _ := json.Marshal(ConfirmedAsset)

	fmt.Println(string(myConfirmedAasetMarshalled))

	myconfirmedasset_marshalled_as_string := string(myConfirmedAasetMarshalled)

	confirmedAssetParameters := []string{myconfirmedasset_marshalled_as_string}

	checkInvokeWithFunc(t, stub, "createconfirmed", confirmedAssetParameters)

	//Transfer parameters
	transferParameters := []string{"aa", "bb", "A1"}

	checkInvokeWithFunc(t, stub, "transfer", transferParameters)

	transferConfirmation := &AssetTransfer{OwnerFrom: "aa",
		OwnerTo:    "bb",
		Identifier: "A1",
		BigChainId: "BIG"}

	transferConfirmationMarshalled, _ := json.Marshal(transferConfirmation)

	fmt.Println(string(transferConfirmationMarshalled))

	transferMarshalledAsString := string(transferConfirmationMarshalled)

	transferConfirmedParameters := []string{transferMarshalledAsString}

	checkInvokeWithFunc(t, stub, "transferconfirmed", transferConfirmedParameters)

	var queryBytes []byte
	queryBytes = checkQueryWithFuncAndResult(t, stub, "assets", "aa")

	fmt.Printf("Result for aa =: %s", string(queryBytes))

	queryBytes = checkQueryWithFuncAndResult(t, stub, "assets", "bb")

	fmt.Printf("Result for bb =: %s", string(queryBytes))

	queryBytes = checkQueryWithFuncAndResult(t, stub, "assets", "TransferTemp")

	fmt.Printf("Result for TransferTemp =: %s", string(queryBytes))
}
