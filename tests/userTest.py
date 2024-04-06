import requests

def getAllUsers():
  r = requests.get("http://127.0.0.1:8080/user/getAll")
  print(r.text)

def main():
  getAllUsers()

if __name__ == "__main__":
  main()
