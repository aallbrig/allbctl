package languages

import (
	"testing"
)

func TestLanguageForFile(t *testing.T) {
	cases := []struct {
		path     string
		expected string
	}{
		// Extension-based detection
		{"main.go", "Go"},
		{"src/app.py", "Python"},
		{"index.js", "JavaScript"},
		{"app.tsx", "TypeScript"},
		{"Main.java", "Java"},
		{"lib.rs", "Rust"},
		{"script.sh", "Shell"},
		{"style.css", "CSS"},
		{"page.html", "HTML"},
		{"config.yml", "YAML"},
		{"data.json", "JSON"},
		{"main.tf", "HCL"},
		{"player.gd", "GDScript"},
		{"shader.glsl", "GLSL"},
		{"query.sql", "SQL"},
		{"doc.md", "Markdown"},
		{"schema.proto", "Protocol Buffers"},

		// Filename-based detection
		{"Dockerfile", "Dockerfile"},
		{"Makefile", "Makefile"},
		{"GNUmakefile", "Makefile"},
		{"Vagrantfile", "Ruby"},
		{"CMakeLists.txt", "CMake"},

		// Filename takes priority over extension
		{"src/Makefile", "Makefile"},

		// Unrecognized files
		{"README", ""},
		{".gitignore", ""},
		{"image.png", ""},
		{"font.woff2", ""},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			got := LanguageForFile(tc.path)
			if got != tc.expected {
				t.Errorf("LanguageForFile(%q) = %q, want %q", tc.path, got, tc.expected)
			}
		})
	}
}

func TestIsVendored(t *testing.T) {
	cases := []struct {
		path     string
		expected bool
	}{
		{"vendor/github.com/pkg/errors/errors.go", true},
		{"node_modules/express/index.js", true},
		{"third_party/lib/code.c", true},
		{"src/vendor/internal.go", true}, // nested vendor
		{"Pods/AFNetworking/lib.m", true},
		{"src/main.go", false},
		{"cmd/root.go", false},
		{"pkg/util/helpers.go", false},
		{"vendored_code.go", false}, // file starts with "vendor" but is not a dir prefix
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			got := isVendored(tc.path)
			if got != tc.expected {
				t.Errorf("isVendored(%q) = %v, want %v", tc.path, got, tc.expected)
			}
		})
	}
}

func TestParseLsTree(t *testing.T) {
	t.Run("basic multi-language repo", func(t *testing.T) {
		input := `100644 blob abc123      1000	main.go
100644 blob abc124       500	pkg/util.go
100644 blob abc125       300	script.py
100644 blob abc126       200	README.md
`
		breakdown, err := parseLsTreeHelper(t, input)
		if err != nil {
			t.Fatal(err)
		}

		if len(breakdown) != 3 {
			t.Fatalf("Expected 3 languages, got %d: %v", len(breakdown), breakdown)
		}

		// Should be sorted by size descending
		if breakdown[0].Name != "Go" {
			t.Errorf("Expected first language to be Go, got %s", breakdown[0].Name)
		}
		if breakdown[0].Size != 1500 {
			t.Errorf("Expected Go size 1500, got %d", breakdown[0].Size)
		}
		if breakdown[1].Name != "Python" {
			t.Errorf("Expected second language to be Python, got %s", breakdown[1].Name)
		}
		if breakdown[2].Name != "Markdown" {
			t.Errorf("Expected third language to be Markdown, got %s", breakdown[2].Name)
		}
	})

	t.Run("skips vendored files", func(t *testing.T) {
		input := `100644 blob abc123      1000	main.go
100644 blob abc124      5000	vendor/github.com/pkg/errors/errors.go
100644 blob abc125      3000	node_modules/express/index.js
`
		breakdown, err := parseLsTreeHelper(t, input)
		if err != nil {
			t.Fatal(err)
		}

		if len(breakdown) != 1 {
			t.Fatalf("Expected 1 language (vendor excluded), got %d: %v", len(breakdown), breakdown)
		}
		if breakdown[0].Name != "Go" {
			t.Errorf("Expected Go, got %s", breakdown[0].Name)
		}
		if breakdown[0].Size != 1000 {
			t.Errorf("Expected 1000, got %d", breakdown[0].Size)
		}
	})

	t.Run("skips unrecognized extensions", func(t *testing.T) {
		input := `100644 blob abc123      1000	main.go
100644 blob abc124       500	image.png
100644 blob abc125       200	.gitignore
`
		breakdown, err := parseLsTreeHelper(t, input)
		if err != nil {
			t.Fatal(err)
		}

		if len(breakdown) != 1 {
			t.Fatalf("Expected 1 language, got %d", len(breakdown))
		}
		if breakdown[0].Percent != 100 {
			t.Errorf("Expected 100%% for sole language, got %d%%", breakdown[0].Percent)
		}
	})

	t.Run("skips submodule entries", func(t *testing.T) {
		input := `100644 blob abc123      1000	main.go
160000 commit def456         -	submodule
`
		breakdown, err := parseLsTreeHelper(t, input)
		if err != nil {
			t.Fatal(err)
		}

		if len(breakdown) != 1 {
			t.Fatalf("Expected 1 language, got %d", len(breakdown))
		}
	})

	t.Run("empty input returns nil", func(t *testing.T) {
		breakdown, err := ParseLsTree("")
		if err != nil {
			t.Fatal(err)
		}
		if breakdown != nil {
			t.Errorf("Expected nil for empty input, got %v", breakdown)
		}
	})

	t.Run("percentage calculation", func(t *testing.T) {
		input := `100644 blob abc123       750	main.go
100644 blob abc124       250	script.py
`
		breakdown, err := parseLsTreeHelper(t, input)
		if err != nil {
			t.Fatal(err)
		}

		if len(breakdown) != 2 {
			t.Fatalf("Expected 2 languages, got %d", len(breakdown))
		}
		if breakdown[0].Percent != 75 {
			t.Errorf("Expected Go at 75%%, got %d%%", breakdown[0].Percent)
		}
		if breakdown[1].Percent != 25 {
			t.Errorf("Expected Python at 25%%, got %d%%", breakdown[1].Percent)
		}
	})

	t.Run("handles Dockerfile and Makefile", func(t *testing.T) {
		input := `100644 blob abc123       100	Dockerfile
100644 blob abc124       200	Makefile
100644 blob abc125       300	main.go
`
		breakdown, err := parseLsTreeHelper(t, input)
		if err != nil {
			t.Fatal(err)
		}

		if len(breakdown) != 3 {
			t.Fatalf("Expected 3 languages, got %d: %v", len(breakdown), breakdown)
		}

		langNames := make(map[string]bool)
		for _, b := range breakdown {
			langNames[b.Name] = true
		}
		if !langNames["Dockerfile"] {
			t.Error("Expected Dockerfile language")
		}
		if !langNames["Makefile"] {
			t.Error("Expected Makefile language")
		}
		if !langNames["Go"] {
			t.Error("Expected Go language")
		}
	})
}

