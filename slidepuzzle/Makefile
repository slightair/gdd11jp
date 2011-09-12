all: Puzzle 

Test: libPuzzle.6 Test.6
	6l -o Test Test.6

Test.6: Test.go
	6g Test.go

Puzzle: libPuzzle.6 Puzzle.6
	6l -o Puzzle Puzzle.6

Puzzle.6: Puzzle.go
	6g Puzzle.go

libPuzzle.6: libPuzzle.go
	6g libPuzzle.go