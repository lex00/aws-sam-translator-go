// Package parser provides YAML/JSON template parsing with intrinsic function detection.
package parser

import (
	"testing"
)

func TestParseYAMLWithRef(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantKey string
		wantVal string
	}{
		{
			name: "short form !Ref",
			yaml: `
Resources:
  MyFunc:
    Type: AWS::Serverless::Function
    Properties:
      Runtime: !Ref RuntimeParam
`,
			wantKey: "Ref",
			wantVal: "RuntimeParam",
		},
		{
			name: "long form Ref",
			yaml: `
Resources:
  MyFunc:
    Type: AWS::Serverless::Function
    Properties:
      Runtime:
        Ref: RuntimeParam
`,
			wantKey: "Ref",
			wantVal: "RuntimeParam",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			result, err := p.ParseRawYAML([]byte(tt.yaml))
			if err != nil {
				t.Fatalf("ParseRawYAML failed: %v", err)
			}

			resources := result["Resources"].(map[string]interface{})
			myFunc := resources["MyFunc"].(map[string]interface{})
			props := myFunc["Properties"].(map[string]interface{})
			runtime := props["Runtime"].(map[string]interface{})

			if _, ok := runtime[tt.wantKey]; !ok {
				t.Errorf("expected key %q not found in runtime: %v", tt.wantKey, runtime)
			}
			if runtime[tt.wantKey] != tt.wantVal {
				t.Errorf("expected value %q, got %v", tt.wantVal, runtime[tt.wantKey])
			}
		})
	}
}

func TestParseYAMLWithSub(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantKey string
	}{
		{
			name: "short form !Sub",
			yaml: `
Resources:
  MyFunc:
    Type: AWS::Serverless::Function
    Properties:
      Handler: !Sub "${Stage}-handler"
`,
			wantKey: "Fn::Sub",
		},
		{
			name: "long form Fn::Sub",
			yaml: `
Resources:
  MyFunc:
    Type: AWS::Serverless::Function
    Properties:
      Handler:
        Fn::Sub: "${Stage}-handler"
`,
			wantKey: "Fn::Sub",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			result, err := p.ParseRawYAML([]byte(tt.yaml))
			if err != nil {
				t.Fatalf("ParseRawYAML failed: %v", err)
			}

			resources := result["Resources"].(map[string]interface{})
			myFunc := resources["MyFunc"].(map[string]interface{})
			props := myFunc["Properties"].(map[string]interface{})
			handler := props["Handler"].(map[string]interface{})

			if _, ok := handler[tt.wantKey]; !ok {
				t.Errorf("expected key %q not found in handler: %v", tt.wantKey, handler)
			}
		})
	}
}

func TestParseYAMLWithGetAtt(t *testing.T) {
	tests := []struct {
		name      string
		yaml      string
		wantKey   string
		wantArray bool
		wantLen   int
	}{
		{
			name: "short form !GetAtt with dot notation",
			yaml: `
Resources:
  MyFunc:
    Type: AWS::Serverless::Function
    Properties:
      Role: !GetAtt MyRole.Arn
`,
			wantKey:   "Fn::GetAtt",
			wantArray: true,
			wantLen:   2,
		},
		{
			name: "long form Fn::GetAtt with array",
			yaml: `
Resources:
  MyFunc:
    Type: AWS::Serverless::Function
    Properties:
      Role:
        Fn::GetAtt:
          - MyRole
          - Arn
`,
			wantKey:   "Fn::GetAtt",
			wantArray: true,
			wantLen:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			result, err := p.ParseRawYAML([]byte(tt.yaml))
			if err != nil {
				t.Fatalf("ParseRawYAML failed: %v", err)
			}

			resources := result["Resources"].(map[string]interface{})
			myFunc := resources["MyFunc"].(map[string]interface{})
			props := myFunc["Properties"].(map[string]interface{})
			role := props["Role"].(map[string]interface{})

			if _, ok := role[tt.wantKey]; !ok {
				t.Errorf("expected key %q not found in role: %v", tt.wantKey, role)
			}

			if tt.wantArray {
				arr, ok := role[tt.wantKey].([]interface{})
				if !ok {
					// Could also be []string from the dot notation conversion
					arrStr, ok := role[tt.wantKey].([]string)
					if !ok {
						t.Errorf("expected array for %s, got %T", tt.wantKey, role[tt.wantKey])
						return
					}
					if len(arrStr) != tt.wantLen {
						t.Errorf("expected array length %d, got %d", tt.wantLen, len(arrStr))
					}
					return
				}
				if len(arr) != tt.wantLen {
					t.Errorf("expected array length %d, got %d", tt.wantLen, len(arr))
				}
			}
		})
	}
}

