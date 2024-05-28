package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"shop-api/user-web/forms"
	"shop-api/user-web/global"
	"shop-api/user-web/global/response"
	"shop-api/user-web/middlewares"
	"shop-api/user-web/models"
	"shop-api/user-web/proto"
	"strconv"
	"strings"
	"time"
)

func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
			return
		}
	}
}

func HandleValidatorError(ctx *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func GetUserList(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	zap.S().Infof("访问用户：%d", claims.(*models.CustomClaims).ID)
	//defer userConn.Close()
	Page := ctx.DefaultQuery("page", "1")
	PageInt, _ := strconv.Atoi(Page)
	PageSize := ctx.DefaultQuery("pageSize", "10")
	PageSizeInt, _ := strconv.Atoi(PageSize)
	rsp, err := global.UserClient.GetUserList(context.Background(), &proto.PageInfo{
		Page:     int32(PageInt),
		PageSize: int32(PageSizeInt),
	})
	if err != nil {
		zap.S().Errorw("获取用户列表失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		//data := make(map[string]interface{})

		user := response.UserResponse{
			Id:       value.Id,
			Nickname: value.Nickname,
			Birthday: response.JsonTime(time.Unix(int64(value.Birthday), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}

		//data["id"] = value.Id
		//data["name"] = value.Nickname
		//data["birthday"] = value.Birthday
		//data["gender"] = value.Gender
		//data["mobile"] = value.Mobile

		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}

func LoginByPassword(ctx *gin.Context) {
	loginForm := forms.LoginForm{}
	if err := ctx.ShouldBind(&loginForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	// 暂时取消验证码验证
	//if store.Verify(loginForm.CaptchaId, loginForm.Captcha, true) {
	//	ctx.JSON(http.StatusBadRequest, gin.H{
	//		"captcha": "验证码错误",
	//	})
	//	return
	//}

	if rsp, err := global.UserClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: loginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, gin.H{
					"msg": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "登录失败，请联系管理员",
				})
			}
			return
		}
	} else {
		if pwdRsp, pwdErr := global.UserClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          loginForm.Password,
			EncryptedPassword: rsp.Password,
		}); pwdErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "登录失败",
			})
		} else {
			if pwdRsp.Result {
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.Nickname,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),              // 签名生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*7, // 过期时间 7 天
						Issuer:    "sean",                         //签名的发行者
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}
				ctx.JSON(http.StatusOK, gin.H{
					"id":        rsp.Id,
					"nickname":  rsp.Nickname,
					"token":     token,
					"expiresAt": claims.StandardClaims.ExpiresAt * 1000,
				})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "密码错误，请重新输入",
				})
			}

		}
	}

}

func Register(ctx *gin.Context) {
	registerForm := forms.RegisterForm{}
	if err := ctx.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.Redis.Host, global.ServerConfig.Redis.Port),
		Password: global.ServerConfig.Redis.Password,
	})
	value, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "未知异常，请联系管理员",
		})
		return
	} else {
		if value != registerForm.Code {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "验证码错误，请重新输入",
			})
			return
		}
	}

	user, _ := global.UserClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Nickname: "用户" + registerForm.Mobile,
		Password: registerForm.Password,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		return
	}

	if err != nil {
		zap.S().Errorw("新建用户失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.Nickname,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),              // 签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*7, // 过期时间 7 天
			Issuer:    "sean",                         //签名的发行者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":        user.Id,
		"nickname":  user.Nickname,
		"token":     token,
		"expiresAt": claims.StandardClaims.ExpiresAt * 1000,
	})

}
