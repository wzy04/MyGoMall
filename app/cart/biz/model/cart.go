package model

import (
	"cart/biz/dal/mysql"
	"github.com/BeroKiTeer/MyGoMall/common/kitex_gen/cart"
	"gorm.io/gorm"
)

// QueryItemsByUser 查询 这个 user 的所有 商品
func QueryItemsByUser(userCart *cart.Cart) error {
	err := mysql.DB.Table("cart").
		Select("product_id", "quantity").
		Where("user_id = ?", userCart.UserId).
		Find(&userCart.Items).Error
	return err
}

// CheckItemsByUserAndProduct 检查商品是否已存在在购物车
func CheckItemsByUserAndProduct(UserID uint32, ProductID int64, targetItemQuantity *int32) error {
	err := mysql.DB.Table("cart").
		Select("quantity").
		Where("product_id = ?", ProductID).
		Where("user_id = ?", UserID).Scan(targetItemQuantity).Error
	return err
}

// CheckItemsByUser 查询该用户的所有商品
func CheckItemsByUser(UserID uint32, targetItemQuantity *int32) error {
	err := mysql.DB.Table("cart").
		Select("SUM(quantity)").
		Where("user_id = ?", UserID).
		Group("user_id").Scan(targetItemQuantity).Error
	return err
}

// AddItem 向用户购物车中添加商品
func AddItem(UserID uint32, ProductID int64, Quantity int32) error {
	err := mysql.DB.Table("cart").Create(&Cart{
		UserID:    UserID,
		ProductID: ProductID,
		Quantity:  Quantity,
	}).Error
	return err
}

// UpdateItem 在已有的商品上面更新数量
func UpdateItem(UserID uint32, ProductID int64, Quantity int32) error {
	err := mysql.DB.Table("cart").
		Select("quantity").
		Where("product_id = ?", ProductID).
		Where("user_id = ?", UserID).
		Update("quantity", gorm.Expr("quantity + ?", Quantity)).Error
	return err
}

// EmptyCart 清空该用户的购物车
func EmptyCart(UserID uint32) error {
	err := mysql.DB.Table("cart").
		Where("user_id = ?", UserID).
		Delete(&UserID).Error
	return err
}
