package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

type User struct {
	MKABID        int    `json:"id"`
	NAME          string `json:"name"`
	FAMILY        string `json:"family"`
	OT            string `json:"ot"`
	SS            string `json:"ss"`
	S_DOC         string `json:"s_doc"`
	N_DOC         string `json:"n_doc"`
	DATE_BD       string `json:"date_bd"`
	ADRES         string `json:"adres"`
	Rf_kl_SexID   string `json:"sexid"`
	ContactEmail  string `json:"contactEmail"`
	ContactMPhone string `json:"contactMPhone"`
}

type File struct {
	Path string `json:"path"`
}

type ResponseFiles struct {
	Files []string `json:"files"`
	Error string   `json:"error,omitempty"`
}

type ResponseUsers struct {
	Users []User `json:"users"`
	Error string `json:"error,omitempty"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("не могу загрузить .env файл, продолжаю с переменными окружения")
	}

	query := url.Values{}
	query.Add("database", os.Getenv("DB_NAME"))

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASS")),
		Host:     net.JoinHostPort(os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		RawQuery: query.Encode(),
	}

	db, err := sql.Open("sqlserver", u.String())
	if err != nil {
		log.Fatalf("не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /list", handleList(db))
	mux.HandleFunc("POST /files", handleFiles(db))

	port := fmt.Sprintf("0.0.0.0:%s", os.Getenv("APP_PORT"))
	fmt.Printf("Сервер запущен на http://%s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}

func handleList(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			Query string `json:"query"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}

		if len(input.Query) < 3 {
			http.Error(w, "Строка запроса должна состоять не меньше чем из 3х символов", http.StatusBadRequest)
			return
		}

		users, err := getUsers(db, input.Query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := ResponseUsers{Users: users}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func getUsers(db *sql.DB, query string) ([]User, error) {
	var users []User

	rows, err := db.Query("SELECT MKABID, FAMILY, NAME, OT, SS, S_DOC, N_DOC, DATE_BD, ADRES, rf_kl_SexID, contactEmail, contactMPhone FROM hlt_MKAB WHERE SS LIKE @p1", "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.MKABID,
			&user.NAME,
			&user.FAMILY,
			&user.OT,
			&user.SS,
			&user.S_DOC,
			&user.N_DOC,
			&user.DATE_BD,
			&user.ADRES,
			&user.Rf_kl_SexID,
			&user.ContactEmail,
			&user.ContactMPhone); err != nil {
			return nil, fmt.Errorf("ошибка при чтении данных: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %v", err)
	}

	return users, nil
}

func handleFiles(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Query string `json:"query"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}

		files, err := getFiles(db, input.Query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := ResponseFiles{Files: files}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func getFiles(db *sql.DB, snils string) ([]string, error) {
	var files []string

	rows, err := db.Query(
		"SELECT atf_fileattachment.path FROM atf_fileattachment JOIN atf_fileinfo ON atf_fileattachment.rf_fileinfoid = atf_fileinfo.fileinfoid JOIN hlt_mkab ON hlt_mkab.uguid = atf_fileinfo.descguid WHERE hlt_mkab.SS = @p1",
		snils,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var file File
		if err := rows.Scan(&file.Path); err != nil {
			return nil, fmt.Errorf("ошибка при чтении данных: %v", err)
		}
		files = append(files, file.Path)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %v", err)
	}

	return files, nil
}
