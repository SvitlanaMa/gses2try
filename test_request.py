import requests

# тест до першого
# url = "http://127.0.0.1:3333/api/rate"
# r = requests.get(url)
# #r = requests.post(url, data={})
# #print(type(r.json()["rate"]))
# print(r.json())

# тест до другого
# url = "http://127.0.0.1:3333/api/subscribe"
# #r = requests.get(url)
# r = requests.post(url, data={"email": " щщщ"}, headers={"Content-Type": "application/x-www-form-urlencoded"})
# print(r.text, r.status_code)
# print(r.json())

# тест до третього
url = "http://127.0.0.1:3333/api/sendEmails"
r = requests.get(url)
# r = requests.post(url, data={})
print(r.text, r.status_code)
print(r.json())



