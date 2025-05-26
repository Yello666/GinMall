package Goods

type Goodsinfo struct {
	ID         int64   `json:"id"`
	Is_deleted bool    `json:"is_deleted" gorm:"default:false"`
	GoodsName  string  `json:"goodsName"`
	Price      float64 `json:"price"`
	//Star      bool    `json:"star" gorm:"default:false"` //是否收藏进了购物车
}
type GoodsResponse struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	ID    int64   `json:"id"`
}
