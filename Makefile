IMAGE_TAG ?= $(shell git describe --tags)

.PHONY: build-gofakeit
build-gofakeit:
	go build -v -o bin/gofakeit -ldflags "-w -s" cmd/gofakeit/gofakeit.go

.PHONY: build-server
build-server:
	go build -v -o bin/gofakeitserver -ldflags "-w -s" cmd/gofakeitserver/main.go

.PHONY: images
images:
	docker build --progress=plain -t yz271544/gofakeit-server:${IMAGE_TAG} -f Dockerfile .

.PHONY: docker-run
docker-run:
	docker run -itd -p 18080:8080 --restart always --name gofaker yz271544/gofakeit-server:${IMAGE_TAG}
