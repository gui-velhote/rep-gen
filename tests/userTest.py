import requests

def getAllUsers():
  r = requests.get("http://127.0.0.1:8080/employee/getAll")
  print(r.text)
  print(r.status_code)

def main():
    getAllUsers()

if __name__ == "__main__":
  main()
