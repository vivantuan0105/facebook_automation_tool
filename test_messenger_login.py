import requests
import re
import time

session = requests.Session()
session.headers.update({
    "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
    "accept-language": "en-US,en;q=0.9,vi;q=0.8",
    "sec-ch-ua": '"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"',
    "sec-ch-ua-mobile": "?0",
    "sec-ch-ua-platform": '"Windows"',
    "sec-fetch-dest": "document",
    "sec-fetch-mode": "navigate",
    "sec-fetch-site": "none",
    "sec-fetch-user": "?1",
    "upgrade-insecure-requests": "1",
    "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
})

r = session.get("https://www.messenger.com/")
lsd = [line.split('value="')[1].split('"')[0] for line in r.text.split('\n') if 'name="lsd"' in line]
lsd = lsd[0] if lsd else "fake"
jazoest = [line.split('value="')[1].split('"')[0] for line in r.text.split('\n') if 'name="jazoest"' in line]
jazoest = jazoest[0] if jazoest else "fake"

print("LSD:", lsd, "Jazoest:", jazoest)

payload = {
    "lsd": lsd,
    "jazoest": jazoest,
    "email": "61574299097448",
    "encpass": f"#PWD_BROWSER:5:{int(time.time())}:S2lj9NlT93",
    "login_source": "comet_header_shortwave",
    "next": ""
}

session.headers.update({
    "content-type": "application/x-www-form-urlencoded",
    "origin": "https://www.messenger.com",
    "referer": "https://www.messenger.com/",
    "sec-fetch-site": "same-origin"
})

resp = session.post("https://www.messenger.com/login/password/", data=payload, allow_redirects=False)

print("Status:", resp.status_code)
print("Cookies:", session.cookies.get_dict())
if "c_user" not in session.cookies.get_dict():
    with open("debug_test_messenger.html", "w", encoding="utf-8") as f:
        f.write(resp.text)
    print("Failed. Wrote debug_test_messenger.html")
