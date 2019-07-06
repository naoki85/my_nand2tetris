package main

type SymbolTable struct {
	table                  map[string]int
	variableAddressCounter int
}

const (
	MinVariableAddress int = 0
)

func initializeSymbolTable() (SymbolTable, error) {
	symbolTable := SymbolTable{
		map[string]int{},
		MinVariableAddress,
	}
	return symbolTable, nil
}

func (s SymbolTable) addEntry(symbol string, address int) {
	s.table[symbol] = address
}

func (s SymbolTable) contains(symbol string) bool {
	_, ok := s.table[symbol]
	return ok
}

func (s SymbolTable) getAddress(symbol string) int {
	v, _ := s.table[symbol]
	return v
}
