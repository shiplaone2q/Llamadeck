# Contributing to Llamadeck

Thank you for considering contributing to Llamadeck! This document outlines the development setup and contribution process.

## Development Setup

### Prerequisites

- Go 1.24 or later
- Node.js 22 or later
- `llama-server` executable (from [llama.cpp](https://github.com/ggml-org/llama.cpp))

### Getting Started

1. **Clone the repository**
   ```bash
   git clone https://github.com/lordmathis/llamadeck.git
   cd llamadeck
   ```

2. **Install dependencies**
   ```bash
   # Go dependencies
   go mod download
   
   # Frontend dependencies
   cd webui && npm ci && cd ..
   ```

3. **Run for development**
   ```bash
   # Start backend server
   go run ./cmd/server

   # Or with automatic rebuild on save (https://github.com/air-verse/air)
   air
