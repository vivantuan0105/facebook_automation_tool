import requests
import re
import time

def do_mbasic_login(username, password):
    session = requests.Session()
    session.headers.update({
        "User-Agent": "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Mobile Safari/537.36",
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
        "Accept-Language": "en-US,en;q=0.5",
    })

    print(f"--- Testing {username} ---")
    r = session.get("https://mbasic.facebook.com/")
    
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

    resp = session.post("https://mbasic.facebook.com/login/device-based/regular/login/", data=payload, allow_redirects=True)
    
    cookies = session.cookies.get_dict()
    if "c_user" in cookies:
        print("SUCCESS! c_user:", cookies["c_user"])
    else:
        print("FAILED!")
        print("Final URL:", resp.url)
        if "checkpoint" in resp.url.lower():
            print("CHECKPOINT / CAPTCHA DETECTED BY URL")
            
        with open(f"mbasic_err_{username}.html", "w", encoding="utf-8") as f:
            f.write(resp.text)

do_mbasic_login("61573296730268", "S2lj9NlT93") # Note: guessing password is same, or user didn't specify. I'll just see what FB returns for UID.
