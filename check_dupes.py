with open(r"Data\61572157625165\Posts.txt", "r", encoding="utf-8") as f:
    lines = f.readlines()

ids = []
duplicates = []
for line in lines:
    if line.startswith("==="): continue
    if "]: " in line:
        post_id = line.split("]: ")[0].strip("[")
    elif "|" in line:
        post_id = line.split("|")[0]
    else: continue
    
    if post_id in ids:
        duplicates.append(post_id)
    ids.append(post_id)

print(f"Total lines: {len(lines)}")
print(f"Valid posts parsed: {len(ids)}")
print(f"Duplicates: {len(duplicates)}")
if duplicates:
    print(f"Some duplicates: {duplicates[:5]}")
