import requests

url = "https://www.facebook.com/api/graphql/"
data = {
    "doc_id": "26505601199092382",
    "variables": '{"UFI2CommentsProvider_commentsKey":"ProfileCometTimelineRoute","beforeTime":null,"count":10,"cursor":null,"dialtone_active":false,"feedLocation":"TIMELINE","feedbackSource":0,"focusCommentID":null,"has_react_native_v2":false,"id":"61572157625165","omitPinnedPost":true,"privacySelectorRenderLocation":"COMET_STREAM","renderLocation":"timeline","scale":1,"useDefaultActor":false}'
}
r = requests.post(url, data=data)
print(r.text)
