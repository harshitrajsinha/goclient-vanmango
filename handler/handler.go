// handler.go
package handler

import (
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"text/template"
	"time"

	"github.com/goclient-vanmango/service"
	"github.com/joho/godotenv"
)

type apiHandler struct {
	apiService *service.Service
}

func NewAPIHandler(service *service.Service) *apiHandler {
	return &apiHandler{
		apiService: service,
	}
}

// Read a specific cookie by name

func isUserAuthorized(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {

		return "", err
	} else {
		return cookie.Value, nil
	}
}

func (h *apiHandler) PageSignInHandler(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	var UserLoggedIn bool

	if r.URL.Path != "/sign-in" {
		tmpl, _ := template.ParseFiles("templates/page-not-found.html")
		_ = tmpl.Execute(w, nil)
		log.Println("PageHomeHandler::Incorrect error path ", r.URL.Path)
		return
	}

	// Read a specific cookie by name
	cookieValue, _ := isUserAuthorized(r)
	if cookieValue != "" {
		UserLoggedIn = true
	} else {
		UserLoggedIn = false
	}

	_ = godotenv.Load()
	authUser := os.Getenv("AUTH_USER")
	authPass := os.Getenv("AUTH_PASS")

	data := struct {
		User     string
		Password string
		LoggedIn bool
	}{
		User:     authUser,
		Password: authPass,
		LoggedIn: UserLoggedIn,
	}

	tmpl, err := template.ParseFiles("templates/sign-in.html")
	if err != nil {
		tmpl, _ := template.ParseFiles("templates/service-unavailable.html")
		_ = tmpl.Execute(w, nil)
		panic(err)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		tmpl, _ := template.ParseFiles("templates/service-unavailable.html")
		_ = tmpl.Execute(w, nil)
		panic(err)
	}
}

func (h *apiHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	if r.URL.Path != "/logout" {
		tmpl, _ := template.ParseFiles("templates/page-not-found.html")
		_ = tmpl.Execute(w, nil)
		log.Println("PageHomeHandler::Incorrect error path ", r.URL.Path)
		return
	}

	// check if user is logged in
	cookieVal, err := isUserAuthorized(r)

	if cookieVal != "" || err == nil {
		// Logout by deleting the cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			Domain:   "",
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   true,
		})
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *apiHandler) AuthorizeUserHandler(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	isSet, err := h.apiService.AuthorizeUser(w, r)

	if err != nil {
		panic(err)
	}

	if isSet == 1 {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
	}

}

func (h *apiHandler) PageHomeHandler(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	var UserLoggedIn bool

	if r.URL.Path != "/" && r.URL.Path != "/favicon.ico" {
		tmpl, _ := template.ParseFiles("templates/page-not-found.html")
		_ = tmpl.Execute(w, nil)
		log.Println("PageHomeHandler::Incorrect error path ", r.URL.Path)
		return
	}

	// Read a specific cookie by name
	cookieValue, _ := isUserAuthorized(r)
	if cookieValue != "" {
		UserLoggedIn = true
	} else {
		UserLoggedIn = false
	}

	allVanData, err := h.apiService.GetAllVans(w, r)
	if err != nil {
		tmpl, _ := template.ParseFiles("templates/service-unavailable.html")
		_ = tmpl.Execute(w, nil)
		panic(err)
	}

	dataToRender := struct {
		LoggedIn   bool
		VanDetails []struct {
			VanCode  string
			Name     string
			Category string
			Price    string
			ImageURL string
		}
	}{
		LoggedIn:   UserLoggedIn,
		VanDetails: *allVanData,
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		tmpl, _ := template.ParseFiles("templates/service-unavailable.html")
		_ = tmpl.Execute(w, nil)
		panic(err)
	}
	err = tmpl.Execute(w, dataToRender)
	if err != nil {
		tmpl, _ := template.ParseFiles("templates/service-unavailable.html")
		_ = tmpl.Execute(w, nil)
		panic(err)
	}
}

