release:
	@go fmt ./...
	@go build -ldflags '-s -w'
debug:
	@go fmt ./...
	@go build -gcflags='-l -N'
clean:
	@rm ./gopass
