package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// テーブル定義
type User struct {
	gorm.Model
	Name string
	Age  int
}

func main() {
	godotenv.Load()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// テーブル作成（なければ）
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatal(err)
	}

	// テーブル一覧
	tables, err := db.Migrator().GetTables()
	if err != nil {
		log.Fatal(err)
	}
	for _, table := range tables {
		fmt.Println("[CREATE Table] name:", table)
	}

	// テーブルを空にする(SERIALもリセット)
	if err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY").Error; err != nil {
		log.Fatal(err)
	}

	// INSERT (1回目)
	var u User

	u = User{
		Name: "田中",
		Age:  20,
	}

	if err := db.Create(&u).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Println("[INSERT] id =", u.ID)

	// INSERT (2回目)
	u = User{
		Name: "鈴木",
		Age:  23,
	}

	if err := db.Create(&u).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Println("[INSERT] id =", u.ID)

	// SELECT (1回目)
	var users []User

	if err := db.Find(&users).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Println("[SELECT]")
	for _, user := range users {
		fmt.Printf(
			"id=%d name=%s age=%d\n",
			user.ID,
			user.Name,
			user.Age,
		)
	}

	// UPDATE
	var id, age int

	id = 1
	age = 21
	fmt.Println("[UPDATE] id =", id, "age =", age)
	if err := db.Model(&User{}).Where("id = ?", id).Update("age", age).Error; err != nil {
		log.Fatal(err)
	}

	// SELECT (2回目)
	if err := db.Find(&users).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Println("[SELECT]")
	for _, user := range users {
		fmt.Printf(
			"id=%d name=%s age=%d\n",
			user.ID,
			user.Name,
			user.Age,
		)
	}

	// DELETE
	id = 2
	fmt.Println("[DELETE] id =", id)
	if err := db.Delete(&User{}, id).Error; err != nil {
		log.Fatal(err)
	}

	// SELECT (3回目)
	if err := db.Find(&users).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Println("[SELECT]")
	for _, user := range users {
		fmt.Printf(
			"id=%d name=%s age=%d\n",
			user.ID,
			user.Name,
			user.Age,
		)
	}

	// Unscoped SELECT (1回目)
	if err := db.Unscoped().Find(&users).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Println("[Unscoped SELECT]")
	for _, user := range users {
		fmt.Printf(
			"id=%d name=%s age=%d deleted_at=%s\n",
			user.ID,
			user.Name,
			user.Age,
			user.DeletedAt.Time.Format("2006-01-02 15:04:05"),
		)
	}

	// Unscoped DELETE
	id = 2
	fmt.Println("[Unscoped DELETE] id =", id)
	if err := db.Unscoped().Delete(&User{}, id).Error; err != nil {
		log.Fatal(err)
	}

	// Unscoped SELECT (2回目)
	if err := db.Unscoped().Find(&users).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Println("[Unscoped SELECT]")
	for _, user := range users {
		fmt.Printf(
			"id=%d name=%s age=%d deleted_at=%s\n",
			user.ID,
			user.Name,
			user.Age,
			user.DeletedAt.Time.Format("2006-01-02 15:04:05"),
		)
	}
}
