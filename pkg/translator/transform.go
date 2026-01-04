package translator

// This file contains type parsing functions that convert raw map[string]interface{}
// properties from JSON/YAML into typed SAM structs for transformation.

import (
	"encoding/json"
	"fmt"

	"github.com/lex00/aws-sam-translator-go/pkg/sam"
)

// parseFunction parses properties into a Function struct.
func (t *Translator) parseFunction(props map[string]interface{}) (*sam.Function, error) {
	if props == nil {
		return nil, fmt.Errorf("function properties cannot be nil")
	}

	fn := &sam.Function{}

	// Required for Zip package type
	if v, ok := props["Handler"].(string); ok {
		fn.Handler = v
	}
	if v, ok := props["Runtime"].(string); ok {
		fn.Runtime = v
	}
	if v, ok := props["CodeUri"]; ok {
		fn.CodeUri = v
	}
	if v, ok := props["ImageUri"]; ok {
		fn.ImageUri = v
	}
	if v, ok := props["PackageType"].(string); ok {
		fn.PackageType = v
	}
	if v, ok := props["Description"].(string); ok {
		fn.Description = v
	}
	if v, ok := props["MemorySize"].(int); ok {
		fn.MemorySize = v
	} else if v, ok := props["MemorySize"].(float64); ok {
		fn.MemorySize = int(v)
	}
	if v, ok := props["Timeout"].(int); ok {
		fn.Timeout = v
	} else if v, ok := props["Timeout"].(float64); ok {
		fn.Timeout = int(v)
	}
	if v, ok := props["Role"]; ok {
		fn.Role = v
	}
	if v, ok := props["Policies"]; ok {
		fn.Policies = v
	}
	if v, ok := props["Environment"].(map[string]interface{}); ok {
		fn.Environment = v
	}
	if v, ok := props["Events"].(map[string]interface{}); ok {
		fn.Events = v
	}
	if v, ok := props["Tags"].(map[string]interface{}); ok {
		fn.Tags = convertToStringMap(v)
	}
	if v, ok := props["Layers"].([]interface{}); ok {
		fn.Layers = v
	}
	if v, ok := props["VpcConfig"].(map[string]interface{}); ok {
		fn.VpcConfig = v
	}
	if v, ok := props["FunctionName"]; ok {
		fn.FunctionName = v
	}
	if v, ok := props["Architectures"].([]interface{}); ok {
		fn.Architectures = convertToStringSlice(v)
	}
	if v, ok := props["AutoPublishAlias"].(string); ok {
		fn.AutoPublishAlias = v
	}
	if v, ok := props["AutoPublishCodeSha256"].(string); ok {
		fn.AutoPublishCodeSha256 = v
	}
	if v, ok := props["DeploymentPreference"].(map[string]interface{}); ok {
		fn.DeploymentPreference = v
	}
	if v, ok := props["ProvisionedConcurrencyConfig"].(map[string]interface{}); ok {
		fn.ProvisionedConcurrencyConfig = v
	}
	if v, ok := props["ReservedConcurrentExecutions"]; ok {
		switch val := v.(type) {
		case int:
			fn.ReservedConcurrentExecutions = &val
		case float64:
			intVal := int(val)
			fn.ReservedConcurrentExecutions = &intVal
		}
	}
	if v, ok := props["Tracing"].(string); ok {
		fn.Tracing = v
	}
	if v, ok := props["DeadLetterQueue"].(map[string]interface{}); ok {
		fn.DeadLetterQueue = v
	}
	if v, ok := props["KmsKeyArn"]; ok {
		fn.KmsKeyArn = v
	}
	if v, ok := props["EphemeralStorage"].(map[string]interface{}); ok {
		fn.EphemeralStorage = v
	}
	if v, ok := props["SnapStart"].(map[string]interface{}); ok {
		fn.SnapStart = v
	}
	if v, ok := props["FileSystemConfigs"].([]interface{}); ok {
		fn.FileSystemConfigs = convertToMapSlice(v)
	}
	if v, ok := props["ImageConfig"].(map[string]interface{}); ok {
		fn.ImageConfig = v
	}
	if v, ok := props["CodeSigningConfigArn"]; ok {
		fn.CodeSigningConfigArn = v
	}
	if v, ok := props["RuntimeManagementConfig"].(map[string]interface{}); ok {
		fn.RuntimeManagementConfig = v
	}
	if v, ok := props["PermissionsBoundary"]; ok {
		fn.PermissionsBoundary = v
	}
	if v, ok := props["FunctionUrlConfig"].(map[string]interface{}); ok {
		fn.FunctionUrlConfig = v
	}
	if v, ok := props["LoggingConfig"].(map[string]interface{}); ok {
		fn.LoggingConfig = v
	}
	if v, ok := props["RecursiveLoop"].(string); ok {
		fn.RecursiveLoop = v
	}

	return fn, nil
}

