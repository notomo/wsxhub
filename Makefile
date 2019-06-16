
test:
	go build -o dist/wsxhub ./main.go

install:
	go install github.com/notomo/wsxhub/cmd/wsxhub
	go install github.com/notomo/wsxhub/cmd/wsxhubd

reup:
	docker-compose down
	docker-compose up -d

start:
	go run main.go server

v=
deploy:
	git tag v${v}
	git push origin v${v}

.PHONY: test
.PHONY: install
.PHONY: reup
.PHONY: start
.PHONY: deploy
