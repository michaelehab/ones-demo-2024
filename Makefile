.PHONY: build test run run_mqtt

MICROSERVICES=cmd/creator/creator-demo \
				cmd/mutator/mutator-demo \
				cmd/transitor/transitor-demo

.PHONY: $(MICROSERVICES)

.PHONY: build
build: $(MICROSERVICES)

.PHONY: cmd/creator/creator-demo
cmd/creator/creator-demo:
	go build -o $@ ./cmd/creator

.PHONY: cmd/mutator/mutator-demo
cmd/mutator/mutator-demo:
	go build -o $@ ./cmd/mutator

.PHONY: cmd/transitor/transitor-demo
cmd/transitor/transitor-demo:
	go build -o $@ ./cmd/transitor

run:
	cd scripts/bin && ./launch.sh

run_mqtt:
	cd scripts/bin && ./launch-mqtt.sh

test:
	go test -cover ./...
	go vet ./...
	gofmt -l .
	[ "`gofmt -l .`" = "" ]