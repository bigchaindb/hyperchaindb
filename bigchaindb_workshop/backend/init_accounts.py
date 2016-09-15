import logging

import bigchaindb.config_utils

from .models.accounts import (
    Account,
    store_account
)
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
    for i in range(NUM_ACCOUNTS):
        account = Account('account_{}'.format(i))
        store_account(account, bigchain, db=APP_DB_NAME)
        logging.info('INIT: account {} added to app: {}'.format(account.name, APP_DB_NAME))

if __name__ == '__main__':
    main()
