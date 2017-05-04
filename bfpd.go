package main

import (
	"bfpd/auth"
	"bfpd/model"
	//"encoding/json"
	"flag"
	"fmt"
	//"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var d = flag.Bool("d", false, "Debug")
var i = flag.Bool("i", false, "Initialize database database")
var c = flag.String("c", "config.yaml", "YAML Config file")

/*type DB struct {
	*gorm.DB
}*/

func Database(cs *bfpd.Config) (*bfpd.DB, error) {
	c := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8&parseTime=True&loc=Local", cs.User, cs.Pw, cs.Db)
	//open a db connection\
	db, err := gorm.Open("mysql", c)
	if err != nil {
		log.Fatal(err)
		panic("failed to connect database")
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	return &bfpd.DB{db}, nil
}

func main() {
	flag.Parse()
	// try to read config from a file
	var cs bfpd.Config
	raw, err := ioutil.ReadFile(*c)
	if err != nil {
		log.Println(err.Error())
	}
	yaml.Unmarshal(raw, &cs)
	// check the environment for one or more config variables (overrides config file)
	if os.Getenv("BFPD_DB") != "" {
		cs.Db = os.Getenv("BFPD_DB")
	}
	if os.Getenv("BFPD_USER") != "" {
		cs.User = os.Getenv("BFPD_USER")
	}
	if os.Getenv("BFPD_PW") != "" {
		cs.Pw = os.Getenv("BFPD_PW")
	}
	db, _ := Database(&cs)
	defer db.Close()
	if *d == true {
		db.LogMode(*d)
	}
	var u *auth.User
	if *i == true {
		log.Printf("Initializing database...")
		db.AutoMigrate(&bfpd.Food{},
			&bfpd.Ingredients{},
			&bfpd.FoodGroup{},
			&bfpd.Weights{},
			&bfpd.Nutrient{},
			&bfpd.Manufacturer{},
			&bfpd.Unit{},
			&bfpd.SourceCode{},
			&bfpd.Derivation{},
			&bfpd.FootNote{},
			&bfpd.NutrientData{},
			&auth.User{},
			&auth.Role{})
		db.Model(&bfpd.Food{}).AddForeignKey("food_group_id", "food_groups(id)", "NO ACTION", "NO ACTION")
		db.Model(&bfpd.Food{}).AddForeignKey("ingredients_id", "ingredients(id)", "CASCADE", "CASCADE")
		db.Model(&bfpd.Food{}).AddForeignKey("manufacturer_id", "manufacturers(id)", "NO ACTION", "NO ACTION")
		db.Model(&bfpd.NutrientData{}).AddForeignKey("food_id", "foods(id)", "NO ACTION", "NO ACTION")
		db.Model(&bfpd.NutrientData{}).AddForeignKey("nutrient_id", "nutrients(id)", "NO ACTION", "NO ACTION")
		db.Model(&bfpd.NutrientData{}).AddForeignKey("derivation_id", "derivations(id)", "NO ACTION", "NO ACTION")
		db.Model(&bfpd.NutrientData{}).AddForeignKey("source_id", "source_codes(id)", "NO ACTION", "NO ACTION")
		u.BootstrapUsers(db)
	}
	// initialize our jwt authentication

	authMiddleware := u.AuthMiddleware(db)
	router := gin.Default()

	v1 := router.Group("/ndb/api/v1")
	{

		v1.POST("/login", authMiddleware.LoginHandler)
		// page through foods
		v1.GET("/food/", func(c *gin.Context) {
			var foods []bfpd.Food
			var _foods []bfpd.TransformedFood
			max, err := strconv.ParseInt(c.Query("max"), 10, 0)
			if err != nil {
				max = 50
			}
			page, err := strconv.ParseInt(c.Query("page"), 10, 0)
			if err != nil {
				log.Printf("Page=%v", c.Query("page"))
				page = 1
			}
			if page <= 0 {
				page = 1
			}
			offset := page * max
			log.Printf("page=%d offset=%d max=%d", page, offset, max)

			db.Preload("Manufacturer").Preload("Measures").Preload("FoodGroup").Preload("Ingredients").Preload("NutrientData").Preload("NutrientData.Nutrient").Preload("NutrientData.Sourcecode").Preload("NutrientData.Derivation").Offset(offset).Limit(max).Find(&foods)
			if foods == nil {
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No food found!"})
				return
			}
			//transforms the foods for building a good response
			for _, item := range foods {
				_foods = append(_foods, transformfood(&item))
			}
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _foods})
		})
		// return a single food by id
		v1.GET("/food/:id", func(c *gin.Context) {
			var food bfpd.Food
			foodId := c.Param("id")
			db.Preload("Manufacturer").Preload("Measures").Preload("Ingredients").Preload("FoodGroup").Preload("NutrientData").Preload("NutrientData.Nutrient").Preload("NutrientData.Sourcecode").Preload("NutrientData.Derivation").First(&food, foodId)
			if food.Id == 0 {
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No food found!"})
				return
			}
			_food := transformfood(&food)
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _food})
		})

		v1.PUT("/food/:id", func(c *gin.Context) {
			var food bfpd.Food
			foodId := c.Param("id")
			completed, _ := strconv.Atoi(c.PostForm("completed"))
			db.First(&food, foodId)
			if food.Id == 0 {
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No food found!"})
				return
			}
			db.Model(&food).Update("Description", c.PostForm("Description"))
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "food updated successfully!", "completed": completed})
		})
		v1.DELETE("/food/:id", func(c *gin.Context) {
			var food bfpd.Food
			foodId := c.Param("id")
			db.Preload("Ingredients").Preload("NutrientData").Preload("Measures").First(&food, foodId)
			if food.Id == 0 {
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No food found!"})
				return
			}
			tx := db.Begin()
			if err := tx.Where("id=?", food.Ingredients.Id).Delete(&bfpd.Ingredients{}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Delete failed!"})
				return
			}
			if err := tx.Where("food_id=?", food.Id).Delete(&bfpd.NutrientData{}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Delete failed!"})
				return
			}
			if err := tx.Where("food_id=?", food.Id).Delete(&bfpd.Weights{}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Delete failed!"})
				return
			}
			if err := tx.Delete(&food).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Delete failed!"})
				return
			}
			tx.Commit()
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "food deleted successfully!"})
		})
		v1.GET("/ndbno/:ndbno", func(c *gin.Context) {
			var food bfpd.Food
			ndbno := c.Param("ndbno")
			db.Where("ndbno = ?", ndbno).Find(&food)
			if food.Id == 0 {
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No food found!"})
				return
			}
			_food := bfpd.TransformedFood{ID: food.Id, Ndbno: food.Ndbno, Description: food.Description}
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _food})
		})
		v1.GET("/nutrient/", func(c *gin.Context) {
			var nutr []bfpd.Nutrient
			db.Preload("Unit").Find(&nutr)
			if len(nutr) <= 0 {
				c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No nutrients found!"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": nutr})
		})
		// create a food
		v1.POST("/food/", func(c *gin.Context) {
			var food bfpd.Food
			var manu bfpd.Manufacturer
			var nutr bfpd.Nutrient
			var derv bfpd.Derivation
			var source bfpd.SourceCode
			var fg bfpd.FoodGroup
			err := c.BindJSON(&food)
			if err == nil {
				// check for  manufacturer ID
				db.Where("name=?", food.Manufacturer.Name).First(&manu)
				if manu.Id != 0 {
					food.Manufacturer.Id = manu.Id
				}
				// food Group
				db.Where("cd=?", food.FoodGroup.Cd).First(&fg)
				if fg.Id != 0 {
					food.FoodGroup = fg
				}
				// set our nutrient ID's and Derivation Codes
				for _, item := range food.NutrientData {
					db.Where("Nutrientno=?", item.NutrientID).First(&nutr)
					item.NutrientID = int64(nutr.Id)
					nutr.Id = 0
					// set derivation
					if item.Derivation.Code != "" {
						db.Where("code=?", item.Derivation.Code).First(&derv)
						item.Derivation = derv
					}
					// set SourceCode
					if item.Sourcecode.Code != "" {
						db.Where("code=?", item.Sourcecode.Code).First(&source)
						//	log.Printf("%d source code=%s ID=%d",i,food.NutrientData[i].SourcID,source.ID)
						item.Sourcecode = source
					}
				}
				db.Create(&food)
				c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Food item created successfully!", "resourceId": food.Id})
			} else {
				log.Printf("err %v\n", err)
				c.JSON(http.StatusCreated, gin.H{"status": http.StatusNotFound, "message": "Can't parse json1!"})
			}
		})

		// create a nutrient
		v1.POST("/nutrient/", func(c *gin.Context) {
			var nutr bfpd.Nutrient
			var unit bfpd.Unit
			if c.BindJSON(&nutr) == nil {
				// check for unit and create if necessary
				db.Where("Unit=?", nutr.Unit.Unit).Find(&unit)
				if unit.Id != 0 {
					nutr.Unit.Id = unit.Id
					db.Save(&unit)
				}
				nutr.Unit = unit
				db.Create(&nutr)
				c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Nutrient item created successfully!", "resourceId": nutr.Id})
			} else {
				c.JSON(http.StatusCreated, gin.H{"status": http.StatusNotFound, "message": "Can't parse json!"})
			}
		})
	}
	router.Run()

}
func transformfood(f *bfpd.Food) bfpd.TransformedFood {
	var t *bfpd.TransformedNutrientData
	return bfpd.TransformedFood{ID: f.Id, Ndbno: f.Ndbno, Description: f.Description, Manufacturer: f.Manufacturer.Name, Source: f.Datasource, FoodGroup: bfpd.TransformedFoodGroup{Code: f.FoodGroup.Cd, Description: f.FoodGroup.Description}, Measures: f.Measures, Nutrients: t.Transform(&(f.NutrientData)), Ingredients: f.Ingredients.Description}
}
