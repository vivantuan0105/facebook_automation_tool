import requests
import re
import urllib.parse
from bs4 import BeautifulSoup

session = requests.Session()
session.headers.update({
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
    "Accept-Language": "en-US,en;q=0.9",
})

resp = session.get("https://mbasic.facebook.com/login/")
html = resp.text

soup = BeautifulSoup(html, 'html.parser')
form = soup.find('form', action=re.compile('login'))

if not form:
    print("No login form found")
    exit(1)

action_url = form.get('action')
if not action_url.startswith('http'):
    action_url = "https://mbasic.facebook.com" + action_url

print("Action URL:", action_url)

payload = {}
for i in form.find_all('input'):
    name = i.get('name')
    if name:
        payload[name] = i.get('value', '')

payload['email'] = "61574299097448"
payload['pass'] = "S2lj9NlT93"

print("Payload:", payload.keys())

session.headers.update({
    "Content-Type": "application/x-www-form-urlencoded",
    "Origin": "https://mbasic.facebook.com",
    "Referer": "https://mbasic.facebook.com/login/",
})

post_resp = session.post(action_url, data=payload)
print("STATUS:", post_resp.status_code)
print("COOKIES:", session.cookies.get_dict())

if "c_user" in session.cookies.get_dict():
    print("SUCCESS!")
else:
    print("FAILED")
    with open("mbasic_debug.html", "w", encoding="utf-8") as f:
        f.write(post_resp.text)
