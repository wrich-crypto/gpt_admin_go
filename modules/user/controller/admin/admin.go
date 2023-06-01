package admin

import (
	"fmt"
	admin_model "gpt_admin_go/modules/admin/model"
	"gpt_admin_go/modules/user/config"
	"gpt_admin_go/modules/user/model"
	"gpt_admin_go/modules/user/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"context"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
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

	user.Password = ""
	responseData := g.Map{"token": token, "user": user}
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

type UploadFileReq struct {
	g.Meta   `path:"/upload" method:"POST"`
	File     *ghttp.UploadFile `p:"file"`
	FileName string            `p:"file_name"`
}

func (c *UserOpenController) HandleUploadToOSS(ctx context.Context, req *UploadFileReq) (res *cool.BaseRes, err error) {
	// 以下为OSS的相关配置
	var (
		endpoint        = config.Config.Oss.Endpoint
		accessKeyID     = config.Config.Oss.AccessKeyID
		accessKeySecret = config.Config.Oss.AccessKeySecret
		bucketName      = config.Config.Oss.BucketName
	)

	g.Log().Info(ctx, endpoint)
	g.Log().Info(ctx, accessKeyID)
	g.Log().Info(ctx, accessKeySecret)
	g.Log().Info(ctx, bucketName)

	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		g.Log().Error(ctx, "Failed to create OSS client: ", err)
		return cool.Fail("Failed to create OSS client"), err
	}

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		g.Log().Error(ctx, "Failed to get bucket: ", err)
		return cool.Fail("Failed to get bucket"), err
	}

	file, err := req.File.Open()
	if err != nil {
		// handle error
	}

	// 上传文件流。
	err = bucket.PutObject(req.FileName, file)
	if err != nil {
		g.Log().Error(ctx, "Failed to upload file: ", err)
		return cool.Fail("Failed to upload file"), err
	}

	fileUrl := bucketName + "." + endpoint + "/" + req.FileName
	return cool.Ok(g.Map{"file_url": fileUrl}), nil
}

type AddCardReq struct {
	g.Meta         `path:"/add_auto" method:"POST"`
	CardNumber     int     `p:"card_number"`
	RechargeAmount float64 `p:"recharge_amount"`
	CreateUser     string  `p:"create_user"`
	Remark         string  `p:"remark"`
}

func (c *RechargeCardController) HandleAddCardAuto(ctx context.Context, req *AddCardReq) (res *cool.BaseRes, err error) {
	user, ok := ctx.Value("user").(*model.Users)
	if !ok {
		return cool.Fail("Invalid user context"), fmt.Errorf("Invalid user context")
	}

	if user.Role != 2 && user.Role != 3 {
		return cool.Fail("Unauthorized"), fmt.Errorf("Unauthorized")
	}

	if user.Role == 2 {
		balance, err := admin_model.GetAgentBalance(user.ReferralCode)
		if err != nil {
			return cool.Fail("Failed to get agent balance"), err
		}

		if balance < req.RechargeAmount*float64(req.CardNumber) {
			return cool.Fail("Insufficient agent balance"), fmt.Errorf("Insufficient agent balance")
		}

		err = admin_model.UpdateAgentBalance(user.ReferralCode, balance-(req.RechargeAmount*float64(req.CardNumber)))
		if err != nil {
			return cool.Fail("Failed to update agent balance"), err
		}
	}

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