func TestFormatBreakdown(t *testing.T) {
	cases := []struct {
		name      string
		breakdown []LanguageBreakdown
		expected  string
	}{
		{
			"empty",
			nil,
			"",
		},
		{
			"single language",
			[]LanguageBreakdown{{Name: "Go", Size: 1024, Percent: 100}},
			"Go: 1.0 KB (100%)",
		},
		{
			"multiple languages",
			[]LanguageBreakdown{
				{Name: "Go", Size: 1536, Percent: 75},
				{Name: "Python", Size: 512, Percent: 25},
			},
			"Go: 1.5 KB (75%) | Python: 512 bytes (25%)",
		},
		{
			"megabyte sized",
			[]LanguageBreakdown{
				{Name: "Java", Size: 2 * 1024 * 1024, Percent: 90},
				{Name: "XML", Size: 200 * 1024, Percent: 10},
			},
			"Java: 2.0 MB (90%) | XML: 200.0 KB (10%)",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := FormatBreakdown(tc.breakdown)
			if got != tc.expected {
				t.Errorf("FormatBreakdown() = %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	cases := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 bytes"},
		{512, "512 bytes"},
		{1023, "1023 bytes"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{2621440, "2.5 MB"},
	}

	for _, tc := range cases {
		t.Run(tc.expected, func(t *testing.T) {
			got := formatBytes(tc.bytes)
			if got != tc.expected {
				t.Errorf("formatBytes(%d) = %q, want %q", tc.bytes, got, tc.expected)
			}
		})
	}
}

func TestBuildBreakdown(t *testing.T) {
	t.Run("empty map returns nil", func(t *testing.T) {
		result := buildBreakdown(map[string]int64{})
		if result != nil {
			t.Errorf("Expected nil, got %v", result)
		}
	})

	t.Run("sorts by size descending", func(t *testing.T) {
		input := map[string]int64{
			"Python": 100,
			"Go":     300,
			"Rust":   200,
		}
		result := buildBreakdown(input)
		if len(result) != 3 {
			t.Fatalf("Expected 3, got %d", len(result))
		}
		if result[0].Name != "Go" {
			t.Errorf("Expected Go first, got %s", result[0].Name)
		}
		if result[1].Name != "Rust" {
			t.Errorf("Expected Rust second, got %s", result[1].Name)
		}
		if result[2].Name != "Python" {
			t.Errorf("Expected Python third, got %s", result[2].Name)
		}
	})

	t.Run("percentages floor correctly", func(t *testing.T) {
		input := map[string]int64{
			"Go":     333,
			"Python": 333,
			"Rust":   334,
		}
		result := buildBreakdown(input)
		totalPercent := 0
		for _, b := range result {
			totalPercent += b.Percent
			if b.Percent != 33 {
				// 333/1000 = 33.3% → floor = 33, 334/1000 = 33.4% → floor = 33
				t.Errorf("Expected 33%% for %s, got %d%%", b.Name, b.Percent)
			}
		}
	})
}

// parseLsTreeHelper is a test helper that calls ParseLsTree and fails the test on error.
func parseLsTreeHelper(t *testing.T, input string) ([]LanguageBreakdown, error) {
	t.Helper()
	return ParseLsTree(input)
}
