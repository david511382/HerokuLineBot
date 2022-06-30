ENV := local

##########################################
## docker-compose
##########################################

DOCKER_SERVICE_LIST := mysql adminer redis redis-commander loki grafana

up: # debug 全開
	docker-compose up -d

up-d: # debug 本地資料環境
	docker-compose up -d \
	$(DOCKER_SERVICE_LIST)

up-s: # stage
	docker-compose \
	-f docker-compose.yml \
	-f deploy/docker/docker-compose.stage.yml \
	up -d

up-p: # prod
	docker-compose \
	-f docker-compose.yml \
	-f deploy/docker/docker-compose.prod.yml \
	up -d

ps: # 查看
	docker-compose ps

down: # 關閉
	docker-compose down

##########################################
## k8s
##########################################

K8S_SERVICE_LIST := mysql adminer redis rediscommander loki grafana linebot

# in shell
kup: # debug 全開
	for %%s in ($(K8S_SERVICE_LIST)) do ( \
		kubectl apply -f deploy/k8s/%%s/$(ENV) \
	) \

kdown: # 關閉
	for %%s in ($(K8S_SERVICE_LIST)) do ( \
		kubectl delete -f deploy/k8s/%%s/$(ENV) \
	) \

##########################################
## test
##########################################

# -p 設定執行續，預設不同 package 會非同步執行
test: # 測試
	go test \
	./src/pkg/util/... \
	./bootstrap/... \
	./src/repo/... \
	./src/logic/... \
	./src/background/... \
	./src/server/... \
	--count=1

mock: # 更新 mock
	mockgen -source="./src/logic/badminton/team.go" -destination="./src/logic/badminton/team_mock.go" -package=badminton
	mockgen -source="./src/logic/badminton/place.go" -destination="./src/logic/badminton/place_mock.go" -package=badminton
	mockgen -source="./src/logic/badminton/court.go" -destination="./src/logic/badminton/court_mock.go" -package=badminton
	mockgen -source="./src/logic/badminton/activity_logic.go" -destination="./src/logic/badminton/activity_logic_mock.go" -package=badminton
	mockgen -source="./src/logic/api/badminton_activity.go" -destination="./src/logic/api/badminton_activity_mock.go" -package=api