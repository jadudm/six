package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/sync/semaphore"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var _DB *gorm.DB

func load_dotenv() {
	// https://github.com/joho/godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// https://stackoverflow.com/questions/64780409/can-i-use-go-gorm-with-mysql-json-type-field
type QueueJob struct {
	// gorm.Model
	Action  string            `json:"action" gorm:"-:all"`
	UUID    string            `json:"uuid"`
	Payload map[string]string `json:"payload" gorm:"serializer:json"`
}

func (qj *QueueJob) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

// Handler functions
func enqueue(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	job := map[string]string{"crawl": domain}

	qj := &QueueJob{
		Action:  "enqueue",
		UUID:    string(uuid.NewString()),
		Payload: job,
	}

	// Returns a DB pointer for chaining.
	_DB.Create(&qj)

	if err := render.Render(w, r, qj); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

}

var sem = semaphore.NewWeighted(int64(1))

func dequeue(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	sem.Acquire(ctx, 1)
	var qj *QueueJob
	result := _DB.Order("time_inserted ASC").Limit(1).Find(&qj)
	if result.Error == nil {
		// Delete if we found no errors
		_DB.Where("uuid = ?", qj.UUID).Delete(&QueueJob{})
	}
	sem.Release(1)
	qj.Action = "dequeue"
	if err := render.Render(w, r, qj); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func load_db() {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_PATH")), &gorm.Config{})
	if err != nil {
		panic("COULD NOT OPEN DATABASE")
	}
	_DB = db
}

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func main() {
	load_dotenv()
	load_db()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// This eats {domain} dots. Eg. `foo.com` becomes `foo`.
	// r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("search six"))
	})

	r.Put("/enqueue/{domain}", enqueue)
	r.Get("/dequeue", dequeue)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}
