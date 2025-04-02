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
	http.HandleFunc("/sign-in", apiHandler.PageSignInHandler)
	http.HandleFunc("/logout", apiHandler.LogoutHandler)
	http.HandleFunc("/authorize", apiHandler.AuthorizeUserHandler)
	http.HandleFunc("/van-details", apiHandler.PageVanDetailsHandler)
	http.HandleFunc("/create-van", apiHandler.PageCreateVanHandler)
	http.HandleFunc("/create-new-van", apiHandler.CreateNewVanHandler)
	http.HandleFunc("/van-details/update-van", apiHandler.PageUpdateVanHandler)
	http.HandleFunc("/update-van/{vanID}", apiHandler.UpdateVanHandler)
	http.HandleFunc("/van-details/delete-van", apiHandler.DeleteVanHandler)

	// Start server
	fmt.Println("Client server is listening on PORT 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
