import requests
import json
with open("Data/61572157625165/Info.txt", "r", encoding="utf-8") as f:
    acc = json.loads(f.read())
cookie = acc["cookie"]
fb_dtsg = 'NAfuhSt1u-1AEvfYMRaLxxxCj1GiIqObinmG-YQR-xbqDJE6WV-Ulcw:7:1776270267'

session = requests.Session()
session.headers.update({
    "Cookie": cookie,
    "Content-Type": "application/x-www-form-urlencoded",
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
    "X-FB-Friendly-Name": "ProfileCometTimelineFeedRefetchQuery"
})

url = "https://www.facebook.com/api/graphql/"
data = "__a=1&__user=61572157625165&av=61572157625165&doc_id=26505601199092382&fb_api_req_friendly_name=ProfileCometTimelineFeedRefetchQuery&fb_dtsg=NAfuhSt1u-1AEvfYMRaLxxxCj1GiIqObinmG-YQR-xbqDJE6WV-Ulcw%3A7%3A1776270267&lsd=teX-LT13Zed-Tk4rpaPnRF&variables=%7B%22UFI2CommentsProvider_commentsKey%22%3A%22ProfileCometTimelineRoute%22%2C%22afterTime%22%3Anull%2C%22beforeTime%22%3Anull%2C%22count%22%3A10%2C%22cursor%22%3Anull%2C%22dialtone_active%22%3Afalse%2C%22feedLocation%22%3A%22TIMELINE%22%2C%22feedbackSource%22%3A0%2C%22focusCommentID%22%3Anull%2C%22has_react_native_v2%22%3Afalse%2C%22id%22%3A%2261572157625165%22%2C%22omitPinnedPost%22%3Atrue%2C%22postedBy%22%3Anull%2C%22privacySelectorRenderLocation%22%3A%22COMET_STREAM%22%2C%22renderLocation%22%3A%22timeline%22%2C%22scale%22%3A1%2C%22useDefaultActor%22%3Afalse%7D"

r = session.post(url, data=data)
print("Base Response:", r.text[:300])
