TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=github.com
NAMESPACE=sharathrnair87
NAME=awx
BINARY=terraform-provider-${NAME}
VERSION=1.1.1
OS_ARCH=linux_amd64

default: install

fmt:
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

build:
	go build -o ${BINARY}_v${VERSION}

release:
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64

install: build
		# Ensure this path is mentioned in the .terraformrc file
		# https://developer.hashicorp.com/terraform/cli/config/config-file#explicit-installation-method-configuration
		cp ${BINARY}_v${VERSION} ~/go/bin/${BINARY}
		mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
		mv ${BINARY}_v${VERSION} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

#test:
#        go test -i $(TEST) || exit 1
#        echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4
#
#testacc:
#        TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m
