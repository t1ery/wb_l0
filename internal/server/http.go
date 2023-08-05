package server

import (
	"encoding/json"
	"github.com/t1ery/wb_l0/internal/service"
	"html/template"
	"net/http"
)

// HTTPServer представляет HTTP-сервер для обработки запросов клиентов.
type HTTPServer struct {
	srv *service.Service
}

// NewHTTPServer создает новый экземпляр HTTP-сервера с заданным сервисом.
func NewHTTPServer(srv *service.Service) *HTTPServer {
	return &HTTPServer{srv: srv}
}

// ServeHTTP обрабатывает входящий HTTP-запрос и маршрутизирует его в соответствии с URL.
func (hs *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		hs.handleIndexPage(w, r)
	case "/get_data":
		hs.handleGetDataByID(w, r)
	default:
		http.NotFound(w, r)
	}
}

// handleIndexPage обрабатывает запрос на главную страницу и отображает HTML-шаблон для ввода данных заказа.
func (hs *HTTPServer) handleIndexPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка при загрузке шаблона", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка при выполнении шаблона", http.StatusInternalServerError)
		return
	}
}

// handleGetDataByID обрабатывает запрос на получение данных заказа по идентификатору и возвращает результат в формате JSON.
func (hs *HTTPServer) handleGetDataByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	data, err := hs.srv.GetDataByID(id)
	if err != nil {
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
