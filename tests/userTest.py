import requests
import json

def getAllUsers():
  r = requests.get("http://127.0.0.1:8080/employee/getAll")
  print(r.text)
  
def addUser(data:dict):
  r = requests.post("http://127.0.0.1:8080/employee/add", data=json.dumps(data))
  print(r.text)
  
def getUserById(data:dict):
  r = requests.get('http://127.0.0.1:8080/employee/getById', data=json.dumps(data))
  print(r.text)
  
def getUserByName(data:dict):
  r = requests.get("http://127.0.0.1:8080/employee/getByName", data=json.dumps(data))
  print(r.text)

def main():
  """ addUser(json.dumps({
    "name" : "Guilherme",
    "privileges" : "admin"
  })) """
  getAllUsers()
  getUserById({"id" : 1})
  getUserByName({"name" : "Guilherme"})

if __name__ == "__main__":
  main()
