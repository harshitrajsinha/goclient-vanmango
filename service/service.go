package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/goclient-vanmango/apiclient"
	"github.com/joho/godotenv"
)

type authData struct {
	Username string
	Password string
}

type vanInfo []struct {
	VanCode  string
	Name     string
	Category string
	Price    string
	ImageURL string
}

type vanDetails struct {
	VanCode     string
	Name        string
	Brand       string
	Description string
	Category    string
	FuelType    string
	EngineID    string
	Price       string
	ImageURL    string
}

type Service struct {
	apiClient *apiclient.APIClient
}

func NewAPIService(apiClient *apiclient.APIClient) *Service {
	return &Service{
		apiClient: apiClient,
	}
}

func validateLoginForm(data *authData) string {

	_ = godotenv.Load()
	authUser := os.Getenv("AUTH_USER")
	authPass := os.Getenv("AUTH_PASS")

	// Name validation
	if data.Username != authUser {
		return "Invalid username"
	}

	// Password validation
	if data.Password != authPass {
		return "Invalid password"
	}

	return ""
}

func (s *Service) GetAllVans(w http.ResponseWriter, r *http.Request) (*vanInfo, error) {

	var err error
	var data *map[string]interface{}
	var vanDataToRender vanInfo
	_ = godotenv.Load()

	baseURL := os.Getenv("BASE_URL")
	requestUrl := baseURL + "/vans"

	data, err = s.apiClient.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	for _, item := range (*data)["data"].([]interface{}) { // item => [name: vanName]
		var vanData [5]string
		for key, value := range item.(map[string]interface{}) { // key,value => name, vanName

			if strings.ToLower(key) == "van-id" {
				vanData[0] = value.(string)
			} else if strings.ToLower(key) == "name" {
				vanData[1] = value.(string)
			} else if strings.ToLower(key) == "category" {
				vanData[2] = value.(string)
			} else if strings.ToLower(key) == "price" {
				floatToInt := int(value.(float64))
				vanData[3] = strconv.Itoa(floatToInt)
			} else if strings.ToLower(key) == "image-url" {
				vanData[4] = value.(string)
			}
		}

		vanDataToRender = append(vanDataToRender, struct {
			VanCode  string
			Name     string
			Category string
			Price    string
			ImageURL string
		}{VanCode: vanData[0], Name: vanData[1], Category: vanData[2], Price: vanData[3], ImageURL: vanData[4]})
	}

	return &vanDataToRender, nil

}

func (s *Service) GetVanByID(w http.ResponseWriter, r *http.Request, vanID string) (*vanDetails, error) {

	var err error
	var data *map[string]interface{}

	var vanDetailInfo vanDetails

	_ = godotenv.Load()

	baseURL := os.Getenv("BASE_URL")
	requestUrl := baseURL + "/van/" + vanID

	data, err = s.apiClient.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	var vanData [9]string

	for _, item := range (*data)["data"].([]interface{}) { // item => [name: vanName]

		for key, value := range item.(map[string]interface{}) { // key,value => name, vanName

			if strings.ToLower(key) == "van-id" {
				vanData[0] = value.(string)
			} else if strings.ToLower(key) == "name" {
				vanData[1] = value.(string)
			} else if strings.ToLower(key) == "brand" {
				vanData[2] = value.(string)
			} else if strings.ToLower(key) == "description" {
				vanData[3] = value.(string)
			} else if strings.ToLower(key) == "category" {
				vanData[4] = value.(string)
			} else if strings.ToLower(key) == "fuel-type" {
				vanData[5] = value.(string)
			} else if strings.ToLower(key) == "engine-id" {
				vanData[6] = value.(string)
			} else if strings.ToLower(key) == "price" {
				floatToInt := int(value.(float64))
				vanData[7] = strconv.Itoa(floatToInt)
			} else if strings.ToLower(key) == "image-url" {
				vanData[8] = value.(string)
			}
		}
	}

	vanDetailInfo = struct {
		VanCode     string
		Name        string
		Brand       string
		Description string
		Category    string
		FuelType    string
		EngineID    string
		Price       string
		ImageURL    string
	}{
		VanCode:     vanData[0],
		Name:        vanData[1],
		Brand:       vanData[2],
		Description: vanData[3],
		Category:    vanData[4],
		FuelType:    vanData[5],
		EngineID:    vanData[6],
		Price:       vanData[7],
		ImageURL:    vanData[8],
	}

	return &vanDetailInfo, nil

}

func (s *Service) UpdateVan(w http.ResponseWriter, r *http.Request, vanID string) (int64, error) {

	// Parse form data

	var err error
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return -1, err
	}

	// Get form values
	formData := vanDetails{
		Name:        r.FormValue("name"),
		Brand:       r.FormValue("brand"),
		Description: r.FormValue("description"),
		Category:    r.FormValue("category"),
		FuelType:    r.FormValue("fuel-type"),
		Price:       r.FormValue("price"),
		ImageURL:    r.FormValue("image-url"),
	}

	// form validation

	// 	// convert price -> string to int
	priceStrToInt, _ := strconv.Atoi(formData.Price)
	var formInput map[string]interface{} = map[string]interface{}{"name": formData.Name, "brand": formData.Brand, "description": formData.Description, "category": formData.Category, "engine-id": "e1f86b1a-0873-4c19-bae2-fc60329d0140", "fuel-type": formData.FuelType, "price": priceStrToInt, "image-url": formData.ImageURL}

	requestToSend, err := json.Marshal(formInput)
	if err != nil {
		return -1, err
	}

	_ = godotenv.Load()

	baseURL := os.Getenv("BASE_URL")
	requestUrl := baseURL + "/van/" + vanID

	_, err = s.apiClient.Put(requestUrl, &requestToSend)
	if err != nil {
		return -1, err
	} else {
		return 1, nil
	}
}

func (s *Service) DeleteVan(w http.ResponseWriter, r *http.Request, vanID string) (int64, error) {

	var err error

	_ = godotenv.Load()

	baseURL := os.Getenv("BASE_URL")
	requestUrl := baseURL + "/van/" + vanID

	data, err := s.apiClient.Delete(requestUrl)
	if err != nil {
		log.Println(data)
		return -1, err
	} else {
		log.Println(data)
		return 1, nil
	}
}

func (s *Service) AuthorizeUser(w http.ResponseWriter, r *http.Request) (int, error) {

	var token interface{}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		return -1, err
	}

	// Get form values
	formData := authData{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	// Get token id

	// Validate data
	isValid := validateLoginForm(&formData)

	if isValid != "" {
		return -1, fmt.Errorf("%s", isValid)
	}

	_ = godotenv.Load()

	baseURL := os.Getenv("BASE_URL")
	requestUrl := baseURL + "/login"

	requestToSend, err := json.Marshal(formData)
	if err != nil {
		return -1, err
	}

	response, err := s.apiClient.Post(requestUrl, &requestToSend)
	if err != nil {
		return -1, err
	}

	if (*response)["code"] == float64(201) {
		data := (*response)["data"].([]interface{})[0].(map[string]interface{})
		token = data["token"]
	} else {
		return -1, fmt.Errorf("%s", "No Auth token received")
	}

	// Create a new cookie
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    token.(string),
		Path:     "/",
		MaxAge:   1800, // 30mins in seconds
		HttpOnly: true,
		Secure:   true, // Only send over HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	// Set the cookie
	http.SetCookie(w, cookie)

	return 1, nil

}
