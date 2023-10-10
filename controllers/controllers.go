package controllers

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/AlwiLion/database"
	"github.com/AlwiLion/models"
	"github.com/AlwiLion/utils"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
		IsKyc:    false,
		Mobile:   data["mobile"],
		Pan:      data["pan"],
	}
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		// for _, err := range err.(validator.ValidationErrors) {
		// 	fmt.Println(err)
		// }
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "request paramater missing",
		})

	} else {
		result := database.DB.Create(&user)

		// Check for errors during insertion
		if result.Error != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "failed to insert record into the database",
			})
		}
		return c.JSON(&user)
		//fmt.Println("Validation passed!")
	}

}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{"message": "User not found"})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	// cookie := fiber.Cookie{
	// 	Name:     "jwt",
	// 	Value:    token,
	// 	Expires:  time.Now().Add(time.Hour * 24),
	// 	HTTPOnly: true,
	// }

	//c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
		"token":   token,
	})
	//return c.JSON(user)

}

func User(c *fiber.Ctx) error {
	//cookie := c.Cookies("jwt")

	token, err := authentication(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)
	fmt.Printf("type of id %T", user.ID)
	return c.JSON(user)
}

func UpdateKyc(c *fiber.Ctx) error {

	token, err := authentication(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)

	var existingUser models.User

	result := database.DB.Where("id = ?", claims.Issuer).First(&existingUser)

	if result.Error != nil {
		return c.JSON(fiber.Map{
			"message": "Data Not found",
		})
	} else {

		var data map[string]bool

		if err := c.BodyParser(&data); err != nil {
			return err
		}

		if value, ok := data["isKyc"]; ok {

			existingUser.IsKyc = value

		} else {
			return c.JSON(fiber.Map{
				"message": "Key Missing",
			})
		}

		if err := database.DB.Save(&existingUser).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
		}

		return c.JSON(existingUser)
	}

}

func DepositWithdrawMoney(c *fiber.Ctx) error {

	token, err := authentication(c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("id = ?", claims.Issuer).Find(&user)

	if user.IsKyc {
		var txn models.UserStatment
		txn.UserID = user.ID
		txn.TxnId = utils.GenerateRandomString(10)
		data := make(map[string]interface{})

		if err := c.BodyParser(&data); err != nil {
			return err
		}

		if value, ok := data["amount"]; ok {
			// Type assertion to float64
			amount, ok := value.(float64)
			if !ok {
				return c.JSON(fiber.Map{
					"message": "Amount Key has an unexpected type",
				})
			}

			// Successfully asserted to float64, now you can use 'amount'
			txn.Amount = amount
		} else {
			return c.JSON(fiber.Map{
				"message": "Amount Key Missing",
			})
		}

		if value, ok := data["isDeposit"]; ok {
			// Type assertion to float64
			isDeposit, ok := value.(bool)
			if !ok {
				return c.JSON(fiber.Map{
					"message": "isDeposit Key has an unexpected type",
				})
			}

			// Successfully asserted to float64, now you can use 'amount'
			txn.IsDeposit = isDeposit
		} else {
			return c.JSON(fiber.Map{
				"message": "isDeposit Key Missing",
			})
		}

		if value, ok := data["comment"]; ok {
			// Type assertion to float64
			comment, ok := value.(string)
			if !ok {
				return c.JSON(fiber.Map{
					"message": "comment Key has an unexpected type",
				})
			}
			// Successfully asserted to float64, now you can use 'amount'
			txn.Comment = comment
		}
		result := database.DB.Create(&txn)

		// Check for errors during insertion
		if result.Error != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "failed to insert record into the database",
			})
		}
		return c.JSON(txn)
	} else {
		return c.JSON(fiber.Map{
			"message": "Kyc Not done",
		})
	}

}

