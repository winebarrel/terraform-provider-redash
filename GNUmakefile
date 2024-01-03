.PHONY: build
build: vet
	go build

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test -v -count=1 ./...

.PHONY: testacc
testacc:
	$(MAKE) test TF_ACC=1

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	rm -f terraform-provider-redash

dev.tfrc:
	sed "s|{{PATH_TO_PROVIDER}}|$(shell pwd)|" etc/dev.tfrc.tpl > dev.tfrc

.PHONY: tf-plan
tf-plan: build dev.tfrc
	TF_CLI_CONFIG_FILE=dev.tfrc terraform plan

.PHONY: tf-apply
tf-apply: build dev.tfrc
	TF_CLI_CONFIG_FILE=dev.tfrc terraform apply -auto-approve

.PHONY: tf-clean
tf-clean: clean
	rm -f dev.tfrc terraform.tfstate*

.PHONY: redash-setup
redash-setup:
	psql -U postgres -h localhost -p 15432 -f etc/redash.sql

.PHONY: redash-create-db
redash-create-db:
	docker compose run --rm server create_db

# cf. https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-documentation-generation
.PHONY: docs
docs:
	go generate ./...
