# Repository generator requisits

## Fucntional requisits

1. The system should manage the `employees` data such as name and their privileges to access data, with methods to add, modify, delete and search for users

2. The system should manage the `clients` data such as name an address, with methods to add, modify, delete and search for clients

3. The system should be able to manage `building` data such as which client is the owner and the status of the building, allowing for supervisors to add, modify and search for buildings.

4. The system should be able to manage `reports` with the team, date, car and which building the work was performed on, this should be done with the methods of adding, modifying, deleting and searching for entries.

5. `Activities, Pendencies and Observations` should be also stored on the database with the methods of adding, deleting, modifying and searching.

## Non-functional requisits

1. None of the supervisors should be able to delete any data from the database, only the database manager can delete entries.

2. Only the database manager can include new users.

3. Employees should have name and privileges.

4. Clients should have name.

5. Buildings should have the owner's id, the building's address, status and the supervisor's id.

6. Reports should have the Date, car, building id and the team performing the work.

7. Activities, pendencies and observations should have their description and the report id.

8. The Visit should have the bu