// parseSimpleTable parses properties into a SimpleTable struct.
func (t *Translator) parseSimpleTable(props map[string]interface{}) (*sam.SimpleTable, error) {
	st := &sam.SimpleTable{}

	if props == nil {
		return st, nil
	}

	if v, ok := props["TableName"]; ok {
		st.TableName = v
	}
	if v, ok := props["PrimaryKey"].(map[string]interface{}); ok {
		pk := &sam.PrimaryKey{}
		if name, ok := v["Name"]; ok {
			pk.Name = name
		}
		if typ, ok := v["Type"].(string); ok {
			pk.Type = typ
		}
		st.PrimaryKey = pk
	}
	if v, ok := props["ProvisionedThroughput"].(map[string]interface{}); ok {
		pt := &sam.ProvisionedThroughput{}
		if rcu, ok := v["ReadCapacityUnits"]; ok {
			pt.ReadCapacityUnits = rcu
		}
		if wcu, ok := v["WriteCapacityUnits"]; ok {
			pt.WriteCapacityUnits = wcu
		}
		st.ProvisionedThroughput = pt
	}
	if v, ok := props["SSESpecification"].(map[string]interface{}); ok {
		sse := &sam.SSESpecification{}
		if enabled, ok := v["SSEEnabled"]; ok {
			sse.SSEEnabled = enabled
		}
		st.SSESpecification = sse
	}
	if v, ok := props["Tags"].(map[string]interface{}); ok {
		st.Tags = convertToStringMap(v)
	}
	if v, ok := props["PointInTimeRecoverySpecification"].(map[string]interface{}); ok {
		pitr := &sam.PointInTimeRecoverySpecification{}
		if enabled, ok := v["PointInTimeRecoveryEnabled"]; ok {
			pitr.PointInTimeRecoveryEnabled = enabled
		}
		st.PointInTimeRecoverySpecification = pitr
	}

	return st, nil
}

// parseLayerVersion parses properties into a LayerVersion struct.
func (t *Translator) parseLayerVersion(props map[string]interface{}) (*sam.LayerVersion, error) {
	lv := &sam.LayerVersion{}

	if props == nil {
		return nil, fmt.Errorf("layer version properties cannot be nil")
	}

	if v, ok := props["ContentUri"]; ok {
		lv.ContentUri = v
	}
	if v, ok := props["LayerName"].(string); ok {
		lv.LayerName = v
	}
	if v, ok := props["Description"].(string); ok {
		lv.Description = v
	}
	if v, ok := props["LicenseInfo"].(string); ok {
		lv.LicenseInfo = v
	}
	if v, ok := props["CompatibleRuntimes"].([]interface{}); ok {
		lv.CompatibleRuntimes = convertToStringSlice(v)
	}
	if v, ok := props["CompatibleArchitectures"].([]interface{}); ok {
		lv.CompatibleArchitectures = convertToStringSlice(v)
	}
	if v, ok := props["RetentionPolicy"]; ok {
		lv.RetentionPolicy = v
	}

	return lv, nil
}

// parseStateMachine parses properties into a StateMachine struct.
func (t *Translator) parseStateMachine(props map[string]interface{}) (*sam.StateMachine, error) {
	sm := &sam.StateMachine{}

	if props == nil {
		return nil, fmt.Errorf("state machine properties cannot be nil")
	}

	if v, ok := props["Definition"].(map[string]interface{}); ok {
		sm.Definition = v
	}
	if v, ok := props["DefinitionUri"]; ok {
		sm.DefinitionUri = v
	}
	if v, ok := props["DefinitionSubstitutions"].(map[string]interface{}); ok {
		sm.DefinitionSubstitutions = v
	}
	if v, ok := props["Name"].(string); ok {
		sm.Name = v
	}
	if v, ok := props["Role"].(string); ok {
		sm.Role = v
	}
	if v, ok := props["RolePath"].(string); ok {
		sm.RolePath = v
	}
	if v, ok := props["Policies"]; ok {
		sm.Policies = v
	}
	if v, ok := props["Logging"].(map[string]interface{}); ok {
		sm.Logging = t.parseLoggingConfig(v)
	}
	if v, ok := props["Tracing"].(map[string]interface{}); ok {
		sm.Tracing = t.parseTracingConfig(v)
	}
	if v, ok := props["Events"].(map[string]interface{}); ok {
		sm.Events = v
	}
	if v, ok := props["Tags"].(map[string]interface{}); ok {
		sm.Tags = v
	}
	if v, ok := props["Type"].(string); ok {
		sm.Type = v
	}
	if v, ok := props["PermissionsBoundary"].(string); ok {
		sm.PermissionsBoundary = v
	}

	return sm, nil
}

