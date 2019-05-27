package main

import (
	"fmt"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"net/http"
	"strconv"
	"time"
)

var db *gorm.DB


func init() {
	//创建一个数据库的连接
	var err error
	db, err = gorm.Open("sqlite3", "demo1.db")
	if err != nil {
		panic("failed to connect database")
		fmt.Println(err)
	}

	//迁移the schema
	db.AutoMigrate(&Sale{})
	db.AutoMigrate(&User{})
}
type Sale struct {
	gorm.Model  // ID 创建时间 修改时间 删除时间
	CreatUser string  `gorm :"type:varchar(10)" json:"creat_user"`  // 创建人
	SaleNumber string `gorm:"type:varchar(10)" json:"sale_number"` // 销售单号
	Client string `grom:"type:varchar(30)" json:"client" form:"client"`// 客户名称
	City string `json:"city"`
	BillingTime time.Time `json:"billing_time"`// 开单时间
	ContractPeriod string `json:"contract_period"` // 合同账期
	AccountPeriod int `json:"account_period"`//账期
	Merchandiser string `json:"merchandiser"`// 跟单员
	Salesman string `json:"salesman"`// 业务员
	Currency string `json:"currency"`// 币种
	UnitPrice float64 `json:"unit_price"`// 单价
	Quantity float64 `json:"quantity"`// 数量
	AmountReceivable float64 `json:"amount_receivable"`// 应收金额
	Invoice int `json:"invoice"`// 发票
	PaidAmount float64 `json:"paid_amount"`// 实收金额
	UncollectedAmount float64 `json:"uncollected_amount"`// 未收金额
	DueDate  time.Time `json:"due_date"`// 应收日期
	CollectionDate  time.Time `json:"collection_date"`// 收款日
	CollectionAmount  float64 `json:"collection_amount"`// 收款金额
	TimeOut time.Time `json:"time_out"`// 超时
	Remarks string `json:"remarks"`// 备注
}

type User struct {
	ID int `json:"id" gorm:"AUTO_INCREMENT"`
	Name string `json:"name"` // 登录用户
	Password string `json:"password"`// 登录密码
	LoginDate time.Time // 登录时间
	Permission int `json:"permission"`// 权限

}


//  返回所有用户列表
func fetchUserList(c *gin.Context){
	var userlist []User
	db.Find(&userlist)
	//fmt.Println(userlist)
	if len(userlist) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message":"没有用户！"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":http.StatusOK,
		"data":userlist,
	})
}


// 创建用户
func createUser(c *gin.Context){
	permission, _ := strconv.Atoi(c.PostForm("permission"))
	user := User{Name: c.PostForm("name"), Password: c.PostForm("password"), Permission:permission}
	db.Save(&user)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!"})

}

// 根据用户名查询
func fetchSingleUser(c *gin.Context){
	var user User
	Uname := c.Param("name")
	db.Where(&User{Name: Uname}).Find(&user)
	fmt.Println(Uname)

	if user.Name == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": user})
}

// 更新匹配到的用户数据
func updateUser(c *gin.Context)  {
	var user User
	uname := c.Param("name")
	db.Where(&User{Name: uname}).Find(&user)
	fmt.Println(user)

	if user.Name == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "未找到!"})
		return
	}
	db.Model(&user).Where("name = ?", uname).Update("password",c.PostForm("password"))
	permission, _ := strconv.Atoi(c.PostForm("permission"))
	db.Model(&user).Where("name = ?", c.Param("name")).Update("permission", permission)
	fmt.Println(c.PostFormArray("password"))
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "用户信息更新成功!"})
}

// 查询订单
func fetchSalesorder(c *gin.Context){
	var sale []Sale
	salenumber := c.Query("sale_number")
	client := strconv.QuoteToASCII(c.Query("client"))
	fmt.Println(salenumber)
	client, _ =strconv.Unquote(client)

	if client == "" {
		db.Where("sale_number = ?", salenumber ).Find(&sale)
		if len(sale) == 0 {
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "没找到该订单!"})
			return
		}
	} else {
		db.Where("client LIKE ?","%"+client+"%").Find(&sale)
		if len(sale) == 0 {
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "没找到该客户订单，请换个描述方式!"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": sale})
}

// 更新订单数据
func updateSales(c *gin.Context)  {
	var sale []Sale
	salenum := c.PostForm("sale_number")
	db.Where(&Sale{SaleNumber: salenum}).Find(&sale)

	if len(sale) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "未找到!"})
		return
	}
	for key, sa := range sale {
		fmt.Println(key, sa.SaleNumber)

		//fmt.Println(c.PostForm("id"))
		//db.Model(&sale).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
		db.Model(&sale[key]).Where("sale_number = ?", sa.SaleNumber).Update("collection_date",
			c.PostFormArray("collection_date")[key])
		//fmt.Println(sa.SaleNumber, sa.CollectionDate)
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "更新成功!"})
}

// 删除订单
func deleteSale(c *gin.Context) {
	var sale Sale
	salenumber := c.Param("sale_number")
	db.Where(&Sale{SaleNumber:salenumber}).First(&sale)
	fmt.Println(sale)

	if sale.SaleNumber == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "未找到!"})
		return
	}

	db.Delete(&sale)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "删除成功！"})
}

//定义一个回调函数，用来决断用户id和密码是否有效
func authCallback(c *gin.Context) (interface{}, error ) {

	var user User
	username := c.PostForm("username")
	password := c.PostForm("password")

	db.Select("id").Where(User{Name : username, Password : password}).First(&user)
	if user.ID > 0 {
		return true, nil
	}

	return nil, jwt.ErrFailedAuthentication
}

