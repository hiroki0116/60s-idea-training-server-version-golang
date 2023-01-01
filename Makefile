run dev:
	docker-compose -f dockers/docker-compose.yml build && docker-compose -f dockers/docker-compose.yml up
run test:
	STAGE=test go test ./test --v     