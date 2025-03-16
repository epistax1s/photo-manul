package model

type User struct {
	ID       int64  `gorm:"column:id;type:bigserial;primaryKey"`
	ChatID   int64  `gorm:"column:chat_id;type:bigint;unique;not null"`
	Username string `gorm:"column:username;type:varchar(255);not null"`
	Name     string `gorm:"column:name;type:varchar(255);not null"`
}

const (
	UserTable          = "tg_users"
	UserChatIDColumn   = "chat_id"
	UserUsernameColumn = "username"
	UserNameColumn     = "name"
)

func (User) TableName() string {
	return UserTable
}
