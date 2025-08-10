package templates

import (
	"fmt"
	"strings"

	"github.com/vibercode/cli/internal/models"
)

// GetTestTemplate returns the test template based on framework and type
func GetTestTemplate(testFile models.TestFile) string {
	switch testFile.Framework {
	case models.TestifyFramework:
		return getTestifyTemplate(testFile)
	case models.GinkgoFramework:
		return getGinkgoTemplate(testFile)
	case models.GoConveyFramework:
		return getGoConveyTemplate(testFile)
	default:
		return getTestifyTemplate(testFile)
	}
}

// getTestifyTemplate generates Testify-based test template
func getTestifyTemplate(testFile models.TestFile) string {
	imports := generateTestImports(testFile.Imports)
	testCases := generateTestifyTestCases(testFile.TestCases, testFile.Package)
	
	packageName := testFile.Package
	if testFile.Type == models.IntegrationTest {
		packageName = "integration"
	}
	
	return fmt.Sprintf(`package %s

%s

%s
`, packageName, imports, testCases)
}

// generateTestifyTestCases generates Testify test cases
func generateTestifyTestCases(testCases []models.TestCase, packageType string) string {
	var cases []string
	
	for _, testCase := range testCases {
		caseCode := fmt.Sprintf(`func %s(t *testing.T) {
	// %s
	%s
	
	// Test cases
	tests := []struct {
		name           string
		%s
		expectedStatus int
		expectedError  string
	}{
		{
			name: "success case",
			expectedStatus: http.StatusOK,
		},
		{
			name: "error case", 
			expectedError: "error message",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			%s
			
			// Assertions
			%s
		})
	}
}`, testCase.Name, testCase.Description, testCase.Setup,
			generateTestCaseFields(packageType),
			generateTestCaseLogic(testCase, packageType),
			strings.Join(testCase.Assertions, "\n\t\t\t"))
		
		cases = append(cases, caseCode)
	}
	
	return strings.Join(cases, "\n\n")
}

