package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"strings"
)

type Urls struct {
	id  int
	URL string
}

func main() {
	ConnectionDB := "user=guodsztyeofrzj password=128219cf5f8b003a4ab5cf3f9ce71ace2f0a48dbc527ab9d40a43dcecf20233a host=ec2-54-247-96-153.eu-west-1.compute.amazonaws.com dbname=d5p8j8r0ll59jk"
	db, err := sql.Open("postgres", ConnectionDB)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nSuccessfully connected to database!\n")
	defer db.Close()

	var value int
	fmt.Println("1. Post\n2. Get")
	fmt.Fscan(os.Stdin, &value)
	switch value {
	case 1:
		fmt.Println("Введите полную ссылку")
		var url string
		fmt.Fscan(os.Stdin, &url)
		fmt.Println("Сокращенная ссылка: ")
		Post(db, url)
		break
	case 2:
		fmt.Println("Введите сокращенную ссылку")
		var url string
		fmt.Fscan(os.Stdin, &url)
		fmt.Println("Полная ссылка: " + Get(db, url))

	}

}

func Post(db *sql.DB, value string) {
	urls := []Urls{}
	result, _ := db.Query("select  url from urls where url = $1", value)
	//Проверка на существующий URL в БД
	for result.Next() {
		url := Urls{}
		err := result.Scan(&url.URL)
		if err != nil {
			fmt.Println(err)
			continue
		}
		urls = append(urls, url)
		for _, url := range urls {
			if url.URL == value {
				value = "Такой URL уже существует"
				fmt.Println(value)
				break
			}
		}
		//Если в БД такого нет, то выводим сокращенную ссылку
		if value != "Такой URL уже существует" {
			_, err := db.Exec("insert into URLs (url) values ($1)", value)
			if err != nil {
				panic(err)
			}
			fmt.Println("Ссылка добавлена")
			test := strings.Split(value, "")
			ar := make([]string, 0)
			for i := 0; i < 10; i++ {
				ar = append(ar, test[i])
			}
			value = strings.Join(ar, "")
			fmt.Println(value)
		}
	}
}

//Нахождение по сокращенной ссылку
func Get(db *sql.DB, value string) string {
	var urls []Urls
	NewValue := "Такого значения нет"
	result, _ := db.Query("select url from URLs")
	for result.Next() {
		url := Urls{}
		err := result.Scan(&url.URL)
		if err != nil {
			fmt.Println(err)
			continue
		}
		urls = append(urls, url)
	}

	for _, url := range urls {
		if strings.Contains(url.URL, value) {
			NewValue = url.URL
			break
		}
	}
	return NewValue
}
