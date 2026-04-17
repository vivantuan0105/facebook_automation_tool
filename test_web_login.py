import requests, re
session = requests.Session()
session.headers.update({
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36',
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
    'Accept-Language': 'en-US,en;q=0.5',
})
r = session.get('https://www.facebook.com/')
def get_val(name):
    m = re.search(f'name="{name}".*?value="([^"]+)"', r.text)
    if not m:
         m = re.search(f'"{name}",\[\],{{"token":"([^"]+)"}}', r.text)
    return m.group(1) if m else ''

lsd = get_val('lsd') or get_val('LSD')
jazoest = get_val('jazoest')

payload = {
    'lsd': lsd,
    'jazoest': jazoest,
    'email': '61568893734847',
    'pass': 'XZbRHwnl83',
    'login_source': 'comet_header_shortwave',
    'next': ''
}
resp = session.post('https://www.facebook.com/login/device-based/regular/login/', data=payload, allow_redirects=True)
if 'c_user' in session.cookies.get_dict():
    print('SUCCESS')
else:
    print('FAILED', resp.url)
    if 'approvals_code' in resp.text:
         print('2FA REQUIRED')
    elif 'checkpoint' in resp.url.lower():
         print('CHECKPOINT URL')
    elif 'incorrect' in resp.text.lower() or 'sai ' in resp.text.lower():
         print('INCORRECT PASSWORD')
    else:
         print('NO C_USER, PROBABLY BLOCKED BY FB SECURITY')
