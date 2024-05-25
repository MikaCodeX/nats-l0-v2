package rest

import (
	"fmt"
	"net/http"

	"awesome/internal/store"

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
	orderID := r1.URL.Query().Get("orderID")

	err, order := r.store.GetOrder(orderID)
	if err != nil {
		log.Print("Ошибка доставки")
		w.WriteHeader(404)
		w.Write(msgNoData)
		return
	}
	fmt.Fprintf(w, "Информация о вашем заказе:\n%s!", order)
}
