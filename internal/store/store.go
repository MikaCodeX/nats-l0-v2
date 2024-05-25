package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v4"
)

type Store struct {
	cache *Cache
}

type OrderUID struct {
	OrderUID string `json:"order_uid"`
}

func New() *Store {
	return &Store{
		cache: NewCache(),
	}
}

func (s *Store) SaveOrder(msg []byte) {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:1234@localhost:5432/db_shop")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	log.Println("Подрубились к дб")
	defer conn.Close(context.Background())

	msg1 := msg
	order_uid := getOrderUidFromJson(msg1)

	_, err = conn.Exec(context.Background(), "INSERT INTO orders (order_uid,order_info) VALUES($1,$2)", order_uid, string(msg))
	log.Printf("Записали заказ: %s в базу данных", order_uid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert data: %v\n", err)
	}
	log.Println("Берем старые заказы из бд и добавляем в кэш")
	s.GetOrderfromDB(conn)
	log.Println("Добавляем новый заказ в кэш")
	s.cache.Set(order_uid, msg)
	log.Printf("Заказ с номером: %s сохранен в кэш", order_uid)
	s.GetOrderfromDB(conn)

}
func (s *Store) GetOrder(id string) ([]byte, error) {
	order, check := s.cache.Get(id)

	if !check {
		return nil, errors.New("такого заказа нету")

	}

	return order, nil
}

func (s *Store) GetOrderfromDB(conn *pgx.Conn) []byte {
	var sliceOrders []byte
	rows, err := conn.Query(context.Background(), "SELECT order_uid, order_info FROM orders")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var orderUID string
		var orderInfo []byte
		err = rows.Scan(&orderUID, &orderInfo)
		if err != nil {
			panic(err)
		}
		s.cache.Set(orderUID, orderInfo)

	}

	return sliceOrders
}

func getOrderUidFromJson(m []byte) string {
	replaceValue := strings.ReplaceAll(string(m), "`", "\"")
	var result OrderUID
	err3 := json.Unmarshal([]byte(replaceValue), &result)
	if err3 != nil {
		log.Fatalf("Error parsing JSON: %s", err3)
	}
	return result.OrderUID
}
