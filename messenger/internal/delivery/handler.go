package delivery

import (
	"github.com/KirillMironov/rapu/messenger/domain"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Handler struct {
	service domain.ClientsService
}

func NewHandler(service domain.ClientsService) *Handler {
	return &Handler{service}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	var m = http.NewServeMux()
	m.HandleFunc("/connect", h.connect)
	return m
}

func (h *Handler) connect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userId := r.URL.Query().Get("userId")
	toUserId := r.URL.Query().Get("toUserId")

	if userId == "" || toUserId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var client = domain.Client{
		UserId:   userId,
		ToUserId: toUserId,
		Conn:     conn,
	}

	h.service.Connect(client)
}
