import urllib.parse
with open("DEBUG_POSTS_PAYLOAD.txt", "r") as f:
    text = f.read().split("\nVARIABLES")[0]
parsed = urllib.parse.parse_qs(text)
for k, v in parsed.items():
    print(f"{k} = {v}")
