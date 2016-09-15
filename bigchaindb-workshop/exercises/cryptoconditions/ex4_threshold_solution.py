import copy
import json

import cryptoconditions as cc

SPACER = '*'*40

# Create some keys and a message to sign
sk1, vk1 = cc.crypto.ed25519_generate_key_pair()
sk2, vk2 = cc.crypto.ed25519_generate_key_pair()
message = 'Hello World! I am a message to sign!'

##############################################################
# Create a Threshold fulfillment and
# add subfulfillments and subconditions as an object or by URI
##############################################################
# Create a 1-of-4 threshold condition (OR gate)
threshold_fulfillment = cc.ThresholdSha256Fulfillment(threshold=1)
# Add a hashlock fulfillment
threshold_fulfillment.add_subfulfillment(cc.PreimageSha256Fulfillment(b'much secret'))
# Add a signature condition
threshold_fulfillment.add_subcondition(cc.Ed25519Fulfillment(public_key=vk1).condition)
# Add a fulfillment URI
threshold_fulfillment.add_subfulfillment_uri('cf:0:')
# Add a condition URI
threshold_fulfillment.add_subcondition_uri('cc:0:3:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU:0')
print(SPACER)
print('Condition URI: {}'.format(threshold_fulfillment.condition_uri))
print('Threshold: {}'.format(threshold_fulfillment.threshold))
print('Is valid fulfillment? {}'.format(threshold_fulfillment.validate(message)))
if threshold_fulfillment.validate():
    print('Fulfillment URI {}'.format(threshold_fulfillment.serialize_uri()))
print('Dict representation:')
print(json.dumps(threshold_fulfillment.to_dict(), sort_keys=True, indent=4, separators=(',', ':')))
print(SPACER)

#######################################################################
# Increase the threshold until the fulfillment doesn't validate anymore
#######################################################################
print(SPACER)
threshold_fulfillment.threshold = 2
print('Threshold: {}'.format(threshold_fulfillment.threshold))
print('Is valid fulfillment? {}'.format(threshold_fulfillment.validate(message)))
threshold_fulfillment.threshold = 3
print('Threshold: {}'.format(threshold_fulfillment.threshold))
print('Is valid fulfillment? {}'.format(threshold_fulfillment.validate(message)))
print(SPACER)

####################################################
# Create a nested Threshold fulfillment and validate
####################################################
print(SPACER)
nested_threshold_fulfillment = copy.deepcopy(threshold_fulfillment)
nested_threshold_fulfillment.threshold = 2
print('Is valid nested fulfillment? {}'.format(nested_threshold_fulfillment.validate(message)))
threshold_fulfillment.add_subfulfillment(nested_threshold_fulfillment)
print('Is valid fulfillment? {}'.format(nested_threshold_fulfillment.validate(message)))
print(json.dumps(threshold_fulfillment.to_dict(), sort_keys=True, indent=4, separators=(',', ':')))
print(SPACER)
