import requests
import re
import time
import urllib.parse

def do_mbasic_login(username, password):
    session = requests.Session()
    # Desktop User Agent!
    session.headers.update({
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
        "Accept-Language": "en-US,en;q=0.5",
    })

    print(f"--- Testing {username} ---")
    r = session.get("https://mbasic.facebook.com/login.php")
    
    lsd = ""
    m = re.search(r'name="lsd"[^>]*?value="([^"]+)"', r.text)
    if m: lsd = m.group(1)

    jazoest = ""
    m = re.search(r'name="jazoest"[^>]*?value="([^"]+)"', r.text)
    if m: jazoest = m.group(1)
    
    m_ts = ""
    m = re.search(r'name="m_ts"[^>]*?value="([^"]+)"', r.text)
    if m: m_ts = m.group(1)
    
    li = ""
    m = re.search(r'name="li"[^>]*?value="([^"]+)"', r.text)
    if m: li = m.group(1)

    print(f"Tokens - lsd: {lsd}, jazoest: {jazoest}, m_ts: {m_ts}, li: {li}")

    payload = {
        "lsd": lsd,
        "jazoest": jazoest,
        "m_ts": m_ts,
        "li": li,
        "try_number": "0",
        "unrecognized_tries": "0",
        "email": username,
        "pass": password,
        "login": "Log In"
    }
    
    # We must use exactly the action URL from the form
    action_url = ""
    m = re.search(r'<form[^>]+action="([^"]+)"', r.text)
    if m:
        action_url = m.group(1).replace("&amp;", "&")
        if not action_url.startswith("http"):
            action_url = "https://mbasic.facebook.com" + action_url
    else:
        action_url = "https://mbasic.facebook.com/login/device-based/regular/login/?refsrc=deprecated&lwv=100"

    print("Action URL:", action_url)

    resp = session.post(action_url, data=payload, allow_redirects=True)
    
    cookies = session.cookies.get_dict()
    if "c_user" in cookies:
        print("SUCCESS! c_user:", cookies["c_user"])
    else:
        print("FAILED!")
        print("Final URL:", resp.url)
        with open("mbasic_desktop_debug.html", "w", encoding="utf-8") as f:
            f.write(resp.text)

do_mbasic_login("61573296730268", "S2lj9NlT93")
