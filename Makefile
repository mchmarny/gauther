GCP_PROJECT=s9-demo
BINARY_NAME=gauther

test:
	go test ./... -v

cover:
	go test ./... -cover
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

cert:
	echo "Downloading latest ca certs from Mozilla..."
	curl -o ./certs/ca-certificates.crt https://curl.haxx.se/ca/cacert.pem

deps:
	go mod tidy

docs:
	godoc -http=:8888 &
	open http://localhost:8888/pkg/github.com/mchmarny/$(BINARY_NAME)/
	# killall -9 godoc

image:
	gcloud builds submit \
		--project $(GCP_PROJECT) \
		--tag gcr.io/$(GCP_PROJECT)/$(BINARY_NAME):latest

docker:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -itP --expose 8080 $(DOCKER_USERNAME)/$(BINARY_NAME):latest

secrets:
	kubectl create secret generic gauther \
		--from-literal=OAUTH_CLIENT_ID=$(GAUTHER_OAUTH_CLIENT_ID) \
		--from-literal=OAUTH_CLIENT_SECRET=$(GAUTHER_OAUTH_CLIENT_SECRET)

service:
	kubectl apply -f deployments/service.yaml

serviceless:
	kubectl delete -f deployments/service.yaml