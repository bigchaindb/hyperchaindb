import binascii, hashlib
import cryptoconditions as cc

SPACER = '*'*40

################################
# Find a good secret and hash it
################################
secret = '42'.encode()  # use bytes only
puzzle = binascii.hexlify(hashlib.sha256(secret).digest())
print(SPACER)
print('The puzzle is {}, guess the secret!'.format(puzzle))
print(SPACER)


###############################################################
# Create a Hashlock fulfillment and broadcast the condition URI
###############################################################
sha256fulfillment = cc.PreimageSha256Fulfillment(preimage=secret)
sha256condition = sha256fulfillment.condition

print(SPACER)
print('Fulfillment URI: {}'.format(sha256fulfillment.serialize_uri()))
print('hashed payload: {}'.format(binascii.hexlify(sha256condition.hash)))
print('Condition URI: {}'.format(sha256fulfillment.condition_uri))
print(SPACER)

##########################
# Validate the fulfillment
##########################
assert sha256fulfillment.validate() and sha256fulfillment.condition_uri == sha256condition.serialize_uri()
