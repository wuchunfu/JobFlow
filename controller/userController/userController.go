package userController

import (
	"gin-vue/common"
	"gin-vue/middleware/jwt"
	"gin-vue/model/loginLogModel"
	"gin-vue/model/tokenModel"
	"gin-vue/model/userModel"
	"gin-vue/utils"
	"gin-vue/utils/datetimeUtils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

//func Register(ctx *gin.Context) {
//	DB := database.GetDB()
//	var requestUser = &userModel.User{}
//	ctx.Bind(requestUser)
//	// 获取参数
//	name := requestUser.Name
//	telephone := requestUser.Telephone
//	password := requestUser.Password
//	response := jsonUtils.JsonResponse{}
//	// 数据验证
//	if len(telephone) != 11 {
//		response.Failure(http.StatusUnprocessableEntity, "手机号必须为11位")
//		return
//	}
//	if len(password) < 6 {
//		response.Failure(http.StatusUnprocessableEntity, "密码不能少于6位")
//		return
//	}
//	// 如果名称为空给一个随机字符串
//	if len(name) == 0 {
//		name = utils.RandomString(10)
//	}
//	if isTelephoneExist(DB, telephone) {
//		response.Failure(http.StatusUnprocessableEntity, "用户已存在")
//		return
//	}
//	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//	if err != nil {
//		response.Failure(http.StatusInternalServerError, "加密错误")
//		return
//	}
//	newUser := &userModel.User{
//		Name:      name,
//		Telephone: telephone,
//		Password:  string(hashPassword),
//	}
//	DB.Create(newUser)
//	// 发送token
//	token, err := jwt.ReleaseToken(newUser)
//	if err != nil {
//		response.Failure(http.StatusInternalServerError, "系统异常")
//		log.Printf("token generate error:%v", err)
//		return
//	}
//	log.Printf("token: %s", token)
//	// 返回结果
//	response.Success("注册成功", map[string]interface{}{"token": token})
//}
//
//func Login(c *gin.Context) {
//	db := database.GetDB()
//	// 获取参数
//	telephone := c.PostForm("telephone")
//	password := c.PostForm("password")
//	response := jsonUtils.JsonResponse{}
//	// 数据验证
//	if len(telephone) != 11 {
//		response.Failure(http.StatusUnprocessableEntity, "手机号必须为11位")
//		return
//	}
//	if len(password) < 6 {
//		response.Failure(http.StatusUnprocessableEntity, "密码不能少于6位")
//		return
//	}
//	// 判断手机号是否存在
//	var user userModel.User
//	db.Where("telephone=?", telephone).First(&user)
//	if user.Id == 0 {
//		response.Failure(http.StatusUnprocessableEntity, "密码不能少于6位")
//		return
//	}
//	// 判断密码是否正确
//	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
//	if err != nil {
//		response.Failure(http.StatusBadRequest, "密码错误")
//		return
//	}
//	// 发送token
//	token, err := jwt.ReleaseToken(&user)
//	if err != nil {
//		response.Failure(http.StatusInternalServerError, "系统异常")
//	}
//	// 返回结果
//	response.Success("登陆成功", map[string]interface{}{"token": token})
//}
//
//func Info(ctx *gin.Context) {
//	user, _ := ctx.Get("user")
//	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(*userModel.User))}})
//}

func Query(ctx *gin.Context) {
	query, _ := ctx.GetQuery("query")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": "hello", "query": query})
}

func Param(ctx *gin.Context) {
	param := ctx.Param("param")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": "hello", "param": param})
}

//func isTelephoneExist(db *gorm.DB, telephone string) bool {
//	var user userModel.User
//	db.Where("telephone=?", telephone).First(&user)
//	if user.Id != 0 {
//		return true
//	}
//	return false
//}

// UserForm 用户表单
type UserForm struct {
	UserId      int64
	UserIds     []int64
	Username    string // 用户名
	Password    string // 密码
	NewPassword string // 确认密码
	Email       string // 邮箱
	IsAdmin     int8   // 是否是管理员 1:管理员 0:普通用户
	Status      int8   // 启用状态 1:启动 0:禁用
	CreateTime  string // 创建时间
	UpdateTime  string // 修改时间
}

func Login(ctx *gin.Context) {
	form := new(UserForm)
	ctx.Bind(form)
	username := strings.TrimSpace(form.Username)
	password := strings.TrimSpace(form.Password)
	user := new(userModel.User)
	usernameCount := user.IsExistsUsername(username)
	if usernameCount <= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "用户名不存在！",
		})
	} else {
		info := user.GetInfoByName(username)
		pwd := utils.Md5(password + user.Salt)
		if info.Password != pwd {
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadRequest,
				"data": nil,
				"msg":  "用户名或密码不正确！",
			})
		} else {
			loginLog := new(loginLogModel.LoginLog)
			loginLog.Username = user.Username
			loginLog.Ip = ctx.ClientIP()
			loginLog.CreateTime = datetimeUtils.FormatDateTime()
			loginLog.Save()
			tokenInfo := jwt.ReleaseToken(user.UserId, user.Username)
			// 保存 session
			session := sessions.Default(ctx)
			session.Set("userId", user.UserId)
			session.Set("username", user.Username)
			session.Save()
			// 保存 token
			tokens := new(tokenModel.UserToken)
			tokens.Save(user.UserId, tokenInfo)
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"data": tokenInfo,
				"msg":  "Token 获取成功！",
			})
		}
	}
}

