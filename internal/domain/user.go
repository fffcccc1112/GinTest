package domain

//	type User struct {
//		UserID   int64  `json:"user_id"`
//		UserName string `json:"user_name"`
//		Age      int8   `json:"age"`
//		Email    string `json:"email"`
//		Password string `json:"password"`
//	}
type User struct {
	// 自定义主键：替换 gorm.Model 的默认 ID
	UserID   string `gorm:"column:user_id;type:varchar(50);primaryKey;comment:用户ID（主键）"`
	Password string `gorm:"column:password;type:varchar(100);not null;comment:加密密码"`
	UserName string `gorm:"column:user_name;type:varchar(50);not null;comment:用户姓名"`
	Gender   int    `gorm:"column:gender;type:int;default:0;comment:性别（0-未知 1-男 2-女）"`
	Account  string `gorm:"column:account;type:varchar(50);not null;unique;comment:登录账号（唯一）"`
	Email    string `gorm:"column:email;type:varchar(100);unique;comment:用户邮箱（唯一）"`
	Role     int    `gorm:"column:role;type:int;default:1;comment:角色（1-普通用户 2-管理员）"`

	// 可选：保留 GORM 的时间字段（如果表有 created_at/updated_at）
	// CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null"`
	// UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null"`
}

func (u *User) TableName() string {
	return "user" // 替换为你的表名（比如 `user`）
}
