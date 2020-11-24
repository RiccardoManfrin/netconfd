GOFILES := $(shell find . -name "*.go" -type f ! -path "./vendor/*")
STATIKTOOL = vendor/github.com/rakyll/statik/statik
OPENAPI = server/

ifeq ($(V),1)
	AT=
else
	AT=@
endif

netconfd: $(GOFILES) $(OPENAPI) swaggerui/statik.go deps
	$(AT) CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o $@
	$(AT) strip $@

netconfd.debug: $(GOFILES) swaggerui/statik.go deps
	$(AT) CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -a -ldflags '-extldflags "-static"' -o $@

$(OPENAPI):
	$(AT) ./gen_templates.sh

test:
	$(AT) go test -v

swaggerui/statik.go: swaggerui/openapi.yaml $(STATIKTOOL)
	$(AT) $(STATIKTOOL) -src=swaggerui -include=*.png,*.yaml,*.html,*.css,*.js,*.json -p=swaggerui
	
$(STATIKTOOL):
	$(AT) cd vendor/github.com/rakyll/statik && go build

deps: swaggerui/statik.go $(OPENAPI)
	$(AT) go get -d -v

clean:
	$(AT) rm -rf netconfd netconfd.debug statik/ swaggerui/statik.go $(STATIKTOOL) $(OPENAPI)

.PHONY: clean debugbuild test deps
