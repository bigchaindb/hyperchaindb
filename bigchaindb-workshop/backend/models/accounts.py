import rethinkdb as r

import bigchaindb.crypto
from tornado.gen import coroutine


class Account:
    def __init__(self, bigchain, name, db):
        self.bigchain = bigchain
        self.db = db
        self.name = name
        self.sk, self.vk = bigchaindb.crypto.generate_key_pair()
        self.save()

    def save(self):
        try:
            r.db_create(self.db).run(self.bigchain.conn)
        except r.ReqlOpFailedError:
            pass

        try:
            r.db(self.db).table_create('accounts').run(self.bigchain.conn)
        except r.ReqlOpFailedError:
            pass

        user_exists = list(r.db(self.db)
                           .table('accounts')
                           .filter(lambda user: (user['name'] == self.name))
                           .run(self.bigchain.conn))
        if not len(user_exists):
            r.db(self.db)\
                .table('accounts')\
                .insert(self.as_dict(), durability='hard')\
                .run(self.bigchain.conn)
        else:
            user_persistent = user_exists[0]
            self.vk = user_persistent['vk']
            self.sk = user_persistent['sk']

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
    accounts = []
    while (yield accounts_feed.fetch_next()):
        account = yield accounts_feed.next()
        accounts.append(account)
    return accounts
