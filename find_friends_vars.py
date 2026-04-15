import requests
import re
import json

cookie = "datr=6M1JaPhd1g6T7V1vZuvpmKC0;xs=1%3A0FwF4gkw-dE3cw%3A2%3A1776045645%3A-1%3A-1;sb=Q07cafyMWDeguux1Y8IExiSZ;fr=1xDujY5gF4655mvfN.AWfs8C-gp2bpncuY2lxPJ3mDzYxyduoF6mR7BaUM96ko8B9pJKw.BobKSj..AAA.0.0.Bp3E5E.AWeSbM-PkJXpFO10lVymdjbJwdY;c_user=61571360758847;pas=61571360758847%3AhAIRYYu4Yh;locale=en_US"
uid = "61571360758847"

session = requests.Session()
session.headers.update({
    "cookie": cookie,
    "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
    "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"
})

resp = session.get(f"https://www.facebook.com/profile.php?id={uid}&sk=friends")
html = resp.text

print("Searching for collection token / variables in HTML:")
# Friends collection id typically starts with app_collection:UID:2356318349
matches = re.finditer(r'(app_collection:[0-9]+:2356318349:[0-9]+)', html)
tokens = set(m.group(1) for m in matches)
print("Collection tokens:", tokens)

import base64
for t in tokens:
    print("Base64 encoded:", base64.b64encode(t.encode()).decode())

# What about exact variables JSON embedded in the HTML?
matches2 = re.finditer(r'"variables":(\{[^}]*"count":\s*\d+[^}]*\})', html)
variables = set(m.group(1) for m in matches2)
print("Found Variables combinations:", len(variables))
for v in list(variables)[:5]:
    if "cursor" in v or "scale" in v:
        print(v[:200])

