import requests
from time import sleep

from bigchaindb import crypto

from bigchaindb_workshop.exercises.bdb.util import (
    b,
    get_account,
    get_accounts,
    post_account,
    post_tx,
    get_tx,
    poll_tx_status_until_valid,
    printd
)

# Create an account
dimi = post_account('dimi').json()
printd(dimi)

# List all accounts
printd(get_accounts().json())

# Select an account
printd(get_account('dimi'))

# Define a digital asset data payload
digital_asset_payload = {'msg': 'Hello BigchainDB!'}

# A create transaction uses the operation `CREATE` and has no inputs
tx = b.create_transaction(b.me, dimi['vk'], None, 'CREATE', payload=digital_asset_payload)

# POST the transaction to BigchainDB
response = post_tx(tx)
tx_received = response.json()
print(response.status_code)

# Poll the transaction status
poll_tx_status_until_valid(tx_received['id'])

tx_received = get_tx(tx_received['id']).json()

# Create an account
mark = post_account('mark').json()

# create a transfer transaction
tx_transfer = b.create_transaction(dimi['vk'], mark['vk'], {'txid': tx_received['id'], 'cid': 0}, 'TRANSFER')
# sign the transaction
tx_transfer_signed = b.sign_transaction(tx_transfer, dimi['sk'])
print(b.validate_transaction(tx_transfer_signed))
# POST the transaction to BigchainDB
response = post_tx(tx_transfer_signed)
tx_received = response.json()
print(response.status_code)

# Poll the transaction status
poll_tx_status_until_valid(tx_received['id'])