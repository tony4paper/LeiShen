GOBIN = bin
GOBUILD_LINUX = env GO111MODULE=on go build -o $(GOBIN)/
GOBUILD_WIN = env CGO_ENABLED=1 GOOS=windows GOARCH=amd64 GO111MODULE=on CC=x86_64-w64-mingw32-gcc go build -o $(GOBIN)/

leishen:
	$(GOBUILD_LINUX) ./cmd/leishen
	$(GOBUILD_WIN) ./cmd/leishen
	@echo "Done building."
