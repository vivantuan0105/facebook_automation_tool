import requests
import re

cookie = "datr=6M1JaPhd1g6T7V1vZuvpmKC0;xs=1%3A0FwF4gkw-dE3cw%3A2%3A1776045645%3A-1%3A-1;sb=Q07cafyMWDeguux1Y8IExiSZ;fr=1xDujY5gF4655mvfN.AWfs8C-gp2bpncuY2lxPJ3mDzYxyduoF6mR7BaUM96ko8B9pJKw.BobKSj..AAA.0.0.Bp3E5E.AWeSbM-PkJXpFO10lVymdjbJwdY;c_user=61571360758847;pas=61571360758847%3AhAIRYYu4Yh;locale=en_US"
uid = "61571360758847"

session = requests.Session()
session.headers.update({
    "cookie": cookie,
    "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
    "accept-language": "en-US,en;q=0.9"
})

resp = session.get(f"https://mbasic.facebook.com/profile.php?id={uid}&v=friends")
print("Response length:", len(resp.text))

# Basic check for friends list in mbasic
matches = re.finditer(r'<a href="/profile\.php\?id=(\d+)[^"]*">([^<]+)</a>', resp.text)
friends = {}
for m in matches:
    if m.group(1) != uid: # Not self
        friends[m.group(1)] = m.group(2)

print("Found friends basic format:", len(friends))

# Sometimes custom URL paths are used (e.g. /john.doe)
matches2 = re.finditer(r'<a href="/([^/"]+)\?(?:fref=pb|hc_location=friends_tab)[^"]*">([^<]+)</a>', resp.text)
for m in matches2:
    if "profile.php?" not in m.group(1):
        friends[m.group(1)] = m.group(2)

print("Total unique exact friends found:", len(friends))
for k, v in list(friends.items())[:5]:
    print(k, "=>", v)
