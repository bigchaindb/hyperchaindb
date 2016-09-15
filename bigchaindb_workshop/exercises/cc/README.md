# Cryptoconditions exercises

[![Join the chat at https://gitter.im/bigchaindb/bigchaindb](https://badges.gitter.im/bigchaindb/bigchaindb.svg)](https://gitter.im/bigchaindb/bigchaindb?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

This repo contains examples and solutions for Interledger Cryptoconditions.

# How to install

We support **Python>=3.4**.

Install from pypi:

```
$ pip install cryptocondtions
```

# Crypto Conditions

This spec is a python port from the [**Interledger Protocol (ILP)**]
(https://interledger.org/five-bells-condition/spec.html)

## Motivation

We would like a way to describe a signed message such that multiple actors in a
distributed system can all verify the same signed message and agree on whether
it matches the description.

This provides a useful primitive for distributed, event-based systems since we
can describe events (represented by signed messages) and therefore define
generic authenticated event handlers.

## Implementations

Cryptoconditions are available in [python](https://github.com/bigchaindb/cryptoconditions), 
[javascript](https://github.com/interledger/five-bells-condition) and Java (coming soon).

This workshop is mainly focussed on Python, but should be easily ported.
Note that the serialized URI's are platform and language independent.

## Exercises:

Each python file contains either an exercise or a fulfletched solution.
The scope covered is:

1. General usage of conditions and fulfillments.
2. Hash locks
3. Asymmetric keys
4. Threshold conditions