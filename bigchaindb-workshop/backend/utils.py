import os

import bigchaindb
from tornado.gen import coroutine

API_BASE_HOST = os.environ.get('DOCKER_MACHINE_IP') or 'localhost'
API_BASE_PORT = int(os.environ.get('API_BASE_PORT', '8888'))
APP_DB_NAME = 'bigchaindb_workshop'
NUM_ACCOUNTS = 3

try:
    CONFIG_FILE = os.environ['BIGCHAINDB_CONFIG']
except KeyError:
    CONFIG_FILE = '.bigchaindb_workshop'


def get_bigchain(conf=CONFIG_FILE):
    if os.path.isfile(conf):
        bigchaindb.config_utils.autoconfigure(filename=conf, force=True)
    return bigchaindb.Bigchain()


@coroutine
def feed_to_list(feed):
    result_list = []
    while (yield feed.fetch_next()):
        result = yield feed.next()
        result_list.append(result)
    return result_list
