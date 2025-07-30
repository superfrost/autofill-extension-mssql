package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

type User struct {
	MKABID            int    `json:"id"`
	NAME              string `json:"name"`
	FAMILY            string `json:"family"`
	OT                string `json:"ot"`
	SS                string `json:"ss"`
	S_DOC             string `json:"s_doc"`
	N_DOC             string `json:"n_doc"`
	DATE_BD           string `json:"date_bd"`
	ADRES             string `json:"adres"`
	Rf_kl_SexID       string `json:"sexid"`
	ContactEmail      string `json:"contactEmail"`
	ContactMPhone     string `json:"contactMPhone"`
	Rf_TYPEDOCID      string `json:"rf_TYPEDOCID"`
	AdresFact         string `json:"adresFact"`
	DateDoc           string `json:"dateDoc"`
	DocIssuedBy       string `json:"docIssuedBy"`
	Rf_MilitaryDutyID string `json:"rf_MilitaryDutyID"`
}

type File struct {
	Path string `json:"path"`
}

type View struct {
	ViewData string `json:"view_data"`
	Sign     string `json:"sign"`
}

type ResponseViews struct {
	Views []View `json:"views"`
	Error string `json:"error,omitempty"`
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
	mux.HandleFunc("POST /views", handleViews(db))
	mux.HandleFunc("POST /addr", handleAddress(db))

	port := fmt.Sprintf("0.0.0.0:%s", os.Getenv("APP_PORT"))
	fmt.Printf("Сервер запущен на http://%s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}

func handleList(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

	rows, err := db.Query("SELECT MKABID, FAMILY, NAME, OT, SS, S_DOC, N_DOC, DATE_BD, ADRES, rf_kl_SexID, contactEmail, contactMPhone, rf_TYPEDOCID, AdresFact, DateDoc, DocIssuedBy, rf_MilitaryDutyID FROM hlt_MKAB WHERE SS LIKE @p1", "%"+query+"%")
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
			&user.ContactMPhone,
			&user.Rf_TYPEDOCID,
			&user.AdresFact,
			&user.DateDoc,
			&user.DocIssuedBy,
			&user.Rf_MilitaryDutyID,
		); err != nil {
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

func handleViews(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Query string `json:"query"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}

		views, err := getViews(db, input.Query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := ResponseViews{Views: views}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func getViews(db *sql.DB, snils string) ([]View, error) {
	var views []View

	rows, err := db.Query(
		"SELECT ViewData, Sign FROM hlt_MedRecord JOIN hlt_VisitHistory ON hlt_MedRecord.rf_VisitHistoryID=hlt_VisitHistory.VisitHistoryID JOIN hlt_MKAB ON hlt_VisitHistory.rf_MKABID=hlt_MKAB.MKABID WHERE hlt_MKAB.SS = @p1;",
		snils,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var view View
		if err := rows.Scan(
			&view.ViewData,
			&view.Sign,
		); err != nil {
			return nil, fmt.Errorf("ошибка при чтении данных: %v", err)
		}
		views = append(views, view)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %v", err)
	}

	return views, nil
}

func handleAddress(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Query string `json:"query"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}

		addr, err := getAddress(db, input.Query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(addr)
	}
}

func getAddress(db *sql.DB, address string) (map[string]string, error) {
	url := os.Getenv("COMPLETION_URL")
	apiKey := os.Getenv("AI_TOKEN")

	requestBody := map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": "Не генерируй ничего кроме json простой строкой без переноса строк и экранирования. Раздели адресс на следующие поля в формате ключь значение в формате json {\"Тип адреса гражданина\": , \"Почтовый индекс\": , \"Субъект Российской Федерации\": , \"Район субъекта РФ\": , \"Наименование населенного пункта\": , \"Улица\": ,  \"Дом (корпус, строение)\": , \"Квартира\": , \"Фактическое место нахождения гражданина\":  } Вот Адрес: " + address,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка при кодировании JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	var openAIResponse map[string]interface{}
	if err := json.Unmarshal(body, &openAIResponse); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании ответа JSON: %v", err)
	}

	choices := openAIResponse["choices"].([]interface{})
	if len(choices) == 0 {
		return nil, fmt.Errorf("нет доступного поля choices в ответе")
	}

	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	content := message["content"].(string)

	var addressInfo map[string]string
	if err := json.Unmarshal([]byte(content), &addressInfo); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании адреса JSON: %v", err)
	}

	return addressInfo, nil
}
