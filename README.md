
rep-gen
=======

A report generator made to ease daily reports.

## About the project

This project was made with the intentions of maing it easier to create reports for my day to day life on the company I work at.

## Table of Contents

- [Requirements](#requirements)
- [Getting Started](#getting-started)
- [Endpoints](#endpoints)
    - [Employee](#employee)
        - [GET](#get-endpoints)
            - [getall](#getall)
            - [getById](#getbyid)
            - [getByName](#getbyname)
        - [POST](#post-endpoints)
            - [add](#add)

## Requirements

This project require you to have installed:
- golang;
- MySql or MariaDb (More options will be developed);
- python (optional for the testings provided).

## Getting Started

For the default installation on linux based distros run:

```
$ ./setup.sh
```

__\*Windows systems are not supported yet\*__

To create the database needed use the sql script provided on `/sql/generate.sql` as root.

Example:  
```
$ mariadb -u root -p < /path/to/sql/generate.sql
```

## Endpoints

### Employee 

All employees methods and endpoints starts with `/employee/`.

### GET endpoints

#### getAll

Returns all employees registered on the database.

Example:

```
[
    {
        "id" : 1,
        "name" : "name1",
        "privileges" : "admin",
    },
    {
        "id" : 2,
        "name" : "name2",
        "privileges" : "employee"
    }
]
```

#### getById

Returns the employee based on the id parsed through the request as a json body.

Exemple:

- python:
```
import requests

data = {"id" : 1}
r = requests.get("http://localhost:8080/employee/getById", data=data)
print(r.text)
```

- returns:
```
{
    "id" : 1,
    "name" : "yourName",
    "privileges" : "admin"
}
```

#### getByName

Returns the employee by the name given by the http body data.

Example:

- python:

```
import requests

data = {"name" : "yourName"}
r = requests.get("http://localhost:8080/employee/getByName", data=data)
print(r.text)
```

- returns:

```
{
    "id" : 1,
    "name" : "yourName",
    "privileges" : "admin"
}
```

### POST endpoints

#### add

Endpoint to add a new user based on the data passed through the http body with the name and privileges of the user.

Example:

- python:

```
import requests

data = {
    "name" : "yourName",
    "privileges" : "admin"
}
r = requests.post("http://localhost:8080/employee/addEmployee", data=data)
```
