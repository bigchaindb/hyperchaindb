
from bigchaindb import Bigchain
from bigchaindb import crypto

b = Bigchain()


# Define a digital asset data payload
digital_asset_payload = {'msg': 'Hello BigchainDB!'}

users = {}



def addUser(userName):

    if  False== users.__contains__(userName):
        user = type('user',(),{'name':userName, 'pub':'','priv':''})

        user.priv,user.pub = crypto.generate_key_pair()

        users[userName]=user

def newAsset(userName, asset):

    if False ==users.__contains__(userName):
        return 'Cant find that user buddy', False
    
    user = users[userName]

    # A create transaction uses the operation `CREATE` and has no inputs
    tx = b.create_transaction(b.me, user.pub, None, 'CREATE', payload=asset)

    # All transactions need to be signed by the user creating the transaction
    tx_signed = b.sign_transaction(tx, b.me_private)

    # Write the transaction to the bigchain.
    # The transaction will be stored in a backlog where it will be validated,
    # included in a block, and written to the bigchain 

    b.write_transaction(tx_signed)

    print('tx_signed[id]', tx_signed['id'])


    tx_retrieved=None

    while  tx_retrieved==None:
        tx_retrieved = b.get_transaction(tx_signed['id'])

    print('tx retrieved', tx_retrieved)


    return (tx_retrieved, True)
    

def transferAsset(fromUserName,toUserName,assetToTransfer):

    print('\nasset to Transfer',assetToTransfer)

    if  False==users.__contains__(fromUserName):
        return 'Cant find the frombUser Name ' + fromUserName, False

    if  False==users.__contains__(toUserName):
        return 'Cant find the to User Name ' + toUserName, False

    fromUser = users[fromUserName]
    toUser = users[toUserName]

    assetIds = b.get_owned_ids(fromUser.pub)

    print('assetIds \n***\n',assetIds)
    print('\n***')

    print(assetIds)

    retrievedAssetToTransfer= [assetId for assetId in assetIds if assetId['txid']==assetToTransfer['id']]

    print('\nretrievedAssetToTransfer',retrievedAssetToTransfer)

    if len(retrievedAssetToTransfer)==1:

        print('\ntransfering:',retrievedAssetToTransfer[0])

        tx_transfer = b.create_transaction(fromUser.pub, toUser.pub, retrievedAssetToTransfer[0], 'TRANSFER')

        # Sign the transaction  
        tx_transfer_signed = b.sign_transaction(tx_transfer, fromUser.priv)

        print('\ntx_transfer_signed:',tx_transfer_signed)

        # Write the transaction
        
        b.write_transaction(tx_transfer_signed)

        tx_transfer_retrieved=None

        while tx_transfer_retrieved==None:
            tx_transfer_retrieved = b.get_transaction(tx_transfer_signed['id'])

        print('tx transfer retrieved:',tx_transfer_retrieved)

        return (tx_transfer_retrieved, True)       

    else:
        return (None, False)

def getAssetIds(userName):

    if  False==users.__contains__(userName):
        return 'Cant find the to User Name ' + userName, False

    user = users[userName]
   
    assetIds = b.get_owned_ids(user.pub)

    return assetIds,True

def getAssets(assetIds):

    assets =[]

    print('\nasset Ids:', assetIds)

    for assetId in assetIds:

        print('\nasset Id:', assetId)

        asset_retrieved = b.get_transaction(assetId['txid'])

        assets.append(asset_retrieved)

    return (assets,True)

def testIt():

    addUser('aa')
    addUser('bb')

    tx,result1 = newAsset('aa',{'Name':'BartCoin2'})

    print('\nnew Asset, Result:', result1, 'tx:', tx)

    txTran, result2= transferAsset('aa','bb',tx)

    print('\ntransfer Asset, Result:', result2, 'txTran:', txTran)


    # print('\ntx',tx,'\nresult1',result1,'\ntxTran',txTran,'\nResult2',result2)

    # assets_for_bb, result = getAssets(getAssetIds('bb')[0])

    # for asset in assets_for_bb:
    #     print('\nAsset for bb:',asset)

                        
