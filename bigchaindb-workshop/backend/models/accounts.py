import rethinkdb as r

import bigchaindb.crypto
from tornado.gen import coroutine

from ..utils import feed_to_list


class Account:
    def __init__(self, name):
        self.name = name
        self.sk, self.vk = bigchaindb.crypto.generate_key_pair()

    def as_dict(self):
        return {
            'name': self.name,
            'sk': self.sk,
            'vk': self.vk
        }


@coroutine
def retrieve_accounts(bigchain, db):
    conn = yield bigchain.conn
    accounts_feed = yield r.db(db).table('accounts').run(conn)
    accounts = yield feed_to_list(accounts_feed)
    return accounts


def store_account(account, bigchain, db):
    conn = bigchain.conn
    try:
        r.db_create(db).run(conn)
    except r.ReqlOpFailedError:
        pass

    try:
        r.db(db).table_create('accounts').run(conn)
    except r.ReqlOpFailedError:
        pass

    user_exists = \
        list(r.db(db)
             .table('accounts')
             .filter(lambda user: (user['name'] == account.name))
             .run(conn))

    if not len(user_exists):
        r.db(db) \
            .table('accounts') \
            .insert(account.as_dict(), durability='hard') \
            .run(conn)


@coroutine
def store_account_async(account, bigchain, db):
    conn = yield bigchain.conn
    try:
        yield r.db_create(db).run(conn)
    except r.ReqlOpFailedError:
        pass

    try:
        yield r.db(db).table_create('accounts').run(conn)
    except r.ReqlOpFailedError:
        pass

    user_exists_feed = \
        yield r.db(db) \
            .table('accounts') \
            .filter(lambda user: (user['name'] == account.name)) \
            .run(conn)
    user_exists = yield feed_to_list(user_exists_feed)

    if not len(user_exists):
        yield r.db(db) \
            .table('accounts') \
            .insert(account.as_dict(), durability='hard') \
            .run(conn)
        return account.as_dict()
    else:
        return user_exists[0]
