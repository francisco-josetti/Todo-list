package main 

import(
	"html/template"
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var store = sessions.NewCookieStore([]byte("secret"))
type Task struct{
	ID int 
	Title string 
	Done bool 
}

type User struct{
	ID int
	Email string
	Password string
}

type AuthPageData struct {
	Error string
}
func getTasks() []Task{
	rows, err := db.Query("SELECT id, title, done FROM tasks")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Done)
		if err != nil {
			panic(err)
		}
		tasks = append(tasks, task)
	}
	return tasks
}

func userTasks(userID int) []Task{
	rows, err := db.Query("SELECT id, title, done FROM tasks WHERE user_id = ?", userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Done)
		if err != nil {
			panic(err)
		}
		tasks = append(tasks, task)
	}
	return tasks
}
func handler(w http.ResponseWriter, r *http.Request){
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	session, _ := store.Get(r, "session")
	userID,ok := session.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var user User
	err := db.QueryRow("SELECT id, email FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/base.html","templates/index.html"))
	tmpl.ExecuteTemplate(w, "base", userTasks(user.ID))
}

func registerHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" {
		email := r.FormValue("user")
		password := r.FormValue("pass")
		confirmedPassword := r.FormValue("confirm-pass")
		if password != confirmedPassword {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_,err = db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", email, hashedPassword)
		if err != nil {
	tmpl := template.Must(template.ParseFiles("templates/base.html","templates/auth.html"))
	data := AuthPageData{
		Error: "Email already exists.",
	}
	tmpl.ExecuteTemplate(w, "base", data)
	return
}
		
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/base.html","templates/auth.html"))
	tmpl.ExecuteTemplate(w, "base", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("user")
		password := r.FormValue("pass")
		var user User
		err := db.QueryRow("SELECT id,password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Password)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
	http.Error(w, "Invalid email or password", http.StatusUnauthorized)
	return
}

		session, _ := store.Get(r, "session")
		session.Values["user_id"] = user.ID
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/base.html","templates/login.html"))
	data := AuthPageData{
	Error: "Email ou senha incorretos.",
}
	tmpl.ExecuteTemplate(w,"base",data)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func isAuthenticated(r *http.Request) bool {
	session, _ := store.Get(r, "session")
	userID := session.Values["user_id"]
	return userID != nil
}
func createHandler(w http.ResponseWriter, r *http.Request){
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	session, _ := store.Get(r, "session")
	userID,ok := session.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method == "POST" {
		title := r.FormValue("title")
		query := `INSERT INTO tasks (title, done, user_id) VALUES (?, ?, ?)`
		_,err := db.Exec(query, title, false,userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/base.html","templates/create.html"))
	tmpl.ExecuteTemplate(w, "base", nil)
}

func doneHandler(w http.ResponseWriter, r *http.Request){
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	id := r.URL.Query().Get("id")
	query := `UPDATE tasks SET done = 1 WHERE id = ?`
	_,err := db.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func removeHandler(w http.ResponseWriter, r *http.Request){
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	query := `DELETE FROM tasks WHERE done = 1`
	_,err := db.Exec(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func main(){
	db, _ = sql.Open("sqlite3", "./tasks.db")
	queryTasks := `CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		done BOOLEAN NOT NULL,
		user_id INTEGER
	);`
	queryUsers := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`
	_, err := db.Exec(queryTasks)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(queryUsers)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	http.HandleFunc("/", handler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/done", doneHandler)
	http.HandleFunc("/remover", removeHandler)
	http.HandleFunc("/auth", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}