// GetIntegrationTestTemplate generates integration test template
func GetIntegrationTestTemplate(testFile models.TestFile) string {
	resourceName := strings.Title(testFile.Name)
	
	switch testFile.Framework {
	case models.TestifyFramework:
		return fmt.Sprintf(`package integration

import (
	"testing"
	"net/http"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
)

// %sAPITestSuite tests %s API endpoints
type %sAPITestSuite struct {
	suite.Suite
	server *TestServer
	client *TestClient
}

// SetupSuite runs before all tests
func (suite *%sAPITestSuite) SetupSuite() {
	suite.server = NewTestServer()
	suite.client = NewTestClient(suite.server.URL)
}

// TearDownSuite runs after all tests
func (suite *%sAPITestSuite) TearDownSuite() {
	suite.server.Close()
}

// SetupTest runs before each test
func (suite *%sAPITestSuite) SetupTest() {
	suite.server.ResetDatabase()
}

// TestCreate%s tests %s creation
func (suite *%sAPITestSuite) TestCreate%s() {
	// Test data
	%sData := map[string]interface{}{
		"name":  "Test %s",
		"email": "test@example.com",
	}
	
	// Make request
	resp, err := suite.client.POST("/api/%s", %sData)
	suite.NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)
	
	// Verify response
	var %s %s
	err = json.NewDecoder(resp.Body).Decode(&%s)
	suite.NoError(err)
	suite.Equal("Test %s", %s.Name)
	suite.NotEmpty(%s.ID)
}

// TestGet%s tests %s retrieval
func (suite *%sAPITestSuite) TestGet%s() {
	// Create test %s
	%s := suite.createTest%s()
	
	// Make request
	resp, err := suite.client.GET("/api/%s/" + %s.ID)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Verify response
	var retrieved%s %s
	err = json.NewDecoder(resp.Body).Decode(&retrieved%s)
	suite.NoError(err)
	suite.Equal(%s.ID, retrieved%s.ID)
	suite.Equal(%s.Name, retrieved%s.Name)
}

// TestUpdate%s tests %s update
func (suite *%sAPITestSuite) TestUpdate%s() {
	// Create test %s
	%s := suite.createTest%s()
	
	// Update data
	updateData := map[string]interface{}{
		"name": "Updated %s",
	}
	
	// Make request
	resp, err := suite.client.PUT("/api/%s/"+%s.ID, updateData)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	
	// Verify update
	resp, err = suite.client.GET("/api/%s/" + %s.ID)
	suite.NoError(err)
	
	var updated%s %s
	err = json.NewDecoder(resp.Body).Decode(&updated%s)
	suite.NoError(err)
	suite.Equal("Updated %s", updated%s.Name)
}

// TestDelete%s tests %s deletion
func (suite *%sAPITestSuite) TestDelete%s() {
	// Create test %s
	%s := suite.createTest%s()
	
	// Make delete request
	resp, err := suite.client.DELETE("/api/%s/" + %s.ID)
	suite.NoError(err)
	suite.Equal(http.StatusNoContent, resp.StatusCode)
	
	// Verify deletion
	resp, err = suite.client.GET("/api/%s/" + %s.ID)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

// Helper method to create test %s
func (suite *%sAPITestSuite) createTest%s() *%s {
	%sData := map[string]interface{}{
		"name":  "Test %s",
		"email": "test@example.com",
	}
	
	resp, err := suite.client.POST("/api/%s", %sData)
	suite.NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)
	
	var %s %s
	err = json.NewDecoder(resp.Body).Decode(&%s)
	suite.NoError(err)
	
	return &%s
}

// TestIntegrationTestSuite runs the test suite
func Test%sAPITestSuite(t *testing.T) {
	suite.Run(t, new(%sAPITestSuite))
}
`, resourceName, strings.ToLower(resourceName), resourceName,
			resourceName, resourceName, resourceName, resourceName, resourceName, resourceName, resourceName,
			strings.ToLower(resourceName), resourceName, strings.ToLower(resourceName), strings.ToLower(resourceName),
			strings.ToLower(resourceName), resourceName, strings.ToLower(resourceName), resourceName, strings.ToLower(resourceName), strings.ToLower(resourceName),
			resourceName, strings.ToLower(resourceName), resourceName, resourceName, strings.ToLower(resourceName), resourceName, strings.ToLower(resourceName),
			resourceName, strings.ToLower(resourceName), strings.ToLower(resourceName), resourceName, strings.ToLower(resourceName), strings.ToLower(resourceName), strings.ToLower(resourceName),
			resourceName, strings.ToLower(resourceName), resourceName, resourceName, strings.ToLower(resourceName), resourceName, resourceName,
			strings.ToLower(resourceName), strings.ToLower(resourceName), strings.ToLower(resourceName), strings.ToLower(resourceName), resourceName, strings.ToLower(resourceName), resourceName, strings.ToLower(resourceName),
			resourceName, strings.ToLower(resourceName), resourceName, resourceName, strings.ToLower(resourceName), resourceName,
			strings.ToLower(resourceName), strings.ToLower(resourceName), resourceName, resourceName, strings.ToLower(resourceName), resourceName, strings.ToLower(resourceName),
			strings.ToLower(resourceName), strings.ToLower(resourceName), resourceName, strings.ToLower(resourceName), strings.ToLower(resourceName),
			resourceName, resourceName)
		
	default:
		return getTestifyTemplate(testFile)
	}
}

// GetMockTemplate generates mock template
func GetMockTemplate(testFile models.TestFile) string {
	if len(testFile.Mocks) == 0 {
		return ""
	}

	mock := testFile.Mocks[0]
	
	return fmt.Sprintf(`package mocks

import (
	"github.com/stretchr/testify/mock"
)

// %s is a mock implementation of %s
type %s struct {
	mock.Mock
}

%s

// NewMock%s creates a new mock instance
func New%s() *%s {
	return &%s{}
}
`, mock.Name, mock.Interface, mock.Name, 
		generateMockMethods(mock.Methods, mock.Name),
		mock.Interface, mock.Name, mock.Name, mock.Name)
}

