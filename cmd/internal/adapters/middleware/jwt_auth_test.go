package middleware_test

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"swapp-go/cmd/internal/adapters/middleware"
	"testing"
	"time"
)

// Helper Functions
func init() {
	_ = os.Setenv("JWT_SECRET", "test_jwt_secret")
}

func jwtSecret() string {
	return os.Getenv("JWT_SECRET")
}

func generateClaims(expiration time.Time) jwt.MapClaims {
	return jwt.MapClaims{
		"email": "test@email.com",
		"sub":   "6e9648ee-fd0b-4267-adcb-0c03b0176277",
		"exp":   expiration.Unix(),
	}
}

func generateTestToken(t *testing.T, claims jwt.MapClaims) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret()))
	assert.NoError(t, err)

	return signedToken
}

//func performRequestWithToken(t *testing.T, token string) *httptest.ResponseRecorder {
//	t.Helper()
//	gin.SetMode(gin.TestMode)
//
//	router := gin.New()
//	router.Use(middleware.JwtAuthMiddleware(jwtSecret()))
//
//	router.GET("/protected", func(context *gin.Context) {
//		userID := context.GetString("userID")
//		email := context.GetString("email")
//
//		context.JSON(http.StatusOK, gin.H{
//			"userID": userID,
//			"email":  email,
//		})
//	})
//
//	request, _ := http.NewRequest(http.MethodGet, "/protected", nil)
//	if token != "" {
//		request.Header.Set("Authorization", "Bearer "+token)
//	}
//
//	response := httptest.NewRecorder()
//
//	router.ServeHTTP(response, request)
//
//	return response
//}

// Tests
//func TestJwtAuthMiddleware_ValidToken(t *testing.T) {
//	claims := generateClaims(time.Now().Add(time.Hour))
//	token := generateTestToken(t, claims)
//	response := performRequestWithToken(t, token)
//
//	assert.Equal(t, http.StatusOK, response.Code)
//	assert.Contains(t, response.Body.String(), "6e9648ee-fd0b-4267-adcb-0c03b0176277")
//	assert.Contains(t, response.Body.String(), "test@email.com")
//}
//
//func TestJwtAuthMiddleware_MissingHeader(t *testing.T) {
//	resp := performRequestWithToken(t, "") // No token
//
//	assert.Equal(t, http.StatusUnauthorized, resp.Code)
//	assert.Contains(t, resp.Body.String(), "Authorization header is missing")
//}
//
//func TestJwtAuthMiddleware_MalformedHeader(t *testing.T) {
//	gin.SetMode(gin.TestMode)
//
//	request, _ := http.NewRequest(http.MethodGet, "/protected", nil)
//	request.Header.Set("Authorization", "missing-bearer-token")
//
//	response := httptest.NewRecorder()
//
//	router := gin.New()
//	router.Use(middleware.JwtAuthMiddleware(os.Getenv("JWT_SECRET")))
//	router.GET("/protected", func(c *gin.Context) {
//		c.JSON(http.StatusOK, gin.H{"message": "OK"})
//	})
//
//	router.ServeHTTP(response, request)
//	assert.Equal(t, http.StatusUnauthorized, response.Code)
//	assert.Contains(t, response.Body.String(), "Authorization header is malformed")
//}
//
//func TestJwtAuthMiddleware_InvalidToken(t *testing.T) {
//	response := performRequestWithToken(t, "invalid-token")
//
//	assert.Equal(t, http.StatusUnauthorized, response.Code)
//	assert.Contains(t, response.Body.String(), "Invalid or expired token")
//}
//
//func TestJwtAuthMiddleware_ExpiredToken(t *testing.T) {
//	claims := generateClaims(time.Now().Add(-time.Hour))
//	token := generateTestToken(t, claims)
//	response := performRequestWithToken(t, token)
//
//	assert.Equal(t, http.StatusUnauthorized, response.Code)
//	assert.Contains(t, response.Body.String(), "Invalid or expired token")
//}

// Table-driven approach
func TestJwtAuthMiddleware_Scenarios(t *testing.T) {
	validClaims := generateClaims(time.Now().Add(time.Hour))
	expiredClaims := generateClaims(time.Now().Add(-time.Hour))

	validToken := generateTestToken(t, validClaims)
	expiredToken := generateTestToken(t, expiredClaims)

	testCases := []struct {
		name                 string
		authorizationHeader  string
		expectedStatusCode   int
		expectedBodyContains string
	}{
		{
			name:                 "Missing Authorization Header",
			authorizationHeader:  "",
			expectedStatusCode:   http.StatusUnauthorized,
			expectedBodyContains: "Authorization header is missing",
		},
		{
			name:                 "Malformed Authorization Header",
			authorizationHeader:  "missing-bearer-token",
			expectedStatusCode:   http.StatusUnauthorized,
			expectedBodyContains: "Authorization header is malformed",
		},
		{
			name:                 "Invalid Token",
			authorizationHeader:  "Bearer invalid-token",
			expectedStatusCode:   http.StatusUnauthorized,
			expectedBodyContains: "Invalid or expired token",
		},
		{
			name:                 "Expired Token",
			authorizationHeader:  "Bearer " + expiredToken,
			expectedStatusCode:   http.StatusUnauthorized,
			expectedBodyContains: "Invalid or expired token",
		},
		{
			name:                 "Valid Token",
			authorizationHeader:  "Bearer " + validToken,
			expectedStatusCode:   http.StatusOK,
			expectedBodyContains: "test@email.com",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			router := gin.New()
			router.Use(middleware.JwtAuthMiddleware(jwtSecret()))
			router.GET("/protected", func(context *gin.Context) {
				userID := context.GetString("userID")
				email := context.GetString("email")

				context.JSON(http.StatusOK, gin.H{
					"userID": userID,
					"email":  email,
				})
			})

			httpRequest, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			if testCase.authorizationHeader != "" {
				httpRequest.Header.Set("Authorization", testCase.authorizationHeader)
			}

			httpResponseRecorder := httptest.NewRecorder()
			router.ServeHTTP(httpResponseRecorder, httpRequest)

			assert.Equal(t, testCase.expectedStatusCode, httpResponseRecorder.Code)
			assert.Contains(t, httpResponseRecorder.Body.String(), testCase.expectedBodyContains)
		})
	}
}
