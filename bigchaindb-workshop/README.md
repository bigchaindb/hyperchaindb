# BigchainDB Workshop

[![Join the chat at https://gitter.im/bigchaindb/bigchaindb](https://badges.gitter.im/bigchaindb/bigchaindb.svg)](https://gitter.im/bigchaindb/bigchaindb?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

This repo contains workshop examples and tutorials for BigchainDB.

### Dependencies

The examples can be run via Docker

 - OS dependencies: see [setup BigchainDB & RethinkDB](https://bigchaindb.readthedocs.io/en/latest/installing-server.html#install-and-run-rethinkdb-server)
 - python>=3.4
 - node>=5.3 using [nvm](https://github.com/creationix/nvm#installation) (**recommended**), or [manually](https://nodejs.org/en/download/)
 - [npm>=3.3](https://docs.npmjs.com/getting-started/installing-node) (should be installed with node)

## Quick start


### Docker

To run via Docker, set up your docker environment as necessary and:

```bash
$ make
```

**Note**: If using docker-machine, you'll have to run `make` with your docker-machine ip:

```bash
$ DOCKER_MACHINE_IP=$(docker-machine ip) make
```

The API will be available at <http://localhost:48888> (replace ``localhost`` with your
docker-machine ip as necessary).

### Locally

If you'd like to run the backend locally (preferably in a virtualenv), you can do so using:

```bash
# Install everything in a virtualenv
$ virtualenv venv -p python3
$ source venv/bin/activate
$ pip install -e .[dev]

# Initialize everything
$ bigchaindb -yc .bigchaindb_workshop configure
$ bigchaindb -c .bigchaindb_workshop init

# Start everything
$ bigchaindb -c .bigchaindb_workshop start
# in an other console run the API server
$ python -m backend.server

# Drop everything
$ bigchaindb -c .bigchaindb_workshop drop
```

The API will be available at <http://localhost:8888>.

## Exercises

### Cryptoconditions

Cryptoconditions is a specification for communicating between different ledgers.
One can think of it as authenticated event handlers or the TCP/IP version for ledgers.

See the [README](exercises/cc/README.md) for details and exercises.

## API endpoints

### Accounts

#### Retrieve account list

##### Request
```http
GET /accounts/ HTTP/1.1
Host: localhost:<PORT>
Content-Type: application/json
```

##### Response
```http
HTTP/1.1 200 OK
[   
    {
        "id": "<uuid>"
        "name": "<string>",
        "sk": "<base58>"
        "vk": "<base58>"
    }, ...
]
```

#### Create account

##### Request
```http
POST /accounts/ HTTP/1.1
Host: localhost:<PORT>
Content-Type: application/json
Body:
    {
        "name": "<string: optional>"
    }
```

##### Response
```http
HTTP/1.1 200 OK
{
    "name": "<string>",
    "sk": "<base58>"
    "vk": "<base58>"
}
```

## Acknowledgements:

Special thanks to the BigchainDB/ascribe.io team for their insights and code contributions:

@r-marques, @vrde, @ttmc, @rhsimplex, @SohKai, @sbellem, @TimDaub
