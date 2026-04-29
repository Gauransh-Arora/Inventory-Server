package repository

import (
	"context"
	"server/internal/models"

	"github.com/jackc/pgx/v5"
)

type ProductRepository struct{
	DB *pgx.Conn
}

func NewProductRepository(db *pgx.Conn) *ProductRepository{
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, p models.Product) error{
	query:=`insert into products (barcode, item_name, mrp) values($1,$2,$3)`
	_,err := r.DB.Exec(ctx,query,p.Barcode, p.ItemName, p.MRP)
	return err
}

func (r *ProductRepository) GetAllProducts(ctx context.Context)([]models.Product,error){
	query:=`select id,barcode,item_name,mrp,created_at from products order by created_at desc`
	rows,err := r.DB.Query(ctx,query)
	if err != nil{
		return nil,err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next(){
		var product models.Product
		err := rows.Scan(&product.ID, &product.Barcode, &product.ItemName, &product.MRP, &product.CreatedAt)
		if err != nil{
			return nil,err
		}
		products = append(products,product)
	}
	return products,nil
}

func(r *ProductRepository) GetProductByBarcode(ctx context.Context, barcode string)([]models.Product, error){
	query:=`select id,barcode,item_name,mrp,created_at from products where barcode = $1`
	rows,err := r.DB.Query(ctx,query,barcode)
	if err != nil{
		return nil,err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next(){
		var product models.Product
		err := rows.Scan(&product.ID, &product.Barcode, &product.ItemName, &product.MRP, &product.CreatedAt)
		if err != nil{
			return nil,err
		}
		products = append(products,product)
	}
	return products,nil
}