// parseTracingConfig parses Tracing configuration for StateMachine.
func (t *Translator) parseTracingConfig(m map[string]interface{}) *sam.TracingConfig {
	cfg := &sam.TracingConfig{}
	if v, ok := m["Enabled"].(bool); ok {
		cfg.Enabled = v
	}
	return cfg
}

// parseLoggingConfig parses Logging configuration for StateMachine.
func (t *Translator) parseLoggingConfig(m map[string]interface{}) *sam.LoggingConfig {
	cfg := &sam.LoggingConfig{}
	if v, ok := m["Level"].(string); ok {
		cfg.Level = v
	}
	if v, ok := m["IncludeExecutionData"].(bool); ok {
		cfg.IncludeExecutionData = v
	}
	if v, ok := m["Destinations"].([]interface{}); ok {
		cfg.Destinations = v
	}
	return cfg
}

// parseApi parses properties into an Api struct.
func (t *Translator) parseApi(props map[string]interface{}) (*sam.Api, error) {
	api := &sam.Api{}

	if props == nil {
		return nil, fmt.Errorf("api properties cannot be nil")
	}

	if v, ok := props["StageName"]; ok {
		api.StageName = v
	}
	if v, ok := props["Name"]; ok {
		api.Name = v
	}
	if v, ok := props["Description"]; ok {
		api.Description = v
	}
	if v, ok := props["DefinitionBody"].(map[string]interface{}); ok {
		api.DefinitionBody = v
	}
	if v, ok := props["DefinitionUri"]; ok {
		api.DefinitionUri = v
	}
	if v, ok := props["CacheClusterEnabled"].(bool); ok {
		api.CacheClusterEnabled = v
	}
	if v, ok := props["CacheClusterSize"].(string); ok {
		api.CacheClusterSize = v
	}
	if v, ok := props["Variables"].(map[string]interface{}); ok {
		api.Variables = v
	}
	if v, ok := props["EndpointConfiguration"].(map[string]interface{}); ok {
		api.EndpointConfiguration = t.parseEndpointConfig(v)
	}
	if v, ok := props["MethodSettings"].([]interface{}); ok {
		api.MethodSettings = t.parseMethodSettings(v)
	}
	if v, ok := props["BinaryMediaTypes"].([]interface{}); ok {
		api.BinaryMediaTypes = v
	}
	if v, ok := props["MinimumCompressionSize"].(int); ok {
		api.MinimumCompressionSize = v
	} else if v, ok := props["MinimumCompressionSize"].(float64); ok {
		api.MinimumCompressionSize = int(v)
	}
	if v, ok := props["Cors"]; ok {
		api.Cors = t.parseCorsConfig(v)
	}
	if v, ok := props["Auth"].(map[string]interface{}); ok {
		api.Auth = t.parseApiAuth(v)
	}
	if v, ok := props["GatewayResponses"].(map[string]interface{}); ok {
		api.GatewayResponses = v
	}
	if v, ok := props["AccessLogSetting"].(map[string]interface{}); ok {
		api.AccessLogSetting = t.parseAccessLogSetting(v)
	}
	if v, ok := props["CanarySetting"].(map[string]interface{}); ok {
		api.CanarySetting = t.parseCanarySettingConfig(v)
	}
	if v, ok := props["TracingEnabled"].(bool); ok {
		api.TracingEnabled = v
	}
	if v, ok := props["OpenApiVersion"].(string); ok {
		api.OpenApiVersion = v
	}
	if v, ok := props["Models"].(map[string]interface{}); ok {
		api.Models = v
	}
	if v, ok := props["Domain"].(map[string]interface{}); ok {
		api.Domain = t.parseDomainConfig(v)
	}
	if v, ok := props["FailOnWarnings"].(bool); ok {
		api.FailOnWarnings = v
	}
	if v, ok := props["DisableExecuteApiEndpoint"].(bool); ok {
		api.DisableExecuteApiEndpoint = v
	}
	if v, ok := props["Tags"].(map[string]interface{}); ok {
		api.Tags = v
	}
	if v, ok := props["ApiKeySourceType"].(string); ok {
		api.ApiKeySourceType = v
	}

	return api, nil
}

