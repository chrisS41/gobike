package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Route struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`  // MongoDB의 기본 ID 필드
	Name          string             `bson:"name"`           // 경로의 이름
	Description   string             `bson:"description"`    // 경로에 대한 설명
	Distance      float64            `bson:"distance"`       // 경로의 총 거리 (킬로미터)
	Duration      time.Duration      `bson:"duration"`       // 예상 소요 시간
	Difficulty    string             `bson:"difficulty"`     // 난이도 (예: 쉬움, 보통, 어려움)
	GPXData       string             `bson:"gpx_data"`       // GPX 형식의 경로 데이터
	CreatedAt     time.Time          `bson:"created_at"`     // 경로 생성 일시
	UpdatedAt     time.Time          `bson:"updated_at"`     // 경로 수정 일시
	StartPoint    GeoPoint           `bson:"start_point"`    // 경로 시작 지점 (위도, 경도)
	EndPoint      GeoPoint           `bson:"end_point"`      // 경로 종료 지점 (위도, 경도)
	ElevationGain float64            `bson:"elevation_gain"` // 총 상승 고도 (미터)
	Tags          []string           `bson:"tags"`           // 경로와 관련된 태그
}

type GeoPoint struct {
	Latitude  float64 `bson:"latitude"`  // 위도
	Longitude float64 `bson:"longitude"` // 경도
}
