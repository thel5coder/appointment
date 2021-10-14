package middleware

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"profira-backend/helpers/functioncaller"
	jwtPkg "profira-backend/helpers/jwt"
	"profira-backend/helpers/logruslogger"
	"profira-backend/helpers/messages"
	"profira-backend/server/handlers"
	"profira-backend/usecase"
	"strings"
	"time"
)

type JwtMiddleware struct {
	*usecase.UcContract
}

//jwt middleware
func (jwtMiddleware JwtMiddleware) New(ctx *fiber.Ctx) (err error) {
	claims := &jwtPkg.CustomClaims{}
	handler := handlers.Handler{UseCaseContract: jwtMiddleware.UcContract}

	//check header is present or not
	header := ctx.Get("Authorization")
	if !strings.Contains(header, "Bearer") {
		logruslogger.Log(logruslogger.WarnLevel, messages.AuthHeaderNotPresent, functioncaller.PrintFuncName(), "middleware-jwt-checkHeader")
		return handler.SendResponseUnauthorized(ctx,errors.New(messages.AuthHeaderNotPresent))
	}

	//check claims and signing method
	token := strings.Replace(header, "Bearer ", "", -1)
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			logruslogger.Log(logruslogger.WarnLevel, messages.UnexpectedSigningMethod, functioncaller.PrintFuncName(), "middleware-jwt-checkSigningMethod")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secret := jwtMiddleware.JwtConfig.SigningKey
		return secret, nil
	})
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, messages.UnexpectedClaims, functioncaller.PrintFuncName(), "middleware-jwt-checkClaims")
		return handler.SendResponseUnauthorized(ctx,errors.New(messages.UnexpectedClaims))
	}

	//check token live time
	if claims.ExpiresAt < time.Now().Unix() {
		logruslogger.Log(logruslogger.WarnLevel, messages.ExpiredToken, functioncaller.PrintFuncName(), "middleware-jwt-checkTokenLiveTime")
		return handler.SendResponseUnauthorized(ctx,errors.New(messages.ExpiredToken))
	}

	//jwe roll back encrypted id
	jweRes, err := jwtMiddleware.Jwe.Rollback(claims.Id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, messages.Unauthorized, functioncaller.PrintFuncName(), "pkg-jwe-rollback")
		return handler.SendResponseUnauthorized(ctx, errors.New(messages.Unauthorized))
	}
	if jweRes == nil {
		logruslogger.Log(logruslogger.WarnLevel, messages.Unauthorized, functioncaller.PrintFuncName(), "pkg-jwe-resultNil")
		return handler.SendResponseUnauthorized(ctx,errors.New(messages.Unauthorized))
	}

	//set id to uce case contract
	claims.Id = fmt.Sprintf("%v", jweRes["id"])
	jwtMiddleware.UcContract.UserID = claims.Id

	userUc := usecase.UserUseCase{UcContract: jwtMiddleware.UcContract}
	user,err := userUc.ReadBy("u.id",jwtMiddleware.UcContract.UserID,"=")
	if err != nil {
		fmt.Println(err.Error())
		return handler.SendResponseUnauthorized(ctx, errors.New(messages.UserNotFound))
	}

	if user.Role.Slug == "admin" {
		staffUc := usecase.StaffUseCase{UcContract:jwtMiddleware.UcContract}
		staff,err := staffUc.ReadBy("s.user_id",jwtMiddleware.UserID,"=")
		if err != nil {
			return handler.SendResponseUnauthorized(ctx, errors.New(messages.UserNotFound))
		}
		jwtMiddleware.ClinicID = staff.Clinics[0].ID
	}

	return ctx.Next()
}