// parseEndpointConfig parses EndpointConfiguration for Api.
func (t *Translator) parseEndpointConfig(m map[string]interface{}) *sam.EndpointConfig {
	cfg := &sam.EndpointConfig{}
	if v, ok := m["Type"].(string); ok {
		cfg.Type = v
	}
	if v, ok := m["VPCEndpointIds"].([]interface{}); ok {
		cfg.VPCEndpointIds = v
	}
	return cfg
}

// parseMethodSettings parses MethodSettings for Api.
func (t *Translator) parseMethodSettings(arr []interface{}) []sam.MethodSettingConfig {
	var result []sam.MethodSettingConfig
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			cfg := sam.MethodSettingConfig{}
			if v, ok := m["HttpMethod"].(string); ok {
				cfg.HttpMethod = v
			}
			if v, ok := m["ResourcePath"].(string); ok {
				cfg.ResourcePath = v
			}
			if v, ok := m["CachingEnabled"].(bool); ok {
				cfg.CachingEnabled = v
			}
			if v, ok := m["CacheTtlInSeconds"].(int); ok {
				cfg.CacheTtlInSeconds = v
			} else if v, ok := m["CacheTtlInSeconds"].(float64); ok {
				cfg.CacheTtlInSeconds = int(v)
			}
			if v, ok := m["DataTraceEnabled"].(bool); ok {
				cfg.DataTraceEnabled = v
			}
			if v, ok := m["LoggingLevel"].(string); ok {
				cfg.LoggingLevel = v
			}
			if v, ok := m["MetricsEnabled"].(bool); ok {
				cfg.MetricsEnabled = v
			}
			if v, ok := m["ThrottlingBurstLimit"].(int); ok {
				cfg.ThrottlingBurstLimit = v
			} else if v, ok := m["ThrottlingBurstLimit"].(float64); ok {
				cfg.ThrottlingBurstLimit = int(v)
			}
			if v, ok := m["ThrottlingRateLimit"].(float64); ok {
				cfg.ThrottlingRateLimit = v
			}
			result = append(result, cfg)
		}
	}
	return result
}

// parseCorsConfig parses Cors configuration for Api.
func (t *Translator) parseCorsConfig(val interface{}) *sam.CorsConfig {
	// Cors can be a string, bool, or map
	if m, ok := val.(map[string]interface{}); ok {
		cfg := &sam.CorsConfig{}
		if v, ok := m["AllowOrigin"]; ok {
			cfg.AllowOrigin = v
		}
		if v, ok := m["AllowMethods"]; ok {
			cfg.AllowMethods = v
		}
		if v, ok := m["AllowHeaders"]; ok {
			cfg.AllowHeaders = v
		}
		if v, ok := m["MaxAge"]; ok {
			cfg.MaxAge = v
		}
		if v, ok := m["AllowCredentials"].(bool); ok {
			cfg.AllowCredentials = v
		}
		return cfg
	}
	// For string values like "'*'", store as AllowOrigin
	if s, ok := val.(string); ok {
		return &sam.CorsConfig{AllowOrigin: s}
	}
	return nil
}

// parseApiAuth parses Auth configuration for Api.
func (t *Translator) parseApiAuth(m map[string]interface{}) *sam.ApiAuth {
	auth := &sam.ApiAuth{}
	if v, ok := m["DefaultAuthorizer"].(string); ok {
		auth.DefaultAuthorizer = v
	}
	if v, ok := m["Authorizers"].(map[string]interface{}); ok {
		auth.Authorizers = v
	}
	if v, ok := m["ApiKeyRequired"].(bool); ok {
		auth.ApiKeyRequired = v
	}
	if v, ok := m["ResourcePolicy"]; ok {
		auth.ResourcePolicy = v
	}
	if v, ok := m["UsagePlan"]; ok {
		auth.UsagePlan = v
	}
	return auth
}

// parseAccessLogSetting parses AccessLogSetting for Api.
func (t *Translator) parseAccessLogSetting(m map[string]interface{}) *sam.AccessLogSetting {
	setting := &sam.AccessLogSetting{}
	if v, ok := m["DestinationArn"]; ok {
		setting.DestinationArn = v
	}
	if v, ok := m["Format"]; ok {
		setting.Format = v
	}
	return setting
}

// parseCanarySettingConfig parses CanarySetting for Api.
func (t *Translator) parseCanarySettingConfig(m map[string]interface{}) *sam.CanarySettingConfig {
	cfg := &sam.CanarySettingConfig{}
	if v, ok := m["PercentTraffic"].(float64); ok {
		cfg.PercentTraffic = v
	}
	if v, ok := m["DeploymentId"]; ok {
		cfg.DeploymentId = v
	}
	if v, ok := m["StageVariableOverrides"].(map[string]interface{}); ok {
		cfg.StageVariableOverrides = v
	}
	if v, ok := m["UseStageCache"].(bool); ok {
		cfg.UseStageCache = v
	}
	return cfg
}

