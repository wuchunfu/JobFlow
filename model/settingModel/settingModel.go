package settingModel

import (
	"encoding/json"
	logger "github.com/sirupsen/logrus"
	"github.com/wuchunfu/JobFlow/middleware/database"
)

type Setting struct {
	Id    int    `gorm:"type:int(20); primary_key; auto_increment; not null"`
	Code  string `gorm:"type:varchar(100); not null"`
	Key   string `gorm:"type:varchar(100); not null"`
	Value string `gorm:"type:varchar(4096); not null; default '' "`
}

const slackTemplate = `
任务ID:  {{.TaskId}}
任务名称: {{.TaskName}}
状态:    {{.Status}}
执行结果: {{.Result}}
`
const emailTemplate = `
任务ID:  {{.TaskId}}
任务名称: {{.TaskName}}
状态:    {{.Status}}
执行结果: {{.Result}}
`
const webhookTemplate = `
{
  "task_id": "{{.TaskId}}",
  "task_name": "{{.TaskName}}",
  "status": "{{.Status}}",
  "result": "{{.Result}}"
}
`

const (
	SlackCode        = "slack"
	SlackUrlKey      = "url"
	SlackTemplateKey = "template"
	SlackChannelKey  = "channel"
)

const (
	MailCode        = "mail"
	MailTemplateKey = "template"
	MailServerKey   = "server"
	MailUserKey     = "user"
)

const (
	WebhookCode        = "webhook"
	WebhookTemplateKey = "template"
	WebhookUrlKey      = "url"
)

// 初始化基本字段 邮件、slack等
func (setting *Setting) InitBasicField() {
	db := database.GetDB()

	setting.Code = SlackCode
	setting.Key = SlackUrlKey
	setting.Value = ""
	db.Model(&setting).Create(&setting)

	setting.Id = 0
	setting.Code = SlackCode
	setting.Key = SlackTemplateKey
	setting.Value = slackTemplate
	db.Model(&setting).Create(&setting)

	setting.Id = 0
	setting.Code = MailCode
	setting.Key = MailServerKey
	setting.Value = ""
	db.Model(&setting).Create(&setting)

	setting.Id = 0
	setting.Code = MailCode
	setting.Key = MailTemplateKey
	setting.Value = emailTemplate
	db.Model(&setting).Create(&setting)

	setting.Id = 0
	setting.Code = WebhookCode
	setting.Key = WebhookTemplateKey
	setting.Value = webhookTemplate
	db.Model(&setting).Create(&setting)

	setting.Id = 0
	setting.Code = WebhookCode
	setting.Key = WebhookUrlKey
	setting.Value = ""
	db.Model(&setting).Create(&setting)
}

// region slack配置
type Slack struct {
	Url      string    `json:"url"`
	Channels []Channel `json:"channels"`
	Template string    `json:"template"`
}

type Channel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (setting *Setting) Slack() Slack {
	db := database.GetDB()
	list := make([]Setting, 0)
	err := db.Model(&setting).Where(&Setting{Code: SlackCode}).Find(&list)
	slack := Slack{}
	if err.Error != nil {
		logger.Error("<<<slack>>> failed to get slack configuration from database, error msg: %v", err.Error)
		return slack
	}

	setting.formatSlack(list, &slack)

	return slack
}

func (setting *Setting) formatSlack(list []Setting, slack *Slack) {
	for _, v := range list {
		switch v.Key {
		case SlackUrlKey:
			slack.Url = v.Value
		case SlackTemplateKey:
			slack.Template = v.Value
		default:
			slack.Channels = append(slack.Channels, Channel{
				v.Id, v.Value,
			})
		}
	}
}

func (setting *Setting) UpdateSlack(url, template string) {
	db := database.GetDB().Begin()

	urlMap := make(map[string]interface{})
	urlMap["code"] = SlackCode
	urlMap["key"] = SlackUrlKey
	urlErr := db.Model(&setting).Where(&Setting{Value: url}).Updates(urlMap)
	if urlErr.Error != nil {
		logger.Errorf("update slack fail, error msg: %v", urlErr.Error)
		db.Rollback()
	}
	db.Commit()

	templateMap := make(map[string]interface{})
	templateMap["code"] = SlackCode
	templateMap["key"] = SlackTemplateKey
	templateErr := db.Model(&setting).Where(&Setting{Value: template}).Updates(templateMap)
	if templateErr.Error != nil {
		logger.Errorf("update slack fail, error msg: %v", templateErr.Error)
		db.Rollback()
	}
	db.Commit()
}

// 创建slack渠道
func (setting *Setting) CreateChannel(channel string) *Setting {
	setting.Code = SlackCode
	setting.Key = SlackChannelKey
	setting.Value = channel

	db := database.GetDB()
	err := db.Model(&setting).Create(&setting)
	if err.Error != nil {
		logger.Errorf("create channel fail, error msg: %v", err.Error)
		return nil
	}
	return setting
}

func (setting *Setting) IsChannelExist(channel string) bool {
	setting.Code = SlackCode
	setting.Key = SlackChannelKey
	setting.Value = channel

	db := database.GetDB()
	var count int64
	err := db.Model(&setting).Find(&setting).Count(&count)
	if err.Error != nil && count <= 0 {
		return false
	} else {
		return count > 0
	}
}

