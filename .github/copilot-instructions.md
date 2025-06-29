# Copilot Instructions

<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

## Project Context
This is a Go benchmark testing project for the `tree-framework` from `github.com/catalinfl/tree-framework`.

## Guidelines
- Focus on performance testing and benchmarking HTTP handlers and routing
- Use Go's built-in `testing` package for benchmarks
- Follow Go naming conventions for benchmark functions (`BenchmarkXxx`)
- Include memory allocation benchmarks where relevant
- Test different payload sizes and request patterns
- Compare performance against standard library implementations when possible

## Code Style
- Use idiomatic Go code
- Include proper error handling
- Add meaningful comments for complex benchmark scenarios
- Use table-driven tests for multiple test cases
- Measure both CPU and memory performance