// parseDomainConfig parses Domain configuration for Api.
func (t *Translator) parseDomainConfig(m map[string]interface{}) *sam.DomainConfig {
	cfg := &sam.DomainConfig{}
	if v, ok := m["DomainName"]; ok {
		cfg.DomainName = v
	}
	if v, ok := m["CertificateArn"]; ok {
		cfg.CertificateArn = v
	}
	if v, ok := m["EndpointConfiguration"]; ok {
		cfg.EndpointConfiguration = v
	}
	if v, ok := m["Route53"]; ok {
		cfg.Route53 = v
	}
	if v, ok := m["BasePath"]; ok {
		cfg.BasePath = v
	}
	return cfg
}

// parseHttpApi parses properties into an HttpApi struct.
func (t *Translator) parseHttpApi(props map[string]interface{}) (*sam.HttpApi, error) {
	api := &sam.HttpApi{}

	if props == nil {
		// HttpApi can have empty properties - uses defaults
		return api, nil
	}

	if v, ok := props["StageName"]; ok {
		api.StageName = v
	}
	if v, ok := props["Name"]; ok {
		api.Name = v
	}
	if v, ok := props["Description"]; ok {
		api.Description = v
	}
	if v, ok := props["DefinitionBody"].(map[string]interface{}); ok {
		api.DefinitionBody = v
	}
	if v, ok := props["DefinitionUri"]; ok {
		api.DefinitionUri = v
	}
	if v, ok := props["StageVariables"].(map[string]interface{}); ok {
		api.StageVariables = v
	}
	if v, ok := props["CorsConfiguration"]; ok {
		api.CorsConfiguration = v
	}
	if v, ok := props["Auth"].(map[string]interface{}); ok {
		api.Auth = t.parseHttpApiAuth(v)
	}
	if v, ok := props["AccessLogSettings"].(map[string]interface{}); ok {
		api.AccessLogSettings = t.parseHttpApiAccessLogSettings(v)
	}
	if v, ok := props["DefaultRouteSettings"].(map[string]interface{}); ok {
		api.DefaultRouteSettings = t.parseHttpApiRouteSettings(v)
	}
	if v, ok := props["RouteSettings"].(map[string]interface{}); ok {
		api.RouteSettings = v
	}
	if v, ok := props["Domain"].(map[string]interface{}); ok {
		api.Domain = t.parseHttpApiDomain(v)
	}
	if v, ok := props["FailOnWarnings"]; ok {
		api.FailOnWarnings = v
	}
	if v, ok := props["DisableExecuteApiEndpoint"]; ok {
		api.DisableExecuteApiEndpoint = v
	}
	if v, ok := props["Tags"].(map[string]interface{}); ok {
		api.Tags = v
	}

	return api, nil
}

// parseHttpApiAuth parses Auth configuration for HttpApi.
func (t *Translator) parseHttpApiAuth(m map[string]interface{}) *sam.HttpApiAuth {
	auth := &sam.HttpApiAuth{}
	if v, ok := m["DefaultAuthorizer"].(string); ok {
		auth.DefaultAuthorizer = v
	}
	if v, ok := m["Authorizers"].(map[string]interface{}); ok {
		auth.Authorizers = v
	}
	if v, ok := m["EnableIamAuthorizer"].(bool); ok {
		auth.EnableIamAuthorizer = v
	}
	return auth
}

// parseHttpApiAccessLogSettings parses AccessLogSettings for HttpApi.
func (t *Translator) parseHttpApiAccessLogSettings(m map[string]interface{}) *sam.HttpApiAccessLogSettings {
	settings := &sam.HttpApiAccessLogSettings{}
	if v, ok := m["DestinationArn"]; ok {
		settings.DestinationArn = v
	}
	if v, ok := m["Format"]; ok {
		settings.Format = v
	}
	return settings
}

// parseHttpApiRouteSettings parses RouteSettings for HttpApi.
func (t *Translator) parseHttpApiRouteSettings(m map[string]interface{}) *sam.HttpApiRouteSettings {
	settings := &sam.HttpApiRouteSettings{}
	if v, ok := m["ThrottlingBurstLimit"]; ok {
		settings.ThrottlingBurstLimit = v
	}
	if v, ok := m["ThrottlingRateLimit"]; ok {
		settings.ThrottlingRateLimit = v
	}
	if v, ok := m["DetailedMetricsEnabled"]; ok {
		settings.DetailedMetricsEnabled = v
	}
	return settings
}

