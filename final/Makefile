all:
	git submodule update --init --recursive
	GOPATH=${PWD} go install ...go-bindata
	cd src/sm/final/assets/ && make
	GOPATH=${PWD} go install sm/...
