GOFLAGS=
TESTFLAGS=
BUILDFLAGS=

gamma: test
	go ${GOFLAGS} build ${BUILDFLAGS} -o gamma

.PHONY: 
test:
	cd sexpr/ && go ${GOFLAGS} test ${TESTFLAGS}
	cd parse/ && go ${GOFLAGS} test ${TESTFLAGS}
	cd interp/ && go ${GOFLAGS} test ${TESTFLAGS}
	go ${GOFLAGS} test ${TESTFLAGS}

.PHONY: fmt
fmt:
	cd sexpr/ && go ${GOFLAGS} fmt
	cd parse/ && go ${GOFLAGS} fmt
	cd interp/ && go ${GOFLAGS} fmt
	go ${GOFLAGS} fmt

.PHONY: vet
vet:
	cd sexpr/ && go ${GOFLAGS} vet
	cd parse/ && go ${GOFLAGS} vet
	cd interp/ && go ${GOFLAGS} vet
	go ${GOFLAGS} vet

clean:
	rm -f gamma
