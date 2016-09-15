import requests

from .driver import (
    b,
    BDB_URL,
    post_account,
    get_account,
    get_accounts,
    post_tx,
    get_tx,
    get_tx_owned,
    poll_tx_status_until_valid,
)
from .util import printd

# Test if the API is running
response = requests.get(BDB_URL)
print('BDB serves at {}'.format(BDB_URL))
printd(response.json())

# Create an account
dimi = post_account('dimi').json()
printd(dimi)
# vk stands for verifying key and equals the public key
# sk stands for signing key and equals the private key

# List all accounts
printd(get_accounts().json())

# Select an account
printd(get_account('dimi'))
# Define a digital asset data payload
digital_asset_payload = {'msg': 'Hello BigchainDB!'}

# A create transaction uses the operation `CREATE` and has no inputs
tx = b.create_transaction(dimi['vk'], None, None, 'CREATE', payload=digital_asset_payload)

# POST the transaction to BigchainDB
response = post_tx(tx)
tx_received = response.json()
print(response.status_code)

# Poll the transaction status
poll_tx_status_until_valid(tx_received['id'])

# GET the transaction by ID
response = get_tx(tx_received['id'])
printd(response.json())

# GET the assets IDs owned by the test user
printd(get_tx_owned(dimi['vk']))
