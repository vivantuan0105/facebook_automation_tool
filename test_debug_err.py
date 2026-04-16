import json
with open("DEBUG_POSTS_RESPONSE.json", "r", encoding="utf-8") as f:
    resp = json.load(f)
print("Keys:", resp.keys())
print("Data exists in dict?", "data" in resp)