// parseHttpApiDomain parses Domain configuration for HttpApi.
func (t *Translator) parseHttpApiDomain(m map[string]interface{}) *sam.HttpApiDomain {
	domain := &sam.HttpApiDomain{}
	if v, ok := m["DomainName"]; ok {
		domain.DomainName = v
	}
	if v, ok := m["CertificateArn"]; ok {
		domain.CertificateArn = v
	}
	if v, ok := m["EndpointConfiguration"]; ok {
		domain.EndpointConfiguration = v
	}
	if v, ok := m["SecurityPolicy"]; ok {
		domain.SecurityPolicy = v
	}
	if v, ok := m["BasePath"]; ok {
		domain.BasePath = v
	}
	if v, ok := m["MutualTlsAuthentication"].(map[string]interface{}); ok {
		domain.MutualTlsAuthentication = v
	}
	if v, ok := m["Route53"].(map[string]interface{}); ok {
		domain.Route53 = v
	}
	return domain
}

// parseApplication parses properties into an Application struct.
func (t *Translator) parseApplication(props map[string]interface{}) (*sam.Application, error) {
	app := &sam.Application{}

	if props == nil {
		return nil, fmt.Errorf("application properties cannot be nil")
	}

	if v, ok := props["Location"]; ok {
		app.Location = v
	}
	if v, ok := props["Parameters"].(map[string]interface{}); ok {
		app.Parameters = v
	}
	if v, ok := props["NotificationArns"].([]interface{}); ok {
		app.NotificationArns = v
	}
	if v, ok := props["Tags"].(map[string]interface{}); ok {
		app.Tags = convertToStringMap(v)
	}
	if v, ok := props["TimeoutInMinutes"]; ok {
		switch val := v.(type) {
		case int:
			app.TimeoutInMinutes = val
		case float64:
			app.TimeoutInMinutes = int(val)
		}
	}

	return app, nil
}

// parseGraphQLApi parses properties into a GraphQLApi struct.
func (t *Translator) parseGraphQLApi(props map[string]interface{}) (*sam.GraphQLApi, error) {
	gql := &sam.GraphQLApi{}

	if props == nil {
		return nil, fmt.Errorf("graphql api properties cannot be nil")
	}

	if v, ok := props["Name"].(string); ok {
		gql.Name = v
	}
	if v, ok := props["SchemaInline"].(string); ok {
		gql.SchemaInline = v
	}
	if v, ok := props["SchemaUri"]; ok {
		gql.SchemaUri = v
	}
	if v, ok := props["Auth"].(map[string]interface{}); ok {
		gql.Auth = t.parseGraphQLAuth(v)
	}
	if v, ok := props["AdditionalAuthenticationProviders"].([]interface{}); ok {
		gql.AdditionalAuthenticationProviders = t.parseGraphQLAuthList(v)
	}
	if v, ok := props["XrayEnabled"].(bool); ok {
		gql.XrayEnabled = v
	}
	if v, ok := props["Logging"].(map[string]interface{}); ok {
		gql.Logging = v
	}
	if v, ok := props["DataSources"].(map[string]interface{}); ok {
		gql.DataSources = t.parseGraphQLDataSources(v)
	}
	if v, ok := props["Functions"].(map[string]interface{}); ok {
		gql.Functions = t.parseGraphQLFunctions(v)
	}
	if v, ok := props["Resolvers"].(map[string]interface{}); ok {
		gql.Resolvers = t.parseGraphQLResolvers(v)
	}
	if v, ok := props["DomainName"].(map[string]interface{}); ok {
		gql.DomainName = v
	}
	if v, ok := props["Tags"].(map[string]interface{}); ok {
		gql.Tags = convertToStringMap(v)
	}
	if v, ok := props["Cache"].(map[string]interface{}); ok {
		gql.Cache = v
	}
	if v, ok := props["ApiKeys"].([]interface{}); ok {
		gql.ApiKeys = t.parseGraphQLApiKeys(v)
	}
	if v, ok := props["Condition"].(string); ok {
		gql.Condition = v
	}
	if v, ok := props["DependsOn"]; ok {
		gql.DependsOn = v
	}
	if v, ok := props["Metadata"].(map[string]interface{}); ok {
		gql.Metadata = v
	}

	return gql, nil
}

