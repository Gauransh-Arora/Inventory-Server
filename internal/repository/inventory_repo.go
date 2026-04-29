package repository

import (
	"context"
	"server/internal/models"

	"github.com/jackc/pgx/v5"
)

type InventoryRepository struct{
	DB *pgx.Conn
}

func NewInventoryRepository(db *pgx.Conn) *InventoryRepository{
	return &InventoryRepository{DB: db}
}

func (r *InventoryRepository) CreateLog(ctx context.Context, log models.InventoryLog) error{
	query := `insert into inventory_logs (product_id,quantity) values($1,$2)`
	_,err := r.DB.Exec(ctx,query,log.ProductID,log.Quantity)
	return err
}

func (r *InventoryRepository) GetAllLogs(ctx context.Context) ([]models.InventoryLog, error){
	query := `select id,product_id,quantity,created_at from inventory_logs order by created_at desc`
	rows, err := r.DB.Query(ctx,query)
	if err != nil{
		return nil,err
	}
	defer rows.Close()

	var logs []models.InventoryLog
	for rows.Next(){
		var log models.InventoryLog
		err := rows.Scan(&log.ID, &log.ProductID, &log.Quantity, &log.CreatedAt)
		if err != nil{
			return nil,err
		}
		logs = append(logs,log)
	}
	return logs,nil
}