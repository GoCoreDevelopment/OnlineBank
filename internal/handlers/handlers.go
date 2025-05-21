package handlers

import (
	"api/internal/middlewarecustom"
	"api/internal/models/transaction"
	"api/internal/models/user"
	"api/internal/servises/jwtservice"
	"api/internal/servises/userservice"
	"api/internal/utils/bind"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService userservice.UserService
	jwtService  jwtservice.JWTService
}

func NewUserHandlers(userServices userservice.UserService, jwtService jwtservice.JWTService) UserHandler {
	return UserHandler{
		userService: userServices,
		jwtService:  jwtService,
	}
}

func (h *UserHandler) registredUserHandler(c echo.Context) error {
	var user user.UserRegistred
	if err := bind.Bind(c, &user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
	}

	id, err := h.userService.Registred(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]int{"id": id})
}

func (h *UserHandler) loginUserHandler(c echo.Context) error {
	var user user.UserLogin
	if err := bind.Bind(c, &user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}

	id, err := h.userService.Login(user)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	token, err := h.jwtService.CreateJWT(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Please retry"})
	}

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Expires:  time.Now().Add(1 * time.Hour),
	})

	return c.JSON(http.StatusOK, map[string]string{"Authorized": "Succsesfull"})
}

func (h *UserHandler) checkBalanceHandler(c echo.Context) error {
	id := c.Get("id").(int)
	balance, err := h.userService.CheckBalance(id)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]int{
		"balance": balance,
	})
}

func (h *UserHandler) transactionHandler(c echo.Context) error {
	senderID := c.Get("id").(int)

	var transactionRequest transaction.TransactionRequest
	if err := bind.Bind(c, &transactionRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	log.Printf("id_sender: %d, email_receiver: %s, amount_transaction: %d", senderID, transactionRequest.To, transactionRequest.Amount)

	receiverID, err := h.userService.GetUserIDByEmail(transactionRequest.To)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err = h.userService.Transfer(senderID, receiverID, transactionRequest.Amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"transaction": "Succsesfull"})
}

func (h *UserHandler) InitHandlers(e *echo.Echo) {
	e.POST("/api/register", h.registredUserHandler)
	e.POST("/api/login", h.loginUserHandler)
	e.GET("/api/balance", h.checkBalanceHandler, middlewarecustom.MiddlewareBalance(h.jwtService))
	e.POST("api/transfer", h.transactionHandler, middlewarecustom.MiddlewareBalance(h.jwtService))
}
