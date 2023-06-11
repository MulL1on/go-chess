package dao

type User struct {
	UserID   int64  `json:"userId" gorm:"column:user_id;primaryKey"`
	Username string `json:"username" gorm:"column:username;unique;not null"`
	Password string `json:"password" gorm:"column:password;not null"`
}
