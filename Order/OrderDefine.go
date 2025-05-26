package Order

type Order struct {
	ID         int64 `json:"id"`
	UserID     uint  `json:"userID"`
	GoodsID    int64 `json:"goodsID"`
	Is_deleted bool  `json:"is_deleted"`
}
