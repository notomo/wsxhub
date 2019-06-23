
test:
	GO111MODULE=on go build -o dist/wsxhub ./main.go
	go test -v github.com/notomo/wsxhub/... -race -coverprofile=coverage.txt -covermode=atomic
	$(MAKE) coverage

coverage:
	go tool cover -html=coverage.txt -o index.html

install:
	go install github.com/notomo/wsxhub

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
.PHONY: coverage
.PHONY: install
.PHONY: reup
.PHONY: start
.PHONY: deploy
