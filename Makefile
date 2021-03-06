.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/productinfo productinfo/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/addproduct addproduct/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getproducts getproducts/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/deleteproduct deleteproduct/main.go
clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
