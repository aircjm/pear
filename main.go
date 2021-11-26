package main

import (
	"context"
	"fmt"

	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	im "github.com/larksuite/oapi-sdk-go/service/im/v1"
)

var conf *config.Config

func init() {
	// 企业自建应用的配置
	// AppID、AppSecret: "开发者后台" -> "凭证与基础信息" -> 应用凭证（App ID、App Secret）
	// EncryptKey、VerificationToken："开发者后台" -> "事件订阅" -> 事件订阅（Encrypt Key、Verification Token）
	// HelpDeskID、HelpDeskToken, 服务台 token：https://open.feishu.cn/document/ukTMukTMukTM/ugDOyYjL4gjM24CO4IjN
	// 更多介绍请看：Github->README.zh.md->如何构建应用配置（AppSettings）
	appSettings := core.NewInternalAppSettings(
		core.SetAppCredentials("AppID", "cli_9e9796203e2fd107"),    // 必需
		core.SetAppEventKey("VerificationToken", "EncryptKey"),     // 非必需，订阅事件、消息卡片时必需
		core.SetHelpDeskCredentials("HelpDeskID", "HelpDeskToken")) // 非必需，使用服务台API时必需

	// 当前访问的是飞书，使用默认的内存存储（app/tenant access token）、默认日志（Error级别）
	// 更多介绍请看：Github->README.zh.md->如何构建整体配置（Config）
	conf = core.NewConfig(core.DomainFeiShu, appSettings, core.SetLoggerLevel(core.LoggerLevelError))
}

func main() {
	imService := im.NewService(conf)
	coreCtx := core.WrapContext(context.Background())
	reqCall := imService.Messages.Create(coreCtx, &im.MessageCreateReqBody{
		ReceiveId: "g878f12d",
		Content:   `{"text":"<at user_id="g878f12d">Tom</at> test content"}`,
		MsgType:   "text",
	})
	reqCall.SetReceiveIdType("open_id")
	message, err := reqCall.Do()
	// 打印 request_id 方便 oncall 时排查问题
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(tools.Prettify(err))
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(tools.Prettify(message))
}
