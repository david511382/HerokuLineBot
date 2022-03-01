# HerokuLineBot

## Start Local

### docker

``` bash
make up-d
```
http://localhost:8882/?pgsql=db&username=root&sql=CREATE%20DATABASE%20club%3B

``` bash
make up
```

### k8s

``` bash
docker build -f ./deploy/docker/without-config.Dockerfile -t line-bot .

make kup
```

http://localhost:31882/?pgsql=postgres&username=root&sql=CREATE%20DATABASE%20club%3B

``` bash
make up
```

## Start Debug

``` bash
make up-d
```

http://localhost:8882/?pgsql=db&username=root&sql=CREATE%20DATABASE%20club%3B

## dev

token:`{"RoleID":1}`

## test

```
go test ./src/util/... --count=1

go test ./bootstrap/... --count=1

go test ./src/repo/database/... --count=1

go test ./src/logic/... --count=1

go test ./src/background/... --count=1

go test ./src/server/... --count=1
```

## db migration

### postgre

```
SELECT SETVAL('activity_id_seq', (SELECT MAX(id) FROM activity));
SELECT SETVAL('activity_finished_id_seq', (SELECT MAX(id) FROM activity_finished));
SELECT SETVAL('income_id_seq', (SELECT MAX(id) FROM income));
SELECT SETVAL('logistic_id_seq', (SELECT MAX(id) FROM logistic));
SELECT SETVAL('member_id_seq', (SELECT MAX(id) FROM member));
SELECT SETVAL('member_activity_id_seq', (SELECT MAX(id) FROM member_activity));
SELECT SETVAL('place_id_seq', (SELECT MAX(id) FROM place));
SELECT SETVAL('rental_court_id_seq', (SELECT MAX(id) FROM rental_court));
SELECT SETVAL('rental_court_detail_id_seq', (SELECT MAX(id) FROM rental_court_detail));
SELECT SETVAL('rental_court_exception_id_seq', (SELECT MAX(id) FROM rental_court_exception));
SELECT SETVAL('rental_court_ledger_id_seq', (SELECT MAX(id) FROM rental_court_ledger));
SELECT SETVAL('rental_court_ledger_court_id_seq', (SELECT MAX(id) FROM rental_court_ledger_court));
SELECT SETVAL('rental_court_refund_ledger_id_seq', (SELECT MAX(id) FROM rental_court_refund_ledger));
SELECT SETVAL('team_id_seq', (SELECT MAX(id) FROM team));
```

## 羽球業務邏輯

1. 訂金跟尾款必須付同樣日期的場地
