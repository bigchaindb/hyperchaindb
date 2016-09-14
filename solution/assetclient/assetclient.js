

var localServer = "http://localhost:7050";

var blockchainServer = localServer;

var chaincodeGuid = "";

function deploy() {

    var request = require('request');
    //Lets configure and request
    request({
        url: blockchainServer + '/chaincode', //URL to hit
        method: 'POST',
        json: {
            "jsonrpc": "2.0",
            "method": "deploy",
            "params": {
                "type": 1,
                "chaincodeID": {
                    "name":"gg",
                    "path": "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_assetmarket"
                },
                "ctorMsg": {
                    "function": "init",
                    "args": ["aa", "bb", "cc"]
                }
            },
            "id": 1
        }
    }, function (error, response, body) {
        if (error) {
            console.log(error);
        } else {
            console.log(response.statusCode, body);
            chaincodeGuid = body.result.message;

            console.log("chaincodeGuid:" + chaincodeGuid)
        }

    });
}

function createasset(){

    var chaincodeName = chaincodeGuid;

    var assetBartCoin = {"Name":"BartCoin",
		"Identifier": "A1",
		"Type":       "Crypto",
		"Owner":      "aa",
		"Price":      3.14,
        "BigChainId": ""};

    var assetBartCoinStringed = JSON.stringify(assetBartCoin)

    createAssetInternal(assetBartCoinStringed,chaincodeName )

}

function createAssetInternal(asset, chaincodeName){

    var request = require('request');

    request({
        url: blockchainServer + '/chaincode', //URL to hit
        method: 'POST',
        json: {
            "jsonrpc": "2.0",
            "method": "invoke",
            "params": {
                "type": 1,
                "chaincodeID": {
                    "name": chaincodeName
                },
                "ctorMsg": {
                    "function": "create",
                    "args": [asset]
                }
            },
            "id": 2
        }
    }, function (error, response, body) {
        if (error) {
            console.log(error);
        } else {
            console.log(response.statusCode, body);
        }
    });
}

function transfer(fromUser, toUser, assetIdentifier){

    var chaincodeName = chaincodeGuid;

    transferInternal(fromUser, toUser, assetIdentifier,chaincodeName )
}

function transferInternal(fromUser, toUser, assetIdentifier, chaincodeName){

    var request = require('request');

    request({
        url: blockchainServer + '/chaincode', //URL to hit
        method: 'POST',
        json: {
            "jsonrpc": "2.0",
            "method": "invoke",
            "params": {
                "type": 1,
                "chaincodeID": {
                    "name": chaincodeName
                },
                "ctorMsg": {
                    "function": "transfer",
                    "args": [fromUser,toUser,assetIdentifier]
                }
            },
            "id": 3
        }
    }, function (error, response, body) {
        if (error) {
            console.log(error);
        } else {
            console.log(response.statusCode, body);
        }
    });
}

function listAssets(userName){

    var chaincodeName = chaincodeGuid;
    listAssetsInternal(userName,chaincodeName)
}

function listAssetsInternal(userName, chaincodeName) {

    var request = require('request');
    //Lets configure and request
    request({
        url: blockchainServer + '/chaincode', //URL to hit
        method: 'POST',
        json: {
            "jsonrpc": "2.0",
            "method": "query",
            "params": {
                "type": 1,
                "chaincodeID": {
                    "name": chaincodeName
                },
                "ctorMsg": {
                    "function": "assets",
                    "args": [userName]
                }
            },
            "id": 4
        }
    }, function (error, response, body) {
        if (error) {
            console.log(error);
        } else {
            console.log(response.statusCode, body);
        }
    });
};



