package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rudiarta/privy-code/database"
	"github.com/rudiarta/privy-code/model"
)

func AddBalanceController(c *gin.Context) {
	userID, exist := c.Get("user_id")
	if !exist {
		c.JSON(500, gin.H{
			"message": "User not found",
		})

		return
	}
	stringID := fmt.Sprintf("%v", userID)
	idUser, _ := strconv.Atoi(stringID)

	db, err := database.InitDatabase()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Database error",
		})

		return
	}
	defer func() {
		db.Close()
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	oldDataUserBalance := model.UserBalance{}
	var count int
	if errResult := getCountAndOldUserBalance(&oldDataUserBalance, &count, db, idUser); errResult != nil {
		c.JSON(422, gin.H{
			"message": errResult.Error(),
		})
	}

	dataUserBalance := model.UserBalance{}
	dataUserBalance.UserID = idUser
	balance, _ := strconv.Atoi(c.PostForm("balance"))
	if balance < 1 {
		c.JSON(422, gin.H{
			"message": "balance < 1",
		})

		return
	}

	{
		var oldBalance int
		oldBalance = balance
		if count != 0 {
			oldBalance = oldDataUserBalance.Balance + balance
		}
		dataUserBalance.Balance = oldBalance
		dataUserBalance.BalanceAchieve = balance
	}

	if err := db.Create(&dataUserBalance).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Insert error",
		})

		return
	}

	dataUserBalanceHistory := model.UserBalanceHistory{}
	dataUserBalanceHistory.UserBalanceID = dataUserBalance.ID
	dataUserBalanceHistory.Type = "credit"
	dataUserBalanceHistory.Activity = `{"user_id":` + stringID + `,"activity":"credit","balance":` + c.PostForm("balance") + `}`
	dataUserBalanceHistory.Author = "user"
	{
		if count == 0 {
			dataUserBalanceHistory.BalanceBefore = 0
		} else {
			dataUserBalanceHistory.BalanceBefore = oldDataUserBalance.Balance
		}
	}

	{
		dataUserBalanceHistory.BalanceAfter = dataUserBalance.Balance

	}

	dataUserBalanceHistory.IP = ""
	dataUserBalanceHistory.Location = ""
	dataUserBalanceHistory.UserAgent = ""

	if err := db.Create(&dataUserBalanceHistory).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Insert error",
			"id":      dataUserBalance,
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "Success",
	})

	return
}

func TakeOutBalanceController(c *gin.Context) {
	userID, exist := c.Get("user_id")
	if !exist {
		c.JSON(500, gin.H{
			"message": "User not found",
		})

		return
	}
	stringID := fmt.Sprintf("%v", userID)
	idUser, _ := strconv.Atoi(stringID)

	db, err := database.InitDatabase()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Database error",
		})

		return
	}
	defer func() {
		db.Close()
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	oldDataUserBalance := model.UserBalance{}
	var count int
	if errResult := getCountAndOldUserBalance(&oldDataUserBalance, &count, db, idUser); errResult != nil {
		c.JSON(422, gin.H{
			"message": errResult.Error(),
		})
	}

	dataUserBalance := model.UserBalance{}
	dataUserBalance.UserID = idUser
	balance, _ := strconv.Atoi(c.PostForm("balance"))
	if balance > -1 {
		c.JSON(422, gin.H{
			"message": "balance > -1",
		})

		return
	}

	{
		var oldBalance int
		oldBalance = balance
		if count != 0 {
			oldBalance = oldDataUserBalance.Balance + balance
		}
		dataUserBalance.Balance = oldBalance
		dataUserBalance.BalanceAchieve = balance
	}

	if err := db.Create(&dataUserBalance).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Insert error",
		})

		return
	}

	dataUserBalanceHistory := model.UserBalanceHistory{}
	dataUserBalanceHistory.UserBalanceID = dataUserBalance.ID
	dataUserBalanceHistory.Type = "debit"
	dataUserBalanceHistory.Activity = `{"user_id":` + stringID + `,"activity":"debit","balance":` + c.PostForm("balance") + `}`
	dataUserBalanceHistory.Author = "user"
	{
		if count == 0 {
			dataUserBalanceHistory.BalanceBefore = 0
		} else {
			dataUserBalanceHistory.BalanceBefore = oldDataUserBalance.Balance
		}
	}

	{
		dataUserBalanceHistory.BalanceAfter = dataUserBalance.Balance

	}

	dataUserBalanceHistory.IP = ""
	dataUserBalanceHistory.Location = ""
	dataUserBalanceHistory.UserAgent = ""

	if err := db.Create(&dataUserBalanceHistory).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Insert error",
			"id":      dataUserBalance,
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "Success",
	})

	return
}

func getCountAndOldUserBalance(oldDataUserBalance *model.UserBalance, count *int, db *gorm.DB, idUser int) error {
	errCount := db.Table("user_balance").Where("user_id = ?", idUser).Count(count).Error
	if errCount != nil {
		return errCount
	}

	errGet := db.Order("created_at desc").Where("user_id = ?", idUser).First(oldDataUserBalance)
	if !errGet.RecordNotFound() && errGet.Error != nil {
		return errGet.Error
	}

	return nil
}
