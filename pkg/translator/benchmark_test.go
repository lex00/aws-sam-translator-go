package translator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// BenchmarkTransformBasicFunction benchmarks transformation of a basic function.
func BenchmarkTransformBasicFunction(b *testing.B) {
	input, err := os.ReadFile(filepath.Join(testdataInputDir, "basic_function.yaml"))
	if err != nil {
		b.Fatalf("Failed to read input file: %v", err)
	}

	tr := New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := tr.TransformBytes(input)
		if err != nil {
			b.Fatalf("Transform failed: %v", err)
		}
	}
}

// BenchmarkTransformFunctionWithEvents benchmarks transformation with multiple events.
func BenchmarkTransformFunctionWithEvents(b *testing.B) {
	// Find a fixture with multiple events
	input, err := os.ReadFile(filepath.Join(testdataInputDir, "function_event_conditions.yaml"))
	if err != nil {
		b.Skipf("Fixture not found: %v", err)
	}

	tr := New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := tr.TransformBytes(input)
		if err != nil {
			b.Fatalf("Transform failed: %v", err)
		}
	}
}

// BenchmarkTransformAPI benchmarks API transformation.
func BenchmarkTransformAPI(b *testing.B) {
	input, err := os.ReadFile(filepath.Join(testdataInputDir, "api_with_cors.yaml"))
	if err != nil {
		b.Skipf("Fixture not found: %v", err)
	}

	tr := New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := tr.TransformBytes(input)
		if err != nil {
			b.Fatalf("Transform failed: %v", err)
		}
	}
}

// BenchmarkTransformStateMachine benchmarks StateMachine transformation.
func BenchmarkTransformStateMachine(b *testing.B) {
	input, err := os.ReadFile(filepath.Join(testdataInputDir, "state_machine_with_definition_substitutions.yaml"))
	if err != nil {
		b.Skipf("Fixture not found: %v", err)
	}

	tr := New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := tr.TransformBytes(input)
		if err != nil {
			b.Fatalf("Transform failed: %v", err)
		}
	}
}

// BenchmarkTransformConnector benchmarks Connector transformation.
func BenchmarkTransformConnector(b *testing.B) {
	input, err := os.ReadFile(filepath.Join(testdataInputDir, "connector_function_to_table.yaml"))
	if err != nil {
		b.Skipf("Fixture not found: %v", err)
	}

	tr := New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := tr.TransformBytes(input)
		if err != nil {
			b.Fatalf("Transform failed: %v", err)
		}
	}
}

// BenchmarkTransformAllPolicyTemplates benchmarks the most complex fixture.
func BenchmarkTransformAllPolicyTemplates(b *testing.B) {
	input, err := os.ReadFile(filepath.Join(testdataInputDir, "all_policy_templates.yaml"))
	if err != nil {
		b.Skipf("Fixture not found: %v", err)
	}

	tr := New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := tr.TransformBytes(input)
		if err != nil {
			b.Fatalf("Transform failed: %v", err)
		}
	}
}

// BenchmarkTransformAllFixtures benchmarks transformation of all success fixtures.
func BenchmarkTransformAllFixtures(b *testing.B) {
	inputFiles, err := filepath.Glob(filepath.Join(testdataInputDir, "*.yaml"))
	if err != nil {
		b.Fatalf("Failed to glob input files: %v", err)
	}

	// Filter out error fixtures
	var successFiles []string
	for _, f := range inputFiles {
		if !strings.HasPrefix(filepath.Base(f), "error_") {
			successFiles = append(successFiles, f)
		}
	}

	// Read all inputs
	inputs := make([][]byte, len(successFiles))
	for i, path := range successFiles {
		input, err := os.ReadFile(path)
		if err != nil {
			b.Fatalf("Failed to read %s: %v", path, err)
		}
		inputs[i] = input
	}

	b.Logf("Benchmarking %d fixtures", len(inputs))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tr := New()
		for _, input := range inputs {
			_, _ = tr.TransformBytes(input)
		}
	}
}

// BenchmarkNewTranslator benchmarks translator creation.
func BenchmarkNewTranslator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New()
	}
}

// BenchmarkTranslatorReuse tests if reusing translator is beneficial.
func BenchmarkTranslatorReuse(b *testing.B) {
	input, err := os.ReadFile(filepath.Join(testdataInputDir, "basic_function.yaml"))
	if err != nil {
		b.Fatalf("Failed to read input file: %v", err)
	}

	b.Run("new-each-time", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tr := New()
			_, err := tr.TransformBytes(input)
			if err != nil {
				b.Fatalf("Transform failed: %v", err)
			}
		}
	})

	b.Run("reuse-translator", func(b *testing.B) {
		tr := New()
		for i := 0; i < b.N; i++ {
			_, err := tr.TransformBytes(input)
			if err != nil {
				b.Fatalf("Transform failed: %v", err)
			}
		}
	})
}

// BenchmarkParallelTransform tests parallel transformation performance.
func BenchmarkParallelTransform(b *testing.B) {
	input, err := os.ReadFile(filepath.Join(testdataInputDir, "basic_function.yaml"))
	if err != nil {
		b.Fatalf("Failed to read input file: %v", err)
	}

	b.RunParallel(func(pb *testing.PB) {
		tr := New()
		for pb.Next() {
			_, err := tr.TransformBytes(input)
			if err != nil {
				b.Fatalf("Transform failed: %v", err)
			}
		}
	})
}
