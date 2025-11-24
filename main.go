package main

import (
    "os"
    
    "github.com/axellelanca/urlshortener/cmd"
    _ "github.com/axellelanca/urlshortener/cmd/cli"
    _ "github.com/axellelanca/urlshortener/cmd/server"
)

func main() {
    if err := cmd.Execute(); err != nil {
        os.Exit(1)
    }
}