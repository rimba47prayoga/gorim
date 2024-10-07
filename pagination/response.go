package pagination

type PaginatedResponse struct {
	Status			bool				`json:"status"`
	Message			string				`json:"message"`
	Pagination		Pagination			`json:"pagination"`
	Data         	interface{} 		`json:"data"`
}
