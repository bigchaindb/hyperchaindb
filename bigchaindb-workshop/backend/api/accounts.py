import random
from json import JSONDecodeError

from tornado import web
from tornado.escape import json_decode
from tornado.gen import coroutine

from ..models.accounts import (
    Account,
    retrieve_accounts
)
from ..utils import (
    get_bigchain,
    APP_DB_NAME
)

bigchain = get_bigchain()


class AccountsHandler(web.RequestHandler):
    @coroutine
    def get(self):
        accounts = yield retrieve_accounts(bigchain, APP_DB_NAME)
        self.write({'result': accounts})

    @coroutine
    def post(self):
        name = 'account_{}'.format(str(random.randint(10, 1000)))
        try:
            data = json_decode(self.request.body)
            name = data.get('name')
        except JSONDecodeError:
            pass

        account = Account(bigchain=bigchain,
                          name=name,
                          db=APP_DB_NAME)
        self.write(account)
#
#
# @api_views.route('/accounts/<account_vk>/assets/')
# def get_assets_for_account(account_vk):
#     query = request.args.get('search')
#
#     result = {
#         'bigchain': assets.get_owned_assets(bigchain, vk=account_vk, query=query),
#         'backlog': assets.get_owned_assets(bigchain, vk=account_vk, query=query, table='backlog')
#     }
#     return flask.jsonify({'assets': result, 'account': account_vk})