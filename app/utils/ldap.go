package utils

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/go-ldap/ldap/v3"
)

type LDAP struct {

}

func NewLDAP() *LDAP {
	return &LDAP{}
}

func (this *LDAP) VerifyLDAP(url, base, username, password string) (success bool, data map[string]interface{}, err error) {
	lc, err := ldap.DialURL(url)
	if err != nil {
		logs.Error("连接 LDAP 失败 ->", err)
		return false, nil, errors.New("连接 LDAP 失败")
	}
	defer lc.Close()
	err = lc.Bind(username, password)
	if err != nil {
		logs.Error("绑定 LDAP 用户失败 ->", err)
		return false, nil, errors.New("绑定 LDAP 用户失败")
	}
	searchRequest := ldap.NewSearchRequest(
		base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		//修改objectClass通过配置文件获取值
		fmt.Sprintf("(&(objectClass=User)(userPrincipalName=%s))", username),
		[]string{"dn", "mail", "displayName", "telephoneNumber", "mobile", "department"},
		nil,
	)
	searchResult, err := lc.Search(searchRequest)
	if err != nil {
		logs.Error("查找 LDAP 用户失败 ->", err)
		return false, nil, errors.New("查找 LDAP 用户失败")
	}
	if len(searchResult.Entries) != 1 {
		logs.Error("查找 LDAP 用户失败 ->ErrLDAPUserNotFoundOrTooMany")
		return false, nil, errors.New("查找 LDAP 用户失败")
	}
	data = map[string]interface{}{
		"given_name": searchResult.Entries[0].GetAttributeValue("displayName"),
		"email":      searchResult.Entries[0].GetAttributeValue("mail"),
		"mobile":     searchResult.Entries[0].GetAttributeValue("mobile"),
		"phone":      searchResult.Entries[0].GetAttributeValue("telephoneNumber"),
		"department": searchResult.Entries[0].GetAttributeValue("department"),
		"position":   "",
		"location":   "",
		"im":         "",
	}
	return true, data, nil
}
