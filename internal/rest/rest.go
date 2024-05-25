package rest

import (
	"awesome/internal/store"
	"net/http"

	"github.com/rs/zerolog/log"
)

var (
	msgNoData = []byte(`{"message": "No data"}`)
)

type Rest struct {
	store *store.Store
}

func New(store *store.Store) *Rest {
	return &Rest{
		store: store,
	}
}
func (r *Rest) Run() {
	http.Handle("/", http.FileServer(http.Dir("../../static")))
	http.HandleFunc("/order", r.handler)
	http.ListenAndServe(":8080", nil)
}

func (r *Rest) handler(w http.ResponseWriter, r1 *http.Request) {
	id := r1.URL.Query().Get("orderID")

	order, err := r.store.GetOrder(id)
	if err != nil {
		log.Err(err).Msg("")
		w.WriteHeader(404)
		w.Write(msgNoData)
		return
	}
	w.Write(order)
}
