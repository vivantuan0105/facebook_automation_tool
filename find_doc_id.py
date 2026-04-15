import requests
import re
import json

cookie = "datr=6M1JaPhd1g6T7V1vZuvpmKC0;xs=1%3A0FwF4gkw-dE3cw%3A2%3A1776045645%3A-1%3A-1;sb=Q07cafyMWDeguux1Y8IExiSZ;fr=1xDujY5gF4655mvfN.AWfs8C-gp2bpncuY2lxPJ3mDzYxyduoF6mR7BaUM96ko8B9pJKw.BobKSj..AAA.0.0.Bp3E5E.AWeSbM-PkJXpFO10lVymdjbJwdY;c_user=61571360758847;pas=61571360758847%3AhAIRYYu4Yh;locale=en_US"
uid = "61571360758847"

session = requests.Session()
session.headers.update({
    "cookie": cookie,
    "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
    "accept": "text/html,application/xhtml+xml"
})

resp = session.get(f"https://www.facebook.com/profile.php?id={uid}&sk=friends")
html = resp.text

with open("friends_page_dump.html", "w", encoding="utf-8") as f:
    f.write(html)

print("Saved HTML!")

# Look for friends in the string:
matches = re.finditer(r'\{"node":\{"__isProfile":"User","id":"(\d+)","name":"([^"]+)"', html)
friends = set()
for m in matches:
    friends.add(m.group(2))
print("Users found:", len(friends))

# Look for doc_id
# Facebook usually defines the graphql queries in RelayStream
matches = re.findall(r'"doc_id":"(\d+)"', html)
if matches:
    print("Found doc_ids in HTML:", set(matches))
