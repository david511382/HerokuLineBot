
##########################################
## docker-compose
##########################################

DATA_SERVICE_LIST := db adminer redis redis-commander loki grafana

up: # debug 全開
	docker-compose up

up-d: # debug 本地資料環境
	docker-compose up \
	-d $(DATA_SERVICE_LIST)

up-s: # stage
	docker-compose \
	-f docker-compose.yml \
	-f docker/docker-compose.stage.yml \
	up

up-p: # prod
	docker-compose \
	-f docker-compose.yml \
	-f docker/docker-compose.prod.yml \
	up

ps: # 查看
	docker-compose ps

down: # 關閉
	docker-compose down