func CheckBalance(c *fiber.Ctx) error {

	token, err := authentication(c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("id = ?", claims.Issuer).Find(&user)

	var totalDepositAmount float64
	database.DB.Model(models.UserStatment{}).Select("COALESCE(SUM(amount), 0) as total_deposit_amount").Where("is_deposit = ? AND user_id = ?", true, user.ID).Scan(&totalDepositAmount)

	var totalWithdrawAmount float64
	database.DB.Model(models.UserStatment{}).Select("COALESCE(SUM(amount), 0) as total_deposit_amount").Where("is_deposit = ? AND user_id = ?", false, user.ID).Scan(&totalWithdrawAmount)

	return c.JSON(fiber.Map{
		"balance": totalDepositAmount - totalWithdrawAmount,
	})
}

func TransferMoney(c *fiber.Ctx) error {

	var wg sync.WaitGroup

	token, err := authentication(c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var senderUser models.User
	var reciverUser models.User

	database.DB.Where("id = ?", claims.Issuer).Find(&senderUser)
	data := make(map[string]interface{})
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if !senderUser.IsKyc {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Kyc pending",
		})
	}

	var amount float64

	if value, ok := data["amount"]; ok {
		// Type assertion to float64
		amnt, ok := value.(float64)
		if !ok {
			return c.JSON(fiber.Map{
				"message": "amount Key has an unexpected type",
			})
		}
		amount = amnt
	} else {
		return c.JSON(fiber.Map{
			"message": "amount Key Missing",
		})
	}

	if value, ok := data["reciverMobileNumber"]; ok {
		// Type assertion to float64
		reciverMobileNumber, ok := value.(string)
		if !ok {
			return c.JSON(fiber.Map{
				"message": "reciverMobileNumber Key has an unexpected type",
			})
		}
		result := database.DB.Model(models.User{}).Where("mobile = ?", reciverMobileNumber).Find(&reciverUser)

		//database.DB.Where("mobile = ?", senderMobileNumber).Find(&senderUser)

		if result.Error != nil {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "user not found",
			})

		} else {
			if senderUser.Mobile == reciverUser.Mobile {
				return c.JSON(fiber.Map{
					"message": "Amount cant be transfer to same number",
				})
			}
			var txnSender models.UserStatment
			var txnReceiver models.UserStatment

			txnSender.UserID = senderUser.ID
			txnReceiver.UserID = reciverUser.ID
			txnSender.IsDeposit = false
			txnReceiver.IsDeposit = true
			txnId := utils.GenerateRandomString(10)
			txnReceiver.TxnId = txnId
			txnSender.TxnId = txnId
			txnReceiver.Amount = amount
			txnSender.Amount = amount
			if value, ok := data["comment"]; ok {
				// Type assertion to float64
				comment, ok := value.(string)
				if !ok {
					return c.JSON(fiber.Map{
						"message": "comment Key has an unexpected type",
					})
				}
				// Successfully asserted to float64, now you can use 'amount'
				txnSender.Comment = comment
				txnReceiver.Comment = comment
			}

			wg.Add(2)
			go func() {
				defer wg.Done()
				// Asynchronously create txnSender
				fmt.Println("txnSender ", txnSender)
				database.DB.Create(&txnSender)
			}()

			go func() {
				defer wg.Done()
				// Asynchronously create txnReceiver
				fmt.Println("txnReciver ", txnReceiver)
				database.DB.Create(&txnReceiver)
			}()

			wg.Wait()

			var totalDepositAmount float64
			database.DB.Model(models.UserStatment{}).Select("COALESCE(SUM(amount), 0) as total_deposit_amount").Where("is_deposit = ? AND user_id = ?", true, senderUser.ID).Scan(&totalDepositAmount)

			var totalWithdrawAmount float64
			database.DB.Model(models.UserStatment{}).Select("COALESCE(SUM(amount), 0) as total_deposit_amount").Where("is_deposit = ? AND user_id = ?", false, senderUser.ID).Scan(&totalWithdrawAmount)

			return c.JSON(fiber.Map{
				"message":          "Balance tarnsfer successfull",
				"txnId":            txnId,
				"availableBalance": totalDepositAmount - totalWithdrawAmount,
			})

		}

	} else {
		return c.JSON(fiber.Map{
			"message": "senderMobileNumber Key Missing",
		})
	}

}
