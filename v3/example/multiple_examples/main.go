package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-swagno/swagno/v3"
	"github.com/go-swagno/swagno/v3/components/endpoint"
	"github.com/go-swagno/swagno/v3/components/http/response"
	"github.com/go-swagno/swagno/v3/components/parameter"
)

// User represents a user in the system
type User struct {
	ID        uint64 `json:"id" example:"12345"`
	Username  string `json:"username" example:"john_doe"`
	Email     string `json:"email" example:"john@example.com"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Username  string `json:"username" example:"new_user"`
	Email     string `json:"email" example:"new@example.com"`
	FirstName string `json:"first_name" example:"New"`
	LastName  string `json:"last_name" example:"User"`
	Password  string `json:"password" example:"securePassword123"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"Invalid request"`
}

func main() {
	// Create a new OpenAPI instance
	openapi := swagno3.New(swagno3.Config{
		Title:       "User API with Multiple Examples",
		Version:     "1.0.0",
		Summary:     "A comprehensive example of OpenAPI with multiple examples support",
		Description: "This API demonstrates how to use multiple examples in OpenAPI 3.0 documentation",
	})

	// Add some servers
	openapi.AddServer("https://api.example.com/v1", "Production server")
	openapi.AddServer("https://sandbox.api.example.com/v1", "Sandbox server")

	// Define endpoints with multiple examples

	// GET /users - List all users with multiple parameter examples
	listUsersEndpoint := endpoint.New(
		endpoint.GET,
		"/users",
		endpoint.WithTags("users"),
		endpoint.WithSummary("List all users"),
		endpoint.WithDescription("Returns a list of all users in the system"),
		endpoint.WithParams(
			parameter.IntParam("limit", parameter.Query,
				parameter.WithDescription("Maximum number of users to return"),
				parameter.WithExample(10),
				parameter.WithExamples(map[string]interface{}{
					"smallPage":  5,
					"mediumPage": 20,
					"largePage":  50,
				}),
				parameter.WithDefault(20),
			),
			parameter.StrParam("cursor", parameter.Query,
				parameter.WithDescription("Pagination cursor"),
				parameter.WithExample("abc123xyz456"),
				parameter.WithExamples(map[string]interface{}{
					"firstPage":  "cursor-001",
					"secondPage": "cursor-002",
					"lastPage":   "cursor-end",
				}),
			),
			parameter.StrParam("role", parameter.Query,
				parameter.WithDescription("Filter by user role"),
				parameter.WithExample("admin"),
				parameter.WithExamples(map[string]interface{}{
					"adminRole": "admin",
					"userRole":  "user",
					"guestRole": "guest",
					"allRoles":  "admin,user,guest",
				}),
			),
		),
		endpoint.WithSuccessfulReturns([]response.Response{
			response.New([]User{}, "200", "Successfully retrieved users").
				WithExample([]User{
					{
						ID:        1,
						Username:  "john_doe",
						Email:     "john@example.com",
						FirstName: "John",
						LastName:  "Doe",
					},
				}).
				WithExamples(map[string]interface{}{
					"emptyResult": []User{},
					"singleUser": []User{
						{
							ID:        1,
							Username:  "john_doe",
							Email:     "john@example.com",
							FirstName: "John",
							LastName:  "Doe",
						},
					},
					"multipleUsers": []User{
						{
							ID:        1,
							Username:  "john_doe",
							Email:     "john@example.com",
							FirstName: "John",
							LastName:  "Doe",
						},
						{
							ID:        2,
							Username:  "jane_smith",
							Email:     "jane@example.com",
							FirstName: "Jane",
							LastName:  "Smith",
						},
					},
				}),
		}),
		endpoint.WithErrors([]response.Response{
			response.New(ErrorResponse{}, "400", "Bad Request").
				WithExample(ErrorResponse{
					Error:   true,
					Message: "Invalid query parameters",
				}).
				WithExamples(map[string]interface{}{
					"invalidLimit": ErrorResponse{
						Error:   true,
						Message: "Limit must be between 1 and 100",
					},
					"invalidCursor": ErrorResponse{
						Error:   true,
						Message: "Invalid cursor format",
					},
					"invalidRole": ErrorResponse{
						Error:   true,
						Message: "Role must be one of: admin, user, guest",
					},
				}),
			response.New(ErrorResponse{}, "401", "Unauthorized").
				WithExample(ErrorResponse{
					Error:   true,
					Message: "Authentication required",
				}).
				WithExamples(map[string]interface{}{
					"missingToken": ErrorResponse{
						Error:   true,
						Message: "Authorization token is missing",
					},
					"invalidToken": ErrorResponse{
						Error:   true,
						Message: "Invalid authorization token",
					},
					"expiredToken": ErrorResponse{
						Error:   true,
						Message: "Authorization token has expired",
					},
				}),
		}),
	)

	// POST /users - Create a new user with multiple request/response examples
	createUserEndpoint := endpoint.New(
		endpoint.POST,
		"/users",
		endpoint.WithTags("users"),
		endpoint.WithSummary("Create a new user"),
		endpoint.WithDescription("Creates a new user in the system"),
		endpoint.WithBody(CreateUserRequest{},
			endpoint.WithBodyDescription("User data for creation"),
			endpoint.WithBodyRequired(true),
			endpoint.WithBodyExample(CreateUserRequest{
				Username:  "new_user_example",
				Email:     "new@example.com",
				FirstName: "New",
				LastName:  "User",
				Password:  "SecurePassword123!",
			}),
			endpoint.WithBodyExamples(map[string]interface{}{
				"minimalUser": CreateUserRequest{
					Username:  "minimal",
					Email:     "minimal@example.com",
					FirstName: "Min",
					LastName:  "User",
					Password:  "Password123!",
				},
				"completeUser": CreateUserRequest{
					Username:  "complete_user",
					Email:     "complete@example.com",
					FirstName: "Complete",
					LastName:  "User",
					Password:  "VerySecurePassword123!",
				},
				"adminUser": CreateUserRequest{
					Username:  "admin_user",
					Email:     "admin@example.com",
					FirstName: "Admin",
					LastName:  "User",
					Password:  "AdminSecurePassword123!",
				},
			}),
		),
		endpoint.WithSuccessfulReturns([]response.Response{
			response.New(User{}, "201", "User created successfully").
				WithExample(User{
					ID:        12345,
					Username:  "new_user_example",
					Email:     "new@example.com",
					FirstName: "New",
					LastName:  "User",
				}).
				WithExamples(map[string]interface{}{
					"regularUser": User{
						ID:        10001,
						Username:  "regular_user",
						Email:     "regular@example.com",
						FirstName: "Regular",
						LastName:  "User",
					},
					"premiumUser": User{
						ID:        20001,
						Username:  "premium_user",
						Email:     "premium@example.com",
						FirstName: "Premium",
						LastName:  "User",
					},
					"adminUser": User{
						ID:        30001,
						Username:  "admin_user",
						Email:     "admin@example.com",
						FirstName: "Admin",
						LastName:  "User",
					},
				}),
		}),
		endpoint.WithErrors([]response.Response{
			response.New(ErrorResponse{}, "400", "Bad Request").
				WithExample(ErrorResponse{
					Error:   true,
					Message: "Invalid user data",
				}).
				WithExamples(map[string]interface{}{
					"missingUsername": ErrorResponse{
						Error:   true,
						Message: "Username is required",
					},
					"invalidEmail": ErrorResponse{
						Error:   true,
						Message: "Email must be a valid email address",
					},
					"weakPassword": ErrorResponse{
						Error:   true,
						Message: "Password must be at least 8 characters and contain uppercase, lowercase, and special characters",
					},
					"usernameTaken": ErrorResponse{
						Error:   true,
						Message: "Username is already taken",
					},
					"emailTaken": ErrorResponse{
						Error:   true,
						Message: "Email is already registered",
					},
				}),
			response.New(ErrorResponse{}, "409", "Conflict").
				WithExample(ErrorResponse{
					Error:   true,
					Message: "Username or email already exists",
				}).
				WithExamples(map[string]interface{}{
					"usernameConflict": ErrorResponse{
						Error:   true,
						Message: "Username 'admin' is already taken",
					},
					"emailConflict": ErrorResponse{
						Error:   true,
						Message: "Email 'admin@example.com' is already registered",
					},
					"bothConflict": ErrorResponse{
						Error:   true,
						Message: "Both username and email are already in use",
					},
				}),
		}),
	)

	// GET /users/{id} - Get a specific user with multiple response examples
	getUserEndpoint := endpoint.New(
		endpoint.GET,
		"/users/{id}",
		endpoint.WithTags("users"),
		endpoint.WithSummary("Get a specific user"),
		endpoint.WithDescription("Returns detailed information about a specific user"),
		endpoint.WithParams(
			parameter.IntParam("id", parameter.Path,
				parameter.WithDescription("User ID"),
				parameter.WithExample(12345),
				parameter.WithExamples(map[string]interface{}{
					"regularUserId": 10001,
					"premiumUserId": 20001,
					"adminUserId":   30001,
					"nonExistentId": 99999,
				}),
				parameter.WithRequired(),
			),
		),
		endpoint.WithSuccessfulReturns([]response.Response{
			response.New(User{}, "200", "Successfully retrieved user").
				WithExample(User{
					ID:        12345,
					Username:  "john_doe",
					Email:     "john@example.com",
					FirstName: "John",
					LastName:  "Doe",
				}).
				WithExamples(map[string]interface{}{
					"regularUser": User{
						ID:        10001,
						Username:  "regular_user",
						Email:     "regular@example.com",
						FirstName: "Regular",
						LastName:  "User",
					},
					"premiumUser": User{
						ID:        20001,
						Username:  "premium_user",
						Email:     "premium@example.com",
						FirstName: "Premium",
						LastName:  "User",
					},
					"adminUser": User{
						ID:        30001,
						Username:  "admin_user",
						Email:     "admin@example.com",
						FirstName: "Admin",
						LastName:  "User",
					},
				}),
		}),
		endpoint.WithErrors([]response.Response{
			response.New(ErrorResponse{}, "404", "Not Found").
				WithExample(ErrorResponse{
					Error:   true,
					Message: "User not found",
				}).
				WithExamples(map[string]interface{}{
					"deletedUser": ErrorResponse{
						Error:   true,
						Message: "User with ID 99999 was not found (may have been deleted)",
					},
					"invalidId": ErrorResponse{
						Error:   true,
						Message: "Invalid user ID format",
					},
					"neverExisted": ErrorResponse{
						Error:   true,
						Message: "No user has ever had ID 99999",
					},
				}),
		}),
	)

	// Add endpoints to OpenAPI
	openapi.AddEndpoints([]*endpoint.EndPoint{listUsersEndpoint, createUserEndpoint, getUserEndpoint})

	// Generate OpenAPI JSON
	jsonData := openapi.MustToJson()

	// Print the JSON
	fmt.Println("Generated OpenAPI specification with multiple examples:")
	fmt.Println("======================================================")

	var prettyJSON map[string]interface{}
	if err := json.Unmarshal(jsonData, &prettyJSON); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	prettyBytes, err := json.MarshalIndent(prettyJSON, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal pretty JSON: %v", err)
	}

	fmt.Println(string(prettyBytes))

	// Export to file
	openapi.ExportOpenAPIDocs("user_api_with_multiple_examples.json")
	fmt.Println("\nOpenAPI specification exported to: user_api_with_multiple_examples.json")
	fmt.Println("This file demonstrates comprehensive multiple examples support for:")
	fmt.Println("- Request body examples (single + multiple)")
	fmt.Println("- Response examples (single + multiple)")
	fmt.Println("- Parameter examples (single + multiple)")
	fmt.Println("- Error response examples (single + multiple)")
}
