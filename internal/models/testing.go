package models

import (
	"fmt"
	"strings"
)

// TestType represents different types of tests
type TestType string

const (
	UnitTest        TestType = "unit"
	IntegrationTest TestType = "integration"
	BenchmarkTest   TestType = "benchmark"
	MockTest        TestType = "mock"
	UtilsTest       TestType = "utils"
)

// TestFramework represents supported testing frameworks
type TestFramework string

const (
	TestifyFramework  TestFramework = "testify"
	GinkgoFramework   TestFramework = "ginkgo"
	GoConveyFramework TestFramework = "goconvey"
)

// TestTarget represents what to generate tests for
type TestTarget string

const (
	ResourceTarget    TestTarget = "resource"
	HandlerTarget     TestTarget = "handler"
	ServiceTarget     TestTarget = "service"
	RepositoryTarget  TestTarget = "repository"
	MiddlewareTarget  TestTarget = "middleware"
	APITarget         TestTarget = "api"
	AllTargets        TestTarget = "all"
)

// TestConfig represents test generation configuration
type TestConfig struct {
	Type        TestType      `json:"type"`
	Framework   TestFramework `json:"framework"`
	Target      TestTarget    `json:"target"`
	Name        string        `json:"name"`
	Package     string        `json:"package"`
	FullSuite   bool          `json:"full_suite"`
	WithMocks   bool          `json:"with_mocks"`
	WithUtils   bool          `json:"with_utils"`
	WithBench   bool          `json:"with_bench"`
	BDDStyle    bool          `json:"bdd_style"`
	Description string        `json:"description"`
}

// TestFile represents a test file to be generated
type TestFile struct {
	Name        string            `json:"name"`
	Package     string            `json:"package"`
	Type        TestType          `json:"type"`
	Framework   TestFramework     `json:"framework"`
	Imports     []string          `json:"imports"`
	TestCases   []TestCase        `json:"test_cases"`
	Mocks       []MockDefinition  `json:"mocks"`
	Utilities   []TestUtility     `json:"utilities"`
	Benchmarks  []BenchmarkCase   `json:"benchmarks"`
}

// TestCase represents a single test case
type TestCase struct {
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Setup           string            `json:"setup"`
	Input           map[string]string `json:"input"`
	ExpectedOutput  string            `json:"expected_output"`
	ExpectedError   string            `json:"expected_error"`
	Assertions      []string          `json:"assertions"`
	MockSetup       []string          `json:"mock_setup"`
	Cleanup         string            `json:"cleanup"`
}

// MockDefinition represents a mock definition
type MockDefinition struct {
	Name        string   `json:"name"`
	Interface   string   `json:"interface"`
	Package     string   `json:"package"`
	Methods     []string `json:"methods"`
}

// TestUtility represents a test utility function
type TestUtility struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Code        string   `json:"code"`
}

// BenchmarkCase represents a benchmark test case
type BenchmarkCase struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Setup       string `json:"setup"`
	Code        string `json:"code"`
}

// GetTestFileName returns the test file name
func (tf *TestFile) GetTestFileName() string {
	switch tf.Type {
	case BenchmarkTest:
		return fmt.Sprintf("%s_bench_test.go", strings.ToLower(tf.Name))
	case MockTest:
		return fmt.Sprintf("%s_mock.go", strings.ToLower(tf.Name))
	default:
		return fmt.Sprintf("%s_test.go", strings.ToLower(tf.Name))
	}
}

// GetTestPackagePath returns the package path for the test file
func (tf *TestFile) GetTestPackagePath() string {
	switch tf.Type {
	case UnitTest:
		return fmt.Sprintf("test/unit/%s", tf.Package)
	case IntegrationTest:
		return "test/integration"
	case BenchmarkTest:
		return "test/benchmark"
	case MockTest:
		return "test/mocks"
	case UtilsTest:
		return "test/utils"
	default:
		return "test"
	}
}

// GetDefaultImports returns default imports for the test framework
func (tf *TestFile) GetDefaultImports() []string {
	var imports []string
	
	// Common imports
	imports = append(imports, "testing")
	
	switch tf.Framework {
	case TestifyFramework:
		imports = append(imports,
			"github.com/stretchr/testify/assert",
			"github.com/stretchr/testify/require",
			"github.com/stretchr/testify/mock",
		)
		if tf.Type == IntegrationTest {
			imports = append(imports, "github.com/stretchr/testify/suite")
		}
		
	case GinkgoFramework:
		imports = append(imports,
			"github.com/onsi/ginkgo/v2",
			"github.com/onsi/gomega",
		)
		
	case GoConveyFramework:
		imports = append(imports,
			"github.com/smartystreets/goconvey/convey",
		)
	}
	
	// Add type-specific imports
	switch tf.Type {
	case UnitTest, IntegrationTest:
		imports = append(imports,
			"net/http",
			"net/http/httptest",
			"bytes",
			"encoding/json",
		)
		if strings.Contains(tf.Package, "handler") {
			imports = append(imports, "github.com/gin-gonic/gin")
		}
		
	case BenchmarkTest:
		imports = append(imports, "runtime")
		
	case MockTest:
		if tf.Framework == TestifyFramework {
			imports = append(imports, "github.com/stretchr/testify/mock")
		}
	}
	
	return imports
}

