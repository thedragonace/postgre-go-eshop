package migrations

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"goshop/dbs"
	"goshop/objects/category"
	"goshop/objects/order"
	"goshop/objects/orderLine"
	"goshop/objects/product"
)

func Migrate() {
	Product := product.Product{}
	Pategory := category.Category{}
	Order := order.Order{}
	OrderLine := orderLine.OrderLine{}

	dbs.Database.AutoMigrate(&Product, &Pategory, &Order, &OrderLine)
	dbs.Database.Model(&Product).AddForeignKey("categ_uuid", "categories(uuid)", "CASCADE", "CASCADE")
}
