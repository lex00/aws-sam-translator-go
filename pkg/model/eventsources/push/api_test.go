package push

import (
	"testing"
)

func TestApi_Validate(t *testing.T) {
	tests := []struct {
		name    string
		api     *Api
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid api event with all required fields",
			api: &Api{
				Path:   "/test",
				Method: "GET",
			},
			wantErr: false,
		},
		{
			name: "valid api event with rest api id",
			api: &Api{
				Path:      "/users",
				Method:    "POST",
				RestApiId: map[string]interface{}{"Ref": "MyApi"},
			},
			wantErr: false,
		},
		{
			name: "missing path",
			api: &Api{
				Method: "GET",
			},
			wantErr: true,
			errMsg:  "Api event source requires Path property",
		},
		{
			name: "missing method",
			api: &Api{
				Path: "/test",
			},
			wantErr: true,
			errMsg:  "Api event source requires Method property",
		},
		{
			name:    "missing both path and method",
			api:     &Api{},
			wantErr: true,
			errMsg:  "Api event source requires Path property",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.api.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Api.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Api.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestApi_EventType(t *testing.T) {
	api := &Api{
		Path:   "/test",
		Method: "GET",
	}

	eventType := api.EventType()
	expected := "Api"

	if eventType != expected {
		t.Errorf("Api.EventType() = %v, want %v", eventType, expected)
	}
}

func TestApi_ToCloudFormation(t *testing.T) {
	tests := []struct {
		name              string
		api               *Api
		functionLogicalId string
		functionRef       interface{}
		stageName         interface{}
		wantErr           bool
		checkResources    func(t *testing.T, resources map[string]interface{})
	}{
		{
			name: "basic api event",
			api: &Api{
				Path:   "/",
				Method: "GET",
			},
			functionLogicalId: "MyFunction",
			functionRef:       map[string]interface{}{"Ref": "MyFunction"},
			stageName:         "Prod",
			wantErr:           false,
			checkResources: func(t *testing.T, resources map[string]interface{}) {
				if len(resources) == 0 {
					t.Error("Expected resources to be generated")
				}
				// Check that we have at least a method and permission
				hasMethod := false
				hasPermission := false
				for key, resource := range resources {
					if resMap, ok := resource.(map[string]interface{}); ok {
						if resType, ok := resMap["Type"].(string); ok {
							if resType == "AWS::ApiGateway::Method" {
								hasMethod = true
							}
							if resType == "AWS::Lambda::Permission" {
								hasPermission = true
							}
						}
					}
					t.Logf("Generated resource: %s", key)
				}
				if !hasMethod {
					t.Error("Expected AWS::ApiGateway::Method resource")
				}
				if !hasPermission {
					t.Error("Expected AWS::Lambda::Permission resource")
				}
			},
		},
		{
			name: "api event with explicit rest api id",
			api: &Api{
				Path:      "/users",
				Method:    "POST",
				RestApiId: map[string]interface{}{"Ref": "MyApi"},
			},
			functionLogicalId: "UserFunction",
			functionRef:       map[string]interface{}{"Ref": "UserFunction"},
			stageName:         "Stage",
			wantErr:           false,
			checkResources: func(t *testing.T, resources map[string]interface{}) {
				if len(resources) == 0 {
					t.Error("Expected resources to be generated")
				}
			},
		},
		{
			name: "api event with auth",
			api: &Api{
				Path:   "/secure",
				Method: "GET",
				Auth: &ApiAuth{
					Authorizer:     "MyAuthorizer",
					ApiKeyRequired: true,
				},
			},
			functionLogicalId: "SecureFunction",
			functionRef:       map[string]interface{}{"Ref": "SecureFunction"},
			stageName:         "Prod",
			wantErr:           false,
			checkResources: func(t *testing.T, resources map[string]interface{}) {
				if len(resources) == 0 {
					t.Error("Expected resources to be generated")
				}
			},
		},
		{
			name: "missing path should error",
			api: &Api{
				Method: "GET",
			},
			functionLogicalId: "MyFunction",
			functionRef:       map[string]interface{}{"Ref": "MyFunction"},
			stageName:         "Prod",
			wantErr:           true,
		},
		{
			name: "missing method should error",
			api: &Api{
				Path: "/test",
			},
			functionLogicalId: "MyFunction",
			functionRef:       map[string]interface{}{"Ref": "MyFunction"},
			stageName:         "Prod",
			wantErr:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resources, err := tt.api.ToCloudFormation(tt.functionLogicalId, tt.functionRef, tt.stageName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Api.ToCloudFormation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.checkResources != nil {
				tt.checkResources(t, resources)
			}
		})
	}
}

func TestApi_getAuthorizationType(t *testing.T) {
	tests := []struct {
		name string
		api  *Api
		want interface{}
	}{
		{
			name: "no auth returns NONE",
			api: &Api{
				Path:   "/",
				Method: "GET",
			},
			want: "NONE",
		},
		{
			name: "with authorizer returns CUSTOM",
			api: &Api{
				Path:   "/",
				Method: "GET",
				Auth: &ApiAuth{
					Authorizer: "MyAuthorizer",
				},
			},
			want: "CUSTOM",
		},
		{
			name: "with invoke role returns AWS_IAM",
			api: &Api{
				Path:   "/",
				Method: "GET",
				Auth: &ApiAuth{
					InvokeRole: "arn:aws:iam::123456789012:role/MyRole",
				},
			},
			want: "AWS_IAM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.api.getAuthorizationType()
			if got != tt.want {
				t.Errorf("Api.getAuthorizationType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApi_buildSourceArn(t *testing.T) {
	tests := []struct {
		name      string
		api       *Api
		restApiId interface{}
		stageName interface{}
		wantType  string
	}{
		{
			name: "builds source arn with ref",
			api: &Api{
				Path:   "/users",
				Method: "GET",
			},
			restApiId: map[string]interface{}{"Ref": "MyApi"},
			stageName: "Prod",
			wantType:  "map",
		},
		{
			name: "builds source arn with string api id",
			api: &Api{
				Path:   "/",
				Method: "POST",
			},
			restApiId: "MyApiId",
			stageName: "Stage",
			wantType:  "map",
		},
		{
			name: "builds source arn with nil stage (uses wildcard)",
			api: &Api{
				Path:   "/test",
				Method: "PUT",
			},
			restApiId: map[string]interface{}{"Ref": "Api"},
			stageName: nil,
			wantType:  "map",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arn := tt.api.buildSourceArn(tt.restApiId, tt.stageName)
			if arn == nil {
				t.Error("Api.buildSourceArn() returned nil")
				return
			}
			// Check that it returns a map (Fn::Sub structure)
			if _, ok := arn.(map[string]interface{}); !ok && tt.wantType == "map" {
				t.Errorf("Api.buildSourceArn() type = %T, want map[string]interface{}", arn)
			}
		})
	}
}

func TestApi_getResourceId(t *testing.T) {
	tests := []struct {
		name      string
		api       *Api
		restApiId interface{}
		path      interface{}
		wantType  string
	}{
		{
			name:      "root path returns root resource id",
			api:       &Api{},
			restApiId: map[string]interface{}{"Ref": "MyApi"},
			path:      "/",
			wantType:  "GetAtt",
		},
		{
			name:      "non-root path returns resource reference",
			api:       &Api{},
			restApiId: map[string]interface{}{"Ref": "MyApi"},
			path:      "/users",
			wantType:  "Ref",
		},
		{
			name:      "non-string path treated as root",
			api:       &Api{},
			restApiId: map[string]interface{}{"Ref": "MyApi"},
			path:      map[string]interface{}{"Ref": "PathParam"},
			wantType:  "GetAtt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceId := tt.api.getResourceId(tt.restApiId, tt.path)
			if resourceId == nil {
				t.Error("Api.getResourceId() returned nil")
				return
			}

			// Check the structure
			if resMap, ok := resourceId.(map[string]interface{}); ok {
				if tt.wantType == "GetAtt" {
					if _, hasGetAtt := resMap["Fn::GetAtt"]; !hasGetAtt {
						t.Errorf("Api.getResourceId() expected Fn::GetAtt, got %v", resMap)
					}
				} else if tt.wantType == "Ref" {
					if _, hasRef := resMap["Ref"]; !hasRef {
						t.Errorf("Api.getResourceId() expected Ref, got %v", resMap)
					}
				}
			}
		})
	}
}
