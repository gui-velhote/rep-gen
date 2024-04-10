import requests, json

def addBuilding(data:dict):
    clientName = {"name" : data.pop("client_name")}
    
    r = requests.get("http://127.0.0.1:8080/client/getByName", data=json.dumps(clientName))
    
    client = json.loads(r.text)

    building = {
        "client_id" : client[0].get("id"),
        "address" : data.get("address"),
        "status" : data.get("status") 
    }
    
    print(building)
    
    r = requests.post("http://127.0.0.1:8080/building/add", data=json.dumps(building))
    print(r.text)
    
    # clientDict = clientTest.getClientByName(client) 
    # print(clientDict)

def getAllBuildings():
    r = requests.get("http://127.0.0.1:8080/building/getAll")
    buildings = json.loads(r.text)
    print(buildings)
    
def getBuildingsByClientId(data:dict):
    r = requests.get("http://127.0.0.1:8080/building/getByClientId", data=json.dumps(data))
    building = json.loads(r.text)
    print(building)
    
def getBuildingById(data:dict):
    r = requests.get("http://127.0.0.1:8080/building/getById", data=json.dumps(data))
    building = json.loads(r.text)
    print(building)

def main():
    """ addBuilding({
        "client_name" : "Camilla",
        "address": "RUA NEBRASKA 929 CASA 1 A",
        "status" : "Finalizando"
    }) """
    getAllBuildings()
    getBuildingsByClientId({"id" : 1})
    getBuildingById({"id" : 1})

if __name__ == "__main__":
    main()