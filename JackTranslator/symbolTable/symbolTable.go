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

// VarCount returns the number of variables of a given kind
func (st *SymbolTable) VarCount(kind string) int {
	return st.indexCounters[kind]
}

// KindOf returns the kind of the named identifier
func (st *SymbolTable) KindOf(name string) string {
	if symbol, ok := st.subroutineScope[name]; ok {
		return symbol.Kind
	}
	if symbol, ok := st.classScope[name]; ok {
		return symbol.Kind
	}
	return "NONE"
}

// TypeOf returns the type of the named identifier
func (st *SymbolTable) TypeOf(name string) string {
	if symbol, ok := st.subroutineScope[name]; ok {
		return symbol.Type
	}
	if symbol, ok := st.classScope[name]; ok {
		return symbol.Type
	}
	return "NONE"
}

// IndexOf returns the index of the named identifier
func (st *SymbolTable) IndexOf(name string) int {
	if symbol, ok := st.subroutineScope[name]; ok {
		return symbol.Index
	}
	if symbol, ok := st.classScope[name]; ok {
		return symbol.Index
	}
	return -1
}

func countFields(classSymbolTable *SymbolTable) int {
	count := 0
	for _, symbol := range classSymbolTable.classScope {
		if symbol.Kind == "field" {
			count++
		}
	}
	return count
}
