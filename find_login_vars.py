import requests
import json
import time

doc_id = "25845944448416411"

session = requests.Session()
r = session.get("https://www.facebook.com/")
print("Fetching LSD...")
lsd = None
for line in r.text.split('\n'):
    if 'name="lsd"' in line:
        lsd = line.split('value="')[1].split('"')[0]
        break
    if '"LSD",[]' in line:
        lsd = line.split('"token":"')[1].split('"')[0]
        break

jazoest = None
for line in r.text.split('\n'):
    if 'name="jazoest"' in line:
        jazoest = line.split('value="')[1].split('"')[0]
        break

print(f"LSD: {lsd}, Jazoest: {jazoest}")

# Encoded password using static values for test
timestamp = int(time.time())
enc_password = f"#PWD_BROWSER:5:{timestamp}:S2lj9NlT93"
username = "61574299097448"

formats = [
    # FB CometLogin 1
    json.dumps({
        "device_id": "",
        "family_device_id": "",
        "machine_id": "",
        "password": enc_password,
        "username": username,
        "credential_type": "password",
        "login_source": "CometSessionCreateMutation",
        "error_detail_type": "button_with_disabled"
    }),
    # FB caa_login_web
    json.dumps({
        "input": {
            "credentials_type": "password",
            "error_detail_type": "button_with_disabled",
            "source": "login",
            "password": enc_password,
            "login_attempt_count": 0,
            "identifier": username,
            "lsd": lsd,
            "jazoest": jazoest
        },
        "scale": 1.5,
        "is_two_factor_auth_in_app_browser": False,
        "is_two_factor_auth_in_app_browser_webview": False,
        "is_two_factor_auth_in_app_browser_webview_bypass": False
    }),
    # FB Login minimal input wrapper
    json.dumps({
        "input": {
            "identifier": username,
            "password": enc_password,
        }
    })
]

for idx, var_str in enumerate(formats):
    payload = {
        "__a": "1",
        "doc_id": doc_id,
        "variables": var_str,
        "lsd": lsd,
        "jazoest": jazoest
    }
    resp = session.post("https://www.facebook.com/api/graphql/", data=payload)
    print(f"--- Format {idx+1} ---")
    data = resp.json()
    if "errors" in data:
        print(" ERR:", [e.get("message", e)[:80] for e in data["errors"]])
    else:
        print(" SUCCESS:", data.get("data"))
        print(" Cookies:", session.cookies.get_dict())
        break