// GetTestSuiteStruct returns test suite struct for integration tests
func GetTestSuiteStruct(target TestTarget, name string) string {
	structName := fmt.Sprintf("%s%sTestSuite", strings.Title(name), strings.Title(string(target)))
	return structName
}

// GetDefaultTestCases returns default test cases for a target
func GetDefaultTestCases(target TestTarget, name string) []TestCase {
	switch target {
	case HandlerTarget:
		return getHandlerTestCases(name)
	case ServiceTarget:
		return getServiceTestCases(name)
	case RepositoryTarget:
		return getRepositoryTestCases(name)
	case MiddlewareTarget:
		return getMiddlewareTestCases(name)
	default:
		return []TestCase{
			{
				Name:        fmt.Sprintf("Test%sSuccess", strings.Title(name)),
				Description: fmt.Sprintf("Test successful %s operation", name),
				Setup:       "// Setup test data",
				ExpectedOutput: "expected result",
				Assertions:  []string{"assert.NoError(t, err)", "assert.NotNil(t, result)"},
			},
		}
	}
}

// getHandlerTestCases returns default test cases for handlers
func getHandlerTestCases(name string) []TestCase {
	resourceName := strings.Title(name)
	return []TestCase{
		{
			Name:        fmt.Sprintf("Test%sHandler_Get%s_Success", resourceName, resourceName),
			Description: fmt.Sprintf("Test successful %s retrieval", name),
			Setup: fmt.Sprintf(`gin.SetMode(gin.TestMode)
	mockService := &Mock%sService{}
	handler := New%sHandler(mockService)`, resourceName, resourceName),
			Input: map[string]string{
				"id": "123",
			},
			ExpectedOutput: fmt.Sprintf(`{"id":"123","name":"Test %s"}`, resourceName),
			MockSetup: []string{
				fmt.Sprintf(`mockService.On("Get%s", "123").Return(&%s{ID: "123", Name: "Test %s"}, nil)`, 
					resourceName, resourceName, resourceName),
			},
			Assertions: []string{
				"assert.Equal(t, http.StatusOK, w.Code)",
				"mockService.AssertExpectations(t)",
			},
			Cleanup: "mockService.AssertExpectations(t)",
		},
		{
			Name:        fmt.Sprintf("Test%sHandler_Get%s_NotFound", resourceName, resourceName),
			Description: fmt.Sprintf("Test %s not found", name),
			Setup: fmt.Sprintf(`gin.SetMode(gin.TestMode)
	mockService := &Mock%sService{}
	handler := New%sHandler(mockService)`, resourceName, resourceName),
			Input: map[string]string{
				"id": "999",
			},
			ExpectedError: "not found",
			MockSetup: []string{
				fmt.Sprintf(`mockService.On("Get%s", "999").Return(nil, errors.New("not found"))`, resourceName),
			},
			Assertions: []string{
				"assert.Equal(t, http.StatusNotFound, w.Code)",
				"mockService.AssertExpectations(t)",
			},
		},
		{
			Name:        fmt.Sprintf("Test%sHandler_Create%s_Success", resourceName, resourceName),
			Description: fmt.Sprintf("Test successful %s creation", name),
			Setup: fmt.Sprintf(`gin.SetMode(gin.TestMode)
	mockService := &Mock%sService{}
	handler := New%sHandler(mockService)`, resourceName, resourceName),
			Input: map[string]string{
				"name":  fmt.Sprintf("New %s", resourceName),
				"email": "test@example.com",
			},
			ExpectedOutput: fmt.Sprintf(`{"id":"new-id","name":"New %s"}`, resourceName),
			MockSetup: []string{
				fmt.Sprintf(`mockService.On("Create%s", mock.AnythingOfType("*models.%s")).Return(nil)`, 
					resourceName, resourceName),
			},
			Assertions: []string{
				"assert.Equal(t, http.StatusCreated, w.Code)",
				"mockService.AssertExpectations(t)",
			},
		},
	}
}

// getServiceTestCases returns default test cases for services
func getServiceTestCases(name string) []TestCase {
	resourceName := strings.Title(name)
	return []TestCase{
		{
			Name:        fmt.Sprintf("Test%sService_Get%s_Success", resourceName, resourceName),
			Description: fmt.Sprintf("Test successful %s retrieval from service", name),
			Setup: fmt.Sprintf(`mockRepo := &Mock%sRepository{}
	service := New%sService(mockRepo)`, resourceName, resourceName),
			Input: map[string]string{
				"id": "123",
			},
			MockSetup: []string{
				fmt.Sprintf(`mockRepo.On("GetByID", "123").Return(&%s{ID: "123", Name: "Test"}, nil)`, resourceName),
			},
			Assertions: []string{
				"assert.NoError(t, err)",
				"assert.NotNil(t, result)",
				`assert.Equal(t, "123", result.ID)`,
				"mockRepo.AssertExpectations(t)",
			},
		},
		{
			Name:        fmt.Sprintf("Test%sService_Create%s_ValidationError", resourceName, resourceName),
			Description: fmt.Sprintf("Test %s creation with validation error", name),
			Setup: fmt.Sprintf(`mockRepo := &Mock%sRepository{}
	service := New%sService(mockRepo)`, resourceName, resourceName),
			Input: map[string]string{
				"name": "", // Invalid empty name
			},
			ExpectedError: "validation error",
			Assertions: []string{
				"assert.Error(t, err)",
				"assert.Contains(t, err.Error(), \"validation\")",
			},
		},
	}
}

