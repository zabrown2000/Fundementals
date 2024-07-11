package symbolTable

/*
//do we need an interface?
constructor
startSubroutine
define(name, type, kind)
varCount(kind) int
kindOf(name) kind
typeOf(name) String
indexOf(name) int
*/

// Symbol represents a single symbol in the symbol table
type Symbol struct {
	Name  string
	Type  string
	Kind  string
	Index int
}

// SymbolTable manages the symbols for a class and subroutine scope
type SymbolTable struct {
	classScope      map[string]Symbol
	subroutineScope map[string]Symbol
	indexCounters   map[string]int
}

// NewSymbolTable creates a new symbol table - sets all indices to 0
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		classScope:      make(map[string]Symbol),
		subroutineScope: make(map[string]Symbol),
		indexCounters:   map[string]int{"static": 0, "field": 0, "argument": 0, "local": 0},
	}
}
