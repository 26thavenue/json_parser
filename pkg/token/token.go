package token 

type TokenType string

type Token struct {	
	Type    TokenType
	Literal string
	Line    int
	Start   int
	End     int
}

const (
	
	Illegal TokenType = "ILLEGAL"

	
	EOF TokenType = "EOF" // End of file

	
	String TokenType = "STRING"
	Number TokenType = "NUMBER"

	
	LeftBrace    TokenType = "{"
	RightBrace   TokenType = "}"
	LeftBracket  TokenType = "["
	RightBracket TokenType = "]"
	Comma        TokenType = ","
	Colon        TokenType = ":"
	QUOTE 		 TokenType	= "'"

	
	True  TokenType = "TRUE"
	False TokenType = "FALSE"
	Null  TokenType = "NULL"
)