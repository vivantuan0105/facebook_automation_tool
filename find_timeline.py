import re
import json
with open('DEBUG_PROFILE_DUMP.html', 'r', encoding='utf-8') as f:
    text = f.read()
matches = re.finditer(r'\{[^{]*"id"\s*:\s*"(\d{15,})"[^{}]*"name"\s*:\s*"([^"]+)"', text)
res = set()
for m in matches:
    res.add(m.group(1) + " " + m.group(2))
for r in res:
    print(r)
