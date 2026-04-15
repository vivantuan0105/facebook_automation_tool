import requests
import hashlib
import uuid
import time
import urllib.parse

def md5(text):
    return hashlib.md5(text.encode('utf-8')).hexdigest()

app_keys = [
    # FB for Android (Old Katana)
    ("882a8490361da98702bf97a021ddc14d", "62f8ce9f74b12f84c123cc23437a4a32"),
    # Messenger for iOS
    ("256002347743983", "374e60f8b9bb6b8cbb30f78030438895"),
    # FB for iPhone
    ("6628568379", "c1e620fa708a1d5696fb991c1bde5662"),
    # FB for Android v2?
    ("350685531728", "907106da44b207f23a63319be83be219"),
    # Instagram for Android
    ("1217981644879628", "2d3e099ee7007f9c2d1b7a2d48c8b6b1")
]

username = "61574299097448"
password = "S2lj9NlT93"

for api_key, secret in app_keys:
    print(f"\n--- Testing App: {api_key} ---")
    data = {
        'api_key': api_key,
        'credentials_type': 'password',
        'email': username,
        'format': 'JSON',
        'generate_machine_id': '1',
        'generate_session_cookies': '1',
        'locale': 'en_US',
        'method': 'auth.login',
        'password': password,
        'return_multiple_errors': 'true',
        'device_id': str(uuid.uuid4()),
        'machine_id': str(uuid.uuid4())
    }
    
    sig_data = ''
    for key in sorted(data.keys()):
        sig_data += f"{key}={data[key]}"
    sig_data += secret
    data['sig'] = md5(sig_data)
    
    headers = {
        "User-Agent": "[FBAN/FB4A;FBAV/417.0.0.33.65;FBBV/480086274;FBDM/{density=3.0,width=1080,height=2132};FBLC/en_US;FBRV/481977799;FBCR/Viettel;FBMF/Xiaomi;FBBD/xiaomi;FBPN/com.facebook.katana;FBDV/M2007J20CG;FBSV/11;FBOP/1;FBCA/arm64-v8a:;]",
        "Content-Type": "application/x-www-form-urlencoded"
    }

    try:
        resp = requests.post('https://b-api.facebook.com/method/auth.login', data=data, headers=headers)
        j = resp.json()
        print("Status:", resp.status_code)
        if 'session_cookies' in j:
            print("SUCCESS!")
        elif 'error_msg' in j:
            print("ERROR MSG:", j['error_msg'])
        else:
            print("Unknown response:", j)
    except Exception as e:
        print("Exception:", e)
