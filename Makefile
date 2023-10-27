build:
	@go build -o bin/api

run: build
	@./bin/api

test:
	@go test -v ./...

seed:
	@go run scripts/seed.go

image:
	echo "building docker image ..."
	@docker build -t api .
	echo "running API inside the container"
	@docker run -p 3000:3000 api