truss:
	truss proto/quotesvc.proto  --pbpkg github.com/cg917658910/fzkj-wallet/quote-service/proto --svcpkg github.com/cg917658910/fzkj-wallet/quote-service --svcout . -v 
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/quote-service ./cmd/cli/main.go
build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/quote-service.exe ./cmd/cli/main.go
