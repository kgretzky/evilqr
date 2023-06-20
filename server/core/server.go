package core

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

const API_TOKEN = "269884a8-db69-4dcd-a47f-003b2498a72e"

type QRCode struct {
	ID         string `json:"id"`
	Source     string `json:"source"`
	UpdateTime int64  `json:"update_time,omitifempty"`
	Authorized bool   `json:"authorized"`
	signal     chan struct{}
}

type HttpServer struct {
	r       *chi.Mux
	QRCodes sync.Map
}

func NewHttpServer() (*HttpServer, error) {
	o := &HttpServer{
		r:       chi.NewRouter(),
		QRCodes: sync.Map{},
	}
	return o, nil
}

func (o *HttpServer) Run(wwwdir string) {
	o.r.Use(middleware.Logger)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, wwwdir))
	o.FileServer(o.r, "/www/*", filesDir)

	o.r.Route("/qrcode", func(r chi.Router) {

		r.Route("/{id}", func(r chi.Router) {
			r.Use(o.qrcodeCtx)
			r.With(Authenticator(API_TOKEN)).Put("/", o.PutQRCode)
			r.Get("/", o.GetQRCode)
		})
	})

	http.ListenAndServe(":35000", o.r)
}

func (o *HttpServer) qrcodeCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var qrcode *QRCode = nil

		id := strings.ToLower(chi.URLParam(r, "id"))
		_, err := uuid.Parse(id)
		if err == nil {
			qrcode = o._getQRCode(id)
		} else {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "qrcode", qrcode)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (o *HttpServer) PutQRCode(w http.ResponseWriter, r *http.Request) {
	id := strings.ToLower(chi.URLParam(r, "id"))

	var data QRCode
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var doSignal bool = false
	qrcode := o._getQRCode(id)
	if qrcode == nil {
		qrcode = &data
		qrcode.signal = make(chan struct{})
	} else {
		doSignal = true
	}

	qrcode.Source = data.Source
	qrcode.UpdateTime = time.Now().UnixMilli()
	o.QRCodes.Store(id, qrcode)

	if doSignal {
		if qrcode.signal != nil {
			close(qrcode.signal)
			qrcode.signal = make(chan struct{})
		}
	}

	render.Render(w, r, NewQRCodeResponse(qrcode))
}

func (o *HttpServer) GetQRCode(w http.ResponseWriter, r *http.Request) {
	var err error
	var fromTime int64

	id := strings.ToLower(chi.URLParam(r, "id"))
	_fromTime := r.URL.Query().Get("t")
	if _fromTime != "" {
		fromTime, err = strconv.ParseInt(_fromTime, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	qrcode := o._getQRCode(id)
	if qrcode == nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	if fromTime > 0 {
		if fromTime >= qrcode.UpdateTime {
			ticker := time.NewTicker(60 * time.Second)

			select {
			case <-ticker.C:
				http.Error(w, "", http.StatusRequestTimeout)
				return
			case <-qrcode.signal:
			}

			qrcode = o._getQRCode(id)
		}
	}

	render.Render(w, r, NewQRCodeResponse(qrcode))
}

func (o *HttpServer) _getQRCode(id string) *QRCode {
	var qrcode *QRCode = nil
	if _qrcode, ok := o.QRCodes.Load(id); ok {
		qrcode = _qrcode.(*QRCode)
	}
	return qrcode
}

func (o *HttpServer) FileServer(r chi.Router, path string, root http.FileSystem) {
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

type QRCodeResponse struct {
	*QRCode
}

func NewQRCodeResponse(qrcode *QRCode) *QRCodeResponse {
	resp := &QRCodeResponse{QRCode: qrcode}
	return resp
}

func (o *QRCodeResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	return nil
}
