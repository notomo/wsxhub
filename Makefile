
test:
	go test github.com/notomo/wsxhub/server -race

install:
	go install github.com/notomo/wsxhub/cmd/wsxhub
	go install github.com/notomo/wsxhub/cmd/wsxhubd

reup:
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

.PHONY: test
.PHONY: install
.PHONY: reup