// parseGraphQLAuth parses a map into a GraphQLAuth struct.
func (t *Translator) parseGraphQLAuth(m map[string]interface{}) *sam.GraphQLAuth {
	auth := &sam.GraphQLAuth{}
	if v, ok := m["Type"].(string); ok {
		auth.Type = v
	}
	if v, ok := m["UserPoolConfig"].(map[string]interface{}); ok {
		auth.UserPoolConfig = v
	}
	if v, ok := m["OpenIDConnectConfig"].(map[string]interface{}); ok {
		auth.OpenIDConnectConfig = v
	}
	if v, ok := m["LambdaAuthorizerConfig"].(map[string]interface{}); ok {
		auth.LambdaAuthorizerConfig = v
	}
	return auth
}

// parseGraphQLAuthList parses a list of GraphQLAuth.
func (t *Translator) parseGraphQLAuthList(arr []interface{}) []sam.GraphQLAuth {
	result := make([]sam.GraphQLAuth, 0, len(arr))
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			auth := t.parseGraphQLAuth(m)
			if auth != nil {
				result = append(result, *auth)
			}
		}
	}
	return result
}

// parseGraphQLDataSources parses a map of data sources.
func (t *Translator) parseGraphQLDataSources(m map[string]interface{}) map[string]sam.GraphQLApiDataSource {
	result := make(map[string]sam.GraphQLApiDataSource)
	for name, item := range m {
		if dsMap, ok := item.(map[string]interface{}); ok {
			ds := sam.GraphQLApiDataSource{}
			if v, ok := dsMap["Type"].(string); ok {
				ds.Type = v
			}
			if v, ok := dsMap["Name"].(string); ok {
				ds.Name = v
			}
			if v, ok := dsMap["Description"].(string); ok {
				ds.Description = v
			}
			if v, ok := dsMap["ServiceRoleArn"]; ok {
				ds.ServiceRoleArn = v
			}
			if v, ok := dsMap["LambdaConfig"].(map[string]interface{}); ok {
				ds.LambdaConfig = v
			}
			if v, ok := dsMap["DynamoDBConfig"].(map[string]interface{}); ok {
				ds.DynamoDBConfig = v
			}
			if v, ok := dsMap["HttpConfig"].(map[string]interface{}); ok {
				ds.HttpConfig = v
			}
			if v, ok := dsMap["ElasticsearchConfig"].(map[string]interface{}); ok {
				ds.ElasticsearchConfig = v
			}
			if v, ok := dsMap["OpenSearchServiceConfig"].(map[string]interface{}); ok {
				ds.OpenSearchServiceConfig = v
			}
			if v, ok := dsMap["RelationalDatabaseConfig"].(map[string]interface{}); ok {
				ds.RelationalDatabaseConfig = v
			}
			if v, ok := dsMap["EventBridgeConfig"].(map[string]interface{}); ok {
				ds.EventBridgeConfig = v
			}
			result[name] = ds
		}
	}
	return result
}

// parseGraphQLFunctions parses a map of functions.
func (t *Translator) parseGraphQLFunctions(m map[string]interface{}) map[string]sam.GraphQLApiFunction {
	result := make(map[string]sam.GraphQLApiFunction)
	for name, item := range m {
		if fnMap, ok := item.(map[string]interface{}); ok {
			fn := sam.GraphQLApiFunction{}
			if v, ok := fnMap["Name"].(string); ok {
				fn.Name = v
			}
			if v, ok := fnMap["Description"].(string); ok {
				fn.Description = v
			}
			if v, ok := fnMap["DataSourceName"].(string); ok {
				fn.DataSourceName = v
			}
			if v, ok := fnMap["RequestMappingTemplate"].(string); ok {
				fn.RequestMappingTemplate = v
			}
			if v, ok := fnMap["ResponseMappingTemplate"].(string); ok {
				fn.ResponseMappingTemplate = v
			}
			if v, ok := fnMap["RequestMappingTemplateS3Location"].(string); ok {
				fn.RequestMappingTemplateS3Location = v
			}
			if v, ok := fnMap["ResponseMappingTemplateS3Location"].(string); ok {
				fn.ResponseMappingTemplateS3Location = v
			}
			if v, ok := fnMap["Runtime"].(map[string]interface{}); ok {
				fn.Runtime = v
			}
			if v, ok := fnMap["Code"].(string); ok {
				fn.Code = v
			}
			if v, ok := fnMap["CodeS3Location"].(string); ok {
				fn.CodeS3Location = v
			}
			if v, ok := fnMap["SyncConfig"].(map[string]interface{}); ok {
				fn.SyncConfig = v
			}
			if v, ok := fnMap["MaxBatchSize"].(int); ok {
				fn.MaxBatchSize = v
			} else if v, ok := fnMap["MaxBatchSize"].(float64); ok {
				fn.MaxBatchSize = int(v)
			}
			result[name] = fn
		}
	}
	return result
}

