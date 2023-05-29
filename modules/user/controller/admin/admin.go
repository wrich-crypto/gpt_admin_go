package admin

import (
	"gpt_admin_go/modules/user/model"
	"gpt_admin_go/modules/user/service"

	"github.com/gogf/gf/v2/frame/g"

	"context"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/cool-team-official/cool-admin-go/cool"
)

type UserOpenController struct {
	*cool.Controller
}

type UserController struct {
	*cool.Controller
}

type RechargeCardController struct {
	*cool.Controller
}

func init() {
	var user_open_controller = &UserOpenController{
		&cool.Controller{
			Perfix:  "/admin/base/open",
			Service: service.NewUserService(),
		},
	}

	var user_controller = &UserController{
		&cool.Controller{
			Perfix:  "/admin/user",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewUserService(),
		},
	}

	var recharge_card_controller = &RechargeCardController{
		&cool.Controller{
			Perfix:  "/admin/recharge_card",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewRechargeCardService(),
		},
	}

	cool.RegisterController(user_controller)
	cool.RegisterController(user_open_controller)
	cool.RegisterController(recharge_card_controller)
}

type LoginReq struct {
	g.Meta   `path:"/login" method:"POST"`
	Username string `p:"username"`
	Password string `p:"password"`
}

// HandleUserLogin 登录接口
func (c *UserOpenController) HandleUserLogin(ctx context.Context, req *LoginReq) (res *cool.BaseRes, err error) {
	if req.Username == "" || req.Password == "" {
		g.Log().Error(ctx, "Invalid parameters")
		return cool.Fail("Invalid parameters"), nil
	}

	hashPassword := hashToken(req.Password)
	token := generateToken(req.Username, hashPassword)

	user := new(model.Users)
	err = g.DB().Model(user.TableName()).Where("username", req.Username).Scan(user)
	if err != nil {
		g.Log().Error(ctx, "HandleUserLogin User.query, username: "+req.Username+" account no exist")
		return cool.Fail("Account does not exist"), err
	}

	if user.Token != token {
		g.Log().Error(ctx, "HandleUserLogin User.query, username: "+req.Username+" password error, token invalid")
		return cool.Fail("Password error"), nil
	}

	responseData := g.Map{"token": token}
	return cool.Ok(responseData), nil
}

func hashToken(context string) string {
	hashedPassword := sha256.Sum256([]byte(context))
	return hex.EncodeToString(hashedPassword[:])
}

func generateToken(username, password string) string {
	salt := "12345678"
	hash := sha256.Sum256([]byte(username + password + salt))
	return hex.EncodeToString(hash[:])
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

type AddCardReq struct {
	g.Meta         `path:"/add_auto" method:"POST"`
	CardNumber     int     `p:"card_number"`
	RechargeAmount float64 `p:"recharge_amount"`
	CreateUser     string  `p:"create_user"`
	Remark         string  `p:"remark"`
}

func (c *RechargeCardController) HandleAddCardAuto(ctx context.Context, req *AddCardReq) (res *cool.BaseRes, err error) {
	for i := 0; i < req.CardNumber; i++ {
		cardAccount := stringWithCharset(10, charset)
		cardPassword := stringWithCharset(8, charset)

		rechargeCard := &model.RechargeCards{
			CardAccount:    cardAccount,
			CardPassword:   cardPassword,
			CreateTime:     time.Now(),
			ExpireTime:     time.Now().AddDate(0, 0, 7), // assuming 7 days of validity for the card
			RechargeAmount: req.RechargeAmount,
			CreateUser:     req.CreateUser,
			Remark:         req.Remark,
			Status:         1, // assuming status 1 for a new card
		}

		_, err := g.DB().Model(rechargeCard.TableName()).Data(rechargeCard).Insert()
		if err != nil {
			g.Log().Error(ctx, "Failed to create a new card: ", err)
			return cool.Fail("Failed to create a new card"), err
		}
	}

	return cool.Ok("Cards added successfully"), nil
}
