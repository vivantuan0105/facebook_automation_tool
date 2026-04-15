import requests
import json
import time

doc_id = "25845944448416411"

session = requests.Session()
r = session.get("https://www.facebook.com/")
lsd = "fake_lsd"
for line in r.text.split('\n'):
    if 'name="lsd"' in line: lsd = line.split('value="')[1].split('"')[0]; break
    if '"LSD",[]' in line: lsd = line.split('"token":"')[1].split('"')[0]; break

username = "61574299097448"
enc_pwd = f"#PWD_BROWSER:5:{int(time.time())}:S2lj9NlT93"

formats = [
    json.dumps({
        "device_id": "",
        "family_device_id": "",
        "machine_id": "",
        "password": enc_pwd,
        "username": username,
        "client_mutation_id": "1",
        "credential_type": "password",
        "login_source": "CometSessionCreateMutation"
    }),
    json.dumps({
        "input": {
            "credentials_type": "password",
            "error_detail_type": "button_with_disabled",
            "source": "login",
            "password": enc_pwd,
            "login_attempt_count": 0,
            "server_timestamps": True,
            "identifier": username,
            "lsd": lsd,
            "jazoest": "29199",
            "openid_tokens": {},
            "client_mutation_id": "1",
            "is_two_factor_auth_in_app_browser": False,
            "is_two_factor_auth_in_app_browser_webview": False,
            "is_two_factor_auth_in_app_browser_webview_bypass": False,
            "is_from_native_webview": False
        },
        "scale": 1
    }),
    json.dumps({
        "input": {
            "password": enc_pwd,
            "login_attempt_count": 0,
            "identifier": username,
            "lsd": lsd,
        },
        "scale": 1
    })
]

for idx, var_str in enumerate(formats):
    payload = {"__a": "1", "doc_id": doc_id, "variables": var_str, "lsd": lsd}
    resp = session.post("https://www.facebook.com/api/graphql/", data=payload)
    print(f"--- Format {idx+1} ---")
    try:
        data = resp.json()
        if "errors" in data and len(data["errors"]) > 0:
            print("  ERR count:", len(data["errors"]))
            print("  ERR 0:", data["errors"][0]["message"])
        else:
            print("  SUCCESS!")
            print(data)
    except: pass
