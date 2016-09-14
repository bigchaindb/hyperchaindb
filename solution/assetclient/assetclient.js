

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
                    "args": ["a", "888", "b", "999"]
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

function invoke() {

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
                    "name": chaincodeGuid
                },
                "ctorMsg": {
                    "function": "invoke",
                    "args": ["a", "b", "10"]
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

};

function deposit(account, amount) {

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
                    "name": "gg"
                },
                "ctorMsg": {
                    "function": "deposit",
                    "args": [account,amount]
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


function query(account){
    var chaincodeName = chaincodeGuid;
    queryinternal(account, chaincodeName);
}

function queryinternal(account, chaincodeName) {

    var request = require('request');
    //Lets configure and request
    request({
        url: blockchainServer + '/chaincode', //URL to hit
        //qs: {from: 'blog example', time: +new Date()}, //Query string data
        method: 'POST',
        //Lets post the following key/values as form
        json: {
            "jsonrpc": "2.0",
            "method": "query",
            "params": {
                "type": 1,
                "chaincodeID": {
                    "name": chaincodeName
                },
                "ctorMsg": {
                    "function": "query",
                    "args": [account]
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
};


function sell(){

    var chaincodeName = chaincodeGuid;

    var assetBartCoin = {Name: "BartCoin",
		Identifier: "A1",
		Type:       "Crypto",
		Owner:      "aa",
		Price:      3.14};

    var assetBartCoinStringed = JSON.stringify(assetBartCoin)

    sellinternal(assetBartCoinStringed,chaincodeName )

}

function sellinternal(asset, chaincodeName){

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
                    "function": "sell",
                    "args": [asset]
                }
            },
            "id": 6
        }
    }, function (error, response, body) {
        if (error) {
            console.log(error);
        } else {
            console.log(response.statusCode, body);
        }
    });
}

function createasset(){

    var chaincodeName = chaincodeGuid;

    var assetBartCoin = {Name: "BartCoin",
		Identifier: "A1",
		Type:       "Crypto",
		Owner:      "aa",
		Price:      3.14};

    var assetBartCoinStringed = JSON.stringify(assetBartCoin)

    sellinternal(assetBartCoinStringed,chaincodeName )

}

function createassetInternal(asset, chaincodeName){

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
            "id": 6
        }
    }, function (error, response, body) {
        if (error) {
            console.log(error);
        } else {
            console.log(response.statusCode, body);
        }
    });
}


function sales(){

    var chaincodeName = chaincodeGuid;
    salesinternal(chaincodeName)
}

function salesinternal(chaincodeName) {

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
                    "function": "sales",
                    "args": ["latest"]
                }
            },
            "id": 7
        }
    }, function (error, response, body) {
        if (error) {
            console.log(error);
        } else {
            console.log(response.statusCode, body);
        }
    });
};

