export GOOS=linux
export GOARCH=amd64
go build -o $1 $2
unset GOOS
unset GOARCH
