from .driver import (
    b,
    get_account,
    get_accounts,
    get_tx_owned,
    post_account,
    post_tx,
    get_tx,
    poll_tx_status_until_valid,
)
from .util import printd

# Create an account
dimi = post_account('dimi').json()
printd(dimi)

# Create another account
mark = post_account('mark').json()

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
printd(get_tx_owned(dimi['vk']))
printd(get_tx_owned(mark['vk']))

tx_received = get_tx(tx_received['id']).json()

# create a transfer transaction
tx_transfer = b.create_transaction(dimi['vk'], mark['vk'], {'txid': tx_received['id'], 'cid': 0}, 'TRANSFER')
# sign the transaction
tx_transfer_signed = b.sign_transaction(tx_transfer, dimi['sk'])
print(b.is_valid_transaction(tx_transfer_signed))

# POST the transaction to BigchainDB
response = post_tx(tx_transfer_signed)
tx_received = response.json()
print(response.status_code)


# Poll the transaction status
poll_tx_status_until_valid(tx_received['id'])

printd(get_tx_owned(dimi['vk']))
printd(get_tx_owned(mark['vk']))
