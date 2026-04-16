import json
with open("DEBUG_POSTS_RESPONSE.json", "r", encoding="utf-8") as f:
    text = f.read()

def find_url(d, p=""):
    results = []
    if isinstance(d, dict):
        for k, v in d.items():
            if "url" in k.lower() or "link" in k.lower():
                results.append(f"{p}.{k} = {v}")
            if isinstance(v, (dict, list)):
                results.extend(find_url(v, f"{p}.{k}" if p else k))
    elif isinstance(d, list):
        for i, v in enumerate(d):
            results.extend(find_url(v, f"{p}[{i}]"))
    return results

lines = text.split("\n")
data = json.loads(lines[0])
edges = data.get("data", {}).get("node", {}).get("timeline_list_feed_units", {}).get("edges", [])
if edges:
    print("URLS IN ROOT EDGE:")
    print("\n".join(find_url(edges[0])[0:20]))

for i, line in enumerate(lines[1:5]):
    if not line.strip(): continue
    chunk = json.loads(line)
    if "node" in chunk.get("data", {}):
        print(f"URLS IN CHUNK {i+1}:")
        print("\n".join(find_url(chunk.get("data", {}).get("node"))[0:20]))
