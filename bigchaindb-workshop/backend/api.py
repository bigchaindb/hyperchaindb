from tornado import web

import rethinkdb as r
import bigchaindb


class GetTestHandler(web.RequestHandler):
    def get(self):
        response = {'id': 'test'}
        self.write(response)
