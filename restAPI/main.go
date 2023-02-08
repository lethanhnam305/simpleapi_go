package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	models "simpleapi_go/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// -------------------------------- DB SET UP----------------------------//
var db *sql.DB
var err error

// --------------------------------- METHOD	-----------------------------	//

// ---------------------------GET ALL--------Í-------------//
func getArticles(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Getting all Article...")
	rows, err := db.Query("SELECT * FROM ARTICLE")
	// check errors
	if err != nil {
		log.Fatal(err)
	}

	var articles []models.Article

	for rows.Next() {
		var (
			Id          string
			Title       string
			Description string
			Content     string
		)

		err = rows.Scan(&Id, &Title, &Description, &Content)

		// check errors
		if err != nil {
			log.Fatal(err)
		}

		articles = append(articles, models.Article{Id: Id, Title: Title, Description: Description, Content: Content})
	}

	var response = models.JsonResponse{Type: "success", Data: articles}

	json.NewEncoder(w).Encode(response)

	fmt.Println("Done!")
}

// --------------------------GET BY ID---------------------//
func getArticleById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	fmt.Println("Finding ID of Article...")
	var articles []models.Article

	result, err := db.Query("SELECT * FROM ARTICLE WHERE ID =?", params["id"])
	if err != nil {
		log.Fatal(err)
	}

	var (
		Id          string
		Title       string
		Description string
		Content     string
	)
	// check errors
	for result.Next() {
		err = result.Scan(&Id, &Title, &Description, &Content)
		if err != nil {
			log.Fatal(err)
		}
		//log.Fatal(err)
		articles = append(articles, models.Article{Id: Id, Title: Title, Description: Description, Content: Content})
	}
	//ID not found
	if len(articles) == 0 {
		var response = models.JsonResponse{Type: "Failed, ID not exists!", Data: nil}
		json.NewEncoder(w).Encode(response)
		return
	}
	//Find ID
	var response = models.JsonResponse{Type: "Success", Data: articles}
	json.NewEncoder(w).Encode(response)

	fmt.Println("Done!")
}

// -------------------------CREATE ARTICLE----------------------//
func createArticle(w http.ResponseWriter, r *http.Request) {

	stmt, err := db.Prepare("INSERT INTO ARTICLE(ID, TITLE, DESCRIPTION, CONTENT) VALUES(?, ?, ?, ?)")

	if err != nil {
		log.Fatal(err)
	}

	var articles []models.Article
	//Value to create
	var article models.Article
	json.NewDecoder(r.Body).Decode(&article)

	//Query ID
	row, err := db.Query("SELECT * FROM ARTICLE WHERE ID = ?", article.Id)

	for row.Next() {
		var (
			Id_t          string
			Title_t       string
			Description_t string
			Content_t     string
		)

		err = row.Scan(&Id_t, &Title_t, &Description_t, &Content_t)
		if err != nil {
			log.Fatal(err)
		}
		articles = append(articles, models.Article{Id: Id_t, Title: Title_t, Description: Description_t, Content: Content_t})
	}
	//Successfully added
	if len(articles) == 0 {
		res, err := stmt.Exec(article.Id, article.Title, article.Description, article.Content)
		_ = res
		if err != nil {
			log.Fatal(err)
		}
		var response = models.JsonResponse{Type: "Successfully created an article", Data: nil}
		json.NewEncoder(w).Encode(response)
		return
	}

	//Failed
	var response = models.JsonResponse{Type: "Failed, existing ID!", Data: nil}
	json.NewEncoder(w).Encode(response)
	return
}

