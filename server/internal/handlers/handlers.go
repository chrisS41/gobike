package handlers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Users  *UserHandler
	Routes *RouteHandler
	Rides  *RideHandler
}

// 파라미터 파싱 헬퍼 함수
func parseParams[T any](c *gin.Context, params ...string) (*T, error) {
	var data T
	contentType := c.GetHeader("Content-Type")

	// 디버깅을 위한 로그 추가
	fmt.Printf("받은 데이터: %+v\n", c.Request.Body)

	switch contentType {
	case "application/json", "": // Content-Type이 없는 경우도 JSON으로 처리
		if err := c.ShouldBindJSON(&data); err != nil {
			return nil, fmt.Errorf("json 파싱 실패: %v", err)
		}
	case "application/x-www-form-urlencoded":
		if err := c.ShouldBind(&data); err != nil {
			return nil, fmt.Errorf("form 데이터 파싱 실패: %v", err)
		}
	default:
		return nil, fmt.Errorf("지원하지 않는 content-type: %s", contentType)
	}

	// 디버깅을 위한 로그 추가
	fmt.Printf("파싱된 데이터: %+v\n", data)

	v := reflect.ValueOf(&data).Elem()
	t := v.Type()

	// 구조체의 모든 필드 정보 로깅
	for i := 0; i < t.NumField(); i++ {
		fmt.Printf("구조체 필드: %s, JSON 태그: %s\n", t.Field(i).Name, t.Field(i).Tag.Get("json"))
	}

	for _, param := range params {
		field := v.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, param)
		})
		if !field.IsValid() {
			return nil, fmt.Errorf("필드 이름 %s를 찾을 수 없습니다", param)
		}
		if field.IsZero() {
			return nil, fmt.Errorf("필수 파라미터 %s가 누락되었습니다", param)
		}
	}

	return &data, nil
}
