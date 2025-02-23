package dto

// PaginationResponse adalah struktur umum untuk respon dengan pagination
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Pagination struct {
		Page       int `json:"page"`
		PageSize   int `json:"pageSize"`
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
	} `json:"pagination"`
}

// Paginasi adalah fungsi untuk membungkus data hasil query dan pagination
func PaginationHelper(data interface{}, totalRecords, page, size int) PaginationResponse {
	totalPages := (totalRecords + size - 1) / size

	response := PaginationResponse{
		Data: data,
		Pagination: struct {
			Page       int `json:"page"`
			PageSize   int `json:"pageSize"`
			Total      int `json:"total"`
			TotalPages int `json:"total_pages"`
		}{
			Page:       page,
			PageSize:   size,
			Total:      totalRecords,
			TotalPages: totalPages,
		},
	}

	return response
}
