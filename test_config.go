package main

import (
	"fmt"
	"log"
	"github.com/goplus/llpkgstore/config"
)

func main() {
	fmt.Println("Testing configuration validation...")
	
	// Test valid configuration
	fmt.Println("\n1. Testing valid configuration:")
	validCfg, err := config.ParseLLPkgConfig("test_llpkg.cfg")
	if err != nil {
		log.Printf("Error parsing valid config: %v", err)
	} else {
		fmt.Printf("Valid config parsed successfully: %+v\n", validCfg)
		
		// Test validation
		if err := validCfg.Llpyg.Validate(); err != nil {
			fmt.Printf("Validation error: %v\n", err)
		} else {
			fmt.Println("Configuration validation passed!")
		}
	}
	
	// Test invalid configuration
	fmt.Println("\n2. Testing invalid configuration:")
	invalidCfg, err := config.ParseLLPkgConfig("invalid_config.cfg")
	if err != nil {
		log.Printf("Error parsing invalid config: %v", err)
	} else {
		fmt.Printf("Invalid config parsed successfully: %+v\n", invalidCfg)
		
		// Test validation
		if err := invalidCfg.Llpyg.Validate(); err != nil {
			fmt.Printf("Validation error (expected): %v\n", err)
		} else {
			fmt.Println("Configuration validation passed (unexpected)!")
		}
	}
}
