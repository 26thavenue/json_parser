package main

import (
	"fmt"

	"github.com/26thavenue/json_parser/pkg/lexer"
	"github.com/26thavenue/json_parser/pkg/parser"
)

func main() {
	input := `{"name": "John", "age": 30, "city": null, "hobbies": ["reading", "swimming"]}`
	l := lexer.New(input)
	p := parser.New(l)

	result, err := p.Parse()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	
	if jsonObject, ok := result.(map[string]interface{}); ok {
		fmt.Printf("%+v\n", jsonObject)
	}
}