func (h *apiHandler) PageVanDetailsHandler(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	var UserLoggedIn bool

	if r.URL.Path != "/van-details" {
		tmpl, _ := template.ParseFiles("templates/page-not-found.html")
		_ = tmpl.Execute(w, nil)
		log.Println("PageHomeHandler::Incorrect error path ", r.URL.Path)
		return
	}

	// Read a specific cookie by name
	cookieValue, _ := isUserAuthorized(r)
	if cookieValue != "" {
		UserLoggedIn = true
	} else {
		UserLoggedIn = false
	}

	vanID := r.URL.Query().Get("van")

	vanDetails, err := h.apiService.GetVanByID(w, r, vanID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		log.Println("Error getting van by id ", err)
		panic(err)
	}

	dataToRender := struct {
		LoggedIn   bool
		VanDetails service.VanDetails
	}{
		LoggedIn:   UserLoggedIn,
		VanDetails: *vanDetails,
	}

	tmpl, err := template.ParseFiles("templates/van-details.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	err = tmpl.Execute(w, dataToRender)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

}

func (h *apiHandler) PageUpdateVanHandler(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	var UserLoggedIn bool

	// 	var err error
	if r.URL.Path != "/van-details/update-van" {
		tmpl, _ := template.ParseFiles("templates/page-not-found.html")
		_ = tmpl.Execute(w, nil)
		log.Println("PageHomeHandler::Incorrect error path ", r.URL.Path)
		return
	}

	vanID := r.URL.Query().Get("van")

	// Read a specific cookie by name
	cookieValue, err := isUserAuthorized(r)
	if err != nil {
		log.Println("No cookie found ", err)
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther) // Redirect to sign-in page
		return
	}
	if cookieValue != "" {
		UserLoggedIn = true
	} else {
		UserLoggedIn = false
	}

	vanData, err := h.apiService.GetVanByID(w, r, vanID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		panic(err)
	}

	dataToRender := struct {
		LoggedIn bool
		VanData  service.VanDetails
	}{
		LoggedIn: UserLoggedIn,
		VanData:  *vanData,
	}

	tmpl, err := template.ParseFiles("templates/update-van.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, dataToRender)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (h *apiHandler) UpdateVanHandler(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	if r := recover(); r != nil {
		log.Println("Error occured ", r)
		debug.Stack()
	}

	var err error
	vanID := strings.TrimPrefix(r.URL.Path, "/update-van/")

	// Read a specific cookie by name
	cookieValue, err := isUserAuthorized(r)
	if err != nil {
		log.Println("No cookie found ", err)
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther) // Redirect to sign-in page
		return
	}

	vanData, err := h.apiService.UpdateVan(w, r, vanID, cookieValue)
	if err != nil {
		log.Println("Error occured in updating van request ", err)
		http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect to home
		panic(err)
	}

	log.Println("Response from server for PUT van ", vanData)
	http.Redirect(w, r, "/", http.StatusFound)

}

func (h *apiHandler) DeleteVanHandler(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	if r := recover(); r != nil {
		log.Println("Error occured ", r)
		debug.Stack()
	}

	var err error
	vanID := r.URL.Query().Get("van")

	// Read a specific cookie by name
	cookieValue, err := isUserAuthorized(r)
	if err != nil {
		log.Println("No cookie found ", err)
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther) // Redirect to sign-in page
		return
	}

	vanData, err := h.apiService.DeleteVan(w, r, vanID, cookieValue)
	if err != nil || vanData != 204 {
		log.Println("Error occured in deleting van request ", err)
		http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect to home
		panic(err)
	}

	log.Println("Response code from server for DELETE van ", vanData)
	http.Redirect(w, r, "/", http.StatusFound)

}

func (h *apiHandler) PageCreateVanHandler(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	var UserLoggedIn bool

	if r.URL.Path != "/create-van" {
		tmpl, _ := template.ParseFiles("templates/page-not-found.html")
		_ = tmpl.Execute(w, nil)
		log.Println("PageHomeHandler::Incorrect error path ", r.URL.Path)
		return
	}

	// Read a specific cookie by name
	cookieValue, err := isUserAuthorized(r)
	if err != nil {
		log.Println("No cookie found ", err)
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther) // Redirect to sign-in page
		return
	}
	if cookieValue != "" {
		UserLoggedIn = true
	} else {
		UserLoggedIn = false
	}

	dataToRender := struct {
		LoggedIn bool
	}{
		LoggedIn: UserLoggedIn,
	}

	tmpl, err := template.ParseFiles("templates/create-van.html")
	if err != nil {
		tmpl, _ := template.ParseFiles("templates/service-unavailable.html")
		_ = tmpl.Execute(w, nil)
		panic(err)
	}
	err = tmpl.Execute(w, dataToRender)
	if err != nil {
		tmpl, _ := template.ParseFiles("templates/service-unavailable.html")
		_ = tmpl.Execute(w, nil)
		panic(err)
	}
}

func (h *apiHandler) CreateNewVanHandler(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	var err error

	// Read a specific cookie by name
	cookieValue, err := isUserAuthorized(r)
	if err != nil {
		log.Println("No cookie found ", err)
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther) // Redirect to sign-in page
		return
	}
	isCreated, err := h.apiService.CreateVan(w, r, cookieValue)

	if err != nil {
		log.Println("Error occured in creating van request ", err)
		http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect to home
		panic(err)
	}

	if isCreated == 1 {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Redirect(w, r, "/create-van", http.StatusSeeOther)
	}
}
