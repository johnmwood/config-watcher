YAML=/Users/jwood/go/src/github.com/johnmwood/config-watcher/kubernetes/deployment.yaml

deploy:
	@echo "Deploying config watcher"
	@echo ""
	@kubectl delete -f $(YAML)
	@kubectl apply -f $(YAML)
	@echo ""
	@sleep 3
	
	@echo "Viewing logs:"
	@echo ""
	kubectl logs -f $$(kubectl get pods | grep watcher | awk 'FNR == 1 {print $$1}')

vendor:
	go mod tidy
	go mod vendor

build:
	minikube image build -t config-watcher:v1 .

build-deploy: build deploy