// getRepositoryTestCases returns default test cases for repositories
func getRepositoryTestCases(name string) []TestCase {
	resourceName := strings.Title(name)
	return []TestCase{
		{
			Name:        fmt.Sprintf("Test%sRepository_GetByID_Success", resourceName),
			Description: fmt.Sprintf("Test successful %s retrieval from repository", name),
			Setup: `db := setupTestDB()
	repo := NewUserRepository(db)`,
			Input: map[string]string{
				"id": "test-id",
			},
			Assertions: []string{
				"assert.NoError(t, err)",
				"assert.NotNil(t, result)",
				`assert.Equal(t, "test-id", result.ID)`,
			},
			Cleanup: "cleanupTestDB(db)",
		},
	}
}

// getMiddlewareTestCases returns default test cases for middleware
func getMiddlewareTestCases(name string) []TestCase {
	middlewareName := strings.Title(name)
	return []TestCase{
		{
			Name:        fmt.Sprintf("Test%sMiddleware_Success", middlewareName),
			Description: fmt.Sprintf("Test %s middleware allows valid request", name),
			Setup: fmt.Sprintf(`gin.SetMode(gin.TestMode)
	middleware := New%sMiddleware()`, middlewareName),
			Assertions: []string{
				"assert.Equal(t, http.StatusOK, w.Code)",
				"assert.True(t, nextCalled)",
			},
		},
		{
			Name:        fmt.Sprintf("Test%sMiddleware_Unauthorized", middlewareName),
			Description: fmt.Sprintf("Test %s middleware blocks invalid request", name),
			Setup: fmt.Sprintf(`gin.SetMode(gin.TestMode)
	middleware := New%sMiddleware()`, middlewareName),
			ExpectedError: "unauthorized",
			Assertions: []string{
				"assert.Equal(t, http.StatusUnauthorized, w.Code)",
				"assert.False(t, nextCalled)",
			},
		},
	}
}

// GetDefaultMocks returns default mock definitions for a target
func GetDefaultMocks(target TestTarget, name string) []MockDefinition {
	resourceName := strings.Title(name)
	
	switch target {
	case HandlerTarget:
		return []MockDefinition{
			{
				Name:      fmt.Sprintf("Mock%sService", resourceName),
				Interface: fmt.Sprintf("%sService", resourceName),
				Package:   "services",
				Methods: []string{
					fmt.Sprintf("Get%s(id string) (*%s, error)", resourceName, resourceName),
					fmt.Sprintf("Create%s(%s *%s) error", resourceName, strings.ToLower(name), resourceName),
					fmt.Sprintf("Update%s(id string, %s *%s) error", resourceName, strings.ToLower(name), resourceName),
					fmt.Sprintf("Delete%s(id string) error", resourceName),
				},
			},
		}
		
	case ServiceTarget:
		return []MockDefinition{
			{
				Name:      fmt.Sprintf("Mock%sRepository", resourceName),
				Interface: fmt.Sprintf("%sRepository", resourceName),
				Package:   "repositories",
				Methods: []string{
					fmt.Sprintf("GetByID(id string) (*%s, error)", resourceName),
					fmt.Sprintf("Create(%s *%s) error", strings.ToLower(name), resourceName),
					fmt.Sprintf("Update(id string, %s *%s) error", strings.ToLower(name), resourceName),
					"Delete(id string) error",
				},
			},
		}
		
	default:
		return []MockDefinition{}
	}
}

// IsValidTestType checks if test type is valid
func IsValidTestType(testType string) bool {
	validTypes := []string{
		string(UnitTest),
		string(IntegrationTest),
		string(BenchmarkTest),
		string(MockTest),
		string(UtilsTest),
	}
	
	for _, validType := range validTypes {
		if testType == validType {
			return true
		}
	}
	return false
}

// IsValidTestFramework checks if test framework is valid
func IsValidTestFramework(framework string) bool {
	validFrameworks := []string{
		string(TestifyFramework),
		string(GinkgoFramework),
		string(GoConveyFramework),
	}
	
	for _, validFramework := range validFrameworks {
		if framework == validFramework {
			return true
		}
	}
	return false
}

// IsValidTestTarget checks if test target is valid
func IsValidTestTarget(target string) bool {
	validTargets := []string{
		string(ResourceTarget),
		string(HandlerTarget),
		string(ServiceTarget),
		string(RepositoryTarget),
		string(MiddlewareTarget),
		string(APITarget),
		string(AllTargets),
	}
	
	for _, validTarget := range validTargets {
		if target == validTarget {
			return true
		}
	}
	return false
}