DATE=`date "+%Y-%m-%d_%H:%M:%S"`
VERSION=$(git describe --tags --always --dirty)
REVISION=$(git rev-parse --short HEAD)
go build \
    -mod=vendor \
    -ldflags "-X main.date=$DATE -X main.version=$VERSION -X main.revision=$REVISION" \
    -o gobike-server ./cmd/main.go