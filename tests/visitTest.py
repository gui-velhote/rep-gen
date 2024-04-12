import requests, json

def addVisit(data:dict):
    r = requests.post("http://127.0.0.1:8080/visit/add", data=json.dumps(data))
    print(r.text)
    
def getAllVisits():
    r = requests.get("http://127.0.0.1:8080/visit/getAll")
    print(r.text)
    # visits = json.loads(r.text)
    # print(visits)
    
def reportTest(data:dict):
    r = requests.post("http://127.0.0.1:8080/visit/report/add", data=json.dumps(data))
    print(r.text)
    print(r.status_code)
    
def getVisit(data:dict) -> dict:
    r = requests.get("http://127.0.0.1:8080/visit/report/get", data=json.dumps(data))
    return json.loads(r.text)

def saveReport(report, fileName):
    
    print(report)
    
    dash = 50 * '-' 

    parsedReport = dash + "\nEquipe: "

    for employee in report.get("team_names"):
        if len(report.get("team_names")) > 1 and report.get("team_names").index(employee) != len(report.get("team_names")) - 1:
            parsedReport = parsedReport + employee + ", "
            
    else:
        parsedReport = parsedReport + employee
            
    
    parsedReport = parsedReport + f"\nCliente: {report.get('client_name')}\nData: {report.get('date')}\nCarro: {report.get('car')}\nAtividades:\n"
    
    for activity in report.get("ACTIVITY"):
        parsedReport = parsedReport + f"- {activity.get('activity_description')}\n"
        
    parsedReport = parsedReport + "Observacoes:\n"
    
    for observation in report.get("OBSERVATION"):
        parsedReport = parsedReport + f"- {observation.get('observation_description')}\n"
        
    parsedReport = parsedReport + "Pendencias\n"
    
    for pendency in report.get("PENDENCY"):
        parsedReport = parsedReport + f"- {pendency.get('pendency_description')}"
        
    parsedReport = parsedReport + "\n" + dash

    print(parsedReport)

    # with open(fileName, 'w') as f:
        

def main():
    # addVisit({"date" : "2024-04-08", "car" : "SEO3J42"})
    getAllVisits()
    """ reportTest({
        "date": "2024-04-09",
        "car" : "SEO3J42",
        "client_id": 1,
        "building_id": 2,
        "team_ids" : [1, 2],
        "activity" : [{"activity_id" : 1, "activity_description" : "Passagem de cabo cat6 no terraco"}],
        "observation" : [{"observation_id": 1, "observation_description" : "Sem observacoes"}],
        "pendency" : [{"pendency_id" : 1, "pendency_description" : "Sem pendencias"}]
    }) """
    report = getVisit({"id" : 9})

    saveReport(report, "reportTest.txt")

if __name__ == "__main__":
    main()
