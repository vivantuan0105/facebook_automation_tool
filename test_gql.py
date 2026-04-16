import requests

url = "https://www.facebook.com/api/graphql/"
data = {
    "doc_id": "26505601199092382",
    "variables": "{}"
}
r = requests.post(url, data=data)
print(r.text)
