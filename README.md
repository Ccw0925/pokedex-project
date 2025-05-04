# 🎮 Pokedex API (Go)

A high-performance Pokémon API built with Go, Gin, and Redis-style caching. Provides Pokémon data, and evolution chains by leveraging the [PokeAPI](https://pokeapi.co/).

My first Go project while learning the language - a simple Pokémon API that fetches data from [PokeAPI](https://pokeapi.co/). 

**Goal**: Get comfortable with Go basics, HTTP handlers, and caching.

![Go Version](https://img.shields.io/badge/Go-1.24+-blue)

## 🌟 Features

- **RESTful API** with Gin framework
- **Smart Caching** using `go-cache` (5-minute TTL)
- **Multi-identifier lookup** (ID or name)
- **Evolution chain** tracking

## 📦 Installation

### Prerequisites
- Go 1.21+

### Local Setup
```bash
# Clone the repository
git clone https://github.com/Ccw0925/pokedex-project.git
cd pokedex-project

# Build and run
go run main.go
```

## 🧑‍💻 Why I Built This

As a Go beginner, I wanted to:
- Practice Go syntax and package management
- Learn how to build REST APIs with Gin
- Understand caching mechanisms
- Work with external APIs
- Get familiar with testing in Go

## 🛠️ What's Inside

### Basic Features
- Get Pokémon details by ID or name
- Simple in-memory caching (5-minute expiry)
- Error handling for missing Pokémon
- Health check endpoint