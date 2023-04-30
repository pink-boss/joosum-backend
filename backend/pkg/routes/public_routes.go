package routes

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"github.com/create-go-app/net_http-go-template/app/controllers"
	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() (*Handler, error) {
	return new(Handler), nil
}

// PublicRoutes func for describe group of public routes.
func PublicRoutes(router *mux.Router) {
	// Routes for GET method:
	router.HandleFunc("/api/v1/user/{id}", controllers.GetUser).Methods(http.MethodGet) // get one user by ID
	router.HandleFunc("/api/v1/users", controllers.GetUsers).Methods(http.MethodGet)    // Get list of all users
}

func (h *Handler) GetMainPage(c *gin.Context) {
	log.Println("Main page....")
	c.String(http.StatusOK, "Main page for secure API!!")
}

func (h *Handler) GetUsers(c *gin.Context) {
	c.String(http.StatusOK, "GetUsers API!!")
}

func (h *Handler) GetUser(c *gin.Context) {
	name := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"id":    name,
		"name":  "홍길동",
		"email": "gildong@gmail.com",
	})
}
