import requests
from time import sleep

from bigchaindb import crypto

from .util import (
    b,
    BDB_URL,
    post_tx,
    get_tx,
    poll_tx_status_until_valid,
    printd
)

# Test if the API is running
response = requests.get(BDB_URL)
print('BDB serves at {}'.format(BDB_URL))
printd(response.json())

# Create a test user
sk_test1, vk_test1 = crypto.generate_key_pair()

# Define a digital asset data payload
digital_asset_payload = {'msg': 'Hello BigchainDB!'}

# A create transaction uses the operation `CREATE` and has no inputs
tx = b.create_transaction(b.me, vk_test1, None, 'CREATE', payload=digital_asset_payload)

# Sign the transaction
tx_signed = b.sign_transaction(tx, b.me_private)

# POST the transaction to BigchainDB
response = post_tx(tx_signed)
tx_received = response.json()
print(response.status_code)

# Poll the transaction status
poll_tx_status_until_valid(tx_received['id'])

# GET the transaction by ID
response = get_tx(tx_received['id'])
printd(response.json())

