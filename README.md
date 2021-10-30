# HerokuLineBot

## test

```
go test ./background/... --count=1

go test ./logic/... --count=1

go test ./storage/database/... --count=1

go test ./bootstrap/... --count=1

go test ./util/... --count=1
```

## db migration

### postgre

```
SELECT SETVAL('activity_id_seq', (SELECT MAX(id) FROM activity));
SELECT SETVAL('income_id_seq', (SELECT MAX(id) FROM income));
SELECT SETVAL('logistic_id_seq', (SELECT MAX(id) FROM logistic));
SELECT SETVAL('member_id_seq', (SELECT MAX(id) FROM member));
SELECT SETVAL('member_activity_id_seq', (SELECT MAX(id) FROM member_activity));
SELECT SETVAL('place_id_seq', (SELECT MAX(id) FROM place));
SELECT SETVAL('rental_court_id_seq', (SELECT MAX(id) FROM rental_court));
SELECT SETVAL('rental_court_detail_id_seq', (SELECT MAX(id) FROM rental_court_detail));
SELECT SETVAL('rental_court_ledger_id_seq', (SELECT MAX(id) FROM rental_court_ledger));
SELECT SETVAL('rental_court_ledger_court_id_seq', (SELECT MAX(id) FROM rental_court_ledger_court));
```
