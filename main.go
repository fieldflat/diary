package main

import (
	"net/http"
	_ "net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// Item は出品アイテム情報
type Item struct {
	gorm.Model
	Title       string // 商品タイトル
	Description string // 商品説明
	Point       string // 商品価格
	CreatedTime string // 作成日時
	UpdatedTime string // 更新日時
}

// DBの初期化処理
func dbInit() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}

	db.AutoMigrate(&Item{})
}

// create関数
func create(title string, description string, point string, createdTime string, updatedTime string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	db.Create(&Item{Title: title, Description: description, Point: point, CreatedTime: createdTime, UpdatedTime: updatedTime})
}

func update(id int, title string, description string, point string, updatedTime string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	var item Item
	db.First(&item, id)
	item.Title = title
	item.Description = description
	item.Point = point
	item.UpdatedTime = updatedTime
	db.Save(&item)
	db.Close()
}

// 全てのItemを取得する
func getAll() []Item {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	var items []Item
	db.Order("created_at desc").Find(&items)
	return items
}

// main関数
func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/main/*")
	r.Static("/assets", "./assets")
	dbInit()

	// 一覧取得
	r.GET("/", func(c *gin.Context) {
		items := getAll()
		c.HTML(200, "index.tmpl", gin.H{
			"items": items,
		})
	})

	// 新規作成
	r.POST("/new", func(c *gin.Context) {
		title := c.PostForm("title")
		description := c.PostForm("description")
		point := c.PostForm("point")

		create(title, description, point, time.Now().Format("2006/1/2 15:04:05"), time.Now().Format("2006/1/2 15:04:05"))
		c.Redirect(302, "/")
	})

	r.POST("/update/:id", func(c *gin.Context) {
		title := c.PostForm("title")
		description := c.PostForm("description")
		point := c.PostForm("point")
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("failed to get id\n")
		}

		update(id, title, description, point, time.Now().Format("2006/1/2 15:04:05"))
		c.Redirect(302, "/")
	})

	// assets フォルダの読み取り
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	r.Run()
}
