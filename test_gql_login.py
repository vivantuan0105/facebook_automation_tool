import requests
import json

doc_id = "25845944448416411"

session = requests.Session()
r = session.get("https://www.facebook.com/")
lsd = "fake_lsd"

formats = [
    # Format 1: my original code
    json.dumps({
        "input": {
            "client_mutation_id": "1",
            "identifier": "test_user",
            "password": "#PWD_BROWSER:5:1700000000:fake",
            "lsd": lsd
        },
        "locale": "vi_VN"
    }),
    # Format 2: no locale, minimal input
    json.dumps({
        "input": {
            "identifier": "test_user",
            "password": "#PWD_BROWSER:5:1700000000:fake",
            "credentials_type": "password",
            "lsd": lsd
        }
    }),
    # Format 3: from facebook js dump
    json.dumps({
        "input": {
            "credentials_type":"password",
            "error_detail_type":"button_with_disabled",
            "source":"login",
            "password":"#PWD_BROWSER:5:1700000000:fake",
            "login_attempt_count":0,
            "openid_tokens":{},
            "identifier":"testuser",
            "lsd":lsd
        },
        "scale":1.5,
        "is_two_factor_auth_in_app_browser":False,
        "is_two_factor_auth_in_app_browser_webview":False,
        "is_two_factor_auth_in_app_browser_webview_bypass":False
    }),
    # Format 4: flat variables without input
    json.dumps({
        "identifier":"testuser",
        "password":"#PWD_BROWSER:5:1700000000:fake",
        "lsd":lsd,
        "credentials_type":"password"
    })
]

for idx, var_str in enumerate(formats):
    payload = {
        "__a": "1",
        "doc_id": doc_id,
        "variables": var_str
    }
    resp = session.post("https://www.facebook.com/api/graphql/", data=payload)
    print(f"--- Format {idx+1} ---")
    try:
        data = resp.json()
        if "errors" in data:
            print("  ERR:", [e.get("message", e)[:50] for e in data["errors"]])
        else:
            print("  SUCCESS:", data.get("data"))
    except Exception as e:
        print("  Parse error:", e, resp.text[:100])
