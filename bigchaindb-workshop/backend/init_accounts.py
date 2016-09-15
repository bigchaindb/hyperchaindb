import logging

import bigchaindb.config_utils

from .models.accounts import Account
from .utils import (
    get_bigchain,
    APP_DB_NAME,
    NUM_ACCOUNTS
)

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

bigchain = get_bigchain()
logging.info('INIT: bigchain initialized with database: {}'.format(bigchaindb.config['database']['name']))


def main():
    accounts = []
    for i in range(NUM_ACCOUNTS):
        account = Account(bigchain=bigchain,
                          name='account_{}'.format(i),
                          db=APP_DB_NAME)
        accounts.append(account)
    logging.info('INIT: {} accounts initialized for app: {}'.format(len(accounts), APP_DB_NAME))


if __name__ == '__main__':
    main()
