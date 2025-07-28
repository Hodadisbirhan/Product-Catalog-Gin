package main

import (
	"catalog-gin/config"
	"catalog-gin/router"
	"catalog-gin/seed"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func init() {
	// Load environment variables from .env file
	// godotenv helps manage environment variables during development.
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found or could not be loaded. Assuming environment variables are set directly.")
	}

	// Retrieve OAuth credentials and callback URL from environment variables.
	// It's crucial to get these from your Google Cloud Console project.
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleCallbackURL := os.Getenv("GOOGLE_CALLBACK_URL") // e.g., http://localhost:8080/auth/google/callback

	// Basic validation to ensure credentials are provided.
	if googleClientID == "" || googleClientSecret == "" || googleCallbackURL == "" {
		log.Fatal("Error: GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, and GOOGLE_CALLBACK_URL environment variables must be set.")
	}

	// Configure Goth to use the Google OAuth provider.
	// We define the scopes (permissions) our application needs.
	// "offline_access" is handled by SetAccessType below, not directly as a scope string.
	googleProvider := google.New(
		googleClientID,
		googleClientSecret,
		googleCallbackURL,
		"email",   // Request access to the user's email address
		"profile", // Request access to basic profile information
		"openid",  // Standard OpenID Connect scope
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	)

	// Explicitly set access type to "offline" to request a refresh token.
	// This is the correct way to ask for a refresh token with Google using Goth.
	googleProvider.SetAccessType("offline")
	// "SetPrompt("consent")" can be used to force the user to consent every time,
	// which ensures a refresh token is always issued, especially useful during development.
	googleProvider.SetPrompt("consent")

	goth.UseProviders(googleProvider)

	// Set up the session store for Gothic.
	// Gothic relies on gorilla/sessions for managing the OAuth state.
	// For production, consider a more robust store like RedisStore or MemcacheStore.
	// The `SESSION_KEY` should be a long, random, and securely generated secret.
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		log.Fatal("Error: SESSION_KEY environment variable is required for session management.")
	}

	maxAge := 86400 * 30                         // Session will expire in 30 days
	isProd := os.Getenv("GIN_MODE") == "release" // Check if running in production mode for secure cookies

	store := sessions.NewCookieStore([]byte(sessionKey))
	store.Options.Path = "/"      // Cookie valid across the entire application
	store.Options.MaxAge = maxAge // Set cookie expiry
	store.Options.HttpOnly = true // Prevents JavaScript access to the cookie
	store.Options.Secure = isProd // Only send cookie over HTTPS in production

	gothic.Store = store // Assign the configured store to Gothic
}

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to DB and run migrations
	config.ConnectDB()

	// Optional: Seed default roles and permissions
	seed.SeedRolesAndPermissions()

	// Start the server
	r := router.SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server is running on port %s\n", port)
	err := r.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
