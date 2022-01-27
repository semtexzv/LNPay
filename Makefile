# Reexport values from terraform environment
include ./vars.env
export

MAKEFLAGS += --no-builtin-rules

.PHONY: infra template build deploy frontend template
all: template build deploy frontend

KUBE_DIR=kube
TPL_DIR=target/tmpl
KBLD_YML=target/built.yaml

KUBE_FILES=$(wildcard $(KUBE_DIR)/*.yaml) $(wildcard $(KUBE_DIR)/*.yml)
TPL_FILES=$(KUBE_FILES:kube/%=$(TPL_DIR)/%)

infra:
	cd infra && terraform apply -auto-approve

$(TPL_DIR):
	mkdir -p $(TPL_DIR)

$(TPL_DIR)/%.yaml: kube/%.yaml $(TPL_DIR) ./vars.env
	envsubst "`env | cut -d= -f1 | sed -e "s/^/$$/"`" < $< > $@

template: $(TPL_FILES)

$(KBLD_YML) build: $(TPL_FILES)
	kbld -f $(TPL_DIR) > $(KBLD_YML)

deploy: $(KBLD_YML)
	kapp deploy -y -a lnpay -f $(KBLD_YML)

frontend:
	cd frontend && yarn build
	az storage blob upload-batch --account-name $(ASSETS_STORAGE_ACCOUNT) -d '$$web' -s frontend/dist
