package sms_provider

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/supabase/gotrue/internal/conf"
)

type AliyunSmsProvider struct {
	Config *conf.AliyunSmsProviderConfiguration
}

// Creates a SmsProvider with the AliyunSms Config
func NewAliyunSmsProvider(config conf.AliyunSmsProviderConfiguration) (SmsProvider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &AliyunSmsProvider{
		Config: &config,
	}, nil
}

func (t *AliyunSmsProvider) SendMessage(phone, message, channel, otp string) (string, error) {
	switch channel {
	case SMSProvider:
		return t.SendSms(phone, otp)
	default:
		return "", fmt.Errorf("channel type %q is not supported for AliyunSms", channel)
	}
}

func (t *AliyunSmsProvider) SendSms(phone, otp string) (string, error) {
	client, err := createClient(&t.Config.AccessKeyId, &t.Config.AccessKeySecret)
	if err != nil {
		return "", err
	}

	templateParam := fmt.Sprintf("{\"code\":\"%s\"}", otp)
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String(t.Config.SignName),
		TemplateCode:  tea.String(t.Config.TemplateCode),
		TemplateParam: tea.String(templateParam),
	}
	runtime := &util.RuntimeOptions{}
	result, err := client.SendSmsWithOptions(sendSmsRequest, runtime)
	if err != nil {
		return "", err
	}

	fmt.Printf("result: %s\n", result)
	if *result.Body.Code != "OK" {
		return "", fmt.Errorf("AliyunSms error: %s", *result.Body.Message)
	}

	return *result.Body.BizId, nil
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func createClient(accessKeyId *string, accessKeySecret *string) (result *dysmsapi20170525.Client, err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	result = &dysmsapi20170525.Client{}
	result, err = dysmsapi20170525.NewClient(config)
	return result, err
}
