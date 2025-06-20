package model

type GoodsInfo struct {
	ID         int64   `json:"id" gorm:"primary_key"`
	Is_deleted bool    `json:"is_deleted" gorm:"default:false"`
	GoodsName  string  `json:"goodsName"`
	Price      float64 `json:"price"`
	Stock      int64   `json:"stock"`
	ImagePath  string  `json:"image_path"`
}
