package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/goclient-vanmango/apiclient"
	"github.com/goclient-vanmango/routes"
	"github.com/goclient-vanmango/service"
	"github.com/joho/godotenv"
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

func Handler(w http.ResponseWriter, r *http.Request) {

	apiclient := apiclient.NewAPIClient()
	apiService := service.NewAPIService(apiclient)
	apiHandler := routes.NewAPIHandler(apiService)

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve asset files
	fsAssets := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fsAssets))

	// Handle routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Accept-Content", "application/json")
		json.NewEncoder(w).Encode(struct {
			Code    int
			Message string
		}{
			Code:    http.StatusOK,
			Message: "Client server is functioning",
		})
		log.Println("Client server is functioning")
	})
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

	_ = godotenv.Load()

	port := os.Getenv("PORT")

	// Start server
	fmt.Println("Client server is listening on PORT ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

// func main() {

// 	apiclient := apiclient.NewAPIClient()
// 	apiService := service.NewAPIService(apiclient)
// 	apiHandler := routes.NewAPIHandler(apiService)

// 	// Serve static files
// 	fs := http.FileServer(http.Dir("./static"))
// 	http.Handle("/static/", http.StripPrefix("/static/", fs))

// 	// Serve asset files
// 	fsAssets := http.FileServer(http.Dir("./assets"))
// 	http.Handle("/assets/", http.StripPrefix("/assets/", fsAssets))

// 	// Handle routes
// 	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		w.Header().Add("Accept-Content", "application/json")
// 		json.NewEncoder(w).Encode(struct {
// 			Code    int
// 			Message string
// 		}{
// 			Code:    http.StatusOK,
// 			Message: "Client server is functioning",
// 		})
// 		log.Println("Client server is functioning")
// 	})
// 	http.HandleFunc("/", apiHandler.PageHomeHandler)
// 	http.HandleFunc("/sign-in", apiHandler.PageSignInHandler)
// 	http.HandleFunc("/logout", apiHandler.LogoutHandler)
// 	http.HandleFunc("/authorize", apiHandler.AuthorizeUserHandler)
// 	http.HandleFunc("/van-details", apiHandler.PageVanDetailsHandler)
// 	http.HandleFunc("/create-van", apiHandler.PageCreateVanHandler)
// 	http.HandleFunc("/create-new-van", apiHandler.CreateNewVanHandler)
// 	http.HandleFunc("/van-details/update-van", apiHandler.PageUpdateVanHandler)
// 	http.HandleFunc("/update-van/{vanID}", apiHandler.UpdateVanHandler)
// 	http.HandleFunc("/van-details/delete-van", apiHandler.DeleteVanHandler)

// 	_ = godotenv.Load()

// 	port := os.Getenv("PORT")

// 	// Start server
// 	fmt.Println("Client server is listening on PORT ", port)
// 	log.Fatal(http.ListenAndServe(":"+port, nil))
// }
