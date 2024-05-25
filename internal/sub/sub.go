package sub

import (
	"log"

	"awesome/internal/store"

	"github.com/nats-io/stan.go"
)

type Sub struct {
	store *store.Store
	conn  stan.Conn
}

const (
	clusterID = "test-cluster"
	clientId  = "subscriber"
)

func New(store *store.Store) *Sub {
	conn, err := stan.Connect(
		clusterID,
		clientId,
		stan.NatsURL("nats://localhost:4223"),
	)
	if err != nil {
		log.Fatalf("Ошибка подключения:%s", err)
	}

	return &Sub{
		store: store,
		conn:  conn,
	}
}
func (s *Sub) Run() {

	_, err := s.conn.Subscribe(
		"Orders",
		s.handler,
		stan.DurableName("service orders"),
	)
	if err != nil {
		log.Fatalf("Подписчик не смог подрубиться:%s", err)
	}
	log.Println("Подписчик на месте")

}
func (s *Sub) Close() {
	s.conn.Close()
}
func (s *Sub) handler(m *stan.Msg) {
	s.store.SaveOrder(m.Data)
}