func TestParseYAMLWithJoin(t *testing.T) {
	yaml := `
Resources:
  MyFunc:
    Type: AWS::Serverless::Function
    Properties:
      Environment:
        Variables:
          PATH: !Join
            - ":"
            - - /usr/bin
              - /opt/bin
`
	p := New()
	result, err := p.ParseRawYAML([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseRawYAML failed: %v", err)
	}

	resources := result["Resources"].(map[string]interface{})
	myFunc := resources["MyFunc"].(map[string]interface{})
	props := myFunc["Properties"].(map[string]interface{})
	env := props["Environment"].(map[string]interface{})
	vars := env["Variables"].(map[string]interface{})
	path := vars["PATH"].(map[string]interface{})

	if _, ok := path["Fn::Join"]; !ok {
		t.Errorf("expected Fn::Join not found: %v", path)
	}
}

func TestParseYAMLWithIf(t *testing.T) {
	yaml := `
Resources:
  MyFunc:
    Type: AWS::Serverless::Function
    Properties:
      Timeout: !If
        - IsProd
        - 30
        - 10
`
	p := New()
	result, err := p.ParseRawYAML([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseRawYAML failed: %v", err)
	}

	resources := result["Resources"].(map[string]interface{})
	myFunc := resources["MyFunc"].(map[string]interface{})
	props := myFunc["Properties"].(map[string]interface{})
	timeout := props["Timeout"].(map[string]interface{})

	if _, ok := timeout["Fn::If"]; !ok {
		t.Errorf("expected Fn::If not found: %v", timeout)
	}

	arr := timeout["Fn::If"].([]interface{})
	if len(arr) != 3 {
		t.Errorf("expected Fn::If to have 3 elements, got %d", len(arr))
	}
}

func TestParseYAMLWithSelect(t *testing.T) {
	yaml := `
Resources:
  MyFunc:
    Type: AWS::Serverless::Function
    Properties:
      Zone: !Select
        - 0
        - !GetAZs ""
`
	p := New()
	result, err := p.ParseRawYAML([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseRawYAML failed: %v", err)
	}

	resources := result["Resources"].(map[string]interface{})
	myFunc := resources["MyFunc"].(map[string]interface{})
	props := myFunc["Properties"].(map[string]interface{})
	zone := props["Zone"].(map[string]interface{})

	if _, ok := zone["Fn::Select"]; !ok {
		t.Errorf("expected Fn::Select not found: %v", zone)
	}
}

func TestParseYAMLWithFindInMap(t *testing.T) {
	yaml := `
Mappings:
  RegionMap:
    us-east-1:
      AMI: ami-12345
Resources:
  MyFunc:
    Type: AWS::Serverless::Function
    Properties:
      ImageId: !FindInMap
        - RegionMap
        - !Ref "AWS::Region"
        - AMI
`
	p := New()
	result, err := p.ParseRawYAML([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseRawYAML failed: %v", err)
	}

	resources := result["Resources"].(map[string]interface{})
	myFunc := resources["MyFunc"].(map[string]interface{})
	props := myFunc["Properties"].(map[string]interface{})
	imageId := props["ImageId"].(map[string]interface{})

	if _, ok := imageId["Fn::FindInMap"]; !ok {
		t.Errorf("expected Fn::FindInMap not found: %v", imageId)
	}

	arr := imageId["Fn::FindInMap"].([]interface{})
	if len(arr) != 3 {
		t.Errorf("expected Fn::FindInMap to have 3 elements, got %d", len(arr))
	}
}

func TestParseYAMLWithBase64(t *testing.T) {
	yaml := `
Resources:
  MyFunc:
    Type: AWS::Serverless::Function
    Properties:
      UserData: !Base64 "echo hello"
`
	p := New()
	result, err := p.ParseRawYAML([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseRawYAML failed: %v", err)
	}

	resources := result["Resources"].(map[string]interface{})
	myFunc := resources["MyFunc"].(map[string]interface{})
	props := myFunc["Properties"].(map[string]interface{})
	userData := props["UserData"].(map[string]interface{})

	if _, ok := userData["Fn::Base64"]; !ok {
		t.Errorf("expected Fn::Base64 not found: %v", userData)
	}
}

func TestParseYAMLWithNestedIntrinsics(t *testing.T) {
	yaml := `
Resources:
  MyFunc:
    Type: AWS::Serverless::Function
    Properties:
      Environment:
        Variables:
          URL: !Sub
            - "https://${Domain}/api"
            - Domain: !Ref DomainName
`
	p := New()
	result, err := p.ParseRawYAML([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseRawYAML failed: %v", err)
	}

	resources := result["Resources"].(map[string]interface{})
	myFunc := resources["MyFunc"].(map[string]interface{})
	props := myFunc["Properties"].(map[string]interface{})
	env := props["Environment"].(map[string]interface{})
	vars := env["Variables"].(map[string]interface{})
	url := vars["URL"].(map[string]interface{})

	subValue, ok := url["Fn::Sub"]
	if !ok {
		t.Fatalf("expected Fn::Sub not found: %v", url)
	}

	// Should be an array with template and variable map
	arr, ok := subValue.([]interface{})
	if !ok {
		t.Fatalf("expected Fn::Sub value to be array, got %T", subValue)
	}

	if len(arr) != 2 {
		t.Fatalf("expected Fn::Sub array to have 2 elements, got %d", len(arr))
	}

	// Check nested !Ref inside the variable map
	varMap, ok := arr[1].(map[string]interface{})
	if !ok {
		t.Fatalf("expected second element to be map, got %T", arr[1])
	}

	domain, ok := varMap["Domain"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected Domain to be map, got %T", varMap["Domain"])
	}

	if _, ok := domain["Ref"]; !ok {
		t.Errorf("expected nested Ref not found: %v", domain)
	}
}

func TestParseJSON(t *testing.T) {
	jsonData := `{
		"Resources": {
			"MyFunc": {
				"Type": "AWS::Serverless::Function",
				"Properties": {
					"Runtime": {"Ref": "RuntimeParam"},
					"Role": {"Fn::GetAtt": ["MyRole", "Arn"]}
				}
			}
		}
	}`

	p := New()
	result, err := p.ParseRawJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseRawJSON failed: %v", err)
	}

	resources := result["Resources"].(map[string]interface{})
	myFunc := resources["MyFunc"].(map[string]interface{})
	props := myFunc["Properties"].(map[string]interface{})

	// Check Ref
	runtime := props["Runtime"].(map[string]interface{})
	if _, ok := runtime["Ref"]; !ok {
		t.Errorf("expected Ref not found in runtime: %v", runtime)
	}

	// Check Fn::GetAtt
	role := props["Role"].(map[string]interface{})
	if _, ok := role["Fn::GetAtt"]; !ok {
		t.Errorf("expected Fn::GetAtt not found in role: %v", role)
	}
}

func TestParseYAMLToTemplate(t *testing.T) {
	yaml := `
AWSTemplateFormatVersion: "2010-09-09"
Transform: AWS::Serverless-2016-10-31
Description: Test template
Parameters:
  Stage:
    Type: String
    Default: dev
Resources:
  MyFunc:
    Type: AWS::Serverless::Function
    Properties:
      Runtime: python3.9
      Handler: index.handler
      CodeUri: ./src
Outputs:
  FuncArn:
    Value: !GetAtt MyFunc.Arn
`
	p := New()
	template, err := p.ParseYAML([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseYAML failed: %v", err)
	}

	if template.AWSTemplateFormatVersion != "2010-09-09" {
		t.Errorf("expected AWSTemplateFormatVersion '2010-09-09', got %q", template.AWSTemplateFormatVersion)
	}

	if template.Transform != "AWS::Serverless-2016-10-31" {
		t.Errorf("expected Transform 'AWS::Serverless-2016-10-31', got %v", template.Transform)
	}

	if template.Description != "Test template" {
		t.Errorf("expected Description 'Test template', got %q", template.Description)
	}

	if len(template.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(template.Parameters))
	}

	if param, ok := template.Parameters["Stage"]; ok {
		if param.Type != "String" {
			t.Errorf("expected parameter type 'String', got %q", param.Type)
		}
	} else {
		t.Error("expected parameter 'Stage' not found")
	}

	if len(template.Resources) != 1 {
		t.Errorf("expected 1 resource, got %d", len(template.Resources))
	}

	if res, ok := template.Resources["MyFunc"]; ok {
		if res.Type != "AWS::Serverless::Function" {
			t.Errorf("expected resource type 'AWS::Serverless::Function', got %q", res.Type)
		}
	} else {
		t.Error("expected resource 'MyFunc' not found")
	}

	if len(template.Outputs) != 1 {
		t.Errorf("expected 1 output, got %d", len(template.Outputs))
	}

	if output, ok := template.Outputs["FuncArn"]; ok {
		if output.Value == nil {
			t.Error("expected output Value not to be nil")
		}
	} else {
		t.Error("expected output 'FuncArn' not found")
	}
}

func TestIsIntrinsicFunction(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  bool
	}{
		{
			name:  "Ref",
			value: map[string]interface{}{"Ref": "MyParam"},
			want:  true,
		},
		{
			name:  "Fn::Sub",
			value: map[string]interface{}{"Fn::Sub": "hello"},
			want:  true,
		},
		{
			name:  "Fn::GetAtt",
			value: map[string]interface{}{"Fn::GetAtt": []interface{}{"Res", "Attr"}},
			want:  true,
		},
		{
			name:  "not intrinsic - multiple keys",
			value: map[string]interface{}{"Ref": "A", "Other": "B"},
			want:  false,
		},
		{
			name:  "not intrinsic - unknown key",
			value: map[string]interface{}{"Unknown": "value"},
			want:  false,
		},
		{
			name:  "not intrinsic - string",
			value: "hello",
			want:  false,
		},
		{
			name:  "not intrinsic - nil",
			value: nil,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsIntrinsicFunction(tt.value)
			if got != tt.want {
				t.Errorf("IsIntrinsicFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIntrinsicName(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  string
	}{
		{
			name:  "Ref",
			value: map[string]interface{}{"Ref": "MyParam"},
			want:  "Ref",
		},
		{
			name:  "Fn::Sub",
			value: map[string]interface{}{"Fn::Sub": "hello"},
			want:  "Fn::Sub",
		},
		{
			name:  "not intrinsic",
			value: map[string]interface{}{"Unknown": "value"},
			want:  "",
		},
		{
			name:  "string",
			value: "hello",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetIntrinsicName(tt.value)
			if got != tt.want {
				t.Errorf("GetIntrinsicName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectIntrinsics(t *testing.T) {
	data := map[string]interface{}{
		"Resources": map[string]interface{}{
			"MyFunc": map[string]interface{}{
				"Type": "AWS::Lambda::Function",
				"Properties": map[string]interface{}{
					"Runtime": map[string]interface{}{"Ref": "RuntimeParam"},
					"Role":    map[string]interface{}{"Fn::GetAtt": []interface{}{"MyRole", "Arn"}},
				},
			},
		},
	}

	intrinsics := DetectIntrinsics(data, "")
	if len(intrinsics) != 2 {
		t.Errorf("expected 2 intrinsics, got %d", len(intrinsics))
	}
}

func TestContainsIntrinsics(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  bool
	}{
		{
			name: "contains Ref",
			value: map[string]interface{}{
				"Runtime": map[string]interface{}{"Ref": "MyParam"},
			},
			want: true,
		},
		{
			name: "nested contains intrinsic",
			value: map[string]interface{}{
				"Nested": map[string]interface{}{
					"Deep": map[string]interface{}{"Fn::Sub": "hello"},
				},
			},
			want: true,
		},
		{
			name: "no intrinsics",
			value: map[string]interface{}{
				"Plain": "value",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsIntrinsics(tt.value)
			if got != tt.want {
				t.Errorf("ContainsIntrinsics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateIntrinsicStructure(t *testing.T) {
	tests := []struct {
		name      string
		intrinsic string
		value     interface{}
		wantErr   bool
	}{
		{
			name:      "valid Ref",
			intrinsic: "Ref",
			value:     "MyParam",
			wantErr:   false,
		},
		{
			name:      "invalid Ref - not string",
			intrinsic: "Ref",
			value:     123,
			wantErr:   true,
		},
		{
			name:      "valid Fn::GetAtt array",
			intrinsic: "Fn::GetAtt",
			value:     []interface{}{"Resource", "Attr"},
			wantErr:   false,
		},
		{
			name:      "valid Fn::GetAtt string",
			intrinsic: "Fn::GetAtt",
			value:     "Resource.Attr",
			wantErr:   false,
		},
		{
			name:      "invalid Fn::GetAtt - wrong array length",
			intrinsic: "Fn::GetAtt",
			value:     []interface{}{"OnlyOne"},
			wantErr:   true,
		},
		{
			name:      "valid Fn::Sub string",
			intrinsic: "Fn::Sub",
			value:     "${AWS::StackName}",
			wantErr:   false,
		},
		{
			name:      "valid Fn::Sub array",
			intrinsic: "Fn::Sub",
			value:     []interface{}{"${Var}", map[string]interface{}{"Var": "value"}},
			wantErr:   false,
		},
		{
			name:      "invalid Fn::Sub - wrong array length",
			intrinsic: "Fn::Sub",
			value:     []interface{}{"only one"},
			wantErr:   true,
		},
		{
			name:      "valid Fn::Join",
			intrinsic: "Fn::Join",
			value:     []interface{}{",", []interface{}{"a", "b"}},
			wantErr:   false,
		},
		{
			name:      "invalid Fn::Join - not array",
			intrinsic: "Fn::Join",
			value:     "not an array",
			wantErr:   true,
		},
		{
			name:      "valid Fn::If",
			intrinsic: "Fn::If",
			value:     []interface{}{"Condition", "True", "False"},
			wantErr:   false,
		},
		{
			name:      "invalid Fn::If - wrong length",
			intrinsic: "Fn::If",
			value:     []interface{}{"Condition", "True"},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIntrinsicStructure(tt.intrinsic, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIntrinsicStructure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTemplate(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid template",
			data: map[string]interface{}{
				"Resources": map[string]interface{}{
					"MyFunc": map[string]interface{}{
						"Type": "AWS::Lambda::Function",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing Resources",
			data: map[string]interface{}{
				"Description": "No resources",
			},
			wantErr: true,
		},
		{
			name: "resource missing Type",
			data: map[string]interface{}{
				"Resources": map[string]interface{}{
					"MyFunc": map[string]interface{}{
						"Properties": map[string]interface{}{},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTemplate(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseWithAutoDetect(t *testing.T) {
	yamlData := `
Resources:
  MyFunc:
    Type: AWS::Lambda::Function
`
	jsonData := `{"Resources": {"MyFunc": {"Type": "AWS::Lambda::Function"}}}`

	p := New()

	// Test YAML
	template, err := p.Parse([]byte(yamlData))
	if err != nil {
		t.Fatalf("Parse(yaml) failed: %v", err)
	}
	if _, ok := template.Resources["MyFunc"]; !ok {
		t.Error("expected resource 'MyFunc' not found in YAML parse")
	}

	// Test JSON
	template, err = p.Parse([]byte(jsonData))
	if err != nil {
		t.Fatalf("Parse(json) failed: %v", err)
	}
	if _, ok := template.Resources["MyFunc"]; !ok {
		t.Error("expected resource 'MyFunc' not found in JSON parse")
	}
}

func TestLocationTracking(t *testing.T) {
	yaml := `
Resources:
  MyFunc:
    Type: AWS::Lambda::Function
    Properties:
      Runtime: python3.9
`
	p := NewWithLocationTracking()
	_, err := p.ParseRawYAML([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseRawYAML failed: %v", err)
	}

	if p.Locations == nil {
		t.Fatal("expected Locations to be set")
	}

	// Check that at least some locations were tracked
	loc, ok := p.Locations.Get("Resources")
	if !ok {
		t.Error("expected 'Resources' location to be tracked")
	}
	if loc.Line == 0 {
		t.Error("expected 'Resources' line to be > 0")
	}
}

func TestParseErrorWithLocation(t *testing.T) {
	err := &ParseError{
		Message:  "test error",
		Location: SourceLocation{Line: 10, Column: 5},
	}

	errStr := err.Error()
	if errStr != "parse error at line 10, column 5: test error" {
		t.Errorf("unexpected error string: %s", errStr)
	}
}

func TestParseErrorWithoutLocation(t *testing.T) {
	err := &ParseError{
		Message: "test error",
	}

	errStr := err.Error()
	if errStr != "parse error: test error" {
		t.Errorf("unexpected error string: %s", errStr)
	}
}

func TestAllShortFormIntrinsics(t *testing.T) {
	intrinsics := []struct {
		tag  string
		key  string
		yaml string
	}{
		{"!Ref", "Ref", "Value: !Ref MyParam"},
		{"!Sub", "Fn::Sub", "Value: !Sub hello"},
		{"!GetAtt", "Fn::GetAtt", "Value: !GetAtt Res.Attr"},
		{"!Join", "Fn::Join", "Value: !Join\n  - ','\n  - - a\n    - b"},
		{"!If", "Fn::If", "Value: !If\n  - Cond\n  - A\n  - B"},
		{"!Select", "Fn::Select", "Value: !Select\n  - 0\n  - - a"},
		{"!FindInMap", "Fn::FindInMap", "Value: !FindInMap\n  - Map\n  - Key1\n  - Key2"},
		{"!Base64", "Fn::Base64", "Value: !Base64 hello"},
		{"!Cidr", "Fn::Cidr", "Value: !Cidr\n  - 10.0.0.0/16\n  - 6\n  - 5"},
		{"!GetAZs", "Fn::GetAZs", "Value: !GetAZs ''"},
		{"!ImportValue", "Fn::ImportValue", "Value: !ImportValue SharedValue"},
		{"!Split", "Fn::Split", "Value: !Split\n  - ','\n  - a,b"},
	}

	for _, tc := range intrinsics {
		t.Run(tc.tag, func(t *testing.T) {
			p := New()
			result, err := p.ParseRawYAML([]byte(tc.yaml))
			if err != nil {
				t.Fatalf("ParseRawYAML failed for %s: %v", tc.tag, err)
			}

			value, ok := result["Value"].(map[string]interface{})
			if !ok {
				t.Fatalf("expected Value to be map for %s, got %T", tc.tag, result["Value"])
			}

			if _, ok := value[tc.key]; !ok {
				t.Errorf("expected key %q not found for %s: %v", tc.key, tc.tag, value)
			}
		})
	}
}
