with open("DEBUG_POSTS_RESPONSE.json", "r", encoding="utf-8") as f:
    text = f.read()
lines = text.split("\n")
print(f"Total lines: {len(lines)}")
for i, line in enumerate(lines[:3]):
    print(f"Line {i}: length {len(line)}, starts with {line[:50]}")
