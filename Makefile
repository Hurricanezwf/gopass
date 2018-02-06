release:
	@go fmt ./...
	@go build
debug:
	@go fmt ./...
	@go build -gcflags='-l -N'
clean:
	@rm ./gopass
