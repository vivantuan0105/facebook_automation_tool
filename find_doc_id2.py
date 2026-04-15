import requests
import re
import json

session = requests.Session()
r = session.get("https://www.facebook.com/")
js_files = re.findall(r'<script.*?src=["\']([^"\']+\.js)[^>]+>', r.text)

print(f"Found {len(js_files)} JS files. Searching for Login mutation ID...")

login_doc_id = None
for js_url in js_files:
    if not js_url.startswith('http'):
        continue
    try:
        resp = session.get(js_url, timeout=5)
        # Search for CometSessionCreateMutation or similar
        m = re.search(r'doc_id\s*:\s*["\'](\d{15,})["\'].*?Login', resp.text, re.IGNORECASE)
        if m:
            print("Found doc_id near 'Login':", m.group(1))
        
        m = re.search(r'path\s*:\s*\["caa_login_web"\].*?doc_id\s*:\s*["\'](\d{15,})["\']', resp.text)
        if m:
            print("Found caa_login_web doc_id:", m.group(1))
            
        m = re.search(r'CometSessionCreateMutation.*?doc_id\s*:\s*["\'](\d{15,})["\']', resp.text)
        if m:
            print("Found CometSessionCreateMutation:", m.group(1))
            
        # generic search for standard mutations
        if '"Login"' in resp.text and 'doc_id' in resp.text:
            pass
            
    except Exception as e:
        pass
        
print("Search complete.")
