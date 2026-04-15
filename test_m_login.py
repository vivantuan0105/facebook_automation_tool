import requests
import re
import urllib.parse

session = requests.Session()
session.headers.update({
    "User-Agent": "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Mobile Safari/537.36",
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
    "Accept-Language": "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7",
    "Sec-Fetch-Dest": "document",
    "Sec-Fetch-Mode": "navigate",
    "Sec-Fetch-Site": "none",
    "Sec-Fetch-User": "?1",
    "Upgrade-Insecure-Requests": "1"
})

resp = session.get("https://m.facebook.com/")
html = resp.text

lsd = re.search(r'name="lsd"\s+value="([^"]+)"', html)
jazoest = re.search(r'name="jazoest"\s+value="([^"]+)"', html)
m_ts = re.search(r'name="m_ts"\s+value="([^"]+)"', html)
li = re.search(r'name="li"\s+value="([^"]+)"', html)
try_number = re.search(r'name="try_number"\s+value="([^"]+)"', html)
unrecognized_tries = re.search(r'name="unrecognized_tries"\s+value="([^"]+)"', html)
bi_xrwh = re.search(r'name="bi_xrwh"\s+value="([^"]+)"', html)

lsd_val = lsd.group(1) if lsd else ""
jazoest_val = jazoest.group(1) if jazoest else ""
m_ts_val = m_ts.group(1) if m_ts else ""
li_val = li.group(1) if li else ""
print(f"LSD: {lsd_val}, Jazoest: {jazoest_val}, m_ts: {m_ts_val}, li: {li_val}")

payload = {
    "lsd": lsd_val,
    "jazoest": jazoest_val,
    "m_ts": m_ts_val,
    "li": li_val,
    "try_number": try_number.group(1) if try_number else "0",
    "unrecognized_tries": unrecognized_tries.group(1) if unrecognized_tries else "0",
    "email": "61574299097448", # Use the UID!!
    "pass": "S2lj9NlT93",
    "login": "Đăng nhập"
}
if bi_xrwh: payload["bi_xrwh"] = bi_xrwh.group(1)

session.headers.update({
    "Content-Type": "application/x-www-form-urlencoded",
    "Origin": "https://m.facebook.com",
    "Referer": "https://m.facebook.com/",
    "Sec-Fetch-Site": "same-origin"
})

post_resp = session.post("https://m.facebook.com/login/device-based/regular/login/?refsrc=deprecated&lwv=100&refid=8", data=payload)
print("POST STATUS:", post_resp.status_code)
print("POST URL:", post_resp.url)
print("COOKIES:", session.cookies.get_dict())

if "c_user" in session.cookies.get_dict():
    print("SUCCESS! WE GOT C_USER!")
else:
    print("FAILED. WRITING HTML")
    with open("m_login_debug.html", "w", encoding="utf-8") as f:
        f.write(post_resp.text)
