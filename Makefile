# Assumes following env vars set
#  * GCP_PROJECT - ID of your project
#  * GAUTHER_OAUTH_CLIENT_ID - Google OAuth2 Client ID
#  * GAUTHER_OAUTH_CLIENT_SECRET - Google OAuth2 Client Secret

test:
	go test ./... -v

cover:
	go test ./... -cover
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

deps:
	go mod tidy

docs:
	godoc -http=:8888 &
	open http://localhost:8888/pkg/github.com/mchmarny/gauther/
	# killall -9 godoc

image:
	gcloud builds submit \
		--project $(GCP_PROJECT) \
		--tag gcr.io/$(GCP_PROJECT)/gauther:latest

docker:
	docker build -t gauther .

secrets:
	# kubectl delete secret gauther -n demo
	kubectl create secret generic gauther -n demo \
		--from-literal=OAUTH_CLIENT_ID=$(DEMO_OAUTH_CLIENT_ID) \
		--from-literal=OAUTH_CLIENT_SECRET=$(DEMO_OAUTH_CLIENT_SECRET)

service:
	kubectl apply -f deployments/service.yaml
	kubectl get pods -n demo

service-clean:
	echo "BEFORE"
	kubectl get pods
	kubectl delete -f deployments/service.yaml
	kubectl apply -f deployments/service.yaml
	echo "AFTER"
	kubectl get pods
