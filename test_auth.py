import json
with open('Data/accounts.json', 'r', encoding='utf-8') as f:
    data = json.load(f)
for acc in data:
    if acc.get('uid') == '61571360758847' or acc.get('uid') == '61572157625165':
        print(acc['cookie'][:20], acc['fb_dtsg'][:20])
