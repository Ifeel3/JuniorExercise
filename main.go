package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var address string

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("error:", err.Error())
		os.Exit(1)
	}
	username := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	db := os.Getenv("POSTGRES_DB")
	address = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, db)
}

type Position struct {
	id         int
	orderId    int
	name       string
	count      int
	base       string
	additional []string
}

func ExitWithError(pool *pgxpool.Pool, err error) {
	fmt.Println("Error:", err.Error())
	pool.Close()
	os.Exit(1)
}

func PrintMap(m map[string][]Position, orders []int) error {
	var builder strings.Builder
	tmp := []string{}
	for _, v := range orders {
		tmp = append(tmp, strconv.Itoa(v))
	}
	if _, err := builder.WriteString("=+=+=+=\nСтраница сборки заказов " + strings.Join(tmp, ",") + "\n"); err != nil {
		return err
	}
	for k, arr := range m {
		if _, err := builder.WriteString("===Стеллаж " + k + "\n"); err != nil {
			return err
		}
		for _, pos := range arr {
			if _, err := builder.WriteString(fmt.Sprintf("%s (id=%d)\nЗаказ %d, %d шт\n", pos.name, pos.id, pos.orderId, pos.count)); err != nil {
				return err
			}
			if len(pos.additional) != 0 {
				if _, err := builder.WriteString(fmt.Sprintf("доп стеллаж: %s\n", strings.Join(pos.additional, ","))); err != nil {
					return err
				}
			}
			if _, err := builder.WriteRune('\n'); err != nil {
				return err
			}
		}

	}
	fmt.Print(builder.String())
	return nil
}

func main() {
	pool, err := pgxpool.New(context.Background(), address)
	if err != nil {
		fmt.Println("Address:", address)
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	defer pool.Close()

	orders := []int{}
	if len(os.Args) < 2 {
		ExitWithError(pool, fmt.Errorf("not enough arguments"))
	}
	for k, str := range os.Args {
		if k == 0 {
			continue
		}
		tmpArr := strings.Split(str, ",")
		for _, v := range tmpArr {
			num, err := strconv.Atoi(v)
			if err != nil {
				ExitWithError(pool, err)
			}
			orders = append(orders, num)
		}
	}

	var builder strings.Builder
	if _, err = builder.WriteString("SELECT products.id AS \"id\", orders.id AS \"OrderID\", products.name AS \"Product\", positions.count AS \"Count\", shelves.name AS \"Base\", products.additional AS \"Additional\" FROM positions JOIN products ON positions.product_id = products.id JOIN shelves ON products.base = shelves.id JOIN orders ON positions.order_id = orders.id "); err != nil {
		ExitWithError(pool, err)
	}

	if len(orders) == 1 {
		if _, err = builder.WriteString("WHERE orders.id=" + strconv.Itoa(orders[0]) + ";"); err != nil {
			ExitWithError(pool, err)
		}
	} else {
		tmp := []string{}
		for _, v := range orders {
			tmp = append(tmp, strconv.Itoa(v))
		}
		if _, err := builder.WriteString("WHERE orders.id=" + strings.Join(tmp, " OR orders.id=") + ";"); err != nil {
			ExitWithError(pool, err)
		}
	}

	m := make(map[string][]Position, 0)
	result, err := pool.Query(context.Background(), builder.String())
	if err != nil {
		ExitWithError(pool, err)
	}
	for result.Next() {
		r := Position{}
		err := result.Scan(&r.id, &r.orderId, &r.name, &r.count, &r.base, &r.additional)
		if err != nil {
			fmt.Println("Error:", err.Error())
		}
		m[r.base] = append(m[r.base], r)
	}
	if err := PrintMap(m, orders); err != nil {
		ExitWithError(pool, err)
	}
}
