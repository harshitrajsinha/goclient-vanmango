// handler.go
package handler

import (
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"text/template"

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

func (h *apiHandler) PageHomeHandler(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		log.Println("PageHomeHandler::Incorrect error path")
		return
	}

	allVanData, err := h.apiService.GetAllVans(w, r)
	if err != nil {
		tmpl, _ := template.ParseFiles("templates/service-unavailable.html")
		_ = tmpl.Execute(w, nil)
		panic(err)
	}

	dataToRender := struct {
		VanDetails []struct {
			VanCode  string
			Name     string
			Category string
			Price    string
			ImageURL string
		}
	}{
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

	if r.URL.Path != "/van-details" {
		http.NotFound(w, r)
		return
	}

	vanID := r.URL.Query().Get("van")

	vanDetails, err := h.apiService.GetVanByID(w, r, vanID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		log.Println("Error getting van by id ", err)
		return
	}

	tmpl, err := template.ParseFiles("templates/van-details.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, vanDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (h *apiHandler) PageUpdateVanHandler(w http.ResponseWriter, r *http.Request) {
	// 	var newVanData *map[string]interface{}

	// 	var err error
	if r.URL.Path != "/van-details/update-van" {
		http.NotFound(w, r)
		return
	}

	vanID := r.URL.Query().Get("van")

	vanData, err := h.apiService.GetVanByID(w, r, vanID)

	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		log.Println("Error getting van by id ", err)
		return
	}

	tmpl, err := template.ParseFiles("templates/update-van.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, vanData)
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

	vanData, err := h.apiService.UpdateVan(w, r, vanID)
	if err != nil {
		log.Println("Error occured in updating van request ", err)
		http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect to home
		panic(err)
	}

	log.Println("Response from server for PUT van ", vanData)
	http.Redirect(w, r, "/", http.StatusFound)

}

func (h *apiHandler) DeleteVanHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("hereeeeeeeeee")
	// panic recovery
	if r := recover(); r != nil {
		log.Println("Error occured ", r)
		debug.Stack()
	}

	var err error
	vanID := r.URL.Query().Get("van")

	vanData, err := h.apiService.DeleteVan(w, r, vanID)
	if err != nil {
		log.Println("Error occured in deleting van request ", err)
		http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect to home
		panic(err)
	}

	log.Println("Response from server for DELETE van ", vanData)
	http.Redirect(w, r, "/", http.StatusFound)

}

func (h *apiHandler) PageSignInHandler(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	if r.URL.Path != "/sign-in" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		log.Println("PageSignInHandler::Incorrect error path")
		return
	}

	_ = godotenv.Load()
	authUser := os.Getenv("AUTH_USER")
	authPass := os.Getenv("AUTH_PASS")

	data := struct {
		User     string
		Password string
	}{
		User:     authUser,
		Password: authPass,
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

func (h *apiHandler) PageCreateVanHandler(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	if r.URL.Path != "/create-van" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		log.Println("PageCreateVanHandler::Incorrect error path")
		return
	}

	// Check a cookie for authorization
	_, err := r.Cookie("session_token")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			log.Println("No cookie found")
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther) // Redirect to sign-in page
			return
		default:
			log.Println("Error reading cookie: ", err)
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther) // Redirect to login page
			panic(err)
		}
	}

	tmpl, err := template.ParseFiles("templates/create-van.html")
	if err != nil {
		tmpl, _ := template.ParseFiles("templates/service-unavailable.html")
		_ = tmpl.Execute(w, nil)
		panic(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		tmpl, _ := template.ParseFiles("templates/service-unavailable.html")
		_ = tmpl.Execute(w, nil)
		panic(err)
	}
}
