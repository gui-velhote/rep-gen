import requests, json

def addEmployee(data:dict):
    r = requests.post("http://127.0.0.1:8080/employee/add", data=json.dumps(data))
    print(f"Status code: {r.status_code} : {r.text}")
    
def getAllEmployees():
    r = requests.get("http://127.0.0.1:8080/employee/getAll")
    print(r.text)

def main():
    getAllEmployees()
    """ addEmployee({
        "name" : "Jederson",
        "privileges" : "admin"
    })"""

if __name__ == "__main__":
    main()