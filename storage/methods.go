package storage

import (
	"database/sql"
	"fmt"
	"log"
)

func (s *Storage) Save(order *Order) error {
	tr, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tr.Rollback()

	err = s.SaveOrder(tr, order)
	if err != nil {
		return err
	}

	err = s.saveDelivery(tr, order)
	if err != nil {
		return err
	}

	err = s.savePayment(tr, order)
	if err != nil {
		return err
	}

	err = s.saveItems(tr, order)
	if err != nil {
		return err
	}

	if err := tr.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction:%w", err)
	}

	log.Printf("Order %s saved", order.OrderUID)
	return nil

}

func (s *Storage) SaveOrder(tr *sql.Tx, order *Order) error {
	query := `
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, 
                           customer_id, delivery_service, "shardkey", sm_id, date_created, "oof_shard")
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        ON CONFLICT (order_uid) DO UPDATE SET
            track_number = EXCLUDED.track_number,
            entry = EXCLUDED.entry,
            locale = EXCLUDED.locale,
            internal_signature = EXCLUDED.internal_signature,
            customer_id = EXCLUDED.customer_id,
            delivery_service = EXCLUDED.delivery_service,
            "shardkey" = EXCLUDED."shardkey",
            sm_id = EXCLUDED.sm_id,
            date_created = EXCLUDED.date_created,
            "oof_shard" = EXCLUDED."oof_shard"
	`

	_, err := tr.Exec(query,
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
		order.OofShard,
	)
	if err != nil {
		return err
	}
	s.UpdateCache(order)
	log.Printf("order %s saved", order.OrderUID)
	return nil
}

func (s *Storage) saveDelivery(tr *sql.Tx, order *Order) error {
	query := `
			INSERT INTO delivery(order_uid,name,phone,zip,city,address,region,email)
			VALUES($1,$2,$3,$4,$5,$6,$7,$8)
			ON CONFLICT (order_uid)DO UPDATE SET
			name = EXCLUDED.name,
			phone = EXCLUDED.phone,
			zip = EXCLUDED.zip,
			city = EXCLUDED.city,
			address = EXCLUDED.address,
			region = EXCLUDED.region,
			email = EXCLUDED.email

	`

	_, err := tr.Exec(query,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)

	return err
}

func (s *Storage) savePayment(tr *sql.Tx, order *Order) error {
	query := `
		INSERT INTO payment (order_uid,transaction,request_id,currency,
								provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT(order_uid) DO UPDATE SET
		transaction = EXCLUDED.transaction,
            request_id = EXCLUDED.request_id,
            currency = EXCLUDED.currency,
            provider = EXCLUDED.provider,
            amount = EXCLUDED.amount,
            payment_dt = EXCLUDED.payment_dt,
            bank = EXCLUDED.bank,
            delivery_cost = EXCLUDED.delivery_cost,
            goods_total = EXCLUDED.goods_total,
            custom_fee = EXCLUDED.custom_fee
	`
	_, err := tr.Exec(query,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	return err
}

func (s *Storage) saveItems(tr *sql.Tx, order *Order) error {
	deleteQ := "DELETE FROM items WHERE order_uid = $1"
	_, err := tr.Exec(deleteQ, order.OrderUID)
	if err != nil {
		return err
	}

	insertQ := `
			INSERT INTO items (order_uid,chrt_id,track_number,price,rid, name,
								sale, size,total_price, nm_id,brand,status)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)					
	`
	for _, item := range order.Items {
		_, err := tr.Exec(insertQ,
			order.OrderUID,
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
			item.Status,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) GetOrderByUID(uid string) (*Order, error) {
	order := &Order{}

	query := `SELECT order_uid, track_number, entry, locale, internal_signature, 
               customer_id, delivery_service, "shardkey", sm_id, date_created, "oof_shard"
        FROM orders 
        WHERE order_uid = $1`

	err := s.db.QueryRow(query, uid).Scan(
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
		&order.OofShard,
	)
	if err != nil {
		return nil, err
	}

	err = s.getDelivery(order)
	if err != nil {
		return nil, err
	}

	err = s.getPayment(order)
	if err != nil {
		return nil, err
	}

	err = s.getItems(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Storage) getDelivery(order *Order) error {
	query := `
		SELECT name,phone,zip,city,address,region,email
		FROM delivery
		WHERE order_uid = $1
	`

	return s.db.QueryRow(query, order.OrderUID).Scan(
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
	)
}

func (s *Storage) getPayment(order *Order) error {
	query := `
		SELECT transaction, request_id,currency, provider,amount,
		payment_dt, bank,delivery_cost,goods_total, custom_fee
		FROM payment
		WHERE order_uid = $1
	`

	return s.db.QueryRow(query, order.OrderUID).Scan(
		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDt,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)
}

func (s *Storage) getItems(order *Order) error {
	query := `
		SELECT chrt_id,track_number,price,rid,name,sale,
				size,total_price,nm_id,brand,status
		FROM items
		WHERE order_uid = $1			
	`

	rows, err := s.db.Query(query, order.OrderUID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return err
		}
		items = append(items, item)
	}

	order.Items = items
	return nil
}

func (s *Storage) AllOrders() ([]Order, error) {
	query := "SELECT order_uid FROM orders"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}

		order, err := s.GetOrderByUID(uid)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}
	return orders, nil
}

func (s *Storage) LoadCache() error {
	orders, err := s.AllOrders()
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range orders {
		s.cache[orders[i].OrderUID] = &orders[i]
	}

	log.Printf("Loaded %d orders into cache", len(orders))
	return nil
}

func (s *Storage) GetOrderCache(uid string) (*Order, error) {
	s.mu.RLock()
	order, exists := s.cache[uid]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("order not found")
	}
	return order, nil
}

func (s *Storage) UpdateCache(order *Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[order.OrderUID] = order
}
