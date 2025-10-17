APP_EXECUTABLE=apostle
OUTPUT_DIR=./bin

build:
	GOARCH=amd64 GOOS=darwin go build -o ${OUTPUT_DIR}/${APP_EXECUTABLE}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ${OUTPUT_DIR}/${APP_EXECUTABLE}-linux main.go
	GOARCH=amd64 GOOS=windows go build -o ${OUTPUT_DIR}/${APP_EXECUTABLE}-windows main.go

run: build
	./${OUTPUT_DIR}/${APP_EXECUTABLE}

clean:
	go clean
	rm ${OUTPUT_DIR}/${APP_EXECUTABLE}-darwin
	rm ${OUTPUT_DIR}/${APP_EXECUTABLE}-linux
	rm ${OUTPUT_DIR}/${APP_EXECUTABLE}-windows