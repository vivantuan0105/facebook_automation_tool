import requests
import re
session = requests.Session()
session.headers.update({
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36',
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
    'Accept-Language': 'en-US,en;q=0.5',
})
r = session.get('https://mbasic.facebook.com/login.php')
lsd = re.search(r'name="lsd"[^>]*?value="([^"]+)"', r.text)
lsd = lsd.group(1) if lsd else ''
jazoest = re.search(r'name="jazoest"[^>]*?value="([^"]+)"', r.text)
jazoest = jazoest.group(1) if jazoest else ''
m_ts = re.search(r'name="m_ts"[^>]*?value="([^"]+)"', r.text)
m_ts = m_ts.group(1) if m_ts else ''
li = re.search(r'name="li"[^>]*?value="([^"]+)"', r.text)
li = li.group(1) if li else ''
action_url = re.search(r'<form[^>]+action="([^"]+)"', r.text)
action_url = action_url.group(1).replace('&amp;', '&') if action_url else 'https://mbasic.facebook.com/login/device-based/regular/login/'
if not action_url.startswith('http'): action_url = 'https://mbasic.facebook.com' + action_url

print('lsd:', lsd, 'jazoest:', jazoest, 'm_ts:', m_ts, 'li:', li)
payload = {'lsd': lsd, 'jazoest': jazoest, 'm_ts': m_ts, 'li': li, 'try_number': '0', 'unrecognized_tries': '0', 'email': '61574299097448', 'pass': 'S2lj9NlT93', 'login': 'Log In'}
resp = session.post(action_url, data=payload, allow_redirects=True)

if 'c_user' in session.cookies.get_dict():
    print('SUCCESS')
else:
    print('FAILED', resp.url)
    with open('mbasic_res.html', 'w', encoding='utf-8') as f: f.write(resp.text)
