import json

def get_keys(d, p=""):
    keys = []
    if isinstance(d, dict):
        for k, v in d.items():
            keys.extend(get_keys(v, f"{p}.{k}" if p else k))
    elif isinstance(d, list):
        for i, v in enumerate(d):
            keys.extend(get_keys(v, f"{p}[{i}]"))
    else:
        keys.append((p, d))
    return keys

with open("DEBUG_POSTS_RESPONSE.json", "r", encoding="utf-8") as f:
    for line in f.read().split("\n"):
        if not line.strip(): continue
        try:
            chunk = json.loads(line)
        except: continue
        node = chunk.get("data", {}).get("node", {})
        if node.get("__typename") == "Story" or node.get("post_id"):
            kv = get_keys(node)
            for path, val in kv:
                if any(x in path.lower() for x in ["time", "count", "photo", "image", "url", "text"]):
                    if isinstance(val, (str, int, bool)) and val:
                        print(f"{path}: {val}")
            break
