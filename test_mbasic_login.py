import requests
import json
import re

session = requests.Session()
session.headers.update({
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
})

print("Fetching login page...")
r = session.get("https://m.facebook.com/login.php")
lsd = re.search(r'name="lsd".*?value="([^"]+)"', r.text)
if lsd: lsd = lsd.group(1)

jazoest = re.search(r'name="jazoest".*?value="([^"]+)"', r.text)
if jazoest: jazoest = jazoest.group(1)

m_ts = re.search(r'name="m_ts".*?value="([^"]+)"', r.text)
li = re.search(r'name="li".*?value="([^"]+)"', r.text)

print(f"Tokens: lsd={lsd}, jazoest={jazoest}")

payload = {
    "lsd": lsd,
    "jazoest": jazoest,
    "m_ts": m_ts.group(1) if m_ts else "",
    "li": li.group(1) if li else "",
    "try_number": "0",
    "unrecognized_tries": "0",
    "email": "61574299097448",
    "pass": "S2lj9NlT93",
    "login": "Đăng nhập"
}

print("Submitting login...")
resp = session.post("https://m.facebook.com/login/device-based/regular/login/", data=payload)
cookies = session.cookies.get_dict()
print("Cookies after login:", cookies)
if "c_user" in cookies:
    print("SUCCESS: c_user", cookies["c_user"])
elif "checkpoint" in resp.url or "v" in cookies:
    print("CHECKPOINT/CAPTCHA!")
else:
    print("FAILED. Current URL:", resp.url)
