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
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()
	// Serve static files under /static
	router.Handle("/static/{file:.*}", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Serve assets files under /assets
	router.Handle("/assets/{file:.*}", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	// Serve static files
	// fs := http.FileServer(http.Dir("./static"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))

	// // Serve asset files
	// fsAssets := http.FileServer(http.Dir("./assets"))
	// http.Handle("/assets/", http.StripPrefix("/assets/", fsAssets))

	// Handle routes

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
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
	router.HandleFunc("/", apiHandler.PageHomeHandler).Methods(http.MethodGet)
	router.HandleFunc("/sign-in", apiHandler.PageSignInHandler).Methods(http.MethodGet)
	router.HandleFunc("/logout", apiHandler.LogoutHandler).Methods(http.MethodGet)
	router.HandleFunc("/authorize", apiHandler.AuthorizeUserHandler).Methods(http.MethodGet)
	router.HandleFunc("/van-details", apiHandler.PageVanDetailsHandler).Methods(http.MethodGet)
	router.HandleFunc("/create-van", apiHandler.PageCreateVanHandler).Methods(http.MethodGet)
	router.HandleFunc("/create-new-van", apiHandler.CreateNewVanHandler).Methods(http.MethodGet)
	router.HandleFunc("/van-details/update-van", apiHandler.PageUpdateVanHandler).Methods(http.MethodGet)
	router.HandleFunc("/update-van/{vanID}", apiHandler.UpdateVanHandler).Methods(http.MethodGet)
	router.HandleFunc("/van-details/delete-van", apiHandler.DeleteVanHandler).Methods(http.MethodGet)

	_ = godotenv.Load()

	port := os.Getenv("PORT")

	// Start server
	fmt.Println("Client server is listening on PORT ", port)
	router.ServeHTTP(w, r)

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
// 	router := mux.NewRouter()

// 	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
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
// 	router.HandleFunc("/", apiHandler.PageHomeHandler).Methods(http.MethodGet)
// 	router.HandleFunc("/sign-in", apiHandler.PageSignInHandler).Methods(http.MethodGet)
// 	router.HandleFunc("/logout", apiHandler.LogoutHandler).Methods(http.MethodGet)
// 	router.HandleFunc("/authorize", apiHandler.AuthorizeUserHandler).Methods(http.MethodGet)
// 	router.HandleFunc("/van-details", apiHandler.PageVanDetailsHandler).Methods(http.MethodGet)
// 	router.HandleFunc("/create-van", apiHandler.PageCreateVanHandler).Methods(http.MethodGet)
// 	router.HandleFunc("/create-new-van", apiHandler.CreateNewVanHandler).Methods(http.MethodGet)
// 	router.HandleFunc("/van-details/update-van", apiHandler.PageUpdateVanHandler).Methods(http.MethodGet)
// 	router.HandleFunc("/update-van/{vanID}", apiHandler.UpdateVanHandler).Methods(http.MethodGet)
// 	router.HandleFunc("/van-details/delete-van", apiHandler.DeleteVanHandler).Methods(http.MethodGet)

// 	_ = godotenv.Load()

// 	port := os.Getenv("PORT")

// 	// Start server
// 	fmt.Println("Client server is listening on PORT ", port)
// 	log.Fatal(http.ListenAndServe(":"+port, router))
// }
