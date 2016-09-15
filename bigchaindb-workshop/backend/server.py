import functools
import os
import logging

from tornado import web, ioloop

from .api.accounts import AccountsHandler
from .websockets import ChangeFeedWebSocket, print_changes


logger = logging.getLogger('tornado')

logger.info('Initializing tornado server')

app = web.Application([
    (r'/accounts', AccountsHandler),
    (r'/accounts/(.*)/changes', ChangeFeedWebSocket)
])


def run_tornado_server():
    tornado_port = int(os.environ.get('TORNADO_PORT', 8888))
    tornado_address = os.environ.get('TORNADO_HOST', '127.0.0.1')
    app.listen(tornado_port, address=tornado_address)

    ioloop.IOLoop.current().add_callback(functools.partial(print_changes, 'backlog'))
    ioloop.IOLoop.current().add_callback(functools.partial(print_changes, 'bigchain'))

    logger.info('Running on http://{}:{}'.format(tornado_address, tornado_port))
    ioloop.IOLoop.instance().start()

if __name__ == '__main__':
    run_tornado_server()
