package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db    *pgxpool.Pool
	cache *Cache
}

type OrderUID struct {
	OrderUID string `json:"order_uid"`
}

func New() *Store {
	db, err := pgxpool.New(context.Background(), "postgres://postgres:1234@localhost:5432/db_shop")
	if err != nil {
		log.Fatalf("Ошибка подключения: %s", err)
	}
	return &Store{
		db:    db,
		cache: NewCache(),
	}
}

func (s *Store) SaveOrder(msg []byte) {
	msg1 := msg
	order_uid := getOrderUidFromJson(msg1)
	_, err := s.db.Exec(context.Background(), "INSERT INTO orders (order_uid,order_info) VALUES($1,$2)", order_uid, string(msg))
	log.Printf("Записали заказ: %s в базу данных", order_uid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert data: %v\n", err)
	}
	log.Println("Добавляем новый заказ в кэш")
	s.cache.Set(order_uid, msg)
	log.Printf("Заказ с номером: %s сохранен в кэш", order_uid)

}
func (s *Store) GetOrder(id string) ([]byte, error) {
	order, check := s.cache.Get(id)

	if !check {
		return nil, errors.New("такого заказа нету")

	}

	return order, nil
}

func (s *Store) GetOrderfromDB() {

	rows, err := s.db.Query(context.Background(), "SELECT order_uid, order_info FROM orders")
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

func (s *Store) Close() {
	s.db.Close()
}
