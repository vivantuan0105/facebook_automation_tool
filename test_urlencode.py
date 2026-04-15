import urllib.parse
payload = {
    "variables": '{"count": 10, "cursor": null, "id": "61571360758847", "scale": 1}'
}
print(urllib.parse.urlencode(payload))
