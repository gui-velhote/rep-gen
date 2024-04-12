import requests, json

def addBuilding(data:dict):
    r = requests.post("http://127.0.0.1:8080/building/add", data=json.dumps(data))     
    print(f"Status code: {r.status_code} : {r.text}")

def main():
    addBuilding({
        "client_id" : 1,
        "address" : "Rua Pirapora 250 Ap. 111",
        "status" : "Finalizando"
    })

if __name__ == "__main__":
    main()