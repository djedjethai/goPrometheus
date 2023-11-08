BINARY=binary

.PHONY: build
build:
	@echo "Building auth binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ./server/bin/${BINARY} ./server
	@echo "Done!"

