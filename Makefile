TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=hashicorp.com
NAMESPACE=eeveebank
NAME=sumologiccse
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=linux_amd64
PROVIDER_PATH=${TF_PLUGIN_CACHE_DIR}/registry.terraform.io/${NAMESPACE}/${NAME}/$(VERSION)/${OS_ARCH}/


default: install

build:
	go build -o ${BINARY}_$(VERSION)

release:
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

install: build
	mkdir -p ${PROVIDER_PATH}
	mv ${BINARY}_$(VERSION) ${PROVIDER_PATH}

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   
