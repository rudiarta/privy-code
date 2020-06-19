package model

import (
	"database/sql/driver"
	"encoding/json"
)

type BaseUserBalanceHistory struct {
	ID int `gorm:"column:id;AUTO_INCREMENT;PRIMARY_KEY;"`
}

type UserBalanceHistory struct {
	BaseUserBalanceHistory
	UserBalanceID int      `gorm:"column:user_balance_id;"`
	BalanceBefore int      `gorm:"column:balance_before;"`
	BalanceAfter  int      `gorm:"column:balance_after;"`
	Activity      string   `gorm:"column:activity"`
	Type          string   `gorm:"column:type;"`
	IP            string   `gorm:"column:ip;"`
	Location      Location `gorm:"column:location;"`
	UserAgent     string   `gorm:"column:user_agent;"`
	Author        string   `gorm:"column:author;"`
}

func (UserBalanceHistory) TableName() string {
	return "user_balance_history"
}

type Location struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	RegionName  string `json:"regionName"`
	Timezone    string `json:"timezone"`
	ISP         string `json:"isp"`
}

func (l Location) Value() (driver.Value, error) {
	location, err := json.Marshal(l)
	return location, err
}

func (l *Location) Scan(src interface{}) error {
	source, err := json.Marshal(src)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(source, &l); err != nil {
		return err
	}

	return nil
}
