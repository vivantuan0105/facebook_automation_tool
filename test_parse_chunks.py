import json
with open("DEBUG_POSTS_RESPONSE.json", "r", encoding="utf-8") as f:
    text = f.read()
lines = text.split("\n")
for i in range(1, len(lines)):
    if lines[i].strip():
        chunk = json.loads(lines[i])
        data = chunk.get("data", {})
        if "page_info" in data:
            print("Found page_info in line", i, data["page_info"])
        elif isinstance(data, list):
            print("Data is list in line", i, len(data))
        elif "node" in data:
            print(f"Node in line {i}", data["node"].keys())
        else:
            print(f"Data struct in line {i}:", data.keys() if isinstance(data, dict) else type(data))
