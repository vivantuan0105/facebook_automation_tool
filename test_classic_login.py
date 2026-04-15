import requests
import re
import time

session = requests.Session()
session.headers.update({
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/124.0.0.0 Safari/537.36",
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
    "Accept-Language": "en-US,en;q=0.5",
    "Sec-Fetch-Dest": "document",
    "Sec-Fetch-Mode": "navigate",
    "Sec-Fetch-Site": "none",
    "Sec-Fetch-User": "?1"
})

r = session.get("https://www.facebook.com/")
lsd = "fake"
match = re.search(r'name="lsd"[\s\S]*?value="([^"]+)"', r.text)
if match: lsd = match.group(1)

jazoest = "fake"
match = re.search(r'name="jazoest"[\s\S]*?value="([^"]+)"', r.text)
if match: jazoest = match.group(1)

print("LSD:", lsd)
print("Jazoest:", jazoest)

encpass = f"#PWD_BROWSER:5:{int(time.time())}:S2lj9NlT93"

payload = {
    "lsd": lsd,
    "jazoest": jazoest,
    "email": "61574299097448",
    "encpass": encpass,
    "login_source": "comet_header_shortwave",
    "next": ""
}

print("Logging in...")
resp = session.post("https://www.facebook.com/login/device-based/regular/login/", data=payload, allow_redirects=False)

print("Status:", resp.status_code)
print("Headers:", resp.headers)
print("Cookies:", session.cookies.get_dict())
