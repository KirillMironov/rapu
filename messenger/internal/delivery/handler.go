package delivery

import (
	"github.com/KirillMironov/rapu/messenger/domain"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Handler struct {
	service domain.ClientsService
	logger  Logger
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

func NewHandler(service domain.ClientsService, logger Logger) *Handler {
	return &Handler{service, logger}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	var m = http.DefaultServeMux
	m.HandleFunc("/connect", h.connect)
	return m
}

func (h *Handler) connect(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	toUserId := r.URL.Query().Get("toUserId")
	if userId == "" || toUserId == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Info("not enough query parameters")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	var client = domain.Client{
		UserId:   userId,
		ToUserId: toUserId,
		Conn:     conn,
	}

	h.service.Connect(client)
}
