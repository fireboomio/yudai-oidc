package object

import (
	sender "github.com/casdoor/go-sms-sender"
	"strings"
)

func getSmsClient(provider *Provider) (sender.SmsClient, error) {
	var client sender.SmsClient
	var err error

	client, err = sender.NewSmsClient(provider.Type, provider.ClientId, provider.ClientSecret, provider.SignName, provider.TemplateCode, "")
	if err != nil {
		return nil, err
	}

	return client, nil
}

func SendSmsCode(provider *Provider, content string, phoneNumbers ...string) error {
	params := map[string]string{}
	params["code"] = content
	return SendSms(provider, params, phoneNumbers...)
}

func SendSms(provider *Provider, params map[string]string, phoneNumbers ...string) error {
	client, err := getSmsClient(provider)
	if err != nil {
		return err
	}

	for i, number := range phoneNumbers {
		phoneNumbers[i] = strings.TrimPrefix(number, "+86")
	}

	err = client.SendMessage(params, phoneNumbers...)
	return err
}
