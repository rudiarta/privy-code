package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rudiarta/privy-code/database"
	"github.com/rudiarta/privy-code/model"
)

func AddBalanceController(c *gin.Context) {

	//Get user agent
	validate := c.Request.Header["User-Agent"]
	if len(validate) == 0 {
		c.JSON(422, gin.H{
			"message": "User agent doesn't exist",
		})

		return
	}
	userAgent := validate[0]
	ipAddr := c.ClientIP()

	id, exist := c.Get("user_id")
	if !exist {
		c.JSON(422, gin.H{
			"message": "user id not found.",
		})
	}
	stringID := fmt.Sprintf("%v", id)
	userID, _ := strconv.Atoi(stringID)
	balance, _ := strconv.Atoi(c.PostForm("balance"))
	var balanceType model.BalanceType
	if err := balanceType.Init("credit"); err != nil {
		c.JSON(422, gin.H{
			"message": err,
		})
	}

	result, err := balanceFlow(userID, balance, userAgent, ipAddr, balanceType)
	if err != nil {
		c.JSON(422, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": result,
	})

	return
}

func TakeOutBalanceController(c *gin.Context) {

	//Get user agent
	validate := c.Request.Header["User-Agent"]
	if len(validate) == 0 {
		c.JSON(422, gin.H{
			"message": "User agent doesn't exist",
		})

		return
	}
	userAgent := validate[0]
	ipAddr := c.ClientIP()

	id, exist := c.Get("user_id")
	if !exist {
		c.JSON(422, gin.H{
			"message": "user id not found.",
		})
	}
	stringID := fmt.Sprintf("%v", id)
	userID, _ := strconv.Atoi(stringID)
	balance, _ := strconv.Atoi(c.PostForm("balance"))
	var balanceType model.BalanceType
	if err := balanceType.Init("debit"); err != nil {
		c.JSON(422, gin.H{
			"message": err,
		})
	}

	result, err := balanceFlow(userID, balance, userAgent, ipAddr, balanceType)
	if err != nil {
		c.JSON(422, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": result,
	})

	return
}

func TransferBalanceController(c *gin.Context) {

	//Get user agent
	validate := c.Request.Header["User-Agent"]
	if len(validate) == 0 {
		c.JSON(422, gin.H{
			"message": "User agent doesn't exist",
		})

		return
	}
	userAgent := validate[0]
	ipAddr := c.ClientIP()

	id, exist := c.Get("user_id")
	if !exist {
		c.JSON(422, gin.H{
			"message": "user id not found.",
		})
	}
	stringID := fmt.Sprintf("%v", id)
	userID, _ := strconv.Atoi(stringID)
	balance, _ := strconv.Atoi(c.PostForm("balance"))
	oldBalance := balance
	balance = balance - (balance * 2)
	var result string

	// Debit the money from user that want transfer his/her money
	{
		response := &result
		var balanceType model.BalanceType
		if err := balanceType.Init("debit"); err != nil {
			c.JSON(422, gin.H{
				"message": err,
			})
		}

		rp, err := balanceFlow(userID, balance, userAgent, ipAddr, balanceType)
		if err != nil {
			c.JSON(422, gin.H{
				"message": err.Error(),
			})

			return
		}
		*response = rp
	}

	// Code below find user id the user that will be receive the money by email
	//---- Code Here
	// Credit the money from user before to another user
	{
		response := &result
		var balanceType model.BalanceType
		if err := balanceType.Init("credit"); err != nil {
			c.JSON(422, gin.H{
				"message": err,
			})
		}

		rp, err := balanceFlow(userID, oldBalance, userAgent, ipAddr, balanceType)
		if err != nil {
			c.JSON(422, gin.H{
				"message": err.Error(),
			})

			return
		}
		*response = rp
	}

	c.JSON(200, gin.H{
		"message": result,
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

func balanceFlow(userID int, balance int, userAgent string, ipAddr string, balanceType model.BalanceType) (string, error) {

	db, err := database.InitDatabase()
	if err != nil {
		return "", err
	}
	defer func() {
		db.Close()
		if r := recover(); r != nil {
			panicMessage := fmt.Sprintf("v%", r)
			fmt.Errorf(panicMessage)
		}
	}()

	//check balanceType is not empty
	if len(balanceType.Type) == 0 {
		return "", fmt.Errorf("You must init balance type first")
	}

	if balanceType.Type == "debit" {
		if balance > -1 {
			return "", fmt.Errorf("You can't add balance on debit")
		}
	}

	if balanceType.Type == "credit" {
		if balance < 1 {
			return "", fmt.Errorf("You can't decrease balance on credit")
		}
	}

	oldDataUserBalance := model.UserBalance{}
	var count int
	if errResult := getCountAndOldUserBalance(&oldDataUserBalance, &count, db, userID); errResult != nil {
		return "", errResult
	}

	dataUserBalance := model.UserBalance{}
	dataUserBalance.UserID = userID

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
		return "", err
	}

	dataUserBalanceHistory := model.UserBalanceHistory{}
	dataUserBalanceHistory.UserBalanceID = dataUserBalance.ID
	dataUserBalanceHistory.Type = balanceType.Type
	stringID := strconv.Itoa(userID)
	stringBalance := strconv.Itoa(balance)
	dataUserBalanceHistory.Activity = `{"user_id":` + stringID + `,"activity":"` + balanceType.Type + `","balance":` + stringBalance + `}`
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

	dataUserBalanceHistory.IP = ipAddr

	//get location from IP
	locationModel := model.Location{
		Country:     "",
		CountryCode: "",
		ISP:         "",
		RegionName:  "",
		Timezone:    "",
	}
	dataUserBalanceHistory.Location = locationModel
	response, err := http.Get("http://ip-api.com/json/" + ipAddr)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		responseString := string(data)
		if err := json.Unmarshal([]byte(responseString), &locationModel); err != nil {
			return "", err
		}
		dataUserBalanceHistory.Location = locationModel
	}

	dataUserBalanceHistory.UserAgent = userAgent

	if err := db.Create(&dataUserBalanceHistory).Error; err != nil {
		return "", err
	}

	return "success", nil
}
