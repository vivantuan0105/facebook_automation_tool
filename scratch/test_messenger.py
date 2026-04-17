import requests
import hashlib

def login_messenger(email, password):
    api_key = '256002347743983'
    app_secret = '374e60f8b9bb6b8cbb30f78030438895'
    
    data = {
        'api_key': api_key,
        'credentials_type': 'password',
        'email': email,
        'format': 'JSON',
        'generate_machine_id': '1',
        'generate_session_cookies': '1',
        'locale': 'en_US',
        'method': 'auth.login',
        'password': password,
        'return_multiple_errors': 'true',
    }
    
    sig = hashlib.md5((''.join(f'{k}={data[k]}' for k in sorted(data.keys())) + app_secret).encode('utf-8')).hexdigest()
    data['sig'] = sig
    
    resp = requests.post('https://b-api.facebook.com/method/auth.login', data=data).json()
    print("Messenger API Response:")
    print(resp)

login_messenger('61573223864739', 'y4O9XUnh0')
