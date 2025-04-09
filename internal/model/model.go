package model

import (
	"go_store/generated/proto/common"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type OrderStatus int

const (
	UNSPECIFIED OrderStatus = iota
	PENDING
	PROCESSING
	COMPLETED
	CANCELLED
)

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

type OrderItem struct {
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

type Order struct {
	ID            string      `json:"id"`
	CustomerName  string      `json:"customer_name"`
	CustomerEmail string      `json:"customer_email"`
	Items         []OrderItem `json:"items"`
	Status        OrderStatus `json:"status"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

func (o *Order) ConvertToMessage() *common.Order {
	items := make([]*common.OrderItem, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, &common.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return &common.Order{
		Id:            o.ID,
		CustomerName:  o.CustomerName,
		CustomerEmail: o.CustomerEmail,
		Items:         items,
		Status:        common.OrderStatus(o.Status),
		CreatedAt:     timestamppb.New(o.CreatedAt),
		UpdatedAt:     timestamppb.New(o.UpdatedAt),
	}
}
