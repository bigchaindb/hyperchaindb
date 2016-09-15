import os
import json
import requests

from time import sleep

from ...backend.utils import get_bigchain
from .util import printd

b = get_bigchain()

BDB_PORT = 9984
BDB_URL = os.environ.get('BIGCHAINDB_BASE_URL',
                         'http://localhost:{}'.format(BDB_PORT))
BDB_API_URL = os.environ.get('BIGCHAINDB_API_ENDPOINT',
                             '{}/api/v1'.format(BDB_URL))
BDB_API_TX = '{}/transactions'.format(BDB_API_URL)
BDB_API_ACCOUNTS = os.environ.get('BDB_API_ACCOUNTS',
                                  'http://localhost:{}/accounts'.format(8888))


def post_tx(tx):
    return requests.post(BDB_API_TX, data=json.dumps(tx))


def get_tx(tx_id):
    return requests.get('{}/{}'.format(BDB_API_TX, tx_id))


def get_tx_status(tx_id):
    return requests.get('{}/{}/status'.format(BDB_API_TX, tx_id))


def get_tx_owned(vk):
    return b.get_owned_ids(vk)


def poll_tx_status_until_valid(tx_id):
    tx_status = None
    # Poll the transaction status
    while not tx_status == 'valid':
        response = get_tx_status(tx_id)
        tx_status = response.json()['status']
        printd(response.json())
        sleep(.5)
    return tx_status


def post_account(name):
    return requests.post(BDB_API_ACCOUNTS, data=json.dumps({'name': name}))


def get_accounts():
    return requests.get(BDB_API_ACCOUNTS)


def get_account(name):
    response = get_accounts()
    try:
        return [account for account in response.json()['result'] if account['name'] == name][0]
    except IndexError:
        return None
