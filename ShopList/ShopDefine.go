package ShopList

//type Shoplist struct {
//	gorm.Model
//	UserID uint    `json:"userID"`
//	List   []int64 `json:"list"` //只储存商品id
//}

type Relation struct {
	ID         uint  `json:"ID"`
	UserID     uint  `json:"userID"`
	GoodsID    int64 `json:"goodsID"`
	Is_deleted bool  `json:"is_deleted"`
}
