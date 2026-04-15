import json

variables_dict = {"count": 10, "cursor": None, "id": "61571360758847", "scale": 1}
variables_str = '{"count": 10, "cursor": null, "id": "61571360758847", "scale": 1}'

print("json.dumps:")
print(repr(json.dumps(variables_dict)))
print("my string:")
print(repr(variables_str))
print("Are they exactly equal?", json.dumps(variables_dict) == variables_str)
