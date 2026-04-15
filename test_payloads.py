import requests
import json

cookie = "datr=6M1JaPhd1g6T7V1vZuvpmKC0;xs=1%3A0FwF4gkw-dE3cw%3A2%3A1776045645%3A-1%3A-1;sb=Q07cafyMWDeguux1Y8IExiSZ;fr=1xDujY5gF4655mvfN.AWfs8C-gp2bpncuY2lxPJ3mDzYxyduoF6mR7BaUM96ko8B9pJKw.BobKSj..AAA.0.0.Bp3E5E.AWeSbM-PkJXpFO10lVymdjbJwdY;c_user=61571360758847;pas=61571360758847%3AhAIRYYu4Yh;locale=en_US"
uid = "61571360758847"
doc_id = "26339944415655338"
fb_dtsg = "NAfs_h4cxt2gwuREZoBJXMV2NOk0KbfePHw3mwV2S7fDFud2axMiiSA:1:1776045645"

session = requests.Session()
session.headers.update({
    "cookie": cookie,
    "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
})

URL = "https://www.facebook.com/api/graphql/"

def try_payload_str(variables_str):
    payload = {
        "av": uid,
        "__user": uid,
        "__a": "1",
        "fb_dtsg": fb_dtsg,
        "doc_id": doc_id,
        "variables": variables_str,
        "fb_api_req_friendly_name": "ProfileCometAppSectionFeedPaginationQuery"
    }
    resp = session.post(URL, data=payload)
    return resp.json()

print("Try GO original string:")
res1 = try_payload_str(f'{{"count":200,"cursor":null,"scale":1,"id":"{uid}"}}')
if "errors" in res1: print(res1["errors"][0]["message"])
else: print("SUCCESS!")

print("Try GO spaced string:")
res2 = try_payload_str(f'{{"count": 10, "cursor": null, "id": "{uid}", "scale": 1}}')
if "errors" in res2: print(res2["errors"][0]["message"])
else: print("SUCCESS!")
