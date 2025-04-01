package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/goclient-vanmango/apiclient"
	"github.com/goclient-vanmango/handler"
	"github.com/goclient-vanmango/service"
)

type FormData struct {
	Username string
	Password string
}

type VanCreationFormData struct {
	Name        string
	Brand       string
	Description string
	Category    string
	FuelType    string
	Price       string
	Image       string
	VanCode     string
}

// var pageRequested string
// var vanDetailsData *map[string]interface{}

func main() {

	apiclient := apiclient.NewAPIClient()
	apiService := service.NewAPIService(apiclient)
	apiHandler := handler.NewAPIHandler(apiService)

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle routes
	http.HandleFunc("/", apiHandler.PageHomeHandler)
	http.HandleFunc("/van-details", apiHandler.PageVanDetailsHandler)
	http.HandleFunc("/van-details/update-van", apiHandler.PageUpdateVanHandler)
	http.HandleFunc("/update-van/{vanID}", apiHandler.UpdateVanHandler)
	http.HandleFunc("/van-details/delete-van", apiHandler.DeleteVanHandler)
	http.HandleFunc("/create-van", apiHandler.PageCreateVanHandler)
	http.HandleFunc("/sign-in", apiHandler.PageSignInHandler)
	http.HandleFunc("/authorize", apiHandler.AuthorizeUserHandler)
	// http.HandleFunc("/create-new-van", createNewVanHandler)

	// Start server
	fmt.Println("Client server is listening on PORT 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// func createNewVanHandler(w http.ResponseWriter, r *http.Request) {

// 	var err error

// 	if r.Method != http.MethodPost {
// 		http.Redirect(w, r, "/create-van", http.StatusSeeOther)
// 		return
// 	}

// 	// if cookie is not present
// 	// Read a specific cookie by name
// 	cookie, err := r.Cookie("session_token")
// 	if err != nil {
// 		switch err {
// 		case http.ErrNoCookie:
// 			fmt.Println("No cookie named 'session_token' found")
// 			pageRequested = "create-van"
// 			http.Redirect(w, r, "/sign-in", http.StatusSeeOther) // Redirect to login page
// 			return
// 		default:
// 			fmt.Println("Error reading cookie: ", err)
// 			pageRequested = "create-van"
// 			http.Redirect(w, r, "/sign-in", http.StatusSeeOther) // Redirect to login page
// 			return
// 		}
// 	}

// 	fmt.Println("Cookie value: ", cookie.Value)

// 	// Parse form data
// 	err = r.ParseForm()
// 	if err != nil {
// 		http.Error(w, "Error parsing form", http.StatusBadRequest)
// 		return
// 	}

// 	// Get form values
// 	formData := VanCreationFormData{
// 		Name:        r.FormValue("name"),
// 		Brand:       r.FormValue("brand"),
// 		Description: r.FormValue("description"),
// 		Category:    r.FormValue("category"),
// 		FuelType:    r.FormValue("fuel-type"),
// 		Price:       r.FormValue("price"),
// 		Image:       r.FormValue("image-url"),
// 	}

// 	// form validation

// 	// convert price -> string to int
// 	priceStrToInt, _ := strconv.Atoi(formData.Price)

// 	var formInput map[string]interface{} = map[string]interface{}{"name": formData.Name, "brand": formData.Brand, "description": formData.Description, "category": formData.Category, "fuel-type": formData.FuelType, "engine-id": "e1f86b1a-0873-4c19-bae2-fc60329d0140", "price": priceStrToInt, "image-url": formData.Image}

// 	requestToSend, _ := json.Marshal(formInput)

// 	response, err := sendRequest("http://localhost:8080/van", &requestToSend, "POST")

// 	if err != nil {
// 		log.Println("Error occured in inserting van request ", err)
// 		http.Redirect(w, r, "/create-van", http.StatusSeeOther) // Redirect to home
// 		return
// 	}

// 	log.Println("Response from server for POST van ", response)
// 	http.Redirect(w, r, "/", http.StatusFound)
// }
