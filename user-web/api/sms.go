package api

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/exp/rand"
	"net/http"
	"shop-api/user-web/forms"
	"shop-api/user-web/global"
	"strings"
	"time"
)

func GenerateSMSCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(uint64(time.Now().UnixNano()))

	var build strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&build, "%d", numeric[rand.Intn(r)])
	}
	return build.String()
}

func SendSMS(ctx *gin.Context) {
	sendSMS := forms.SendSMSForm{}
	if err := ctx.ShouldBind(&sendSMS); err != nil {
		HandleValidatorError(ctx, err)
		return

	}

	//	阿里云短信服务
	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", global.ServerConfig.AliSMS.ApiKey, global.ServerConfig.AliSMS.ApiSecret)
	if err != nil {
		panic(err)
	}
	smsCode := GenerateSMSCode(6)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = sendSMS.Mobile                //手机号
	request.QueryParams["SignName"] = "Go物商城"                           //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = "SMS_467375091"               //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}" //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	response, err := client.ProcessCommonRequest(request)
	fmt.Print(client.DoAction(request, response))
	//  fmt.Print(response)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
	// 验证码保存至redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.Redis.Host, global.ServerConfig.Redis.Port),
		Password: global.ServerConfig.Redis.Password,
	})
	rdb.Set(context.Background(), sendSMS.Mobile, smsCode, time.Minute*5)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "发送成功",
	})

}