// --------------------------DELETE ARTICLE-----------------------//
func deleteArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Deleting an article...")

	params := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM ARTICLE WHERE ID = ?")

	row, err := db.Query("SELECT * FROM ARTICLE WHERE ID = ?", params["id"])
	if err != nil {
		log.Fatal(err)
	}

	var articles []models.Article
	for row.Next() {
		var (
			Id_t          string
			Title_t       string
			Description_t string
			Content_t     string
		)

		err = row.Scan(&Id_t, &Title_t, &Description_t, &Content_t)
		if err != nil {
			log.Fatal(err)
		}
		articles = append(articles, models.Article{Id: Id_t, Title: Title_t, Description: Description_t, Content: Content_t})
	}
	//Success
	if len(articles) != 0 {
		res, err := stmt.Exec(params["id"])
		_ = res
		if err != nil {
			log.Fatal(err)
		}
		var response = models.JsonResponse{Type: "Successfully deleted an article!", Data: nil}
		json.NewEncoder(w).Encode(response)
		return
	}
	//Failed
	var response = models.JsonResponse{Type: "Failed, ID not exist!", Data: nil}
	json.NewEncoder(w).Encode(response)
}

// -------------------------------DELETE ALL ------------------//
func deleteAllArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete all article!")
	stmt, err := db.Exec("DELETE FROM ARTICLE")
	_ = stmt
	if err != nil {
		log.Fatal(err)
	}
	var response = models.JsonResponse{Type: "Deleted all article!", Data: nil}
	json.NewEncoder(w).Encode(response)
}

// -------------------------------- UPDATE --------------------//
func updateArticle(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Updating an article!")

	var article models.Article
	json.NewDecoder(r.Body).Decode(&article)

	stmt, err := db.Prepare("UPDATE ARTICLE SET TITLE = ?,  DESCRIPTION = ?, CONTENT = ? WHERE ID = ?")
	if err != nil {
		log.Fatal(err)
	}

	row, err := db.Query("SELECT * FROM ARTICLE WHERE ID = ?", article.Id)
	if err != nil {
		log.Fatal(err)
	}

	var articles []models.Article
	for row.Next() {
		var (
			Id_t          string
			Title_t       string
			Description_t string
			Content_t     string
		)
		err = row.Scan(&Id_t, &Title_t, &Description_t, &Content_t)
		if err != nil {
			log.Fatal(err)
		}

		articles = append(articles, models.Article{Id: Id_t, Title: Title_t, Description: Description_t, Content: Content_t})
	}
	//Update Success
	if len(articles) != 0 {
		res, err := stmt.Exec(article.Title, article.Description, article.Content, article.Id)
		_ = res
		if err != nil {
			log.Fatal(err)
		}
		var response = models.JsonResponse{Type: "Successfully updated article!", Data: nil}
		json.NewEncoder(w).Encode(response)
		return
	}
	//Update failed
	var response = models.JsonResponse{Type: "Failed, ID not exist!", Data: nil}
	json.NewEncoder(w).Encode(response)
}

// ----------------------------- MAIN -----------------------------//
func main() {
	db_user := "root"
	db_pass := ""
	db_name := "simpleapi_go"

	db, err = sql.Open("mysql", db_user+":"+db_pass+"@tcp(127.0.0.1:3306)/"+db_name)
	if err = db.Ping(); err != nil {
		fmt.Println("open database fail")
		log.Fatal(err)
		return
	}

	fmt.Println("We're in the DB!")
	defer db.Close()
	// Init the mux router
	router := mux.NewRouter()

	// Get all article
	router.HandleFunc("/articles", getArticles).Methods("GET")
	// Get by Id
	router.HandleFunc("/articles/{id}", getArticleById).Methods("GET")
	// Create an article
	router.HandleFunc("/articles", createArticle).Methods("POST")
	// Delete a specific article by the Id
	router.HandleFunc("/articles/{id}", deleteArticle).Methods("DELETE")
	// Delete all
	router.HandleFunc("/articles", deleteAllArticle).Methods("DELETE")
	//Update an article
	router.HandleFunc("/articles", updateArticle).Methods("PUT")
	// serve the app
	http.ListenAndServe(":8001", router)
}
