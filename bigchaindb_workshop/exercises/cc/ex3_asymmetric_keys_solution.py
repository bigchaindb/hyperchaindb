import base58

import cryptoconditions as cc

SPACER = '*'*40

#############################
# Generate a random key pair
#############################
sk_base58, vk_base58 = cc.crypto.ed25519_generate_key_pair()

# create Key objects
sk = cc.crypto.Ed25519SigningKey(vk_base58)
vk = sk.get_verifying_key()
print(SPACER)
print('Signing key')
print('Bytes: {}'.format(sk.to_bytes()))
print('Hex: {}'.format(sk.to_ascii(encoding='hex')))
print('Base58: {}'.format(sk.to_ascii(encoding='base58')))
print('Base64: {}'.format(sk.to_ascii(encoding='base64')))
print(SPACER)

print(SPACER)
print('Verifying key')
print('Bytes: {}'.format(vk.to_bytes()))
print('Hex: {}'.format(vk.to_ascii(encoding='hex')))
print('Base58: {}'.format(vk.to_ascii(encoding='base58')))
print('Base64: {}'.format(vk.to_ascii(encoding='base64')))
print(SPACER)

####################################
# Create an ED25519-SHA256 Condition
####################################
ed25519_fulfillment = cc.Ed25519Fulfillment(public_key=vk)
ed25519_condition_uri = ed25519_fulfillment.condition_uri
print(SPACER)
print('Dict representation: {}'.format(ed25519_fulfillment.to_dict()))
print('Condition URI: {}'.format(ed25519_condition_uri))
print('Hashed payload: {}'.format(base58.b58encode(ed25519_fulfillment.condition.hash)))
print('Verifying key: {}'.format(vk.to_ascii().decode()))
print('Is valid fulfillment? {}'.format(ed25519_fulfillment.validate()))
try:
    ed25519_fulfillment.serialize_uri()
except TypeError:
    print('Fulfillment URI threw an exception')
print(SPACER)

#####################################
# Fulfill an ED25519-SHA256 condition
#####################################
message = 'Hello World! I am a message to sign!'
ed25519_fulfillment.sign(message, sk)
print(SPACER)
print('Fulfillment URI: {}'.format(ed25519_fulfillment.serialize_uri()))
print('Condition URI: {}'.format(ed25519_fulfillment.condition_uri))
print('Is valid fulfillment? {}'.format(ed25519_fulfillment.validate()))
print('Is valid fulfillment? {}'.format(ed25519_fulfillment.validate('Not the message to sign')))
print('Is valid fulfillment? {}'.format(ed25519_fulfillment.validate(message)))
print(SPACER)

assert ed25519_fulfillment.validate(message) and ed25519_fulfillment.condition_uri == ed25519_condition_uri

###########################################
# Fulfill an ED25519-SHA256 Fulfillment URI
###########################################
example_fulfillment_uri = \
    'cf:4:7Bcrk61eVjv0kyxw4SRQNMNUZ-8u_U1k6_gZaDRn4r-2IpH62UMvjymL' \
    'nEpIldvik_b_2hpo2t8Mze9fR6DHISpf6jzal6P0wD6p8uisHOyGpR1FISer26CdG28zHAcK'
parsed_ed25519_fulfillment = cc.Ed25519Fulfillment.from_uri(example_fulfillment_uri)

print(SPACER)
print('Fulfillment URI: {}'.format(parsed_ed25519_fulfillment.serialize_uri()))
print('Condition URI: {}'.format(parsed_ed25519_fulfillment.condition_uri))
print('Is valid fulfillment? {}'.format(ed25519_fulfillment.validate(message)))
print(SPACER)
