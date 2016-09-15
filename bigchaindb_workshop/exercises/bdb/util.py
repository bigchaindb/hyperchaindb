import json


def printd(a_dict):
    print(json.dumps(a_dict, sort_keys=True, indent=4, separators=(',', ':')))
