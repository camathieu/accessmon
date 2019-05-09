all: build

###
# Run linters
###
lint:
	@FAIL=0 ;echo -n " - go fmt :" ; OUT=`gofmt -l . | grep -v ^vendor` ; \
	if [[ -z "$$OUT" ]]; then echo " OK" ; else echo " FAIL"; echo "$$OUT"; FAIL=1 ; fi ;\
	echo -n " - go vet :" ; OUT=`go vet ./...` ; \
	if [[ -z "$$OUT" ]]; then echo " OK" ; else echo " FAIL"; echo "$$OUT"; FAIL=1 ; fi ;\
	echo -n " - go lint :" ; OUT=`golint ./... | grep -v ^vendor` ; \
	if [[ -z "$$OUT" ]]; then echo " OK" ; else echo " FAIL"; echo "$$OUT"; FAIL=1 ; fi ;\
	test $$FAIL -eq 0

###
# Run tests
###
test:
	@GORACE="halt_on_error=1" go test -race -cover -p 1 -count=1 ./... 2>&1 | grep -v "no test files"; test $${PIPESTATUS[0]} -eq 0

###
# Build command line
###

build:
	@cd cmd && go build -o accessmon

###
# Build docker
###
docker: build
	docker build -t camathieu/accessmon .