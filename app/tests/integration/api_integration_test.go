package integration

import (
	"net/http"
	"testing"

	"base-golang-restful-app/config"
	"base-golang-restful-app/database"
	"base-golang-restful-app/models"
	helperPkg "base-golang-restful-app/testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APIIntegrationTestSuite struct {
	suite.Suite
	helper *helperPkg.TestHelper
	router *gin.Engine
	token  string
	userID uuid.UUID
}

func (suite *APIIntegrationTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// Load test config
	cfg, err := config.Load("../../.env")
	if err != nil {
		suite.T().Skip("Skipping integration tests - no database config")
	}

	// Connect to test database
	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.Name + "_test",
		SSLMode:  cfg.Database.SSLMode,
	}

	if err := database.Connect(dbConfig); err != nil {
		suite.T().Skip("Skipping integration tests - cannot connect to database")
	}

	// Run migrations
	database.GetDB().AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.Product{},
	)

	// Setup router (would need actual router setup)
	suite.router = gin.New()
	suite.helper = helperPkg.NewTestHelper(suite.T(), suite.router)
}

func (suite *APIIntegrationTestSuite) TearDownSuite() {
	database.Close()
}

func (suite *APIIntegrationTestSuite) SetupTest() {
	// Clean database before each test
	database.GetDB().Exec("TRUNCATE users, roles, permissions, products CASCADE")
}

func (suite *APIIntegrationTestSuite) TestAuthFlow() {
	// Test registration
	suite.Run("User Registration", func() {
		w := suite.helper.POST("/api/v1/auth/register", map[string]interface{}{
			"email":      "test@example.com",
			"password":   "Password123!",
			"first_name": "Test",
			"last_name":  "User",
		})

		suite.helper.AssertStatus(w, http.StatusCreated)

		var response map[string]interface{}
		suite.helper.ParseJSON(w, &response)

		assert.Contains(suite.T(), response, "data")
	})

	// Test login
	suite.Run("User Login", func() {
		// Register user first
		suite.helper.POST("/api/v1/auth/register", map[string]interface{}{
			"email":      "login@example.com",
			"password":   "Password123!",
			"first_name": "Login",
			"last_name":  "Test",
		})

		// Login
		w := suite.helper.POST("/api/v1/auth/login", map[string]interface{}{
			"email":    "login@example.com",
			"password": "Password123!",
		})

		suite.helper.AssertStatus(w, http.StatusOK)

		var response map[string]interface{}
		suite.helper.ParseJSON(w, &response)

		data := response["data"].(map[string]interface{})
		assert.Contains(suite.T(), data, "access_token")
		assert.Contains(suite.T(), data, "refresh_token")
	})
}

func (suite *APIIntegrationTestSuite) TestProductCRUD() {
	// Setup: Create user and get token
	authHelper := suite.helper.Auth()
	token := authHelper.Login("admin@example.com", "admin123")

	suite.Run("Create Product", func() {
		w := suite.helper.POST("/api/v1/products", map[string]interface{}{
			"name":        "Test Product",
			"description": "Test Description",
			"price":       99.99,
			"stock":       100,
			"category":    "Electronics",
		}, map[string]string{
			"Authorization": "Bearer " + token,
		})

		suite.helper.AssertStatus(w, http.StatusCreated)
	})

	suite.Run("List Products", func() {
		w := suite.helper.GET("/api/v1/products")
		suite.helper.AssertStatus(w, http.StatusOK)
	})
}

func TestAPIIntegrationSuite(t *testing.T) {
	suite.Run(t, new(APIIntegrationTestSuite))
}
