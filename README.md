
# Pokedex CLI

A command-line interface (CLI) Pokedex application written in Go. This project allows you to explore Pokemon locations, catch Pokemon, and inspect your Pokedex, all from the comfort of your terminal.

## Features

*   **Explore Locations:** Discover different location areas and the Pokemon that inhabit them.
*   **Catch Pokemon:** Attempt to catch wild Pokemon and add them to your Pokedex.
*   **Inspect Pokedex:** View detailed information about the Pokemon you've caught.
*   **Command-Line Interface:** Interact with the Pokedex through simple and intuitive commands.
*   **Caching:** Implements caching to reduce API calls and improve performance.

## Getting Started

### Prerequisites

*   [Go](https://go.dev/dl/) installed on your system (version 1.20 or higher recommended).
*   A terminal or command prompt.

### Installation

1.  **Clone the repository:**

    ```bash
    git clone <your_repository_url>
    cd pokedexcli
    ```

2.  **Build the application:**

    ```bash
    go build -o pokedex .
    ```

### Usage

Run the application:

```bash
./pokedex
