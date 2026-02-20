# AI Resource Compiler (Go)

Compile [AI Resource Specification](https://github.com/jomadu/ai-resource-spec) resources to tool-specific formats.

## Overview

`ai-resource-compiler-go` transforms validated AI resources into formats consumed by AI coding tools:

- **Kiro CLI** - AWS AI assistant format
- **Cursor** - `.cursorrules` format
- **Claude Code** - Project instructions
- **GitHub Copilot** - Workspace instructions

Built on [ai-resource-core-go](https://github.com/jomadu/ai-resource-core-go) for parsing and validation.

## Installation

bash
go get github.com/jomadu/ai-resource-
compiler-go

## Usage

go
import (
   "github.com/jomadu/ai-resource-core-go/
pkg/airesource"
   "github.com/jomadu/ai-resource-compiler
-go/cursor"
)

// Load resource
prompt,  := airesource.LoadPrompt(
"prompt.yml")

// Compile to Cursor format
output := cursor.Compile(prompt)

// Write to .cursorrules
os.WriteFile(".cursorrules", []byte(output)
, 0644)

## CLI

bash
# Install
go install github.com/jomadu/ai-resource-
compiler-go/cmd/arc@latest

# Compile to Cursor
arc compile --target cursor prompts.yml

# Compile to multiple targets
arc compile --target 
kiro,cursor,claude,copilot prompts.yml

## Supported Targets

- `kiro` - Kiro CLI context format
- `cursor` - Cursor IDE rules
- `claude` - Claude Code project instructions
- `copilot` - GitHub Copilot instructions

## Architecture

Core handles validation. Compiler handles transformation. Clean separation of concerns.

## License

GNU General Public License v2.0