func Logout(ctx *gin.Context) {
	// 获取 session
	session := sessions.Default(ctx)
	userId := session.Get("userId")
	username := session.Get("username")
	session.Delete("userId")
	session.Delete("username")
	session.Clear()
	// 重新生成 token
	tokenInfo := jwt.ReleaseToken(userId.(int64), username.(string))
	tokens := new(tokenModel.UserToken)
	tokens.Update(userId.(int64), tokenInfo)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": tokenInfo,
		"msg":  "Token 获取成功！",
	})
}

// Index 用户列表页
func Index(ctx *gin.Context) {
	page, pageSize := common.ParseQueryParams(ctx)
	username := ctx.Query("username")
	user := new(userModel.User)
	dataList, count := user.List(page, pageSize, username)
	ctx.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"data":     &dataList,
		"msg":      "获取数据成功！",
		"page":     page,
		"pageSize": pageSize,
		"total":    count,
	})
}

// Detail 用户详情
func Detail(ctx *gin.Context) {
	userId, _ := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	user := new(userModel.User)
	detail := user.Detail(userId)
	if detail.UserId >= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": &detail,
			"msg":  "获取数据成功！",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": &detail,
			"msg":  "获取数据失败！",
		})
	}
}

func Save(ctx *gin.Context) {
	form := new(UserForm)
	ctx.Bind(form)
	username := strings.TrimSpace(form.Username)
	password := strings.TrimSpace(form.Password)
	email := strings.TrimSpace(form.Email)
	isAdmin := form.IsAdmin
	status := form.Status

	user := new(userModel.User)
	usernameCount := user.IsExistsUsername(username)
	emailCount := user.IsExistsEmail(email)
	if usernameCount > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "用户名已存在！",
		})
	} else if emailCount > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "邮箱已存在！",
		})
	} else {
		user.Username = username
		user.Salt = utils.RandomString(6)
		user.Password = utils.Md5(password + user.Salt)
		user.Email = email
		user.IsAdmin = isAdmin
		user.Status = status
		user.CreateTime = datetimeUtils.FormatDateTime()
		user.UpdateTime = datetimeUtils.FormatDateTime()
		user.Save()
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": nil,
			"msg":  "保存成功！",
		})
	}
}

func Update(ctx *gin.Context) {
	form := new(UserForm)
	ctx.Bind(form)
	userId := form.UserId
	username := strings.TrimSpace(form.Username)
	email := strings.TrimSpace(form.Email)
	isAdmin := form.IsAdmin
	status := form.Status

	user := new(userModel.User)
	userMap := make(map[string]interface{})
	userMap["username"] = username
	userMap["email"] = email
	userMap["is_admin"] = isAdmin
	userMap["status"] = status
	userMap["update_time"] = datetimeUtils.FormatDateTime()
	user.Update(userId, userMap)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": nil,
		"msg":  "修改成功！",
	})
}

func ChangePassword(ctx *gin.Context) {
	form := new(UserForm)
	ctx.Bind(form)
	userId := form.UserId
	username := strings.TrimSpace(form.Username)
	password := strings.TrimSpace(form.Password)

	if password == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "请输入密码！",
		})
	} else {
		//salt := utils.RandomString(6)
		user := new(userModel.User)
		userMap := make(map[string]interface{})
		userMap["username"] = username
		userMap["password"] = utils.Md5(password + user.Salt)
		userMap["update_time"] = datetimeUtils.FormatDateTime()

		user.ChangePassword(userId, userMap)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": nil,
			"msg":  "密码修改成功！",
		})
	}
}

func ChangeLoginPassword(ctx *gin.Context) {
	form := new(UserForm)
	ctx.Bind(form)
	userId := form.UserId
	username := strings.TrimSpace(form.Username)
	oldPassword := strings.TrimSpace(form.Password)
	newPassword := strings.TrimSpace(form.NewPassword)

	user := new(userModel.User)
	detail := user.Detail(userId)
	password := utils.Md5(oldPassword + user.Salt)
	if detail.Password != password {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "原密码不正确！",
		})
	} else if oldPassword == "" || newPassword == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"data": nil,
			"msg":  "请输入密码！",
		})
	} else {
		userMap := make(map[string]interface{})
		userMap["username"] = username
		userMap["password"] = utils.Md5(newPassword + user.Salt)
		userMap["update_time"] = datetimeUtils.FormatDateTime()

		user.ChangePassword(userId, userMap)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": nil,
			"msg":  "密码修改成功！",
		})
	}
}

func Delete(ctx *gin.Context) {
	form := new(UserForm)
	ctx.Bind(form)
	user := new(userModel.User)
	for _, id := range form.UserIds {
		user.Delete(id)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": nil,
		"msg":  "删除成功！",
	})
}

//设置默认路由当访问一个错误网站时返回
func NotFound(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusNotFound,
		"data": nil,
		"msg":  "404 ,page not exists!",
	})
}
