GOFILES := $(shell find . -name "*.go" -type f ! -path "./vendor/*")
STATIKTOOL = vendor/github.com/rakyll/statik/statik

ifeq ($(V),1)
	AT=
else
	AT=@
endif

netconfd: $(GOFILES) swaggerui/statik.go schemas/statik.go deps
	$(AT) CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o $@
	$(AT) strip $@

netconfd.debug: $(GOFILES) swaggerui/statik.go schemas/statik.go deps
	$(AT) CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -a -ldflags '-extldflags "-static"' -o $@

test:
	$(AT) go test ./crm -v

swaggerui/statik.go: swaggerui/swagger.yaml $(STATIKTOOL)
	$(AT) $(STATIKTOOL) -src=swaggerui -include=*.png,*.yaml,*.html,*.css,*.js -p=swaggerui
	
schemas/statik.go: $(STATIKTOOL)
	$(AT) $(STATIKTOOL) -src=schemas -include=*.json -p=schemas

$(STATIKTOOL):
	$(AT) cd vendor/github.com/rakyll/statik && go build

deps: swaggerui/statik.go swaggerui/statik.go schemas/statik.go
	$(AT) go get -d -v

clean:
	$(AT) rm -rf netconfd netconfd.debug statik/ swaggerui/statik.go schemas/statik.go

.PHONY: clean debugbuild test deps
