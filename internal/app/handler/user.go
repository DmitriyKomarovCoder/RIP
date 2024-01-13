package handler

import (
	"RIP/internal/app/ds"
	"RIP/internal/app/role"
	"RIP/internal/app/utils"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Register godoc
// @Summary Регистрация пользователя
// @Description Регистрация нового пользователя.
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param request body ds.RegisterReq true "Детали регистрации"
// @Router /api/v3/users/sign_up [post]
func (h *Handler) Register(ctx *gin.Context) {
	req := &ds.RegisterReq{}

	err := json.NewDecoder(ctx.Request.Body).Decode(req)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if req.Password == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("pass is empty"))
		return
	}

	if req.Login == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("name is empty"))
		return
	}

	if err = h.Repository.Register(&ds.User{
		Name:     req.Name,
		Role:     role.Buyer,
		Login:    req.Login,
		Password: generateHashString(req.Password),
	}); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &ds.RegisterResp{
		Ok: true,
	})
}

// Login godoc
// @Summary Аутентификация пользователя
// @Description Вход нового пользователя.
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param request body ds.RegisterReq true "Детали входа"
// @Success 200 {object} ds.LoginSwaggerResp "Успешная аутентификация"
// @Failure 400 {object} errorResp "Неверный запрос"
// @Failure 401 {object} errorResp "Неверные учетные данные"
// @Failure 500 {object} errorResp "Внутренняя ошибка сервера"
// @Router /api/v3/users/login [post]
func (h *Handler) Login(ctx *gin.Context) {
	cfg := h.Config
	req := &ds.LoginReq{}

	if err := json.NewDecoder(ctx.Request.Body).Decode(req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	user, err := h.Repository.GetUserByLogin(req.Login)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if req.Login == user.Login && user.Password == generateHashString(req.Password) {
		token := jwt.NewWithClaims(cfg.JWT.SigningMethod, &ds.JWTClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(cfg.JWT.ExpiresIn).Unix(),
				IssuedAt:  time.Now().Unix(),
				Issuer:    "bitop-admin",
			},
			UserID: uint(user.UserId),
			Role:   user.Role,
		})

		if token == nil {
			h.errorHandler(ctx, http.StatusInternalServerError, errors.New("token is nil"))
			return
		}

		strToken, err := token.SignedString([]byte(cfg.JWT.Token))
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, errors.New("cannot create str token"))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"expires_in":   cfg.JWT.ExpiresIn,
			"access_token": strToken,
			"token_type":   "Bearer",
			"role":         user.Role,
			"name":         user.Name,
		})
		return
	}

	h.errorHandler(ctx, http.StatusBadRequest, errors.New("incorrect login or password"))
}

func (h *Handler) UsersList(ctx *gin.Context) {
	users, err := h.Repository.UsersList()
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "users", users)
}

// Logout godoc
// @Summary Выход пользователя
// @Description Завершение сеанса текущего пользователя.
// @Tags Пользователи
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {string} string "Успешный выход"
// @Failure 400 {object} errorResp "Неверный запрос"
// @Failure 401 {object} errorResp "Неверные учетные данные"
// @Failure 500 {object} errorResp "Внутренняя ошибка сервера"
// @Router /api/v3/users/logout [get]
func (h *Handler) Logout(ctx *gin.Context) {
	jwtStr := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	jwtStr = jwtStr[len(jwtPrefix):]

	_, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWT.Token), nil
	})
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Redis.WriteJWTToBlacklist(ctx.Request.Context(), jwtStr, h.Config.JWT.ExpiresIn)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// MARK: - Inner functions

func generateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

const (
	ServerToken = "qwerty"
	ServiceUrl  = "http://127.0.0.1:8000/addRequest/"
)

func (h *Handler) UserRequest(c *gin.Context) {

	var request ds.RequestAsyncService
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("неверный формат"))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		h.errorHandler(c, http.StatusUnauthorized, errors.New("user_id not found"))
		return
	}

	//request.Token = ServerToken
	//var err1 error
	//err1, request.RequestId = h.Repository.GetTenderByUser(userID.(uint))
	//if err1 != nil {
	//	h.errorHandler(c, http.StatusBadRequest, err1)
	//}

	err, reqT := h.Repository.GetTenderByID(userID.(uint), request.RequestId)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, errors.New("request not found"))
	}

	if reqT.Status != utils.Draft && reqT.Status != "сформирован" {
		h.errorHandler(c, http.StatusBadRequest, errors.New("нельзя менять завершенные и отклоненные заявки"))

	}

	body, _ := json.Marshal(request)

	client := &http.Client{}
	req, err := http.NewRequest("PUT", ServiceUrl, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	//err, _ := h.Repository.FormTenderRequestByIDAsynce(request.RequestId, userID.(uint))

	if resp.StatusCode == 200 {
		err, _ := h.Repository.FormTenderRequestByIDAsynce(request.RequestId, userID.(uint))
		if err != nil {
			h.errorHandler(c, http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "заявка принята в обработку"})
		return
	}
	c.AbortWithError(http.StatusInternalServerError, errors.New("заявка не принята в обработку"))
}

// ручка вызывается сервисом на python
func (h *Handler) FinishUserRequest(c *gin.Context) {
	var request ds.RequestAsyncService
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)
		return
	}

	// сохраняем в базу
	err := h.Repository.SaveRequest(request)
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "данные сохранены"})
}
