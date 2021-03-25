go:
	cd api && bash -c  "go run main.go"
go-build:
	cd api && bash -c  "go build -o api"
run-build:
	cd api && ./api
dockerize:
	docker-compose -f docker-compose.yml up -d --build
react:
	cd client && yarn install && yarn start