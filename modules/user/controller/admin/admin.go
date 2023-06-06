package admin

import (
	"fmt"
	admin_model "gpt_admin_go/modules/admin/model"
	"gpt_admin_go/modules/user/config"
	"gpt_admin_go/modules/user/model"
	"gpt_admin_go/modules/user/service"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"context"
	cryptoRand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"path/filepath"
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
	g.Meta `path:"/upload" method:"POST"`
	File   *ghttp.UploadFile `p:"file"`
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
		return cool.Fail("Failed to upload file"), err
	}

	// Get the extension of the uploaded file.
	g.Log().Info(ctx, req.File.Filename)
	fileExtension := filepath.Ext(req.File.Filename)

	// Create a new random UUID for filename.
	uuid := newUUID() + fileExtension

	g.Log().Info(ctx, uuid)

	// 上传文件流。
	err = bucket.PutObject(uuid, file)
	if err != nil {
		g.Log().Error(ctx, "Failed to upload file: ", err)
		return cool.Fail("Failed to upload file"), err
	}

	fileUrl := bucketName + "." + endpoint + "/" + uuid
	return cool.Ok(g.Map{"file_url": fileUrl}), nil
}

func newUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := cryptoRand.Read(b)
	if err != nil {
		panic(err)
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}

type AddCardReq struct {
	g.Meta         `path:"/add_auto" method:"POST"`
	CardNumber     int     `p:"card_number"`
	RechargeAmount float64 `p:"recharge_amount"`
	Remark         string  `p:"remark"`
}

func (c *RechargeCardController) HandleAddCardAuto(ctx context.Context, req *AddCardReq) (res *cool.BaseRes, err error) {
	user, ok := ctx.Value("user").(*model.Users)
	if !ok {
		return cool.Fail("Invalid user context"), fmt.Errorf("Invalid user context")
	}

	if !user.HasRole(model.ROLE_ADMIN) && !user.HasRole(model.ROLE_AGENT) {
		return cool.Fail("Unauthorized"), fmt.Errorf("Unauthorized")
	}

	if user.HasRole(model.ROLE_AGENT) {
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
			ExpireTime:     time.Now().AddDate(0, 0, 7), // assuming 7 days of validity for the card
			RechargeAmount: req.RechargeAmount,
			CreateUser:     user.Username,
			CreateUserId:   int(user.ID),
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

type UserListReq struct {
	g.Meta       `path:"/user_list" method:"POST"`
	Size         int    `p:"size" v:"max:100#最多每页显示100条数据"`
	Page         int    `p:"page"`
	Order        string `p:"order"`
	Sort         string `p:"sort"`
	KeyWordField string `p:"key_word_field"`
	KeyWord      string `p:"key_word"`
}

func (c *UserController) UserList(ctx context.Context, req *UserListReq) (res *cool.BaseRes, err error) {
	user, ok := ctx.Value("user").(*model.Users)
	if !ok {
		return cool.Fail("Invalid user context"), fmt.Errorf("Invalid user context")
	}

	var users []*model.Users

	db := g.DB().Model(model.NewUsers().TableName())

	if req.KeyWord != "" && req.KeyWordField != "" {
		fields := strings.Split(req.KeyWordField, ",")
		for _, field := range fields {
			db = db.Where(field+" LIKE ?", "%"+req.KeyWord+"%")
		}
	}

	if user.HasRole(model.ROLE_ADMIN) {
		// Administrator can see all users
	} else if user.HasRole(model.ROLE_AGENT) {
		// Agent can only see users they invited
		db = db.Where("invitation_user_id = ?", user.ID)
	} else {
		return cool.Fail("User role not authorized"), nil
	}

	if req.Order != "" {
		if req.Sort == "" {
			req.Sort = "asc"
		}
		db = db.Order(req.Order + " " + req.Sort)
	}

	if req.Page != 0 {
		if req.Size == 0 {
			req.Size = 10
		}
		db = db.Page(req.Page, req.Size)
	}

	all, err := db.All()
	if err != nil {
		g.Log().Error(ctx, "UserList error fetching users", err)
		return cool.Fail("Error fetching users"), err
	}

	if err = all.Structs(&users); err != nil {
		g.Log().Error(ctx, "UserList error scanning users", err)
		return cool.Fail("Error scanning users"), err
	}

	return cool.Ok(users), nil
}

type RechargeCardListReq struct {
	g.Meta       `path:"/recharge_card_list" method:"POST"`
	Size         int    `p:"size" v:"max:100#最多每页显示100条数据"`
	Page         int    `p:"page"`
	Order        string `p:"order"`
	Sort         string `p:"sort"`
	KeyWordField string `p:"key_word_field"`
	KeyWord      string `p:"key_word"`
}

func (c *RechargeCardController) RechargeCardList(ctx context.Context, req *RechargeCardListReq) (res *cool.BaseRes, err error) {
	user, ok := ctx.Value("user").(*model.Users)
	if !ok {
		return cool.Fail("Invalid user context"), fmt.Errorf("Invalid user context")
	}

	var cards []*model.RechargeCards

	db := g.DB().Model(model.NewRechargeCards().TableName())

	if req.KeyWord != "" && req.KeyWordField != "" {
		fields := strings.Split(req.KeyWordField, ",")
		for _, field := range fields {
			db = db.Where(field+" LIKE ?", "%"+req.KeyWord+"%")
		}
	}

	if user.HasRole(model.ROLE_ADMIN) {
		// Administrator can see all recharge cards
	} else if user.HasRole(model.ROLE_AGENT) {
		// Agent can only see recharge cards they created
		db = db.Where("create_user_id = ?", user.ID)
	} else {
		return cool.Fail("User role not authorized"), nil
	}

	if req.Order != "" {
		if req.Sort == "" {
			req.Sort = "asc"
		}
		db = db.Order(req.Order + " " + req.Sort)
	}

	if req.Page != 0 {
		if req.Size == 0 {
			req.Size = 10
		}
		db = db.Page(req.Page, req.Size)
	}

	all, err := db.All()
	if err != nil {
		g.Log().Error(ctx, "RechargeCardList error fetching cards", err)
		return cool.Fail("Error fetching recharge cards"), err
	}

	if err = all.Structs(&cards); err != nil {
		g.Log().Error(ctx, "RechargeCardList error scanning cards", err)
		return cool.Fail("Error scanning recharge cards"), err
	}

	return cool.Ok(cards), nil
}
