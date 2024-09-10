package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nestik55/task0/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	*pgxpool.Pool
}

func NewPostgres() *Postgres {

	db, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		"dmitry_db", "dmitry_db", "localhost", "5432", "dmitry_db"))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if err := db.Ping(context.Background()); err != nil {
		fmt.Println(err)
		return nil
	}

	postgres := Postgres{Pool: db}
	return &postgres
}

func (p *Postgres) SetOrder(ctx context.Context, order *model.Order) {

	query := `INSERT INTO "order" (
		"order_uid",
		"track_number",
		"entry",
		"locale",
		"internal_signature",
		"customer_id",
		"delivery_service",
		"shardkey",
		"sm_id",
		"date_created",
		"oof_shard") VALUES ($1, $2, $3, $4,$5,$6,$7,$8,$9,$10,$11)`

	if _, err := p.Pool.Exec(ctx, query,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard); err != nil {
		fmt.Println(err)
		return
	}

	if err := p.setDelivery(ctx, &order.Delivery, order.OrderUID); err != nil {
		fmt.Println(err)
		return
	}

	if err := p.setPayment(ctx, &order.Payment, order.OrderUID); err != nil {
		fmt.Println(err)
		return
	}

	if err := p.setItems(ctx, &order.Items, order.OrderUID); err != nil {
		fmt.Println(err)
		return
	}
}

func (p *Postgres) setDelivery(ctx context.Context, delivery *model.Delivery, orderUID string) error {

	query := `INSERT INTO "delivery"(
		"order_id",
		"phone",
		"name",
		"zip",
		"city",
		"address",
		"region",
		"email") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	if _, err := p.Pool.Exec(ctx, query,
		orderUID,
		delivery.Phone,
		delivery.Name,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email); err != nil {
		return errors.New("не удалось внести данные о Delivery")
	}

	return nil
}

func (p *Postgres) setPayment(ctx context.Context, payment *model.Payment, orderUID string) error {

	query := `INSERT INTO "payment"(
					  "order_id",
                      "transaction",
                      "request_id",
                      "currency",
                      "provider",
                      "amount",
                      "payment_dt",
                      "bank",
                      "delivery_cost",
                      "goods_total",
                      "custom_fee") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`

	if _, err := p.Pool.Exec(ctx, query,
		orderUID,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDt,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee); err != nil {
		return errors.New("не удалось внести данные о Payment")

	}
	return nil
}

func (p *Postgres) setItems(ctx context.Context, items *model.Items, orderUID string) error {
	query := `INSERT INTO "item"(
		"order_id",
		"chrt_id",
		"track_number",
		"price",
		"rid",
		"name",
		"sale",
		"size",
		"total_price",
		"nm_id",
		"brand",
		"status") VALUES ($1, $2, $3, $4,$5,$6,$7,$8,$9,$10,$11,$12)`

	for _, item := range *items {
		if _, err := p.Pool.Exec(ctx, query,
			orderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status); err != nil {
			return errors.New("не удалось внести данные об Item")
		}
	}

	return nil
}

func (p *Postgres) GetOrder(ctx context.Context, order *model.Order, orderUID string) error {
	if err := p.getDelivery(ctx, &order.Delivery, orderUID); err != nil {
		return err
	}

	if err := p.getPayment(ctx, &order.Payment, orderUID); err != nil {
		return err
	}

	order.Items, _ = p.getItems(ctx, orderUID)

	query := `SELECT 
		"order_uid",
		"track_number",
		"entry",
		"locale",
		"internal_signature",
		"customer_id",
		"delivery_service",
		"shardkey",
		"sm_id",
		"date_created",
		"oof_shard"
		FROM "order" WHERE ("order_uid"=$1)`

	if err := p.Pool.QueryRow(ctx, query, orderUID).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard); err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("строка не нашлась (заказ)")
		} else {
			return errors.New("не выполнился запрос поиска строки (заказ)")
		}
	}

	return nil
}

func (p *Postgres) getPayment(ctx context.Context, payment *model.Payment, orderUID string) error {
	query := `SELECT
		"transaction",
		"request_id",
		"currency",
		"provider",
		"amount",
		"payment_dt",
		"bank",
		"delivery_cost",
		"goods_total",
		"custom_fee"
		FROM "payment" WHERE ("order_id"=$1)`

	if err := p.Pool.QueryRow(ctx, query, orderUID).Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee); err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("строка не нашлась (payment)")
		} else {
			return errors.New("не выполнился запрос поиска строки (payment)")
		}
	}

	return nil
}

func (p *Postgres) getDelivery(ctx context.Context, delivery *model.Delivery, orderUID string) error {
	query := `SELECT
		"phone",
		"name",
		"zip",
		"city",
		"address",
		"region",
		"email"
		FROM "delivery" WHERE ("order_id"=$1)`

	if err := p.Pool.QueryRow(ctx, query, orderUID).Scan(
		&delivery.Phone,
		&delivery.Name,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email); err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("строка не нашлась (delivery)")
		} else {
			return err //errors.New("не выполнился запрос поиска строки (delivery)")
		}
	}

	return nil
}

func (p *Postgres) getItems(ctx context.Context, orderUID string) (model.Items, error) {
	var items model.Items

	query := `SELECT
		"chrt_id",
		"track_number",
		"price",
		"rid",
		"name",
		"sale",
		"size",
		"total_price",
		"nm_id",
		"brand",
		"status"
		FROM "item" WHERE ("order_id=$1")`

	rows, err := p.Pool.Query(ctx, query, orderUID)
	if err == pgx.ErrNoRows {
		return nil, errors.New("строка не нашлась (items)")
	} else if err != nil {
		return nil, errors.New("не выполнился запрос поиска строки (items)")
	}

	item := make(model.Items, 1)
	for rows.Next() {
		if err := rows.Scan(&item[0]); err != nil {
			return nil, errors.New("ошибка при сканировании (items)")
		}
		items = append(items, item[0])
	}

	return items, nil
}
