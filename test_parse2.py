import json
with open("DEBUG_POSTS_RESPONSE.json", "r", encoding="utf-8") as f:
    text = f.read()
lines = text.split("\n")
if len(lines) > 1 and lines[1].strip():
    line1 = json.loads(lines[1])
    edges1 = line1.get("data", {}).get("edges", [])
    print(f"Number of edges in line 1: {len(edges1)}")
    if edges1:
        print(edges1[0].keys())
