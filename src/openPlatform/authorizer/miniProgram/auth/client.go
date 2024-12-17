package auth

import (
	"context"
	"fmt"
	"reflect"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform/auth"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform/authorizer/miniProgram/auth/response"
)

type Client struct {
	*kernel.BaseClient

	// PowerWechat\OpenPlatform\Application
	component kernel.ApplicationInterface
}

func NewClient(app kernel.ApplicationInterface, component kernel.ApplicationInterface) (*Client, error) {
	baseClient, err := kernel.NewBaseClient(&app, nil)
	if err != nil {
		return nil, err
	}
	return &Client{
		baseClient,
		component,
	}, nil
}

// 小程序登录
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/others/WeChat_login.html
func (comp *Client) Session(ctx context.Context, code string) (*response.ResponseSession, error) {

	result := &response.ResponseSession{}

	config := (*comp.App).GetConfig()
	componentConfig := comp.component.GetConfig()
	component := comp.component.GetComponent("AccessToken")

	// 打印类型和包路径信息
	fmt.Printf("Type: %T\n", component)
	fmt.Printf("Value: %v\n", component)
	if reflect.TypeOf(component).Kind() == reflect.Ptr {
		fmt.Printf("Type (reflect): %s\n", reflect.TypeOf(component).Elem().PkgPath())
	} else {
		fmt.Printf("Type (reflect): %s\n", reflect.TypeOf(component).PkgPath())
	}
	fmt.Printf("Type (reflect): %s\n", reflect.TypeOf(component).String())

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	token, ok := component.(*auth.AccessToken)
	if !ok {
		return nil, fmt.Errorf("failed to cast component to *auth.AccessToken, actual type: %T", component)
	}
	componentToken, err := token.GetToken(ctx, false)

	query := &object.StringMap{
		"appid":                  config.GetString("app_id", ""),
		"js_code":                code,
		"grant_type":             "authorization_code",
		"component_appid":        componentConfig.GetString("app_id", ""),
		"component_access_token": componentToken.ComponentAccessToken,
	}
	_, err = comp.BaseClient.HttpGet(ctx, "sns/component/jscode2session", query, nil, result)

	return result, err

}
