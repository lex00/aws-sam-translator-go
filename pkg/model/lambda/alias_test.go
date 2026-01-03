package lambda

import (
	"testing"
)

func TestNewAlias(t *testing.T) {
	alias := NewAlias("my-function", "1", "prod")

	if alias.FunctionName != "my-function" {
		t.Errorf("expected FunctionName 'my-function', got %v", alias.FunctionName)
	}
	if alias.FunctionVersion != "1" {
		t.Errorf("expected FunctionVersion '1', got %v", alias.FunctionVersion)
	}
	if alias.Name != "prod" {
		t.Errorf("expected Name 'prod', got %s", alias.Name)
	}
}

func TestNewAliasWithDescription(t *testing.T) {
	alias := NewAliasWithDescription("my-function", "1", "prod", "Production alias")

	if alias.Description != "Production alias" {
		t.Errorf("expected Description 'Production alias', got %s", alias.Description)
	}
}

func TestAliasWithProvisionedConcurrency(t *testing.T) {
	alias := NewAlias("my-function", "1", "prod").WithProvisionedConcurrency(100)

	if alias.ProvisionedConcurrencyConfig == nil {
		t.Fatal("expected ProvisionedConcurrencyConfig to be set")
	}
	if alias.ProvisionedConcurrencyConfig.ProvisionedConcurrentExecutions != 100 {
		t.Errorf("expected ProvisionedConcurrentExecutions 100, got %d",
			alias.ProvisionedConcurrencyConfig.ProvisionedConcurrentExecutions)
	}
}

func TestAliasWithRoutingConfig(t *testing.T) {
	weights := []VersionWeight{
		{FunctionVersion: "2", FunctionWeight: 0.1},
	}
	alias := NewAlias("my-function", "1", "prod").WithRoutingConfig(weights)

	if alias.RoutingConfig == nil {
		t.Fatal("expected RoutingConfig to be set")
	}
	if len(alias.RoutingConfig.AdditionalVersionWeights) != 1 {
		t.Errorf("expected 1 version weight, got %d", len(alias.RoutingConfig.AdditionalVersionWeights))
	}
	if alias.RoutingConfig.AdditionalVersionWeights[0].FunctionWeight != 0.1 {
		t.Errorf("expected FunctionWeight 0.1, got %f", alias.RoutingConfig.AdditionalVersionWeights[0].FunctionWeight)
	}
}

func TestAliasAddVersionWeight(t *testing.T) {
	alias := NewAlias("my-function", "1", "prod").
		AddVersionWeight("2", 0.1).
		AddVersionWeight("3", 0.05)

	if alias.RoutingConfig == nil {
		t.Fatal("expected RoutingConfig to be set")
	}
	if len(alias.RoutingConfig.AdditionalVersionWeights) != 2 {
		t.Errorf("expected 2 version weights, got %d", len(alias.RoutingConfig.AdditionalVersionWeights))
	}
}

func TestAliasToCloudFormation_Minimal(t *testing.T) {
	alias := NewAlias("my-function", "1", "prod")

	result := alias.ToCloudFormation()

	if result["Type"] != ResourceTypeAlias {
		t.Errorf("expected Type %s, got %v", ResourceTypeAlias, result["Type"])
	}

	props := result["Properties"].(map[string]interface{})
	if props["FunctionName"] != "my-function" {
		t.Errorf("expected FunctionName in properties")
	}
	if props["FunctionVersion"] != "1" {
		t.Errorf("expected FunctionVersion in properties")
	}
	if props["Name"] != "prod" {
		t.Errorf("expected Name in properties")
	}
}

func TestAliasToCloudFormation_Full(t *testing.T) {
	alias := NewAliasWithDescription("my-function", "1", "prod", "Production alias").
		WithProvisionedConcurrency(50).
		AddVersionWeight("2", 0.1)

	result := alias.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	if props["Description"] != "Production alias" {
		t.Errorf("expected Description 'Production alias', got %v", props["Description"])
	}

	provConfig := props["ProvisionedConcurrencyConfig"].(map[string]interface{})
	if provConfig["ProvisionedConcurrentExecutions"] != 50 {
		t.Errorf("expected ProvisionedConcurrentExecutions 50, got %v", provConfig["ProvisionedConcurrentExecutions"])
	}

	routingConfig := props["RoutingConfig"].(map[string]interface{})
	weights := routingConfig["AdditionalVersionWeights"].([]map[string]interface{})
	if len(weights) != 1 {
		t.Errorf("expected 1 weight, got %d", len(weights))
	}
	if weights[0]["FunctionVersion"] != "2" {
		t.Errorf("expected FunctionVersion '2', got %v", weights[0]["FunctionVersion"])
	}
}

func TestAliasWithIntrinsicReferences(t *testing.T) {
	fnRef := map[string]interface{}{"Ref": "MyLambdaFunction"}
	versionRef := map[string]interface{}{"Fn::GetAtt": []string{"MyLambdaVersion", "Version"}}

	alias := NewAlias(fnRef, versionRef, "live")

	result := alias.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	fnName := props["FunctionName"].(map[string]interface{})
	if fnName["Ref"] != "MyLambdaFunction" {
		t.Errorf("expected Ref to MyLambdaFunction, got %v", fnName)
	}

	fnVersion := props["FunctionVersion"].(map[string]interface{})
	if fnVersion["Fn::GetAtt"] == nil {
		t.Error("expected Fn::GetAtt in FunctionVersion")
	}
}

func TestAliasRoutingConfigToMap_Empty(t *testing.T) {
	config := &AliasRoutingConfig{}
	result := config.toMap()

	if result != nil {
		t.Errorf("expected nil for empty routing config, got %v", result)
	}
}

func TestAliasRoutingConfigToMap_WithWeights(t *testing.T) {
	config := &AliasRoutingConfig{
		AdditionalVersionWeights: []VersionWeight{
			{FunctionVersion: "2", FunctionWeight: 0.2},
			{FunctionVersion: "3", FunctionWeight: 0.1},
		},
	}

	result := config.toMap()
	weights := result["AdditionalVersionWeights"].([]map[string]interface{})

	if len(weights) != 2 {
		t.Errorf("expected 2 weights, got %d", len(weights))
	}
	if weights[0]["FunctionWeight"] != 0.2 {
		t.Errorf("expected first weight 0.2, got %v", weights[0]["FunctionWeight"])
	}
}

func TestAliasProvisionedConcurrencyConfigToMap(t *testing.T) {
	config := &AliasProvisionedConcurrencyConfig{
		ProvisionedConcurrentExecutions: 75,
	}

	result := config.toMap()
	if result["ProvisionedConcurrentExecutions"] != 75 {
		t.Errorf("expected 75, got %v", result["ProvisionedConcurrentExecutions"])
	}
}

func TestAliasNoRoutingConfigWhenEmpty(t *testing.T) {
	alias := NewAlias("my-function", "1", "prod")

	result := alias.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	if props["RoutingConfig"] != nil {
		t.Error("expected RoutingConfig to be nil when no weights are set")
	}
}