//定义一个回调函数，用来决断用户在认证成功的前提下，是否有权限对资源进行访问
func authPrivCallback(User interface{}, c *gin.Context) bool {
	if v, ok := User.(int); ok && v == 1 {
		return true
	}
	return false
}

//定义一个函数用来处理，认证不成功的情况
func unAuthFunc(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

// jwt 生成
//func usermMiddleware(){
//	// the jwt middleware
//	authMiddleware := &jwt.GinJWTMiddleware{
//		//Realm name to display to the User1. Required.
//		//必要项，显示给用户看的域
//		Realm: "zone name just for test",
//		//Secret key used for signing. Required.
//		//用来进行签名的密钥，就是加盐用的
//		Key: []byte("secret key salt"),
//		//Duration that a jwt token is valid. Optional, defaults to one hour
//		//JWT 的有效时间，默认为一小时
//		Timeout: time.Hour,
//		// This field allows clients to refresh their token until MaxRefresh has passed.
//		// Note that clients can refresh their token in the last moment of MaxRefresh.
//		// This means that the maximum validity timespan for a token is MaxRefresh + Timeout.
//		// Optional, defaults to 0 meaning not refreshable.
//		//最长的刷新时间，用来给客户端自己刷新 token 用的
//		MaxRefresh: time.Hour,
//		// Callback function that should perform the authentication of the User1 based on userID and
//		// password. Must return true on success, false on failure. Required.
//		// Option return User1 data, if so, User1 data will be stored in Claim Array.
//		//必要项, 这个函数用来判断 User1 信息是否合法，如果合法就反馈 true，否则就是 false, 认证的逻辑就在这里
//		Authenticator: authCallback,
//		// Callback function that should perform the authorization of the authenticated User1. Called
//		// only after an authentication success. Must return true on success, false on failure.
//		// Optional, default to success
//		//可选项，用来在 Authenticator 认证成功的基础上进一步的检验用户是否有权限，默认为 success
//		Authorizator: authPrivCallback,
//		// User1 can define own Unauthorized func.
//		//可以用来息定义如果认证不成功的的处理函数
//		Unauthorized: unAuthFunc,
//		// TokenLookup is a string in the form of "<source>:<name>" that is used
//		// to extract token from the request.
//		// Optional. Default value "header:Authorization".
//		// Possible values:
//		// - "header:<name>"
//		// - "query:<name>"
//		// - "cookie:<name>"
//		//这个变量定义了从请求中解析 token 的格式
//		TokenLookup: "header: Authorization, query: token, cookie: jwt",
//		// TokenLookup: "query:token",
//		// TokenLookup: "cookie:token",
//
//		// TokenHeadName is a string in the header. Default value is "Bearer"
//		//TokenHeadName 是一个头部信息中的字符串
//		TokenHeadName: "Bearer",
//
//		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
//		//这个指定了提供当前时间的函数，也可以自定义
//		TimeFunc: time.Now,
//	}
//}

func main(){

	// the jwt middleware
	authMiddleware := &jwt.GinJWTMiddleware{
		//Realm name to display to the User1. Required.
		//必要项，显示给用户看的域
		Realm: "zone name just for test",
		//Secret key used for signing. Required.
		//用来进行签名的密钥，就是加盐用的
		Key: []byte("secret key salt"),
		//Duration that a jwt token is valid. Optional, defaults to one hour
		//JWT 的有效时间，默认为一小时
		Timeout: time.Hour,
		// This field allows clients to refresh their token until MaxRefresh has passed.
		// Note that clients can refresh their token in the last moment of MaxRefresh.
		// This means that the maximum validity timespan for a token is MaxRefresh + Timeout.
		// Optional, defaults to 0 meaning not refreshable.
		//最长的刷新时间，用来给客户端自己刷新 token 用的
		MaxRefresh: time.Hour,
		// Callback function that should perform the authentication of the User1 based on userID and
		// password. Must return true on success, false on failure. Required.
		// Option return User1 data, if so, User1 data will be stored in Claim Array.
		//必要项, 这个函数用来判断 User1 信息是否合法，如果合法就反馈 true，否则就是 false, 认证的逻辑就在这里
		Authenticator: authCallback,
		// Callback function that should perform the authorization of the authenticated User1. Called
		// only after an authentication success. Must return true on success, false on failure.
		// Optional, default to success
		//可选项，用来在 Authenticator 认证成功的基础上进一步的检验用户是否有权限，默认为 success
		Authorizator: authPrivCallback,
		// User1 can define own Unauthorized func.
		//可以用来息定义如果认证不成功的的处理函数
		Unauthorized: unAuthFunc,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		//这个变量定义了从请求中解析 token 的格式
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		//TokenHeadName 是一个头部信息中的字符串
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		//这个指定了提供当前时间的函数，也可以自定义
		TimeFunc: time.Now,
	}

	router := gin.Default()

	router.POST("/longin", authMiddleware.LoginHandler)
	auth := router.Group("/user")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/", fetchUserList)
		auth.GET("/:name", fetchSingleUser)
		auth.POST("/", createUser)
		auth.PUT(":name",updateUser)

	}
	sa := router.Group("/sale")
	{
		sa.GET("/q", fetchSalesorder)
		sa.PUT("/:client", updateSales)
		sa.DELETE("/:sale_number", deleteSale)
	}

	router.Run(":8081")
}