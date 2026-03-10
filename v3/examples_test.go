package swagno3

import (
	"testing"

	"github.com/go-swagno/swagno/v3/components/definition"
	"github.com/go-swagno/swagno/v3/components/endpoint"
	"github.com/go-swagno/swagno/v3/components/http/response"
	"github.com/go-swagno/swagno/v3/components/mime"
	"github.com/go-swagno/swagno/v3/components/parameter"
	"github.com/go-swagno/swagno/v3/components/security"
	"github.com/go-swagno/swagno/v3/components/tag"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type TestRequest struct {
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

type TestResponse struct {
	ID    uint64 `json:"id" example:"12345"`
	Name  string `json:"name" example:"Test Item"`
	Email string `json:"email" example:"test@example.com"`
}

func TestExamplesCombined(t *testing.T) {
	// Common options for cmp.Diff
	cmpOpts := []cmp.Option{
		cmpopts.IgnoreUnexported(OpenAPI{}, endpoint.EndPoint{}, endpoint.JsonEndPoint{}),
		cmpopts.EquateEmpty(),
		cmpopts.IgnoreFields(definition.SchemaProperty{}, "IsRequired"),
		cmpopts.SortSlices(func(a, b string) bool { return a < b }),
	}

	t.Run("BasicExamples", func(t *testing.T) {
		// Create a test endpoint with examples
		testEndpoint := endpoint.New(
			endpoint.POST,
			"/test",
			endpoint.WithBody(TestRequest{}, endpoint.WithBodyExample(map[string]interface{}{
				"name":  "Jane Smith",
				"email": "jane@example.com",
			})),
			endpoint.WithSuccessfulReturns([]response.Response{
				response.New(TestResponse{}, "200", "Success").WithExample(map[string]interface{}{
					"id":    uint64(67890),
					"name":  "Example Response",
					"email": "response@example.com",
				}),
			}),
			endpoint.WithParams(
				parameter.StrParam("api_key", parameter.Query,
					parameter.WithDescription("API Key"),
					parameter.WithExample("abc123-xyz456"),
					parameter.WithRequired(),
				),
			),
		)

		// Create OpenAPI instance
		openapi := New(Config{
			Title:       "Test API",
			Version:     "1.0.0",
			Summary:     "Test API with examples",
			Description: "This is a test API to verify example functionality",
		})

		openapi.AddEndpoint(testEndpoint)
		openapi.generateOpenAPIJson()

		expected := &OpenAPI{
			OpenAPI: "3.0.3",
			Info: Info{
				Title:       "Test API",
				Version:     "1.0.0",
				Summary:     "Test API with examples",
				Description: "This is a test API to verify example functionality",
			},
			Servers: []Server{{URL: "/"}},
			Paths: map[string]endpoint.PathItem{
				"/test": {
					Post: &endpoint.JsonEndPoint{
						OperationId: "post-_test",
						Consume:     []mime.MIME{mime.JSON},
						Produce:     []mime.MIME{mime.JSON},
						Parameters: []parameter.JsonParameter{
							{
								Name:        "api_key",
								In:          "query",
								Description: "API Key",
								Required:    true,
								Schema: &parameter.JsonResponseSchema{
									Type:    "string",
									Example: "abc123-xyz456",
								},
								Example: "abc123-xyz456",
							},
						},
						RequestBody: &endpoint.RequestBody{
							Description: "Request body",
							Required:    true,
							Content: map[string]endpoint.MediaType{
								"application/json": {
									Schema: &parameter.JsonResponseSchema{
										Ref: "#/components/schemas/swagno3.TestRequest",
									},
									Example: map[string]interface{}{
										"name":  "Jane Smith",
										"email": "jane@example.com",
									},
								},
							},
						},
						Responses: map[string]endpoint.JsonResponse{
							"200": {
								Description: "Success",
								Content: map[string]endpoint.MediaType{
									"application/json": {
										Schema: &parameter.JsonResponseSchema{
											Ref: "#/components/schemas/swagno3.TestResponse",
										},
										Example: map[string]interface{}{
											"id":    uint64(67890),
											"name":  "Example Response",
											"email": "response@example.com",
										},
									},
								},
							},
						},
					},
				},
			},
			Components: &Components{
				Schemas: map[string]definition.Schema{
					"swagno3.TestRequest": {
						Type: "object",
						Properties: map[string]definition.SchemaProperty{
							"email": {
								Type:    "string",
								Example: "john@example.com",
							},
							"name": {
								Type:    "string",
								Example: "John Doe",
							},
						},
						Required: []string{"name", "email"},
					},
					"swagno3.TestResponse": {
						Type: "object",
						Properties: map[string]definition.SchemaProperty{
							"email": {
								Type:    "string",
								Example: "test@example.com",
							},
							"id": {
								Type:    "integer",
								Example: float64(12345),
							},
							"name": {
								Type:    "string",
								Example: "Test Item",
							},
						},
						Required: []string{"id", "name", "email"},
					},
				},
				SecuritySchemes: make(map[security.SecuritySchemeName]SecurityScheme),
			},
			Tags: []tag.Tag{},
		}

		if diff := cmp.Diff(expected, openapi, cmpOpts...); diff != "" {
			t.Errorf("OpenAPI mismatch (-expected +got):\n%s", diff)
		}
	})

	t.Run("MultipleExamples", func(t *testing.T) {
		// Create a test endpoint with multiple examples
		testEndpoint := endpoint.New(
			endpoint.POST,
			"/test",
			endpoint.WithBody(TestRequest{},
				endpoint.WithBodyExample(map[string]interface{}{
					"name":  "Single Example",
					"email": "single@example.com",
				}),
				endpoint.WithBodyExamples(map[string]interface{}{
					"example1": map[string]interface{}{
						"name":  "First Multiple Example",
						"email": "first@example.com",
					},
					"example2": map[string]interface{}{
						"name":  "Second Multiple Example",
						"email": "second@example.com",
					},
				}),
			),
			endpoint.WithSuccessfulReturns([]response.Response{
				response.New(TestResponse{}, "200", "Success").
					WithExample(map[string]interface{}{
						"id":    uint64(12345),
						"name":  "Single Response Example",
						"email": "single-response@example.com",
					}).
					WithExamples(map[string]interface{}{
						"successExample1": map[string]interface{}{
							"id":    uint64(67890),
							"name":  "First Success Example",
							"email": "first-success@example.com",
						},
						"successExample2": map[string]interface{}{
							"id":    uint64(54321),
							"name":  "Second Success Example",
							"email": "second-success@example.com",
						},
					}),
			}),
			endpoint.WithParams(
				parameter.StrParam("api_key", parameter.Query,
					parameter.WithDescription("API Key"),
					parameter.WithExample("single-api-key-example"),
					parameter.WithExamples(map[string]interface{}{
						"userKey":    "user-api-key-123",
						"adminKey":   "admin-api-key-456",
						"serviceKey": "service-api-key-789",
					}),
					parameter.WithRequired(),
				),
			),
		)

		// Create OpenAPI instance
		openapi := New(Config{Title: "Multiple Examples Test API", Version: "1.0.0"})

		openapi.AddEndpoint(testEndpoint)
		openapi.generateOpenAPIJson()

		expected := &OpenAPI{
			OpenAPI: "3.0.3",
			Info: Info{
				Title:   "Multiple Examples Test API",
				Version: "1.0.0",
			},
			Servers: []Server{{URL: "/"}},
			Paths: map[string]endpoint.PathItem{
				"/test": {
					Post: &endpoint.JsonEndPoint{
						OperationId: "post-_test",
						Consume:     []mime.MIME{mime.JSON},
						Produce:     []mime.MIME{mime.JSON},
						Parameters: []parameter.JsonParameter{
							{
								Name:        "api_key",
								In:          "query",
								Description: "API Key",
								Required:    true,
								Schema: &parameter.JsonResponseSchema{
									Type:    "string",
									Example: nil,
								},
								Example: nil,
								Examples: map[string]parameter.ComponentExample{
									"userKey":    {Value: "user-api-key-123"},
									"adminKey":   {Value: "admin-api-key-456"},
									"serviceKey": {Value: "service-api-key-789"},
								},
							},
						},
						RequestBody: &endpoint.RequestBody{
							Description: "Request body",
							Required:    true,
							Content: map[string]endpoint.MediaType{
								"application/json": {
									Schema: &parameter.JsonResponseSchema{
										Ref: "#/components/schemas/swagno3.TestRequest",
									},
									Example: nil,
									Examples: map[string]parameter.ComponentExample{
										"example1": {
											Value: map[string]interface{}{
												"name":  "First Multiple Example",
												"email": "first@example.com",
											},
										},
										"example2": {
											Value: map[string]interface{}{
												"name":  "Second Multiple Example",
												"email": "second@example.com",
											},
										},
									},
								},
							},
						},
						Responses: map[string]endpoint.JsonResponse{
							"200": {
								Description: "Success",
								Content: map[string]endpoint.MediaType{
									"application/json": {
										Schema: &parameter.JsonResponseSchema{
											Ref: "#/components/schemas/swagno3.TestResponse",
										},
										Example: nil,
										Examples: map[string]parameter.ComponentExample{
											"successExample1": {
												Value: map[string]interface{}{
													"id":    uint64(67890),
													"name":  "First Success Example",
													"email": "first-success@example.com",
												},
											},
											"successExample2": {
												Value: map[string]interface{}{
													"id":    uint64(54321),
													"name":  "Second Success Example",
													"email": "second-success@example.com",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Components: &Components{
				Schemas: map[string]definition.Schema{
					"swagno3.TestRequest": {
						Type: "object",
						Properties: map[string]definition.SchemaProperty{
							"email": {
								Type:    "string",
								Example: "john@example.com",
							},
							"name": {
								Type:    "string",
								Example: "John Doe",
							},
						},
						Required: []string{"email", "name"},
					},
					"swagno3.TestResponse": {
						Type: "object",
						Properties: map[string]definition.SchemaProperty{
							"email": {
								Type:    "string",
								Example: "test@example.com",
							},
							"id": {
								Type:    "integer",
								Example: float64(12345),
							},
							"name": {
								Type:    "string",
								Example: "Test Item",
							},
						},
						Required: []string{"id", "name", "email"},
					},
				},
				SecuritySchemes: make(map[security.SecuritySchemeName]SecurityScheme),
			},
			Tags: []tag.Tag{},
		}

		if diff := cmp.Diff(expected, openapi, cmpOpts...); diff != "" {
			t.Errorf("OpenAPI mismatch (-expected +got):\n%s", diff)
		}
	})
}
