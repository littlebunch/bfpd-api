package auth

import (
	"bfpd-api/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/appleboy/gin-jwt.v2"
	"log"
	"time"
)

type User struct {
	Id       int64
	Name     string `gorm:"not null"`
	Password string `gorm:"not null"`
	Roles    []Role `gorm:"many2many:user_roles"`
}
type Role struct {
	Id   int64
	Role string `gorm:"not null"`
}

//initialize our jwt components
func (u *User) AuthMiddleware(db *bfpd.DB) *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			var u User
			db.Where("Name=?", userId).First(&u)
			if u.Id > 0 {
				if CheckPasswordHash(password, u.Password) == true {
					return u.Name, true
				}
			}
			return userId, false
		},
		// load user's roles into a []string and put into the claim
		PayloadFunc: func(userId string) map[string]interface{} {
			var u User
			var roles []string
			db.Preload("Roles").Where("Name=?", userId).First(&u)
			for _, r := range u.Roles {
				roles = append(roles, r.Role)
			}
			return map[string]interface{}{"ROLES": roles}
		},
		// must have ADMIN role from the roles assigned to the current user
		Authorizator: func(userId string, c *gin.Context) bool {
			jwtClaims := jwt.ExtractClaims(c)
			for _, v := range jwtClaims["ROLES"].([]interface{}) {
				if v.(string) == "ADMIN" {
					return true
				}
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header:Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}
}

// add a user named 'littlebunch' with ADMIN and USER roles
func (u *User) BootstrapUsers(db *bfpd.DB) {
	var user User
	user.Name = "littlebunch"
	password, err := HashPassword("littlebunch")
	if err != nil {
		log.Fatal(err)
	}
	user.Password = password
	var role, role2 Role
	role.Role = "ADMIN"
	db.Where("role=?", role.Role).First(&role)
	if db.NewRecord(&role) == true {
		db.Create(&role)
	}
	role2.Role = "USER"
	db.Where("role=?", role2.Role).First(&role2)
	if db.NewRecord(&role2) == true {
		db.Create(&role2)
	}
	db.Where("name=?", user.Name).First(&user)
	if db.NewRecord(&user) == true {
		roles := []Role{role, role2}
		db.Create(&user)
		db.Model(&user).Association("Roles").Append(roles)
	}
	return
}

// generates an encrypted password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// compares a plain text password with a hash and returns true for matches
// otherwise false
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
