import cryptoconditions as cc
from bigchaindb.util import get_hash_data

from .driver import (
    b,
    get_tx_owned,
    post_account,
    post_tx,
    poll_tx_status_until_valid,
)
from .util import printd

# Create some accounts
quizmaster = post_account('quizmaster').json()
dimi = post_account('dimi').json()
mark = post_account('mark').json()

# The quizmaster receives a Question Token
tx = b.create_transaction(quizmaster['vk'], None, None, 'CREATE')
response = post_tx(tx)
tx_received = response.json()
print(response.status_code)
poll_tx_status_until_valid(tx_received['id'])

# The quizmaster creates a question
question = {'question': 'What is the answer to life the universe and everything?'}
question_tx = b.create_transaction(quizmaster['vk'], None, {'txid': tx_received['id'], 'cid': 0}, 'TRANSFER')
# Define a secret that will be hashed - fulfillments need to guess the secret
answer = b'42'
hashlock_tx_condition = cc.PreimageSha256Fulfillment(preimage=answer)

# The conditions list is empty, so we need to append a new condition
question_tx['transaction']['conditions'].append({
    'condition': {
        'uri': hashlock_tx_condition.condition.serialize_uri()
    },
    'cid': 0,
    'owners_after': None
})

# Conditions have been updated, so hash needs updating
question_tx['id'] = get_hash_data(question_tx)
question_tx_signed = b.sign_transaction(question_tx, quizmaster['sk'])

# POST the transaction to BigchainDB
b.is_valid_transaction(question_tx_signed)
response = post_tx(question_tx_signed)
tx_received = response.json()
print(response.status_code)

# Poll the transaction status
poll_tx_status_until_valid(question_tx['id'])
printd(get_tx_owned(dimi['vk']))
printd(get_tx_owned(mark['vk']))

# create an empty transaction for dimi
wrong_answer_tx = b.create_transaction(None, dimi['vk'], {'txid': question_tx['id'], 'cid': 0}, 'TRANSFER')

# dimi provides the wrong answer
wrong_answer_tx['transaction']['fulfillments'][0]['fulfillment'] = \
    cc.PreimageSha256Fulfillment(preimage=b'43').serialize_uri()

# POST the transaction to BigchainDB
response = post_tx(wrong_answer_tx)
print(response.status_code)  # 400 : wrong answer

# create an empty transaction for mark
right_answer_tx = b.create_transaction(None, mark['vk'], {'txid': question_tx['id'], 'cid': 0}, 'TRANSFER')

# mark provides the right answer
right_answer_tx['transaction']['fulfillments'][0]['fulfillment'] = \
    cc.PreimageSha256Fulfillment(preimage=b'42').serialize_uri()

# POST the transaction to BigchainDB
response = post_tx(right_answer_tx)
print(response.status_code)

# Poll the transaction status
poll_tx_status_until_valid(right_answer_tx['id'])
printd(get_tx_owned(dimi['vk']))
printd(get_tx_owned(mark['vk']))
