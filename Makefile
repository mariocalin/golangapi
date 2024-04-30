MAIN_PACKAGE_PATH := ./cmd/server
MAIN_SERVER_FILE := api.go
BINARY_NAME := server
API_DOCS=./docs

# ----
# Run
# ----
run: build run/bin

run/bin:
	bin/${BINARY_NAME}

# ----
# Build
# ----

build: build/pre build/do

build/pre:
	mkdir -p bin

build/do:
	go build -o bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}/${MAIN_SERVER_FILE}


# ---
# TEST
# ---
test/only:
	go test -v ./...

test/cover:
	mkdir -p tmp
	go test -v -coverprofile=tmp/coverage.out ./...
	go tool cover -func=tmp/coverage.out

# ---
# UTIL
# ---
api-spec:
	swag init -g ${MAIN_PACKAGE_PATH}/${MAIN_SERVER_FILE} -o ${API_DOCS}
