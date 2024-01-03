package object

import (
	"strings"
)

type Provider struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk unique" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	Type         string `xorm:"varchar(100)" json:"type"`
	ClientId     string `xorm:"varchar(100)" json:"clientId"`
	ClientSecret string `xorm:"varchar(2000)" json:"clientSecret"`
	SignName     string `xorm:"varchar(100)" json:"signName"`
	TemplateCode string `xorm:"varchar(100)" json:"templateCode"`
}

func GetSmsProvider(id string) (*Provider, bool, error) {
	owner, name, _ := strings.Cut(id, "/")
	provider := Provider{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&provider)
	return &provider, existed, err
}

func UpdateSmsProvider(provider *Provider) (rows int64, err error) {
	if provider.Owner == "" {
		provider.Owner = "fireboom"
	}

	return adapter.Engine.Where("owner=? and name=?", provider.Owner, provider.Name).Update(provider)
}
