package version

import (
	"fmt"
	"runtime"
)

var (
	Version   string                                               // 빌드 버전
	Revision  string                                               // 커밋 버전
	BuildDate string                                               // 빌드 날짜
	GoVersion = runtime.Version()                                  // Go 버전
	Compiler  = runtime.Compiler                                   // 컴파일러
	Platform  = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH) // 플랫폼
)
