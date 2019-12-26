package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/phachon/mm-wiki/app/models"
	"github.com/phachon/mm-wiki/app/utils"
	"github.com/astaxie/beego"
)

type AuthorController struct {
	BaseController
}

// login index
func (this *AuthorController) Index() {

	// is open auth login
	ssoOpen := "0"
	config, err := models.ConfigModel.GetConfigByKey(models.Config_Key_AuthLogin)
	if err == nil && len(config) > 0 && config["value"] == "1" {
		ssoOpen = "1"
	}

	this.Data["sso_open"] = ssoOpen
	this.viewLayout("author/login", "author")
}

// login
func (this *AuthorController) Login() {

	if !this.IsPost() {
		this.ViewError("请求方式有误！")
	}
	username := strings.TrimSpace(this.GetString("username"))
	password := strings.TrimSpace(this.GetString("password"))

	if username == "" {
		this.jsonError("系统用户名不能为空！")
	}
	if strings.Contains(username, "_") {
		this.jsonError("系统用户名不合法！")
	}
	if password == "" {
		this.jsonError("密码不能为空！")
	}
	verifyOver := false
	if strings.Contains(username,"@"){
		usernamesplit := strings.Split(username, "@")
		userdomain := usernamesplit[1]
		// get auth login config
		authLogin, err := models.AuthModel.GetAuthByName(userdomain)
		if err != nil || len(authLogin) == 0 {
			this.jsonError("没有对应的认证配置！")
		}
		if authLogin["is_used"] != "1"{
			this.jsonError("对应的认证没有启用！")
		}
		var authData map[string]interface{}
		if strings.HasPrefix(authLogin["url"],"ldap://") || strings.HasPrefix(authLogin["url"],"ldaps://") {
			var success bool
			success, authData, err = utils.NewLDAP().VerifyLDAP(authLogin["url"], authLogin["ext_data"], username, password)
			if err != nil || !success {
				this.jsonError("LDAP认证失败！")
			}
		}else if strings.HasPrefix(authLogin["url"],"http://") || strings.HasPrefix(authLogin["url"],"https://") {
			queryValue := map[string]string{
				"username": username,
				"password": password,
				"ext_data": authLogin["ext_data"],
			}
			// request auth login api
			body, code, err := utils.Request.HttpPost(authLogin["url"], queryValue, nil)
			if err != nil {
				this.jsonError("登录认证失败：" + err.Error())
			}
			if len(body) == 0 {
				this.jsonError("登录认证失败：" + fmt.Sprintf("%d", code))
			}
			v := map[string]interface{}{}
			err = json.Unmarshal(body, &v)
			if err != nil {
				this.jsonError("登录认证失败!" + err.Error())
			}

			if v["message"].(string) != "" {
				this.jsonError("登录失败：" + v["message"].(string))
			}
			authData = v["data"].(map[string]interface{})
		}else{
			this.jsonError("认证配置错误!")
		}

		passwordEncode := models.UserModel.EncodePassword(password)
		userValue := map[string]interface{}{
			"username":   username,
			"given_name": authData["given_name"],
			"password":   passwordEncode,
			"email":      authData["email"],
			"mobile":     authData["mobile"],
			"phone":      authData["phone"],
			"department": authData["department"],
			"position":   authData["position"],
			"location":   authData["location"],
			"im":         authData["im"],
			"last_time":  time.Now().Unix(),
			"last_ip":    this.GetClientIp(),
		}
		ok, err := models.UserModel.HasUsername(username)
		if err != nil {
			this.jsonError("登录失败!")
		}
		if ok {
			// update user info
			_, err = models.UserModel.UpdateUserByUsername(userValue)
		} else {
			// insert user info
			userValue["role_id"] = models.Role_Default_Id
			_, err = models.UserModel.Insert(userValue)
		}
		if err != nil {
			this.jsonError("登录失败！" + err.Error())
		}
		verifyOver = true
	}

	user, err := models.UserModel.GetUserByUsername(username)
	if err != nil {
		this.jsonError("登录出错")
	}
	if len(user) == 0 {
		this.jsonError("用户名或密码错误!")
	}
	if user["is_forbidden"] == fmt.Sprintf("%d", models.User_Forbidden_True) {
		this.jsonError("用户已被禁用!")
	}

	password = utils.Encrypt.Md5Encode(password)
	if !verifyOver {
		if user["password"] != password {
			this.jsonError("用户名或密码错误!")
		}
	}

	// update last_ip and last_login_time
	updateValue := map[string]interface{}{
		"last_time": time.Now().Unix(),
		"last_ip":   this.GetClientIp(),
	}
	_, err = models.UserModel.Update(user["user_id"], updateValue)
	if err != nil {
		this.jsonError("登录出错")
	}

	// save session
	this.SetSession("author", user)
	// save cookie
	identify := utils.Encrypt.Md5Encode(this.Ctx.Request.UserAgent() + this.GetClientIp() + password)
	passportValue := utils.Encrypt.Base64Encode(username + "@@" + identify)
	passport := beego.AppConfig.String("author::passport")
	cookieExpired, _ := beego.AppConfig.Int64("author::cookie_expired")
	this.Ctx.SetCookie(passport, passportValue, cookieExpired)

	this.Ctx.Request.PostForm.Del("password")

	this.InfoLog("登录成功")
	this.jsonSuccess("登录成功！", nil, "/main/index")
}


//logout
func (this *AuthorController) Logout() {
	this.InfoLog("退出成功")
	passport := beego.AppConfig.String("author::passport")
	this.Ctx.SetCookie(passport, "")
	this.SetSession("author", nil)
	this.DelSession("author")

	this.Redirect("/", 302)
}
