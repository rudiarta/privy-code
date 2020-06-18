package model

type BaseUserBalanceHistory struct {
	ID int `gorm:"column:id;AUTO_INCREMENT;PRIMARY_KEY;"`
}

type UserBalanceHistory struct {
	BaseUserBalanceHistory
	UserBalanceID int    `gorm:"column:user_balance_id;"`
	BalanceBefore int    `gorm:"column:balance_before;"`
	BalanceAfter  int    `gorm:"column:balance_after;"`
	Activity      string `gorm:"column:activity"`
	Type          string `gorm:"column:type;"`
	IP            string `gorm:"column:ip;"`
	Location      string `gorm:"column:location;"`
	UserAgent     string `gorm:"column:user_agent;"`
	Author        string `gorm:"column:author;"`
}

func (UserBalanceHistory) TableName() string {
	return "user_balance_history"
}
