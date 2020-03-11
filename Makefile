air_path=$(HOME)/go/bin/air

run-local:
	export GIN_ENV= && $(air_path)

run-dev:
	export GIN_ENV=development && $(air_path)

run:
	go build -o bin/main main.go
	export GIN_ENV=production && ./bin/main

compile:
	echo "Compiling for every OS and Platform"
	go build -o bin/main main.go
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 main.go
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 main.go
