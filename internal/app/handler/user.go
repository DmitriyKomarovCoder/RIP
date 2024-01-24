package handler

import (
	"RIP/internal/app/ds"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SignUp godoc
// @Summary      Sign up a new user
// @Description  Creates a new user account
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body  ds.UserSignUp  true  "User information"
// @Success      201  {object}  map[string]any
// @Failure      400  {object}  error
// @Failure      409  {object}  error
// @Failure      500  {object}  error
// @Router       /api/user/signUp [post]
func (h *Handler) SignUp(c *gin.Context) {
	var newClient ds.UserSignUp
	var err error

	if err = c.BindJSON(&newClient); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных о новом пользователе"})
		return
	}

	if newClient.Password, err = h.Hasher.Hash(newClient.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "неверный формат пароля"})
		return
	}

	if err = h.Repository.SignUp(c.Request.Context(), ds.User{
		Login:    newClient.Login,
		Name:     newClient.Name,
		Password: newClient.Password,
	}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "нельзя создать пользователя с таким логином"})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "пользователь успешно создан"})
}

// SignIn godoc
// @Summary      User sign-in
// @Description  Authenticates a user and generates an access token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body  ds.UserLogin  true  "User information"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Failure      401  {object}  error
// @Failure      500  {object}  error
// @Router       /api/user/signIn [post]
func (h *Handler) SignIn(c *gin.Context) {
	var clientInfo ds.UserLogin
	var err error

	if err = c.BindJSON(&clientInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "неверный формат данных")
		return
	}

	if clientInfo.Password, err = h.Hasher.Hash(clientInfo.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "неверный формат пароля"})
		return
	}

	user, err := h.Repository.GetByCredentials(c.Request.Context(), ds.User{Password: clientInfo.Password, Login: clientInfo.Login})
	if err != nil {
		if errors.Is(err, ds.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка авторизации"})
		return
	}

	token, err := h.TokenManager.NewJWT(int(user.UserId), user.IsAdmin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка при формировании токена"})
		return
	}

	c.SetCookie("AccessToken", "Bearer "+token, 0, "/", "127.0.0.1:8080", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "клиент успешно авторизован"})
}

// Logout godoc
// @Summary      Logout
// @Description  Logs out the user by blacklisting the access token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Router       /api/user/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	jwtStr, err := c.Cookie("AccessToken")
	if !strings.HasPrefix(jwtStr, jwtPrefix) || err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	jwtStr = jwtStr[len(jwtPrefix):]

	_, _, err = h.TokenManager.Parse(jwtStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)
		return
	}

	err = h.Redis.WriteJWTToBlacklist(c.Request.Context(), jwtStr, time.Hour)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
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

	request.Token = ServerToken

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

	if resp.StatusCode == 200 {
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
