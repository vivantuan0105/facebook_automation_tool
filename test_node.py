import json

def find_keys(obj, target_key, current_path=""):
    results = []
    if isinstance(obj, dict):
        for k, v in obj.items():
            new_path = f"{current_path}.{k}" if current_path else k
            if target_key.lower() in k.lower():
                results.append((new_path, type(v).__name__))
            results.extend(find_keys(v, target_key, new_path))
    elif isinstance(obj, list):
        for i, item in enumerate(obj):
            new_path = f"{current_path}[{i}]"
            results.extend(find_keys(item, target_key, new_path))
    return results

with open("DEBUG_POSTS_RESPONSE.json", "r", encoding="utf-8") as f:
    text = f.read()
lines = text.split("\n")
data = json.loads(lines[0])

edges = data.get("data", {}).get("node", {}).get("timeline_list_feed_units", {}).get("edges", [])
if edges:
    node = edges[0].get("node", {})
    
    # Try to find time, image, like, comment
    print("TIME:", [k[0] for k in find_keys(node, "time")])
    print("IMAGE:", [k[0] for k in find_keys(node, "image")])
    print("PHOTO:", [k[0] for k in find_keys(node, "photo")])
    print("COUNT:", [k[0] for k in find_keys(node, "count")])
    
