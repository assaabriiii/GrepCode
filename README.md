# Grepcode 🔍

I've wrote this out of nowhere and i was bored so don't be a genius and add brick to this code, i've just wanted to help myself finding code pieces among other code pieces and yeah here it is.

![Grepcode Demo](https://github.com/assaabriiii/GrepCode/demo.png)

## Features ✨

- 🔎 Regex pattern searching (`-regex` flag)
- 🎨 Colorized output (files, line numbers, matches)
- 📂 File extension filtering (`-ext`)
- 🚫 Directory exclusion (`-exclude-dir`)
- 📜 Context lines around matches (`-context`)
- 📊 Search statistics (`-stats`)
- 🖥 Cross-platform support
- ⚡ Blazing fast performance

## Installation ⚙️

### Prerequisites

- [Go 1.16+](https://go.dev/dl/)

### Quick Install

```bash
go install github.com/fatih/color@latest  # Dependency
go install github.com/yourusername/Grepcode@latest
```

### From Source

```bash
git clone https://github.com/assaabriiii/GrepCode
cd GrepCode
go build -o GrepCode grepCode.go
mv GrepCode /usr/local/bin/  # Or add to your $PATH
```

## Usage 🚀

```bash
Grepcode [flags] -pattern <search_pattern>
```

### Basic Examples

```bash
# Search for "func main" in Go files
Grepcode -pattern "func main" -ext go

# Case-sensitive search in JavaScript/TypeScript files
Grepcode -pattern "TODO" -case -ext js,ts

# Regex search with 2 lines of context
Grepcode -pattern "\b\d{3}-\d{3}-\d{4}\b" -regex -context 2

# Search in specific directory excluding node_modules
Grepcode -pattern "error" -dir ./src -exclude-dir node_modules,vendor
```

### All Flags
  
- Flag	Description	Default
- `-pattern`	Search pattern (required)
- `-dir`	Directory to search	Current dir
- `-case`	Case-sensitive search	false
- `-regex`	Use regular expressions	false
- `-ext`	Comma-separated file extensions (e.g., go,js)	All files
- `-context`	Show N lines of context around match	0
- `-stats`	Show search statistics	false
- `-exclude-dir`	Comma-separated directories to exclude	

## Output Formatting 🎨

- 🔵 Cyan: File paths
- 🟡 Yellow: Line numbers
- 🔴 Red: Matching lines
- 📈 Statistics: Execution time, files scanned, matches found

## Advanced Usage 🧠

### Pipe to Other Commands

```bash
Grepcode -pattern "FIXME" -ext py | sort | uniq
```

### Search Hidden Files

```bash
Grepcode -pattern "password" -dir ~/.config
```

### Time-bound Search

```bash
Grepcode -pattern "deprecated" -stats | grep "Execution time"
```

## FAQ ❓

### How do I search for special characters?

#### Use regex escaping or quotes:

```bash
Grepcode -pattern "\[WARNING\]" -regex
```

#### Can I search binary files?

Grepcode intentionally skips binary files for safety.

#### How to make it faster?

- Use specific file extensions (`-ext`)
- Exclude large directories (`-exclude-dir`)

## Contributing 🤝

Don't.

