
test:
	go build -o dist/wsxhub ./cmd/wsxhub 
	go build -o dist/wsxhubd ./cmd/wsxhubd 
	go test github.com/notomo/wsxhub/server -race

install:
	go install github.com/notomo/wsxhub/cmd/wsxhub
	go install github.com/notomo/wsxhub/cmd/wsxhubd

reup:
	docker-compose down
	docker-compose up -d

v=
deploy:
	git tag v${v}
	git push origin v${v}

.PHONY: test
.PHONY: install
.PHONY: reup
.PHONY: deploy
