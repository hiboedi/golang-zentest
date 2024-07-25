package controllers

import (
	"net/http"

	"zen-test/app/auth"
	"zen-test/app/helpers"
	"zen-test/app/middleware"
	"zen-test/app/web"
	"zen-test/app/web/models"
	"zen-test/app/web/services"
)

type UserControllerImpl struct {
	UserService services.UserService
}

type UserController interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
}

func NewUserController(userService services.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

// SignUp godoc
// @Summary Sign up a new user
// @Description Create a new user account
// @Tags User
// @Accept json
// @Produce json
// @Param user body models.UserCreate true "User Sign Up"
// @Success 200 {object} web.WebResponse{data=models.UserResponse}
// @Failure 400 {object} web.WebResponse
// @Router /users/signup [post]
func (c *UserControllerImpl) SignUp(w http.ResponseWriter, r *http.Request) {
	userSignUp := models.UserCreate{}
	helpers.ToRequestBody(r, &userSignUp)

	userResponse := c.UserService.Register(r.Context(), userSignUp)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   userResponse,
	}

	helpers.WriteResponseBody(w, webResponse)
}

// Login godoc
// @Summary Log in a user
// @Description Authenticate a user and set a session cookie
// @Tags User
// @Accept json
// @Produce json
// @Param user body models.UserLogin true "User Login"
// @Success 200 {object} web.WebResponse{data=models.UserResponse}
// @Failure 401 {object} web.WebResponse
// @Router /users/login [post]
func (c *UserControllerImpl) Login(w http.ResponseWriter, r *http.Request) {
	userLogin := models.UserLogin{}
	helpers.ToRequestBody(r, &userLogin)

	userResponse, loggedIn := c.UserService.Login(r.Context(), userLogin)
	if loggedIn {
		helpers.SetCookie(w, r, helpers.RefreshToken, userResponse.RefreshToken)
		webResponse := web.WebResponse{
			Code:   http.StatusOK,
			Status: "Ok",
			Data:   userResponse,
		}
		helpers.WriteResponseBody(w, webResponse)
	} else {
		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "Your enail or password is not valid",
		}
		helpers.WriteResponseBody(w, webResponse)
	}
}

// Update godoc
// @Summary Update user for the user
// @Description Update user for the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Param user body models.UserUpdate true "User update"
// @Param userId path string true "User ID"
// @Success 200 {object} web.WebResponse{data=models.UserResponse}
// @Failure 401 {object} web.WebResponse
// @Router /users/{userId} [put]
// @Security BearerAuth
func (c *UserControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	userUpdateRequest := models.UserUpdate{}
	helpers.ToRequestBody(r, &userUpdateRequest)

	userId := middleware.GetUserID(r)
	user := c.UserService.Update(r.Context(), userUpdateRequest, userId)
	response := web.WebResponse{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   user,
	}

	helpers.WriteResponseBody(w, response)
}

// Logout godoc
// @Summary Logout for the user
// @Description Logout for the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} web.WebResponse
// @Failure 401 {object} web.WebResponse
// @Router /users/logout [post]
// @Security BearerAuth
func (c *UserControllerImpl) Logout(w http.ResponseWriter, r *http.Request) {
	helpers.DeleteCookieHandler(w, r, "refresh_token")
	response := web.WebResponse{
		Code:   http.StatusOK,
		Status: "Logout success",
	}
	helpers.WriteResponseBody(w, response)
}

// Refresh Token godoc
// @Summary Refresh Token for the user
// @Description Refresh Token for the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} web.WebResponse
// @Failure 401 {object} web.WebResponse
// @Router /users/refresh-token [post]
// @Security BearerAuth
func (c *UserControllerImpl) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := helpers.GetCookie(w, r, helpers.RefreshToken)
	if err != nil {
		http.Error(w, "Refresh token not found", http.StatusUnauthorized)
		return
	}

	claims, err := auth.VerifyToken(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	id := claims["id"].(string)
	newAccessToken, err := auth.CreateToken(id)
	if err != nil {
		http.Error(w, "Failed to create new access token", http.StatusInternalServerError)
		return
	}

	response := web.WebResponse{
		Code:   http.StatusOK,
		Status: "Token refreshed",
		Data: map[string]string{
			"access_token": newAccessToken,
		},
	}
	helpers.WriteResponseBody(w, response)
}
