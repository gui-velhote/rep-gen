SELECT c.name, v.date, e.id FROM Visit as v INNER JOIN Makes_visit as mk_v ON v.id = mk_v.visit_id INNER JOIN Employee as e ON mk_v.employee_id = e.id INNER JOIN Recieves_visit as rv_v ON v.id = rv_v.id INNER JOIN Building as b ON rv_v.building_id = b.id INNER JOIN Client as c ON b.client_id = c.id;

 INNER JOIN Activity as a ON v.id = a.visit_id INNER JOIN Observation as o ON v.id = o.visit_id INNER JOIN Pendency as p ON V.id = p.visit_id;