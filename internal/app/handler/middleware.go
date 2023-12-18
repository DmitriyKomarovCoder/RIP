package handler

import (
	"RIP/internal/app/ds"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"strings"
)

const (
	jwtPrefix = "Bearer "
	userCtx   = "UserId"
	adminCtx  = "IsAdmin"
)

func (h *Handler) WithAuthCheck(assignedRoles []ds.Role) func(ctx *gin.Context) {
	return func(gCtx *gin.Context) {
		jwtStr, err := gCtx.Cookie("AccessToken")
		if err != nil {
			gCtx.AbortWithStatus(http.StatusForbidden) // отдаем что нет доступа
			return
		}

		if !strings.HasPrefix(jwtStr, jwtPrefix) { // если нет префикса то нас дурят!
			gCtx.AbortWithStatus(http.StatusForbidden) // отдаем что нет доступа
			return
		}

		jwtStr = jwtStr[len(jwtPrefix):]

		err = h.Redis.CheckJWTInBlacklist(gCtx.Request.Context(), jwtStr)
		if err == nil { // значит что токен в блеклисте
			gCtx.AbortWithStatus(http.StatusForbidden)

			return
		}
		if !errors.Is(err, redis.Nil) {
			gCtx.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		h.TokenManager.Parse(jwtStr)

		userId, isAdmin, err := h.TokenManager.Parse(jwtStr)

		if len(assignedRoles) == 1 {
			if !isAdmin && assignedRoles[0] == 1 || isAdmin && assignedRoles[0] == 0 {
				gCtx.AbortWithStatus(http.StatusForbidden)
				log.Printf("user %v is not admin", userId)
				return
			}
		}

		gCtx.Set(userCtx, userId)
		gCtx.Set(adminCtx, isAdmin)
		gCtx.Next()
	}
}
