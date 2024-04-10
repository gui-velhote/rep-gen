import requests, json

def addVisit(data:dict):
    r = requests.post("http://127.0.0.1:8080/visit/add", json.dumps(data))
    print(r.text)
    
def getAllVisits():
    r = requests.get("http://127.0.0.1:8080/visit/getAll")
    visits = json.loads(r.text)
    print(visits)
    
def reportTest(data:dict):
    r = requests.post("http://127.0.0.1:8080/report/test", json.dumps(data))
    print(r.text)

def main():
    # addVisit({"date" : "2024-04-08", "car" : "SEO3J42"})
    getAllVisits()
    reportTest({
        "id" : 1,
        "date": "2024-04-08",
        "car" : "SEO3J42",
        "client_id": 1,
        "team_ids" : [1],
        "activity" : [{"activity_id" : 1, "activity_description" : "Passagem de cabo cat6 no terraco"}],
        "observation" : [{"observation_id": 1, "observation_description" : "Sem observacoes"}],
        "pendency" : [{"pendency_id" : 1, "pendency_description" : "Sem pendencias"}]
    })

if __name__ == "__main__":
    main()