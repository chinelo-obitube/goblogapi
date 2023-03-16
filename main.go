package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type BlogPost struct {
	ID         int       `gorm:"primary key" json:"id"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	// CategoryID int       `json:"category_id"`
	Category   string    `json:"category" gorm:"foreignkey:CategoryID"`
	// UserID     string    `json:"user_id"`
	User       string    `json:"user" gorm:"foreignkey:UserID"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Category struct {
	ID        int       `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        int       `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var db *gorm.DB
var err error

func main() {
	db, _ = gorm.Open("sqlite3", "./blog.db")
	if err != nil {
		fmt.Println("Failed to connect the database")
	}
	defer db.Close()

	db.AutoMigrate(&BlogPost{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Category{})

	router := gin.Default()

	router.GET("/", HomePage)
	router.GET("/posts/", GetAllPosts)
	router.POST("/posts/", CreatePost)
	router.GET("/posts/:id", GetPostById)
	router.PUT("/posts/:id", UpdatePost)
	router.DELETE("/posts/:id", DeletePost)

	router.GET("/users", GetAllUsers)
	router.POST("/users/", CreateUser)
	router.GET("/users/:id", GetUserById)
	

	router.GET("/categories", GetAllCategories)
	router.POST("/categories/", CreateCategory)
	router.GET("/categories/:id", GetCategoryById)
	

	router.Run(":4000")
}

func HomePage(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, gin.H{"message": "welcome to v1"})
}

func GetAllPosts(context *gin.Context) {
	var posts []BlogPost
	if err := db.Preload("User").Preload("Category").Find(&posts).Error; err != nil {
		context.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	} else {
		context.JSON(http.StatusOK, posts)
	}
}

func GetPostById(context *gin.Context) {
	id := context.Param("id")
	var post BlogPost
	if err := db.Preload("User").Preload("Category").Where("id = ?", id).First(&post).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "post not found"})
		return
	} else {
		context.JSON(http.StatusOK, post)
	}
}

func CreatePost(context *gin.Context) {
	var post BlogPost
	context.BindJSON(&post)

	db.Create(&post)
	context.JSON(http.StatusOK, post)
}

func UpdatePost(context *gin.Context) {
	var post BlogPost
	id := context.Param("id")

	if err := db.Where("id = ?", id).First(&post).Error; err != nil {
		context.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	}
	context.BindJSON(&post)
	db.Save(&post)
	context.JSON(http.StatusOK, post)
}

func DeletePost(context *gin.Context) {
	id := context.Param("id")
	var post BlogPost

	if err := db.First(&post, id).Error; err != nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	db.Delete(&post)
	context.Status(http.StatusOK)

}

func GetAllUsers(context *gin.Context) {
	var users []User

	db.Find(&users)
	context.JSON(http.StatusOK, users)

}

func GetUserById(context *gin.Context) {
	id := context.Param("id")
	var user User

	if err := db.First(&user, id).Error; err != nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	context.JSON(http.StatusOK, user)

}

func CreateUser(context *gin.Context) {
	var user User

	context.BindJSON(&user)
	db.Create(&user)
	context.JSON(http.StatusOK, user)

}

func GetAllCategories(context *gin.Context) {
	var categories []Category

	db.Find(&categories)
	context.JSON(http.StatusOK, categories)

}

func GetCategoryById(context *gin.Context) {
	id := context.Param("id")
	var category Category

	if err := db.First(&category, id).Error; err != nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	context.JSON(http.StatusOK, category)

}

func CreateCategory(context *gin.Context) {
	var category Category

	context.BindJSON(&category)
	db.Create(&category)
	context.JSON(http.StatusOK, category)

}
