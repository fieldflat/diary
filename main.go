package main

import (
	"log"
	"net/http"
	_ "net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/go-playground/validator.v9"
)

// Item は出品アイテム情報
type Item struct {
	gorm.Model
	Title       string `validate:"required"` // 商品タイトル
	Description string `validate:"required"` // 商品説明
	Point       string `validate:"required"` // 商品価格
	CreatedTime string `validate:"required"` // 作成日時
	UpdatedTime string `validate:"required"` // 更新日時
}

// Validate about Item structure.
func (form *Item) Validate() (ok bool, result map[string]string) {
	result = make(map[string]string)
	// 構造体のデータをタグで定義した検証方法でチェック
	// err := validator.New().Struct(*form)
	validate := validator.New()
	// validate.RegisterValidation("is_tarou", tarou) //第一引数をvalidateタグで設定した名前に合わせる
	err := validate.Struct(*form)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		if len(errors) != 0 {
			for i := range errors {
				// フィールドごとに、検証
				switch errors[i].StructField() {
				case "Title":
					result["Title"] = "タイトルの入力は必須です．"
				case "Description":
					result["Description"] = "本文の入力は必須です．"
				case "Point":
					result["Point"] = "点数の入力は必須です．"
				}
			}
		}
		return false, result
	}
	return true, result
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
func create(title string, description string, point string, createdTime string, updatedTime string) map[string]string {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}

	// 処理を追加
	form := Item{
		Title:       title,
		Description: description,
		Point:       point,
		CreatedTime: createdTime,
		UpdatedTime: updatedTime,
	}

	// バリデーションの検証を行う
	ok, errorMessages := form.Validate()
	if !ok {
		log.Print("入力エラーあり")
		log.Print(errorMessages)
		return errorMessages
	}

	log.Print("入力エラーなし！！")
	db.Create(&Item{Title: title, Description: description, Point: point, CreatedTime: createdTime, UpdatedTime: updatedTime})
	return errorMessages
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

		errors := create(title, description, point, time.Now().Format("2006/1/2 15:04:05"), time.Now().Format("2006/1/2 15:04:05"))
		c.Redirect(302, "/", gin.H{
			"errors": errors,
		}))
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
