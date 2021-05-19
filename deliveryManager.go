package delivery_manager

type DMProduct struct {
	Id          int    `json:"id" binding:"numeric" db:"id"`
	Title       string `json:"title" binding:"required" db:"title"`
	Description string `json:"description" db:"description"`
	Price       int    `json:"price" binding:"numeric,min=0" db:"price"`
}

type DMList struct {
	Id          int    `json:"id" binding:"numeric" db:"id"`
	Title       string `json:"title" binding:"required" db:"title"`
	Description string `json:"description" db:"description"`
}

type DMProductSearchQuery struct {
	MatchAllFields     bool   `json:"strict"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	TitleOrDescription string `json:"any"`
}

type DMProductIndex struct {
	Id    int `json:"id"`
	Count int `json:"count"`
}
