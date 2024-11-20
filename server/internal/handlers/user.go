package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/chrisS41/gobike-server/internal/config"
	"github.com/chrisS41/gobike-server/internal/database"
	"github.com/chrisS41/gobike-server/internal/errors"
	"github.com/chrisS41/gobike-server/internal/logger"
	"github.com/chrisS41/gobike-server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	users *database.Collection
	log   *logger.Log
}

func NewUserHandler(users *database.Collection, log *logger.Log) *UserHandler {
	return &UserHandler{users: users, log: log}
}

// 회원가입
func (h *UserHandler) Register(c *gin.Context) {

	user, err := parseParams[models.User](c, "email", "password", "name")
	h.log.Info("회원가입 요청 파라미터", user)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewErrorResponseWithMessage(
				errors.ErrInvalidUserInput,
				err.Error(),
			),
		)
		return
	}

	// 비밀번호 해싱
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(errors.ErrFailedToHashPassword),
		)
		return
	}
	user.Password = string(hashedPassword)

	// 생성 시간 설정
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	user.ID, err = h.users.Create(user)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(errors.ErrFailedToCreateUser),
		)
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(user))
}

// 로그인
func (h *UserHandler) Login(c *gin.Context) {

	loginInput, err := parseParams[models.User](c, "email", "password")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewErrorResponseWithMessage(errors.ErrInvalidUserInput, err.Error()),
		)
		return
	}

	// 사용자 조회
	var user models.User
	err = h.users.ReadOne(bson.M{"email": loginInput.Email}, &user)
	if err != nil {
		c.JSON(
			http.StatusNotFound,
			models.NewErrorResponse(errors.ErrUserNotFound),
		)
		return
	}

	// 비밀번호 검증
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password)); err != nil {
		c.JSON(
			http.StatusUnauthorized,
			models.NewErrorResponse(errors.ErrUnauthorized),
		)
		return
	}

	// 마지막 로그인 시간 업데이트
	user.LastLoginAt = time.Now()
	if err := h.users.Update(
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"last_login_at": user.LastLoginAt}},
	); err != nil {
		log.Printf("Failed to update last login time: %v", err)
	}

	// JWT 토큰 생성
	token, err := h.generateJWT(user)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(errors.ErrFailedToGenerateToken),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		models.NewSuccessResponse(gin.H{
			"token":         token,
			"name":          user.Name,
			"role":          user.Role,
			"last_login_at": user.LastLoginAt,
		}),
	)
}

func (h *UserHandler) generateJWT(user models.User) (string, error) {
	cfg := config.GetConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(cfg.JWTSecret))
}

// 사용자 조회
func (h *UserHandler) GetUser(c *gin.Context) {
	// TODO: 사용자 조회 로직 구현
	c.JSON(http.StatusOK, gin.H{"message": "user details"})
}

// 사용자 업데이트
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// TODO: 사용자 업데이트 로직 구현
	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

// 친구 추가
func (h *UserHandler) AddFriend(c *gin.Context) {
	userID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var friendID primitive.ObjectID
	if err := c.ShouldBindJSON(&friendID); err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewErrorResponseWithMessage(errors.ErrInvalidUserInput, err.Error()),
		)
		return
	}

	err := h.users.Update(
		bson.M{"_id": userID},
		bson.M{"$addToSet": bson.M{"friends": friendID}},
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(errors.ErrFailedToAddFriend),
		)
		return
	}

	c.Status(http.StatusOK)
}

// 친구 조회
func (h *UserHandler) GetFriends(c *gin.Context) {
	// TODO: 친구 조회 로직 구현
	c.JSON(
		http.StatusOK,
		models.NewSuccessResponse("friends list"),
	)
}

func (h *UserHandler) UpdateSubscription(c *gin.Context) {
	// TODO: 구독 정보 업데이트 로직 구현
	c.JSON(
		http.StatusOK,
		models.NewSuccessResponse("subscription updated"),
	)
}

func (h *UserHandler) GetSubscription(c *gin.Context) {
	// TODO: 구독 정보 조회 로직 구현
	c.JSON(
		http.StatusOK,
		models.NewSuccessResponse("subscription details"),
	)
}
