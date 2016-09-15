import binascii
import cryptoconditions as cc

SPACER = '*'*40

###############################
# Parse a condition from a URI
###############################

example_condition_uri = 'cc:0:3:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU:0'
parsed_condition = cc.Condition.from_uri(example_condition_uri)

# check if it's actually a condition
assert isinstance(parsed_condition, cc.Condition)

# print out all the internals of the condition
print(SPACER)
print('URI: {}'.format(parsed_condition.serialize_uri()))
print('class: {}'.format(parsed_condition.__class__.__name__))
print('type ID: {}'.format(parsed_condition.type_id))
print('type: {}'.format(cc.TypeRegistry.get_class_from_type_id(parsed_condition.type_id).__name__))
print('bitmask: {}'.format(parsed_condition.bitmask))
print('maximum fulfillment length: {}'.format(parsed_condition.max_fulfillment_length))
print('hashed payload: {}'.format(binascii.hexlify(parsed_condition.hash)))
print('dict representation: {}'.format(parsed_condition.to_dict()))
print(SPACER)

#################################
# Parse a fulfillment from a URI
#################################
example_fulfillment_uri = 'cf:0:'
parsed_fulfillment = cc.Fulfillment.from_uri(example_fulfillment_uri)

# check if it's actually a condition
assert isinstance(parsed_fulfillment, cc.Fulfillment)

# print out all the internals of the condition
print(SPACER)
print('URI: {}'.format(parsed_fulfillment.serialize_uri()))
print('class: {}'.format(parsed_fulfillment.__class__.__name__))
print('type ID: {}'.format(parsed_fulfillment.type_id))
print('type: {}'.format(cc.TypeRegistry.get_class_from_type_id(parsed_fulfillment.type_id).__name__))
print('bitmask: {}'.format(parsed_fulfillment.bitmask))
print('payload: {}'.format(binascii.hexlify(parsed_fulfillment.preimage)))
print('dict representation: {}'.format(parsed_fulfillment.to_dict()))


# Retrieve the condition of the fulfillment
print('condition URI: {}'.format(parsed_fulfillment.condition_uri))
# prints 'cc:0:3:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU:0'

# Validate a fulfillment
assert parsed_fulfillment.validate()
assert parsed_fulfillment.condition_uri == example_condition_uri

print(SPACER)
