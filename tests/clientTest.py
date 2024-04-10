import requests
import json

def getAllClients():
    r = requests.get("http://127.0.0.1:8080/client/getAll")
    print(r.text)
    
def getClientByName(data:dict) -> dict:
    r = requests.get("http://127.0.0.1:8080/client/getByName", data=json.dumps(data))
    print(r.text)
    return json.loads(r.text)
    
def addClient(data:dict):
    r = requests.post("http://127.0.0.1:8080/client/add", data=json.dumps(data))
    print(r.text)
    
def getClientById(data:dict):
    r = requests.get("http://127.0.0.1:8080/client/getByName", data=json.dumps(data))
    print(r.text)

def main():
    # addClient({"name" : "Camilla Lunardelli"})
    getAllClients()
    getClientByName({"name" : "camilla"})
    getClientById({"id" : 1})

if __name__ == "__main__":
    main()