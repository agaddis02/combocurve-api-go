# combocurve-api-go
Simple go project, meant to replicate the light ComboCurve API Python SDK to allow for auth and pagination


### Imports 
these are the main imports needed to complete authorization, and pagination with combocurve

```
"github.com/agaddis02/combocurve-api-go/auth/combocurve_auth" 
"github.com/agaddis02/combocurve-api-go/auth/service_account"
"github.com/agaddis02/combocurve-api-go/pagination"
"github.com/agaddis02/combocurve-api-go/models"
```


### Example Usage
```
import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/agaddis02/combocurve-api-go/auth/combocurve_auth"
	"github.com/agaddis02/combocurve-api-go/auth/service_account"
	"github.com/agaddis02/combocurve-api-go/pagination"
)

func main() {

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	xApiKey := os.Getenv("X_API_KEY")

	// Load in Service Account Json - Currently only supports from file
	sa, err := service_account.FromFile("./api-service-account-key.json")
	if err != nil {
		log.Fatalf("Failed to load service account: %v", err)
	}

	// Create ComboCurveAuth
	auth, err := combocurve_auth.NewComboCurveAuth(&sa, xApiKey, 60, 60*60)
	if err != nil {
		log.Fatalf("Failed to create ComboCurveAuth: %v", err)
	}

	// Get auth headers
	headers, err := auth.GetAuthHeaders()
	if err != nil {
		log.Fatalf("Failed to get auth headers: %v", err)
	}
	fmt.Printf("Auth headers: %v\n", headers)

	// Make request
	url := "https://api.combocurve.com/v1/monthly-productions?take=1"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-api-key", headers["x-api-key"])
	req.Header.Add("Authorization", headers["Authorization"])

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	// Print response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

	// Get next page url - similar to python library, this will just be an empty string is nothing further availble to proceed
	nextLink := pagination.GetNextPageURL(res.Header)

	fmt.Println(nextLink)
}

```