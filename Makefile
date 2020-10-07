GOFILES := $(shell find . -name "*.go" -type f ! -path "./vendor/*")

ifeq ($(V),1)
	AT=
else
	AT=@
endif

netconfd: $(GOFILES) statik/statik.go deps
	$(AT) CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o $@
	$(AT) strip $@



netconfd.debug: $(GOFILES) statik/statik.go deps
	$(AT) CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -a -ldflags '-extldflags "-static"' -o $@

test:
	go test ./crm -v

statik/statik.go: swaggerui/swagger.yaml
	@ cd vendor/github.com/rakyll/statik && go build
	@ vendor/github.com/rakyll/statik/statik -src=swaggerui -include=*.png,*.yaml,*.html,*.css,*.js

deps: statik/statik.go
	$(AT) go get -d -v

clean:
	@ rm -rf netconfd netconfd.debug statik/

.PHONY: clean debugbuild test deps
