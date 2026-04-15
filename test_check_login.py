import requests
import re
import time

session = requests.Session()
session.headers.update({
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/124.0.0.0 Safari/537.36",
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
})

r = session.get("https://www.facebook.com/")
print("Initial URL:", r.url)

lsd = None
m = re.search(r'name="lsd"[^>]*?value="([^"]+)"', r.text)
if not m: m = re.search(r'"LSD",\[\],\{"token":"([^"]+)"}', r.text)
if m: lsd = m.group(1)

jazoest = None
m = re.search(r'name="jazoest"[^>]*?value="([^"]+)"', r.text)
if m: jazoest = m.group(1)

encpass = f"#PWD_BROWSER:5:{int(time.time())}:S2lj9NlT93"

payload = {
    "lsd": lsd,
    "jazoest": jazoest,
    "email": "61574299097448",
    "encpass": encpass,
    "login_source": "comet_header_shortwave",
    "next": ""
}

print(f"LSD: {lsd}, Jazoest: {jazoest}")

resp = session.post("https://www.facebook.com/login/device-based/regular/login/", data=payload, allow_redirects=True)

print("Final URL:", resp.url)
print("Cookies:", session.cookies.get_dict())
with open("test_checkpoint.html", "w", encoding="utf-8") as f:
    f.write(resp.text)

if 'checkpoint' in resp.url.lower():
    print("REDIRECTED TO CHECKPOINT URL")
elif 'checkpoint' in resp.text.lower():
    print("CHECKPOINT FOUND IN HTML TEXT BUT NOT URL")
