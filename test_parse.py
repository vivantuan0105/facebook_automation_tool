import json
with open("DEBUG_POSTS_RESPONSE.json", "r", encoding="utf-8") as f:
    text = f.read()
lines = text.split("\n")
first_line = json.loads(lines[0])
edges = first_line.get("data", {}).get("node", {}).get("timeline_list_feed_units", {}).get("edges", [])
print(f"Number of edges in first line: {len(edges)}")
