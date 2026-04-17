import requests, re
session = requests.Session()
session.headers.update({
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36',
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
    'Accept-Language': 'en-US,en;q=0.5',
})
r = session.get('https://mbasic.facebook.com/login.php')
def get_val(name):
    m = re.search(f'name="{name}"[^>]*?value="([^"]+)"', r.text)
    return m.group(1) if m else ''

action_url = re.search(r'<form[^>]+action="([^"]+)"', r.text)
action_url = action_url.group(1).replace('&amp;', '&') if action_url else 'https://mbasic.facebook.com/login/device-based/regular/login/'
if not action_url.startswith('http'): action_url = 'https://mbasic.facebook.com' + action_url

payload = {
    'lsd': get_val('lsd'),
    'jazoest': get_val('jazoest'),
    'm_ts': get_val('m_ts'),
    'li': get_val('li'),
    'try_number': '0',
    'unrecognized_tries': '0',
    'email': '61568893734847',
    'pass': 'XZbRHwnl83',
    'login': 'Log In'
}
resp = session.post(action_url, data=payload, allow_redirects=True)
if 'c_user' in session.cookies.get_dict():
    print('SUCCESS', session.cookies.get_dict())
else:
    print('FAILED', resp.url)
    if 'approvals_code' in resp.text:
         print('2FA REQUIRED')
    elif 'checkpoint' in resp.url.lower():
         print('CHECKPOINT URL')
    elif 'incorrect' in resp.text.lower() or 'sai' in resp.text.lower():
         print('INCORRECT PASSWORD / SAI MẬT KHẨU')
    else:
         print('UNKNOWN ERROR RESPONSE')