// GetBenchmarkTemplate generates benchmark test template
func GetBenchmarkTemplate(testFile models.TestFile) string {
	resourceName := strings.Title(testFile.Name)
	packageType := testFile.Package
	
	return fmt.Sprintf(`package benchmark

import (
	"testing"
	"runtime"
)

// Benchmark%s%s benchmarks %s %s operations
func Benchmark%s%s(b *testing.B) {
	// Setup
	%s
	
	// Reset timer after setup
	b.ResetTimer()
	
	// Run benchmark
	for i := 0; i < b.N; i++ {
		// Benchmark code here
		_ = i
	}
}

// Benchmark%s%sMemory benchmarks %s %s memory usage
func Benchmark%s%sMemory(b *testing.B) {
	// Setup
	%s
	
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	// Reset timer
	b.ResetTimer()
	
	// Run benchmark
	for i := 0; i < b.N; i++ {
		// Benchmark code here
		_ = i
	}
	
	runtime.GC()
	runtime.ReadMemStats(&m2)
	
	b.Logf("Total Alloc: %%d bytes", m2.TotalAlloc-m1.TotalAlloc)
	b.Logf("Mallocs: %%d", m2.Mallocs-m1.Mallocs)
}

// Benchmark%s%sConcurrent benchmarks concurrent %s %s operations
func Benchmark%s%sConcurrent(b *testing.B) {
	// Setup
	%s
	
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Concurrent benchmark code here
		}
	})
}
`, resourceName, strings.Title(packageType), strings.ToLower(resourceName), packageType,
		resourceName, strings.Title(packageType), generateBenchmarkSetup(packageType, resourceName),
		resourceName, strings.Title(packageType), strings.ToLower(resourceName), packageType,
		resourceName, strings.Title(packageType), generateBenchmarkSetup(packageType, resourceName),
		resourceName, strings.Title(packageType), strings.ToLower(resourceName), packageType,
		resourceName, strings.Title(packageType), generateBenchmarkSetup(packageType, resourceName))
}

// GetTestUtilityTemplate generates test utility templates
func GetTestUtilityTemplate(utilityName string, framework models.TestFramework) string {
	switch utilityName {
	case "test_database.go":
		return getTestDatabaseUtility()
	case "test_server.go":
		return getTestServerUtility()
	case "test_client.go":
		return getTestClientUtility()
	case "test_factories.go":
		return getTestFactoriesUtility()
	default:
		return ""
	}
}

// Helper functions for template generation

func generateTestImports(imports []string) string {
	if len(imports) == 0 {
		return ""
	}
	
	var importLines []string
	for _, imp := range imports {
		importLines = append(importLines, fmt.Sprintf("\t\"%s\"", imp))
	}
	
	return fmt.Sprintf("import (\n%s\n)", strings.Join(importLines, "\n"))
}

func generateTestCaseFields(packageType string) string {
	switch packageType {
	case "handlers":
		return `input map[string]string
		setupMocks func()`
	case "services":
		return `input interface{}
		setupMocks func()`
	case "repositories":
		return `input interface{}`
	default:
		return `input interface{}`
	}
}

func generateTestCaseLogic(testCase models.TestCase, packageType string) string {
	switch packageType {
	case "handlers":
		return `// Setup mocks
			if tt.setupMocks != nil {
				tt.setupMocks()
			}
			
			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			
			// Set request parameters
			for key, value := range tt.input {
				c.Param(key, value)
			}
			
			// Execute handler
			// handler.Method(c)`
	case "services":
		return `// Setup mocks
			if tt.setupMocks != nil {
				tt.setupMocks()
			}
			
			// Execute service method
			// result, err := service.Method(tt.input)`
	default:
		return `// Execute test logic
			// result, err := target.Method(tt.input)`
	}
}

func generateMockMethods(methods []string, mockName string) string {
	var mockMethods []string
	
	for _, method := range methods {
		// Extract method name from signature
		parts := strings.Split(method, "(")
		if len(parts) > 0 {
			methodName := parts[0]
			mockMethod := fmt.Sprintf(`// %s mocks the %s method
func (m *%s) %s {
	args := m.Called(%s)
	%s
}`, methodName, methodName, mockName, method,
				generateMockCallArgs(method),
				generateMockReturn(method))
			
			mockMethods = append(mockMethods, mockMethod)
		}
	}
	
	return strings.Join(mockMethods, "\n\n")
}

func generateMockCallArgs(methodSig string) string {
	// Simple implementation - extract parameter names
	if strings.Contains(methodSig, "(") && strings.Contains(methodSig, ")") {
		params := strings.Split(strings.Split(methodSig, "(")[1], ")")[0]
		if params == "" {
			return ""
		}
		
		// Extract parameter names (simplified)
		var argNames []string
		paramList := strings.Split(params, ",")
		for _, param := range paramList {
			param = strings.TrimSpace(param)
			parts := strings.Fields(param)
			if len(parts) >= 2 {
				argNames = append(argNames, parts[0])
			}
		}
		
		return strings.Join(argNames, ", ")
	}
	return ""
}

