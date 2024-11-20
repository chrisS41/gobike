package errors

// Error codes
const (
	SUCCESS = 200

	// HTTP related errors (1000-1999)
	ErrInvalidMethod = 1001
	ErrPathNotFound  = 1002
	ErrMissingParams = 1003
	// Database errors (5000-5999)
	ErrDatabaseConn  = 5001
	ErrDatabaseQuery = 5002

	// Auth related errors (6000-6999)
	ErrInvalidToken          = 6001
	ErrTokenExpired          = 6002
	ErrUnauthorized          = 6003
	ErrFailedToGenerateToken = 6004

	// User related errors (7000-7999)
	ErrUserNotFound         = 7001
	ErrInvalidUserInput     = 7002
	ErrDuplicateEmail       = 7003
	ErrFailedToHashPassword = 7004
	ErrFailedToCreateUser   = 7005
	ErrFailedToAddFriend    = 7006

	// Route related errors (8000-8999)
	ErrRouteNotFound       = 8001
	ErrInvalidRoute        = 8002
	ErrFailedToCreateRoute = 8003
	ErrFailedToUpdateRoute = 8004
	ErrFailedToFetchRoutes = 8005

	// Ride related errors (9000-9999)
	ErrFailedToCreateRide = 9001
)

// GetErrorMessage returns predefined error message for error code
func GetErrorMessage(code int) string {
	switch code {
	case SUCCESS:
		return "성공"

	// HTTP errors
	case ErrInvalidMethod:
		return "메서드가 유효하지 않습니다"
	case ErrPathNotFound:
		return "경로를 찾을 수 없습니다"
	case ErrMissingParams:
		return "필수 파라미터가 누락되었습니다"

	// Database errors
	case ErrDatabaseConn:
		return "데이터베이스 연결 오류"
	case ErrDatabaseQuery:
		return "데이터베이스 쿼리 오류"

	// Auth errors
	case ErrInvalidToken:
		return "유효하지 않은 토큰입니다"
	case ErrTokenExpired:
		return "만료된 토큰입니다"
	case ErrUnauthorized:
		return "인증되지 않은 접근입니다"
	case ErrFailedToGenerateToken:
		return "토큰 생성에 실패했습니다"

	// User errors
	case ErrUserNotFound:
		return "사용자를 찾을 수 없습니다"
	case ErrInvalidUserInput:
		return "잘못된 사용자 입력입니다"
	case ErrDuplicateEmail:
		return "이미 존재하는 이메일입니다"
	case ErrFailedToHashPassword:
		return "비밀번호 해싱에 실패했습니다"
	case ErrFailedToCreateUser:
		return "사용자 생성에 실패했습니다"
	case ErrFailedToAddFriend:
		return "친구 추가에 실패했습니다"

	// Route errors
	case ErrRouteNotFound:
		return "경로를 찾을 수 없습니다"
	case ErrInvalidRoute:
		return "잘못된 경로 정보입니다"
	case ErrFailedToCreateRoute:
		return "경로 생성에 실패했습니다"
	case ErrFailedToUpdateRoute:
		return "경로 업데이트에 실패했습니다"
	case ErrFailedToFetchRoutes:
		return "경로 조회에 실패했습니다"

	// Ride errors
	case ErrFailedToCreateRide:
		return "라이드 생성에 실패했습니다"

	default:
		return "내부 서버 오류가 발생했습니다"
	}
}