// parseGraphQLResolvers parses a map of resolvers.
func (t *Translator) parseGraphQLResolvers(m map[string]interface{}) map[string]sam.GraphQLApiResolver {
	result := make(map[string]sam.GraphQLApiResolver)
	for name, item := range m {
		if resMap, ok := item.(map[string]interface{}); ok {
			res := sam.GraphQLApiResolver{}
			if v, ok := resMap["FieldName"].(string); ok {
				res.FieldName = v
			}
			if v, ok := resMap["TypeName"].(string); ok {
				res.TypeName = v
			}
			if v, ok := resMap["DataSourceName"].(string); ok {
				res.DataSourceName = v
			}
			if v, ok := resMap["Kind"].(string); ok {
				res.Kind = v
			}
			if v, ok := resMap["RequestMappingTemplate"].(string); ok {
				res.RequestMappingTemplate = v
			}
			if v, ok := resMap["ResponseMappingTemplate"].(string); ok {
				res.ResponseMappingTemplate = v
			}
			if v, ok := resMap["RequestMappingTemplateS3Location"].(string); ok {
				res.RequestMappingTemplateS3Location = v
			}
			if v, ok := resMap["ResponseMappingTemplateS3Location"].(string); ok {
				res.ResponseMappingTemplateS3Location = v
			}
			if v, ok := resMap["PipelineConfig"].(map[string]interface{}); ok {
				res.PipelineConfig = v
			}
			if v, ok := resMap["CachingConfig"].(map[string]interface{}); ok {
				res.CachingConfig = v
			}
			if v, ok := resMap["SyncConfig"].(map[string]interface{}); ok {
				res.SyncConfig = v
			}
			if v, ok := resMap["MaxBatchSize"].(int); ok {
				res.MaxBatchSize = v
			} else if v, ok := resMap["MaxBatchSize"].(float64); ok {
				res.MaxBatchSize = int(v)
			}
			if v, ok := resMap["Runtime"].(map[string]interface{}); ok {
				res.Runtime = v
			}
			if v, ok := resMap["Code"].(string); ok {
				res.Code = v
			}
			if v, ok := resMap["CodeS3Location"].(string); ok {
				res.CodeS3Location = v
			}
			result[name] = res
		}
	}
	return result
}

// parseGraphQLApiKeys parses a list of API keys.
func (t *Translator) parseGraphQLApiKeys(arr []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(arr))
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			result = append(result, m)
		}
	}
	return result
}

// parseConnector parses properties into a Connector struct.
func (t *Translator) parseConnector(props map[string]interface{}) (*sam.Connector, error) {
	conn := &sam.Connector{}

	if props == nil {
		return nil, fmt.Errorf("connector properties cannot be nil")
	}

	if v, ok := props["Source"].(map[string]interface{}); ok {
		conn.Source = t.parseConnectorEndpoint(v)
	}
	if v, ok := props["Destination"].(map[string]interface{}); ok {
		conn.Destination = t.parseConnectorEndpoint(v)
	}
	if v, ok := props["Permissions"].([]interface{}); ok {
		conn.Permissions = convertToStringSlice(v)
	}

	return conn, nil
}

// parseConnectorEndpoint parses a map into a ConnectorEndpoint struct.
func (t *Translator) parseConnectorEndpoint(m map[string]interface{}) sam.ConnectorEndpoint {
	endpoint := sam.ConnectorEndpoint{}
	if v, ok := m["Id"].(string); ok {
		endpoint.ID = v
	}
	if v, ok := m["Type"].(string); ok {
		endpoint.Type = v
	}
	if v, ok := m["Arn"]; ok {
		endpoint.Arn = v
	}
	if v, ok := m["RoleName"]; ok {
		endpoint.RoleName = v
	}
	if v, ok := m["QueueUrl"]; ok {
		endpoint.QueueUrl = v
	}
	if v, ok := m["Name"]; ok {
		endpoint.Name = v
	}
	if v, ok := m["ResourceId"]; ok {
		endpoint.ResourceId = v
	}
	if v, ok := m["Qualifier"]; ok {
		endpoint.Qualifier = v
	}
	return endpoint
}

// Helper functions for type conversion

func convertToStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		if s, ok := v.(string); ok {
			result[k] = s
		} else {
			// Convert other types to JSON string
			if bytes, err := json.Marshal(v); err == nil {
				result[k] = string(bytes)
			}
		}
	}
	return result
}

func convertToStringSlice(arr []interface{}) []string {
	result := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func convertToMapSlice(arr []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(arr))
	for _, v := range arr {
		if m, ok := v.(map[string]interface{}); ok {
			result = append(result, m)
		}
	}
	return result
}
