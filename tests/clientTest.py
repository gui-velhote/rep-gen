import requests, json

def addClient(data:dict):
    r = requests.post("http://127.0.0.1:8080/client/add", data=json.dumps(data))     
    print(f"Status code: {r.status_code} : {r.text}")

def main():
    addClient({
        "name" : "Ricardo Leão",
    })

if __name__ == "__main__":
    main()