package order_pg

import (
	"database/sql"
	"assignment2/models"
	"assignment2/repository/order_repository"
)

type orderPG struct {
	db *sql.DB
}

const (
	createOrderQuery = `
	INSERT INTO "orders" ("ordered_at", "customer_name)
	VALUES ($1, $2)
	RETURNING "order_id"
	`

	creaeteItemQuery = `
	INSERT INTO "items" ("item_code", "description", "quantity", "order_id")
	VALUES ($1, $2, $3, $4)
	`
)

func NewOrderPG(db *sql.DB) order_repository.Repository {
	return &orderPG{db: db}
}

func (orderPG *orderPG) CreateOrder(orderPayload models.Order, itemPayload []models.Item) error {
	dbTransaction, err := orderPG.db.Begin()
	if err != nil {
		return err
	}
	var orderId int

	orderRow := dbTransaction.QueryRow(createOrderQuery, orderPayload.OrderedAt, orderPayload.CustomerName)
	
	err = orderRow.Scan(&orderId)
	if err != nil {
		dbTransaction.Rollback()
		return err
	}

	for _, eachItem := range itemPayload {
		_, err := dbTransaction.Exec(creaeteItemQuery, eachItem.ItemCode, eachItem.Description, eachItem.Quantity, orderId)
		if err != nil {
			dbTransaction.Rollback()
			return err
		}
	}
	err = dbTransaction.Commit()
	if err != nil {
		dbTransaction.Rollback()
		return err
	}
	return nil
}