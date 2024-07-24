package parser

import (
	"fmt"
	"strconv"

	"github.com/26thavenue/json_parser/pkg/lexer"
	"github.com/26thavenue/json_parser/pkg/token"
)

type Parser struct {
	l     *lexer.Lexer
	token token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.token = p.l.NextToken()
}

func (p *Parser) Parse() (interface{}, error) {
	result := p.parseValue()
	if p.token.Type != token.EOF {
		return nil, fmt.Errorf("unexpected token: %s", p.token.Literal)
	}
	return result, nil
}

func (p *Parser) parseValue() interface{} {
	switch p.token.Type {
	case token.LeftBrace:
		return p.parseObject()
	case token.LeftBracket:
		return p.parseArray()
	case token.String:
		value := p.token.Literal
		p.nextToken()
		return value
	case token.Number:
		num, err := strconv.ParseFloat(p.token.Literal, 64)
		if err != nil {
			panic(fmt.Errorf("invalid number: %s", p.token.Literal))
		}
		p.nextToken()
		return num
	case token.True:
		p.nextToken()
		return true
	case token.False:
		p.nextToken()
		return false
	case token.Null:
		p.nextToken()
		return nil
	default:
		panic(fmt.Errorf("unexpected token: %s", p.token.Literal))
	}
}

func (p *Parser) parseObject() map[string]interface{} {
	obj := make(map[string]interface{})
	p.nextToken() // consume '{'

	for p.token.Type != token.RightBrace {
		if p.token.Type != token.String {
			panic(fmt.Errorf("expected string key, got %s", p.token.Literal))
		}
		key := p.token.Literal
		p.nextToken()

		if p.token.Type != token.Colon {
			panic(fmt.Errorf("expected ':', got %s", p.token.Literal))
		}
		p.nextToken()

		value := p.parseValue()
		obj[key] = value

		if p.token.Type == token.Comma {
			p.nextToken()
		} else if p.token.Type != token.RightBrace {
			panic(fmt.Errorf("expected ',' or '}', got %s", p.token.Literal))
		}
	}

	p.nextToken() // consume '}'
	return obj
}

func (p *Parser) parseArray() []interface{} {
	var arr []interface{}
	p.nextToken() // consume '['

	for p.token.Type != token.RightBracket {
		value := p.parseValue()
		arr = append(arr, value)

		if p.token.Type == token.Comma {
			p.nextToken()
		} else if p.token.Type != token.RightBracket {
			panic(fmt.Errorf("expected ',' or ']', got %s", p.token.Literal))
		}
	}

	p.nextToken() // consume ']'
	return arr
}
