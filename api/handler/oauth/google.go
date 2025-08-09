package oauth

import (
	"anchor-blog/api/handler"
	"anchor-blog/config"
	usersvc "anchor-blog/internal/service/user"
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig *oauth2.Config

const oauthStateCookieName = "oauthstate"

func InitializeGoogleOAuthConfig(cfg *config.Config) {
	googleOAuthConfig = &oauth2.Config{
		RedirectURL:  cfg.OAuth.Google.RedirectURI,
		ClientID:     cfg.OAuth.Google.ClientID,
		ClientSecret: cfg.OAuth.Google.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

type OAuthHandler struct {
	userService *usersvc.UserServices
}

func NewOAuthHandler(us *usersvc.UserServices) *OAuthHandler {
	return &OAuthHandler{userService: us}
}

// GoogleLogin initiates the Google OAuth2 login flow.
func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
	// Generate a random state string for CSRF protection.
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Set the state in a short-lived cookie.
	c.SetCookie(oauthStateCookieName, state, 3600, "/", "localhost", false, true)

	// Redirect the user to Google's consent page.
	url := googleOAuthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the callback from Google after user authentication.
func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	// Compare the state from the cookie with the state from the query parameter.
	cookieState, err := c.Cookie(oauthStateCookieName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state cookie not found"})
		return
	}
	if c.Query("state") != cookieState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state token"})
		return
	}

	code := c.Query("code")
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange code for token"})
		return
	}

	// Fetch user info from Google.
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read user info response"})
		return
	}

	result, err := h.userService.HandleGoogleLogin(c.Request.Context(), contents)
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
