import requests
import json
import re

cookie = "datr=6M1JaPhd1g6T7V1vZuvpmKC0;xs=1%3A0FwF4gkw-dE3cw%3A2%3A1776045645%3A-1%3A-1;sb=Q07cafyMWDeguux1Y8IExiSZ;fr=1xDujY5gF4655mvfN.AWfs8C-gp2bpncuY2lxPJ3mDzYxyduoF6mR7BaUM96ko8B9pJKw.BobKSj..AAA.0.0.Bp3E5E.AWeSbM-PkJXpFO10lVymdjbJwdY;c_user=61571360758847;pas=61571360758847%3AhAIRYYu4Yh;locale=en_US"
uid = "61571360758847"
doc_id = "26339944415655330"

session = requests.Session()
session.headers.update({
    "cookie": cookie,
    "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
})

print("Fetching profile to get fb_dtsg...")
resp = session.get(f"https://www.facebook.com/profile.php?id={uid}")
html = resp.text

# Extract fb_dtsg
match = re.search(r'"DTSGInitialData",\[\],\{"token":"([^"]+)"', html)
if not match:
    print("Could not find fb_dtsg in HTML.")
    exit(1)
fb_dtsg = match.group(1)
print("fb_dtsg:", fb_dtsg)

URL = "https://www.facebook.com/api/graphql/"
variables = {
    "count": 20,
    "id": uid,
    "scale": 1
}

payload = {
    "av": uid,
    "__user": uid,
    "__a": "1",
    "fb_dtsg": fb_dtsg,
    "variables": json.dumps(variables),
    "fb_api_req_friendly_name": "ProfileCometAppSectionFeedPaginationQuery"
}

print("Executing GraphQL request...")
post_resp = session.post(URL, data=payload)
try:
    data = post_resp.json()
    with open("python_debug_fb.json", "w", encoding="utf-8") as f:
        json.dump(data, f, indent=2, ensure_ascii=False)
    print("Saved response to python_debug_fb.json")
except Exception as e:
    print("Failed to parse JSON:", post_resp.text[:500])
