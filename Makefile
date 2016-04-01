.PHONY: tfall tfplan tfapply help

.DEFAULT_GOAL := help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

tfall: tfplan tfapply tfplan ## Run all terraform tasks

tfplan: ## Make a terraform plan
	@bin/terraform plan -var-file terraform/terraform.tfvars -out terraform/terraform.plan -state terraform/terraform.tfstate terraform

tfapply: ## Apply terraform plan
	@bin/terraform apply -state terraform/terraform.tfstate terraform/terraform.plan