func generateMockReturn(methodSig string) string {
	// Extract return types and generate appropriate return statements
	if strings.Contains(methodSig, ")") {
		parts := strings.Split(methodSig, ")")
		if len(parts) > 1 {
			returnPart := strings.TrimSpace(parts[1])
			if returnPart == "" {
				return "return"
			}
			
			// Handle common return patterns
			if strings.Contains(returnPart, "error") {
				if strings.Contains(returnPart, "*") || strings.Contains(returnPart, "interface") {
					return "return args.Get(0), args.Error(1)"
				}
				return "return args.Error(0)"
			} else if strings.Contains(returnPart, "*") || strings.Contains(returnPart, "interface") {
				return "return args.Get(0)"
			}
		}
	}
	return "return"
}

func generateBenchmarkSetup(packageType, resourceName string) string {
	switch packageType {
	case "handlers":
		return fmt.Sprintf(`gin.SetMode(gin.TestMode)
	handler := New%sHandler(nil)`, resourceName)
	case "services":
		return fmt.Sprintf(`service := New%sService(nil)`, resourceName)
	case "repositories":
		return fmt.Sprintf(`repo := New%sRepository(nil)`, resourceName)
	default:
		return "// Setup benchmark environment"
	}
}

// Specific utility templates

func getTestDatabaseUtility() string {
	return `package utils

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	
	_ "github.com/lib/pq"
)

// TestDatabase manages test database operations
type TestDatabase struct {
	DB   *sql.DB
	name string
}

// NewTestDatabase creates a new test database
func NewTestDatabase(t *testing.T) *TestDatabase {
	// Create unique test database name
	dbName := fmt.Sprintf("test_%s_%d", t.Name(), time.Now().Unix())
	
	// Connect to postgres to create test database
	masterDB, err := sql.Open("postgres", "postgres://localhost/postgres?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to master database: %v", err)
	}
	defer masterDB.Close()
	
	// Create test database
	_, err = masterDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	
	// Connect to test database
	testDB, err := sql.Open("postgres", fmt.Sprintf("postgres://localhost/%s?sslmode=disable", dbName))
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	
	return &TestDatabase{
		DB:   testDB,
		name: dbName,
	}
}

// Reset clears all data from test database
func (td *TestDatabase) Reset() {
	// Get all table names
	rows, err := td.DB.Query("SELECT tablename FROM pg_tables WHERE schemaname = 'public'")
	if err != nil {
		log.Printf("Failed to get table names: %v", err)
		return
	}
	defer rows.Close()
	
	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			continue
		}
		tables = append(tables, tableName)
	}
	
	// Truncate all tables
	for _, table := range tables {
		_, err := td.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			log.Printf("Failed to truncate table %s: %v", table, err)
		}
	}
}

// Close closes test database and cleans up
func (td *TestDatabase) Close() {
	td.DB.Close()
	
	// Connect to master database to drop test database
	masterDB, err := sql.Open("postgres", "postgres://localhost/postgres?sslmode=disable")
	if err != nil {
		log.Printf("Failed to connect to master database: %v", err)
		return
	}
	defer masterDB.Close()
	
	// Drop test database
	_, err = masterDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", td.name))
	if err != nil {
		log.Printf("Failed to drop test database: %v", err)
	}
}`
}

func getTestServerUtility() string {
	return `package utils

import (
	"net/http/httptest"
	"github.com/gin-gonic/gin"
)

// TestServer wraps httptest.Server for API testing
type TestServer struct {
	*httptest.Server
	Router *gin.Engine
	DB     *TestDatabase
}

// NewTestServer creates a new test server
func NewTestServer() *TestServer {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	
	// Setup test database
	// db := NewTestDatabase()
	
	// Setup routes
	setupRoutes(router)
	
	server := httptest.NewServer(router)
	
	return &TestServer{
		Server: server,
		Router: router,
		// DB:     db,
	}
}

// ResetDatabase clears all test data
func (ts *TestServer) ResetDatabase() {
	if ts.DB != nil {
		ts.DB.Reset()
	}
}

// Close shuts down the test server
func (ts *TestServer) Close() {
	ts.Server.Close()
	if ts.DB != nil {
		ts.DB.Close()
	}
}

func setupRoutes(router *gin.Engine) {
	// Setup your API routes here
	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
		
		// Add your routes here
		// users := api.Group("/users")
		// {
		//     users.GET("", getUsersHandler)
		//     users.POST("", createUserHandler)
		//     users.GET("/:id", getUserHandler)
		//     users.PUT("/:id", updateUserHandler)
		//     users.DELETE("/:id", deleteUserHandler)
		// }
	}
}`
}

