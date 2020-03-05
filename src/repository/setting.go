package repository

import (
	"github.com/icowan/blog/src/repository/types"
	"github.com/jinzhu/gorm"
	"time"
)

type SettingKey string

const (
	SettingSiteDomain                 SettingKey = "site-domain"                   // 主站地址
	SettingSiteName                   SettingKey = "site-name"                     // 主站名称
	SettingSiteDomainBeian            SettingKey = "site-domain-beian"             // 主站备案号
	SettingSiteDescription            SettingKey = "site-description"              // 在head里的DESCRIPTION信息
	SettingSiteKeywords               SettingKey = "site-keywords"                 // 在head里的keywords信息
	SettingSiteIcon                   SettingKey = "site-icon"                     // 站点的ICON地址 todo: cdn地址或本地相对路径
	SettingViewTemplate               SettingKey = "site-view-template"            // 模板路径
	SettingGlobalFoobarCode           SettingKey = "site-global-foobar-code"       // 全局foobar代码 可以设置成百度统计啥的
	SettingGlobalHeaderCode           SettingKey = "site-global-header-code"       // 全局header代码 可以设置成百度统计啥的
	SettingGlobalDomainImage          SettingKey = "site-global-domain-image"      // 图片cdn地址
	SettingSiteAboutName              SettingKey = "site-about-name"               // 站点用户名
	SettingSiteAboutDesc              SettingKey = "site-about-desc"               // 站点用户介绍
	SettingSiteAboutContent           SettingKey = "site-about-content"            // 站点用户介绍
	SettingSiteAboutAvatar            SettingKey = "site-about-avatar"             // 站点用户头像
	SettingSiteUserGithub             SettingKey = "site-user-github"              // 站点用户github
	SettingSiteUserOccupation         SettingKey = "site-user-occupation"          // 站点用户职业
	SettingSiteUserCity               SettingKey = "site-user-city"                // 站点用户城市
	SettingSiteUserEmail              SettingKey = "site-user-email"               // 站点用户邮箱
	SettingSiteUserQQ                 SettingKey = "site-user-qq"                  // 站点用户QQ
	SettingSiteUserWechatImage        SettingKey = "site-user-wecaht-image"        // 站点用户微信二维码
	SettingSiteGitalkSetting          SettingKey = "site-gtialk-setting"           // gitalk的配置 // todo: 是否需要开启的功能呢？
	SettingSiteShareSetting           SettingKey = "site-share-setting"            // 分享给件的配置
	SettingSiteMediaUploadPath        SettingKey = "site-media-upload-path"        // 上传资源所存放的路径
	SettingWechatOfficialAccountName  SettingKey = "wechat-official-account-name"  // 微信公众号名称
	SettingWechatOfficialAccountDesc  SettingKey = "wechat-official-account-desc"  // 微信公众号描述
	SettingWechatOfficialAccountImage SettingKey = "wechat-official-account-image" // 微信公众号名称 todo: cdn地址或本地相对路径
)

func (s SettingKey) String() string {
	return string(s)
}

type SettingRepository interface {
	Add(key SettingKey, value, desc string) (err error)
	Delete(key SettingKey) (err error)
	Update(setting *types.Setting) (err error)
	List() (res []*types.Setting, err error)
	Find(key SettingKey) (res types.Setting, err error)
}

type setting struct {
	db *gorm.DB
}

func (s *setting) Find(key SettingKey) (res types.Setting, err error) {
	err = s.db.Where("`key` = ?", key).First(&res).Error
	return
}

func (s *setting) Add(key SettingKey, value, desc string) (err error) {
	return s.db.Save(&types.Setting{
		Key:         key.String(),
		Value:       value,
		Description: desc,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}).Error
}

func (s *setting) Delete(key SettingKey) (err error) {
	return s.db.Where("`key` = ?", key).Delete(&types.Setting{}).Error
}

func (s *setting) Update(setting *types.Setting) (err error) {
	return s.db.Model(setting).Where("`key` = ?", setting.Key).Update(setting).Error
}

func (s *setting) List() (res []*types.Setting, err error) {
	err = s.db.Find(&res).Error
	return
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &setting{db: db}
}
