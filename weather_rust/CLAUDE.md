# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build/Run/Test Commands

- Build: `cargo build`
- Run: `OPENWEATHER_API_KEY=your_key cargo run`
- Check: `cargo check`
- Lint: `cargo clippy -- -D warnings`
- Format: `cargo fmt --all -- --check`
- Test: `cargo test`
- Test single: `cargo test test_name`

## Code Style Guidelines

- Use 2018 edition Rust
- Error handling: Use anyhow::Result for general error returns, McpError for API errors
- Struct fields use snake_case
- Use proper typing with Serde for JSON serialization/deserialization
- Doc comments should be present for public interfaces
- Follow standard Rust formatting conventions via rustfmt
- Organize imports: std first, then external crates alphabetically
- For HTTP requests, always handle response status and provide meaningful errors
- Prefer async/await for asynchronous code

