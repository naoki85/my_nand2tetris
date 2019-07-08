package main

type SymbolTable struct {
	table                  map[string]int
	variableAddressCounter int
}

const (
	MinVariableAddress int = 0
)

func InitializeSymbolTable() (SymbolTable, error) {
	symbolTable := SymbolTable{
		map[string]int{},
		MinVariableAddress,
	}
	return symbolTable, nil
}

func (s SymbolTable) AddEntry(symbol string, address int) {
	s.table[symbol] = address
}

func (s SymbolTable) Contains(symbol string) bool {
	_, ok := s.table[symbol]
	return ok
}

func (s SymbolTable) GetAddress(symbol string) int {
	v, _ := s.table[symbol]
	return v
}
