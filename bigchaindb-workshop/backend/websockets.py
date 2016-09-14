from tornado import websocket
from tornado.gen import coroutine

import rethinkdb as r

from .server import get_bigchain

bigchain = get_bigchain()
clients = []


# from http://blog.hiphipjorge.com/django-and-realtime-using-django-with-tornado-and-rethinkdb/
r.set_loop_type('tornado')


@coroutine
def print_changes(db_table):
    conn = yield bigchain.conn
    feed = yield r.table(db_table).changes().run(conn)
    while (yield feed.fetch_next()):
        change = yield feed.next()
        block = get_block_from_change(change, db_table)
        for client in clients:
            for tx in block:
                # TODO: use REQL for filtering
                if tx_contains_vk(tx['transaction'], client.username):
                    msg = {'change': change,
                           'client': client.username}
                    client.write_message(msg)
                    break


def get_block_from_change(change, db_table):
    block = []
    if db_table in ['backlog', 'bigchain'] and (change['old_val'] or change['new_val']):
        block_data = change['old_val'] if change['old_val'] else change['new_val']
        if db_table == 'bigchain':
            block = block_data['block']['transactions']
        else:
            block.append(block_data)
    return block


def tx_contains_vk(tx, vk):
    for condition in tx['conditions']:
        if vk in condition['new_owners']:
            return True
    for fullfillment in tx['fulfillments']:
        if vk in fullfillment['current_owners']:
            return True


class ChangeFeedWebSocket(websocket.WebSocketHandler):
    username = None

    def check_origin(self, origin):
        return True

    def open(self, username):
        if self not in clients:
            self.username = username
            clients.append(self)
        print('ws: open (Pool: {} connections)'.format(len(clients)))

    def on_message(self, message):
        pass

    def on_close(self):
        for i, client in enumerate(clients):
            if client is self:
                clients.remove(self)
                print('ws: close (Pool: {} connections)'.format(len(clients)))
                return