package communication

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type API interface {
	Run(port int)
	RunWithHTTPS(port int)
	Database() *sql.DB
	Router() *http.ServeMux
	SetDatabase(db *sql.DB)
	SetRouter(router *http.ServeMux)
	CloseDatabase() error
	AddCustomHandler(route string, handler func(w http.ResponseWriter, r *http.Request))
	initialize()
}

type RESTInterface struct {
	router *http.ServeMux
	db     *sql.DB
}

func NewRestInterface(db *sql.DB) *RESTInterface {
	return &RESTInterface{
		router: http.NewServeMux(),
		db:     db,
	}
}

func (rest *RESTInterface) Run(port int) {
	rest.initialize()
	log.Printf("<REST>: HTTP serve at %d\n", port)
	addr := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(addr, rest.router)
	if err != nil {
		panic(err)
	}
}

func (rest *RESTInterface) RunWithHTTPS(port int) {
	rest.initialize()
	log.Printf("<REST>: HTTPS serve at %d\n", port)
	addr := fmt.Sprintf(":%d", port)
	err := http.ListenAndServeTLS(addr, "server.crt", "server.key", rest.router)
	if err != nil {
		panic(err)
	}
}

func (rest *RESTInterface) Database() *sql.DB {
	return rest.db
}

func (rest *RESTInterface) Router() *http.ServeMux {
	return rest.router
}

func (rest *RESTInterface) SetDatabase(db *sql.DB) {
	rest.db = db
}

func (rest *RESTInterface) SetRouter(router *http.ServeMux) {
	rest.router = router
}

func (rest *RESTInterface) CloseDatabase() error {
	err := rest.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response, _ := json.Marshal(payload)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func (rest *RESTInterface) AddCustomHandler(route string, handler func(w http.ResponseWriter, r *http.Request)) {
	rest.router.HandleFunc(route, handler)
}

func (rest *RESTInterface) initialize() {}
