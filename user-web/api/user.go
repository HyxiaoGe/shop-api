package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"shop-api/user-web/forms"
	"shop-api/user-web/global"
	"shop-api/user-web/global/response"
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

func HandlerGrpcErrorToHttp(err error, c *gin.Context) {
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

func HandlerValidatorError(ctx *gin.Context, err error) {
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
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserConfig.Host, global.ServerConfig.UserConfig.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorf("连接用户服务失败: %v", err)
		return
	}
	//defer userConn.Close()
	userClient := proto.NewUserClient(userConn)
	Page := ctx.DefaultQuery("page", "1")
	PageInt, _ := strconv.Atoi(Page)
	PageSize := ctx.DefaultQuery("pageSize", "10")
	PageSizeInt, _ := strconv.Atoi(PageSize)
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Page:     int32(PageInt),
		PageSize: int32(PageSizeInt),
	})
	if err != nil {
		zap.S().Errorw("获取用户列表失败", "msg", err.Error())
		HandlerGrpcErrorToHttp(err, ctx)
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
		HandlerValidatorError(ctx, err)
		return
	}

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserConfig.Host, global.ServerConfig.UserConfig.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorf("连接用户服务失败: %v", err)
		return
	}
	//defer userConn.Close()
	userClient := proto.NewUserClient(userConn)
	if rsp, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
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
		if pwdRsp, pwdErr := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          loginForm.Password,
			EncryptedPassword: rsp.Password,
		}); pwdErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "登录失败",
			})
		} else {
			if pwdRsp.Result {
				ctx.JSON(http.StatusOK, gin.H{
					"msg": "登录成功",
				})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "密码错误，请重新输入",
				})
			}

		}
	}

}
