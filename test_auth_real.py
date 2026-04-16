import requests
import json
with open("Data/61572157625165/Info.txt", "r", encoding="utf-8") as f:
    text = f.read()
acc = json.loads(text)
cookie = acc["cookie"]
fb_dtsg = 'NAfuhSt1u-1AEvfYMRaLxxxCj1GiIqObinmG-YQR-xbqDJE6WV-Ulcw:7:1776270267'

session = requests.Session()
session.headers.update({
    "Cookie": cookie,
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
})

url = "https://www.facebook.com/api/graphql/"
data = {
    "av": "61572157625165",
    "__user": "61572157625165",
    "__a": "1",
    "fb_dtsg": fb_dtsg,
    "doc_id": "26505601199092382",
    "variables": '{"UFI2CommentsProvider_commentsKey":"ProfileCometTimelineRoute","beforeTime":null,"count":10,"cursor":null,"dialtone_active":false,"feedLocation":"TIMELINE","feedbackSource":0,"focusCommentID":null,"has_react_native_v2":false,"id":"61572157625165","omitPinnedPost":true,"privacySelectorRenderLocation":"COMET_STREAM","renderLocation":"timeline","scale":1,"useDefaultActor":false}'
}

r = session.post(url, data=data)
print("Base Response:")
print(r.text[:500])
