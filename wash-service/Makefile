truss:
	truss proto/washsvc.proto  --pbpkg github.com/cg917658910/fzkj-wallet/wash-service/proto --svcpkg github.com/cg917658910/fzkj-wallet/wash-service --svcout . -v 
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/wash-service ./cmd/cli/main.go
build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/wash-service.exe ./cmd/cli/main.go
