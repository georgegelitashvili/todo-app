package controllers

import (
	"encoding/json"
	"net/http"
	"todo-app/models"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

type CategoryController struct {
	session *gocql.Session
}

func NewCategoryController(session *gocql.Session) *CategoryController {
	return &CategoryController{session: session}
}

func (c *CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := category.Create(c.session); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, Response{
		Status: "success",
		Data:   category,
	})
}

func (c *CategoryController) GetCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := gocql.ParseUUID(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := models.GetCategoryByID(c.session, id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Category not found")
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status: "success",
		Data:   category,
	})
}

func (c *CategoryController) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := models.GetAllCategories(c.session)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status: "success",
		Data:   categories,
	})
}

func (c *CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := gocql.ParseUUID(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	category.CategoryID = id
	if err := category.Update(c.session); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update category")
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Category updated successfully",
		Data:    category,
	})
}

func (c *CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := gocql.ParseUUID(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := models.DeleteCategoryByID(c.session, id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete category")
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Category deleted successfully",
	})
}
