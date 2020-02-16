APPNAME = blog
BIN = $(GOPATH)/bin
GOCMD = /usr/local/go/bin/go
GOBUILD = $(GOCMD) build
GOINSTALL = $(GOCMD) install
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GORUN = $(GOCMD) run
BINARY_UNIX = $(BIN)/$(APPNAME)
PID = .pid
HUB_ADDR = hub.kpaas.nsini.com
DOCKER_USER =
DOCKER_PWD =
TAG = v0.0.01-test
NAMESPACE = app
PWD = $(shell pwd)

start:
	$(BIN)/$(APPNAME) start -p :8080 -c /etc/blog/app.cfg & echo $$! > $(PID)

restart:
	@echo restart the app...
	@kill `cat $(PID)` || true
	$(BIN)/$(APPNAME) start -p :8080 -c /etc/blog/app.cfg & echo $$! > $(PID)

install:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOINSTALL) -v

stop:
	@kill `cat $(PID)` || true

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

login:
	docker login -u $(DOCKER_USER) -p $(DOCKER_PWD) $(HUB_ADDR)

build:
#	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o kplcloud -v ./main.go
	docker build --rm -t $(APPNAME):$(TAG) .

docker-run:
	docker run -it --rm -p 8080:8080 -v $(PWD)/app.cfg:/etc/app.cfg$(APPNAME):$(TAG)

push:
	docker image tag $(APPNAME):$(TAG) $(HUB_ADDR)/$(NAMESPACE)/$(APPNAME):$(TAG)
	docker push $(HUB_ADDR)/$(NAMESPACE)/$(APPNAME):$(TAG)

run:
	GO111MODULE=on $(GORUN) ./main.go start -p :8080 -c ./app.cfg