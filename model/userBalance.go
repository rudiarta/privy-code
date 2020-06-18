package model

type BaseUserBalance struct {
	ID int `gorm:"column:id;AUTO_INCREMENT;PRIMARY_KEY"`
}

type UserBalance struct {
	BaseUserBalance
	UserID         int `gorm:"column:user_id;"`
	Balance        int `gorm:"column:balance;"`
	BalanceAchieve int `gorm:"column:balance_achieve;"`
}

func (UserBalance) TableName() string {
	return "user_balance"
}
