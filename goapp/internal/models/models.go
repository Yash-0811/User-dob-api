package models

// CreateUserRequest is the payload for POST /users.
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	DOB  string `json:"dob"  validate:"required,datetime=2006-01-02"`
}

// UpdateUserRequest is the payload for PUT /users/:id.
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	DOB  string `json:"dob"  validate:"required,datetime=2006-01-02"`
}

// UserResponse is returned for create / update operations (no age).
type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
}

// UserDetailResponse is returned when fetching a single user or listing users
// (includes dynamically-calculated age).
type UserDetailResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  int    `json:"age"`
}

// PaginatedUsersResponse wraps a paginated list of users.
type PaginatedUsersResponse struct {
	Data       []UserDetailResponse `json:"data"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// ErrorResponse is used for all error payloads.
type ErrorResponse struct {
	Message string `json:"message"`
}
