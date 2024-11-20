DATE=`date "+%Y-%m-%d_%H:%M:%S"`
VERSION=$(git describe --tags --always --dirty)
REVISION=$(git rev-parse --short HEAD)
go build \
    -mod=vendor \
    -ldflags "-X github.com/chrisS41/gobike-server/internal/version.Version=$VERSION \
              -X github.com/chrisS41/gobike-server/internal/version.Revision=$REVISION \
              -X github.com/chrisS41/gobike-server/internal/version.BuildDate=$DATE" \
    -o gobike-server ./cmd/main.go