GOFILES := $(shell find . -name "*.go" -type f ! -path "./vendor/*")
OPENAPI = server/

ifeq ($(V),1)
	AT=
else
	AT=@
endif

netconfd: $(GOFILES) $(OPENAPI) deps swaggerui/openapi.yaml
	$(AT) CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o $@
	$(AT) strip $@

netconfd.debug: $(GOFILES) deps swaggerui/openapi.yaml
	$(AT) CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -a -ldflags '-extldflags "-static"' -o $@

$(OPENAPI):
	$(AT) ./gen_templates.sh

test:
	$(AT) go test -v

deps: $(OPENAPI)
	$(AT) go get -d -v

clean:
	$(AT) rm -rf netconfd netconfd.debug $(OPENAPI)

.PHONY: clean debugbuild test deps
