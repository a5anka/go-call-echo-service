package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"

    "golang.org/x/oauth2/clientcredentials"
)

func callEchoService(w http.ResponseWriter, r *http.Request) {
    // Retrieve environment variables (same as before)
    serviceURL := os.Getenv("SERVICE_URL")
    clientID := os.Getenv("OAUTH_CLIENT_ID")
    clientSecret := os.Getenv("OAUTH_CLIENT_SECRET")
    tokenURL := os.Getenv("OAUTH_TOKEN_URL")


    // Print the read environment variables
    fmt.Println("Environment variables:")
    fmt.Println("SERVICE_URL:", serviceURL)
    fmt.Println("OAUTH_CLIENT_ID:", clientID)
    fmt.Println("OAUTH_CLIENT_SECRET:", clientSecret)
    fmt.Println("OAUTH_TOKEN_URL:", tokenURL)
    
    // OAuth2 configuration (same as before)
    config := &clientcredentials.Config{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        TokenURL:     tokenURL,
    }

    // Get an OAuth2 token (same as before)
    _, err := config.Token(context.Background())
    if err != nil {
        http.Error(w, "Error getting OAuth2 token", http.StatusInternalServerError)
        log.Println("Error getting OAuth2 token:", err)
        return
    }

    // Create an HTTP client with the OAuth2 token (same as before)
    client := config.Client(context.Background())

    // Make the request to the echo microservice
    resp, err := client.Get(serviceURL + "/echo?input=hello") 
    if err != nil {
        http.Error(w, "Error making request to echo service", http.StatusInternalServerError)
        log.Println("Error making request to echo service:", err)
        return
    }
    defer resp.Body.Close()

    // Copy the response from the echo service to the client's response
    _, err = io.Copy(w, resp.Body)
    if err != nil {
        http.Error(w, "Error reading response from echo service", http.StatusInternalServerError)
        log.Println("Error reading response from echo service:", err)
        return
    }
}

func main() {
    http.HandleFunc("/call-echo", callEchoService)

    fmt.Println("Client service listening on port 9000...")
    log.Fatal(http.ListenAndServe(":9000", nil))
}