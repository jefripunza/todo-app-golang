package module

import (
	"fmt"
	"time"
	"todolist/database"
	"todolist/dto"
	"todolist/model"
	"todolist/util"

	"github.com/gofiber/fiber/v2"
)

type Todo struct{}

func (ref Todo) Route(api fiber.Router) {
	handler := TodoHandler{}
	route := api.Group("/tasks")

	route.Post("/", handler.Create)
	route.Get("/", handler.GetAll)
	route.Get("/:id", handler.GetByID)
	route.Put("/:id", handler.Update)
	route.Delete("/:id", handler.Delete)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type TodoHandler struct{}

func (handler TodoHandler) HelloWorld(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Hello World",
	})
}

// CreateTodo creates a new todo
//
//	@Summary		Create Todo
//	@Description	create todo
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.TodoCreateBody	true	"Todo Create Body"
//	@Success		200		{object}	dto.TodoCreateResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		404		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/tasks [post]
func (handler TodoHandler) Create(c *fiber.Ctx) error {
	var body dto.TodoCreateBody
	if is_error, code, message := util.BodyValidator(c, &body); is_error {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    code,
			"message": message,
		})
	}

	// Parse time directly in Asia/Jakarta timezone
	loc, _ := time.LoadLocation("Asia/Jakarta")
	due_date, _ := time.ParseInLocation("2006-01-02 15:04:05", body.DueDate, loc)
	todo := model.Todo{
		Title:       body.Title,
		Description: body.Description,
		Status:      body.Status,
		DueDate:     due_date,
	}

	title_exist := database.Postgres.Where("title = ?", body.Title).First(&model.Todo{})
	if title_exist.Error == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "BAD_REQUEST",
			"message": "title already exists",
		})
	}

	create := database.Postgres.Create(&todo)
	if create.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "failed to create todo",
		})
	}

	last_id := database.Postgres.Last(&model.Todo{})
	var exist model.Todo
	last_id.Scan(&exist)
	todo.ID = exist.ID

	response := dto.TodoCreateResponse{
		Message: "Todo created successfully",
		Task:    todo,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// GetAllTodos lists all existing todos
//
//	@Summary		List Todos with pagination
//	@Description	get todos with pagination
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//	@Param			status	query		string	false	"(optional): Filter tasks by status (pending/completed)."
//	@Param			page	query		string	false	"(optional): Page number for pagination."
//	@Param			limit	query		string	false	"(optional): Number of tasks per page."
//	@Param			search	query		string	false	"(optional): Search term to filter tasks by title or description."
//	@Success		200	{object}	dto.TodoGetAllResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/tasks [get]
func (handler TodoHandler) GetAll(c *fiber.Ctx) error {
	status := c.Query("status", "")
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")
	search := c.Query("search", "")

	if status != "" {
		if status != "pending" && status != "completed" {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Code:    "BAD_REQUEST",
				Message: "Invalid status",
			})
		}
	}

	pageNum, err := util.StringToInt(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	limitNum, err := util.StringToInt(limit)
	if err != nil || limitNum < 1 || limitNum > 100 {
		limitNum = 10
	}

	offset := (pageNum - 1) * limitNum

	query := database.Postgres.Model(&model.Todo{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", searchTerm, searchTerm)
	}

	// order by id
	query = query.Order("id ASC")

	var totalTasks int64
	query.Count(&totalTasks)

	var tasks []model.Todo
	result := query.Offset(offset).Limit(limitNum).Find(&tasks)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to fetch tasks",
		})
	}

	totalPages := int(totalTasks) / limitNum
	if int(totalTasks)%limitNum != 0 {
		totalPages++
	}

	response := dto.TodoGetAllResponse{
		Tasks: tasks,
		Pagination: dto.TodoGetAllPagination{
			CurrentPage: pageNum,
			TotalPages:  totalPages,
			TotalTasks:  int(totalTasks),
		},
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// GetTodoByID lists a specific todo
//
//	@Summary		Get one Todo
//	@Description	get one todo
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Task ID"
//	@Success		200	{object}	model.Todo
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/tasks/{id} [get]
func (handler TodoHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	fmt.Printf("id: %v\n", id)
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid ID",
		})
	}
	idNum, err := util.StringToInt(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid ID",
		})
	}

	var todo model.Todo
	result := database.Postgres.First(&todo, idNum)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "Todo not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(todo)
}

// UpdateTodoByID updates a specific todo
//
//	@Summary		Update Todo
//	@Description	update todo
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Task ID"
//	@Param			body	body		dto.TodoUpdateBody	true	"Todo Update Body"
//	@Success		200	{object}	dto.TodoUpdateResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/tasks/{id} [put]
func (handler TodoHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid ID",
		})
	}
	idNum, err := util.StringToInt(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid ID",
		})
	}

	var todo model.Todo
	result := database.Postgres.First(&todo, idNum)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "Todo not found",
		})
	}

	var body dto.TodoUpdateBody
	if is_error, code, message := util.BodyValidator(c, &body); is_error {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    code,
			"message": message,
		})
	}

	todo.Title = body.Title
	todo.Description = body.Description
	todo.Status = body.Status
	todo.DueDate, _ = time.Parse("2006-01-02 15:04:05", body.DueDate)

	update := database.Postgres.Save(&todo)
	if update.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to update todo",
		})
	}

	response := dto.TodoUpdateResponse{
		Message: "Todo updated successfully",
		Task:    todo,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteTodoByID deletes a specific todo
//
//	@Summary		Delete Todo
//	@Description	delete todo
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Task ID"
//	@Success		200	{object}	dto.TodoDeleteResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/tasks/{id} [delete]
func (handler TodoHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	delete := database.Postgres.Delete(&model.Todo{}, id)
	if delete.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to delete todo",
		})
	}

	response := dto.TodoDeleteResponse{
		Message: "Todo deleted successfully",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