func getTestClientUtility() string {
	return `package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TestClient provides HTTP client for API testing
type TestClient struct {
	baseURL string
	client  *http.Client
}

// NewTestClient creates a new test client
func NewTestClient(baseURL string) *TestClient {
	return &TestClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// GET makes a GET request
func (tc *TestClient) GET(path string) (*http.Response, error) {
	return tc.client.Get(tc.baseURL + path)
}

// POST makes a POST request with JSON body
func (tc *TestClient) POST(path string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	
	return tc.client.Post(
		tc.baseURL+path,
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
}

// PUT makes a PUT request with JSON body
func (tc *TestClient) PUT(path string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest("PUT", tc.baseURL+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	
	return tc.client.Do(req)
}

// DELETE makes a DELETE request
func (tc *TestClient) DELETE(path string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", tc.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	
	return tc.client.Do(req)
}

// PATCH makes a PATCH request with JSON body
func (tc *TestClient) PATCH(path string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest("PATCH", tc.baseURL+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	
	return tc.client.Do(req)
}

// DecodeJSON decodes JSON response body
func (tc *TestClient) DecodeJSON(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

// ReadBody reads response body as string
func (tc *TestClient) ReadBody(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}`
}

func getTestFactoriesUtility() string {
	return `package utils

import (
	"fmt"
	"time"
)

// UserFactory creates test user data
type UserFactory struct{}

// NewUserFactory creates a new user factory
func NewUserFactory() *UserFactory {
	return &UserFactory{}
}

// CreateUser creates a test user
func (f *UserFactory) CreateUser(overrides ...map[string]interface{}) map[string]interface{} {
	user := map[string]interface{}{
		"name":       "Test User",
		"email":      fmt.Sprintf("test-%d@example.com", time.Now().Unix()),
		"age":        25,
		"created_at": time.Now(),
	}
	
	// Apply overrides
	for _, override := range overrides {
		for key, value := range override {
			user[key] = value
		}
	}
	
	return user
}

// CreateUsers creates multiple test users
func (f *UserFactory) CreateUsers(count int) []map[string]interface{} {
	var users []map[string]interface{}
	
	for i := 0; i < count; i++ {
		user := f.CreateUser(map[string]interface{}{
			"name":  fmt.Sprintf("Test User %d", i+1),
			"email": fmt.Sprintf("test-%d-%d@example.com", i+1, time.Now().Unix()),
		})
		users = append(users, user)
	}
	
	return users
}

// Example factory for other entities
// Add more factories as needed:
//
// type ProductFactory struct{}
// 
// func (f *ProductFactory) CreateProduct(overrides ...map[string]interface{}) map[string]interface{} {
//     product := map[string]interface{}{
//         "name":        "Test Product",
//         "description": "A test product",
//         "price":       99.99,
//         "created_at":  time.Now(),
//     }
//     
//     for _, override := range overrides {
//         for key, value := range override {
//             product[key] = value
//         }
//     }
//     
//     return product
// }`
}

// getGinkgoTemplate generates Ginkgo/Gomega BDD-style tests
func getGinkgoTemplate(testFile models.TestFile) string {
	imports := generateTestImports(testFile.Imports)
	
	return fmt.Sprintf(`package %s

%s

var _ = Describe("%s", func() {
	var (
		// Setup variables here
	)
	
	BeforeEach(func() {
		// Setup before each test
	})
	
	AfterEach(func() {
		// Cleanup after each test
	})
	
	Context("when testing %s functionality", func() {
		It("should handle success case", func() {
			// Test implementation
			Expect(true).To(BeTrue())
		})
		
		It("should handle error case", func() {
			// Test implementation
			Expect(false).To(BeFalse())
		})
	})
})
`, testFile.Package, imports, strings.Title(testFile.Name), strings.ToLower(testFile.Name))
}

// getGoConveyTemplate generates GoConvey-style tests
func getGoConveyTemplate(testFile models.TestFile) string {
	imports := generateTestImports(testFile.Imports)
	
	return fmt.Sprintf(`package %s

%s

func Test%s(t *testing.T) {
	Convey("Given a %s", t, func() {
		// Setup
		
		Convey("When performing an operation", func() {
			// Action
			
			Convey("Then it should succeed", func() {
				// Assertion
				So(true, ShouldBeTrue)
			})
		})
		
		Convey("When encountering an error", func() {
			// Error scenario
			
			Convey("Then it should handle the error properly", func() {
				// Error assertion
				So(false, ShouldBeFalse)
			})
		})
	})
}
`, testFile.Package, imports, strings.Title(testFile.Name), strings.ToLower(testFile.Name))
}