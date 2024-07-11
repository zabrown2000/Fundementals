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

// StartSubroutine resets the subroutine scope
func (st *SymbolTable) StartSubroutine(isMethod bool, className string) {
	st.subroutineScope = make(map[string]Symbol)
	st.indexCounters["argument"] = 0
	st.indexCounters["local"] = 0
	if isMethod {
		st.Define("this", className, "argument")
	}
}

// Define adds a new variable to the symbol table
func (st *SymbolTable) Define(name, typ, kind string) {
	index := st.indexCounters[kind]
	symbol := Symbol{Name: name, Type: typ, Kind: kind, Index: index}
	if kind == "static" || kind == "field" {
		st.classScope[name] = symbol
	} else {
		st.subroutineScope[name] = symbol
	}
	st.indexCounters[kind]++
}
