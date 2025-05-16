# Build the binary (we use linux AMD64)
go build 

# To cross-compile for Windows ARM64
GOOS=windows GOARCH=arm64 go build