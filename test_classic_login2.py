import requests
import re
import time

session = requests.Session()
session.headers.update({
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/124.0.0.0 Safari/537.36",
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
})

r = session.get("https://www.facebook.com/")
lsd = None
if 'name="lsd"' in r.text:
    lsd = r.text.split('name="lsd"')[1].split('value="')[1].split('"')[0]
elif '"LSD",[]' in r.text: 
    lsd = r.text.split('"LSD",[],{"token":"')[1].split('"')[0]

jazoest = None
if 'name="jazoest"' in r.text:
    jazoest = r.text.split('name="jazoest"')[1].split('value="')[1].split('"')[0]
elif '"jazoest":"' in r.text:
    jazoest = r.text.split('"jazoest":"')[1].split('"')[0]

print("LSD:", lsd, "Jazoest:", jazoest)

encpass = f"#PWD_BROWSER:5:{int(time.time())}:S2lj9NlT93"

payload = {
    "lsd": lsd,
    "jazoest": jazoest,
    "email": "61574299097448",
    "encpass": encpass,
    "login_source": "comet_header_shortwave",
    "next": ""
}

resp = session.post("https://www.facebook.com/login/device-based/regular/login/", data=payload, allow_redirects=False)

cookies = session.cookies.get_dict()
print("Cookies:", cookies)
if "c_user" in cookies:
    print("SUCCESS: c_user", cookies["c_user"])
elif "checkpoint" in resp.headers.get("Location", ""):
    print("CHECKPOINT REDIRECT!")
else:
    print("Location:", resp.headers.get("Location"))
    if "approvals_code" in resp.headers.get("Location", ""):
        print("2FA REQUIRED!")