// 删除slack渠道
func (setting *Setting) RemoveChannel(id int) {
	setting.Code = SlackCode
	setting.Key = SlackChannelKey
	setting.Id = id

	db := database.GetDB().Begin()
	err := db.Model(&setting).Delete(&setting)
	if err.Error != nil {
		logger.Errorf("delete channel fail, error msg: %v", err.Error)
		db.Rollback()
	}
	db.Commit()
}

type Mail struct {
	Host      string     `json:"host"`
	Port      int        `json:"port"`
	User      string     `json:"user"`
	Password  string     `json:"password"`
	MailUsers []MailUser `json:"mailUsers"`
	Template  string     `json:"template"`
}

type MailUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// region 邮件配置
func (setting *Setting) Mail() Mail {
	db := database.GetDB()
	list := make([]Setting, 0)
	err := db.Model(&setting).Where(&Setting{Code: MailCode}).Find(&list)
	mail := Mail{MailUsers: make([]MailUser, 0)}
	if err.Error != nil {
		logger.Errorf("<<<mail>>> failed to get mail configuration from database, error msg: %v", err.Error)
		return mail
	}

	setting.formatMail(list, &mail)

	return mail
}

func (setting *Setting) formatMail(list []Setting, mail *Mail) {
	mailUser := MailUser{}
	for _, v := range list {
		switch v.Key {
		case MailServerKey:
			json.Unmarshal([]byte(v.Value), mail)
		case MailUserKey:
			json.Unmarshal([]byte(v.Value), &mailUser)
			mailUser.Id = v.Id
			mail.MailUsers = append(mail.MailUsers, mailUser)
		case MailTemplateKey:
			mail.Template = v.Value
		}
	}
}

func (setting *Setting) UpdateMail(config, template string) {
	db := database.GetDB().Begin()

	configMap := make(map[string]interface{})
	configMap["code"] = MailCode
	configMap["key"] = MailServerKey
	configErr := db.Model(&setting).Where(&Setting{Value: config}).Updates(configMap)
	if configErr.Error != nil {
		logger.Errorf("update mail fail, error msg: %v", configErr.Error)
		db.Rollback()
	}
	db.Commit()

	templateMap := make(map[string]interface{})
	templateMap["code"] = MailCode
	templateMap["key"] = MailTemplateKey
	templateErr := db.Model(&setting).Where(&Setting{Value: template}).Updates(templateMap)
	if templateErr.Error != nil {
		logger.Errorf("update mail fail, error msg: %v", templateErr.Error)
		db.Rollback()
	}
	db.Commit()
}

func (setting *Setting) CreateMailUser(username, email string) {
	setting.Code = MailCode
	setting.Key = MailUserKey
	mailUser := MailUser{0, username, email}
	jsonByte, err := json.Marshal(mailUser)
	if err != nil {
		logger.Errorf("create mail user marshal fail, error msg: %v", err.Error())
	}
	setting.Value = string(jsonByte)

	db := database.GetDB()
	createErr := db.Model(&setting).Create(&setting)
	if createErr.Error != nil {
		logger.Errorf("create mail user fail, error msg: %v", createErr.Error)
	}
}

func (setting *Setting) RemoveMailUser(id int) {
	setting.Code = MailCode
	setting.Key = MailUserKey
	setting.Id = id

	db := database.GetDB().Begin()
	err := db.Model(&setting).Delete(&setting)
	if err.Error != nil {
		logger.Errorf("delete mail user fail, error msg: %v", err.Error)
		db.Rollback()
	}
	db.Commit()
}

type WebHook struct {
	Url      string `json:"url"`
	Template string `json:"template"`
}

func (setting *Setting) Webhook() WebHook {
	list := make([]Setting, 0)
	webHook := WebHook{}

	db := database.GetDB()
	err := db.Model(&setting).Where(&Setting{Code: WebhookCode}).Find(&list)
	if err.Error != nil {
		logger.Errorf("<<<webhook>>> failed to get webhook configuration from database, error msg: %v", err.Error)
		return webHook
	}

	setting.formatWebhook(list, &webHook)

	return webHook
}

func (setting *Setting) formatWebhook(list []Setting, webHook *WebHook) {
	for _, v := range list {
		switch v.Key {
		case WebhookUrlKey:
			webHook.Url = v.Value
		case WebhookTemplateKey:
			webHook.Template = v.Value
		}
	}
}

func (setting *Setting) UpdateWebHook(url, template string) {
	db := database.GetDB().Begin()
	urlMap := make(map[string]interface{})
	urlMap["code"] = WebhookCode
	urlMap["key"] = WebhookUrlKey
	urlErr := db.Model(&setting).Where(&Setting{Value: url}).Updates(urlMap)
	if urlErr.Error != nil {
		logger.Errorf("update webhook fail, error msg: %v", urlErr.Error)
		db.Rollback()
	}
	db.Commit()

	templateMap := make(map[string]interface{})
	templateMap["code"] = WebhookCode
	templateMap["key"] = WebhookTemplateKey
	templateErr := db.Model(&setting).Where(&Setting{Value: template}).Updates(templateMap)
	if templateErr.Error != nil {
		logger.Errorf("update webhook fail, error msg: %v", templateErr.Error)
		db.Rollback()
	}
	db.Commit()
}
