package object

import (
	"strings"
)

type Provider struct {
	Owner        string `xorm:"owner varchar(100) notnull pk" json:"owner"`
	Name         string `xorm:"name varchar(100) notnull pk unique" json:"name"`
	CreatedAt    string `xorm:"created_at datetime notnull" json:"createdAt"`
	Type         string `xorm:"type varchar(100)" json:"type"`
	ClientId     string `xorm:"client_id varchar(100) notnull" json:"clientId"`
	ClientSecret string `xorm:"client_secret varchar(2000) notnull" json:"clientSecret"`
	SignName     string `xorm:"sign_name varchar(100)" json:"signName"`
	TemplateCode string `xorm:"template_code varchar(100)" json:"templateCode"`
}

func GetSmsProvider(id string) (*Provider, bool, error) {
	owner, name, _ := strings.Cut(id, "/")
	provider := Provider{Owner: owner, Name: name}
	existed, err := engine.Get(&provider)
	return &provider, existed, err
}

func UpdateSmsProvider(provider *Provider) (rows int64, err error) {
	if provider.Owner == "" {
		provider.Owner = "fireboom"
	}

	return engine.Where("owner=? and name=?", provider.Owner, provider.Name).Update(provider)
}
