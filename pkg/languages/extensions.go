package languages

// extensionToLanguage maps file extensions (with leading dot) to language names.
// Modeled after GitHub's linguist for the most common programming languages.
var extensionToLanguage = map[string]string{
	// Go
	".go": "Go",

	// Python
	".py":  "Python",
	".pyi": "Python",
	".pyx": "Python",

	// JavaScript
	".js":  "JavaScript",
	".mjs": "JavaScript",
	".cjs": "JavaScript",
	".jsx": "JavaScript",

	// TypeScript
	".ts":  "TypeScript",
	".tsx": "TypeScript",
	".mts": "TypeScript",
	".cts": "TypeScript",

	// Java
	".java": "Java",

	// Kotlin
	".kt":  "Kotlin",
	".kts": "Kotlin",

	// C
	".c": "C",
	".h": "C",

	// C++
	".cpp": "C++",
	".cxx": "C++",
	".cc":  "C++",
	".hpp": "C++",
	".hxx": "C++",
	".hh":  "C++",

	// C#
	".cs": "C#",

	// Rust
	".rs": "Rust",

	// Ruby
	".rb":  "Ruby",
	".erb": "Ruby",

	// PHP
	".php": "PHP",

	// Swift
	".swift": "Swift",

	// Objective-C
	".m":  "Objective-C",
	".mm": "Objective-C",

	// Scala
	".scala": "Scala",
	".sc":    "Scala",

	// Perl
	".pl": "Perl",
	".pm": "Perl",

	// Lua
	".lua": "Lua",

	// R
	".r": "R",
	".R": "R",

	// Julia
	".jl": "Julia",

	// Haskell
	".hs":  "Haskell",
	".lhs": "Haskell",

	// Elixir
	".ex":  "Elixir",
	".exs": "Elixir",

	// Erlang
	".erl": "Erlang",
	".hrl": "Erlang",

	// Clojure
	".clj":  "Clojure",
	".cljs": "Clojure",
	".cljc": "Clojure",

	// Dart
	".dart": "Dart",

	// Shell
	".sh":   "Shell",
	".bash": "Shell",
	".zsh":  "Shell",
	".fish": "Shell",

	// PowerShell
	".ps1":  "PowerShell",
	".psm1": "PowerShell",
	".psd1": "PowerShell",

	// Batch
	".bat": "Batch",
	".cmd": "Batch",

	// HTML
	".html": "HTML",
	".htm":  "HTML",

	// CSS
	".css": "CSS",

	// SCSS/Sass
	".scss": "SCSS",
	".sass": "SCSS",

	// Less
	".less": "Less",

	// Vue
	".vue": "Vue",

	// Svelte
	".svelte": "Svelte",

	// SQL
	".sql": "SQL",

	// Dockerfile
	// (handled specially by filename, see detect.go)

	// Makefile handled by filename in detect.go

	// YAML
	".yml":  "YAML",
	".yaml": "YAML",

	// JSON
	".json": "JSON",

	// TOML
	".toml": "TOML",

	// XML
	".xml":  "XML",
	".xsl":  "XML",
	".xslt": "XML",

	// Markdown
	".md":       "Markdown",
	".markdown": "Markdown",

	// reStructuredText
	".rst": "reStructuredText",

	// Protocol Buffers
	".proto": "Protocol Buffers",

	// Terraform
	".tf":     "HCL",
	".tfvars": "HCL",
	".hcl":    "HCL",

	// Nix
	".nix": "Nix",

	// Zig
	".zig": "Zig",

	// V
	".v": "V",

	// Nim
	".nim": "Nim",

	// OCaml
	".ml":  "OCaml",
	".mli": "OCaml",

	// F#
	".fs":  "F#",
	".fsx": "F#",
	".fsi": "F#",

	// Groovy
	".groovy": "Groovy",
	".gradle": "Groovy",

	// GDScript (Godot)
	".gd": "GDScript",

	// GLSL / Shaders
	".glsl":     "GLSL",
	".vert":     "GLSL",
	".frag":     "GLSL",
	".shader":   "GLSL",
	".gdshader": "GLSL",

	// Assembly
	".asm": "Assembly",
	".s":   "Assembly",
	".S":   "Assembly",

	// Fortran
	".f":   "Fortran",
	".f90": "Fortran",
	".f95": "Fortran",

	// COBOL
	".cob": "COBOL",
	".cbl": "COBOL",

	// Verilog / SystemVerilog
	".sv": "SystemVerilog",

	// VHDL
	".vhd":  "VHDL",
	".vhdl": "VHDL",
}

// filenameToLanguage maps exact filenames (without path) to language names.
var filenameToLanguage = map[string]string{
	"Dockerfile":     "Dockerfile",
	"Makefile":       "Makefile",
	"GNUmakefile":    "Makefile",
	"Vagrantfile":    "Ruby",
	"Gemfile":        "Ruby",
	"Rakefile":       "Ruby",
	"CMakeLists.txt": "CMake",
	"Justfile":       "Just",
	"justfile":       "Just",
}

// LanguageForFile returns the language name for a given file path, or empty
// string if the file type is not recognized. It checks exact filename first,
// then falls back to extension matching.
func LanguageForFile(filePath string) string {
	base := fileBase(filePath)

	if lang, ok := filenameToLanguage[base]; ok {
		return lang
	}

	ext := fileExt(filePath)
	if lang, ok := extensionToLanguage[ext]; ok {
		return lang
	}

	return ""
}
