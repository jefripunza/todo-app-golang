package dto

import "todolist/model"

type TodoCreateBody struct {
	Title       string `json:"title" validate:"required,min=5,max=30" example:"account title"`
	Description string `json:"description" validate:"required,min=5,max=100" example:"account description"`
	Status      string `json:"status" validate:"required,oneof=pending completed" example:"pending"`
	DueDate     string `json:"due_date" validate:"required,datetime=2006-01-02 15:04:05" example:"2025-07-02 15:04:05"`
}
type TodoCreateResponse struct {
	Message string     `json:"message"`
	Task    model.Todo `json:"task"`
}

type TodoGetAllResponse struct {
	Tasks      []model.Todo         `json:"tasks"`
	Pagination TodoGetAllPagination `json:"pagination"`
}
type TodoGetAllPagination struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalTasks  int `json:"total_tasks"`
}

type TodoUpdateBody struct {
	Title       string `json:"title" validate:"min=5,max=30" example:"account title updated"`
	Description string `json:"description" validate:"min=5,max=100" example:"account description updated"`
	Status      string `json:"status" validate:"oneof=pending completed" example:"completed"`
	DueDate     string `json:"due_date" validate:"datetime=2006-01-02 15:04:05" example:"2025-07-05 22:19:00"`
}
type TodoUpdateResponse struct {
	Message string     `json:"message"`
	Task    model.Todo `json:"task"`
}

type TodoDeleteResponse struct {
	Message string `json:"message"`
}
