package bf

import (
	"fmt"
	"os"
	"testing"
)

func TestCNF(t *testing.T) {
	f := And(Or(Var("a"), Var("b")), Var("i"), Or(Var("g"), Var("h"), And(Var("c"), Or(Var("d"), Var("e")), Var("f"))))
	sat, _, err := Solve(f)
	if err != nil {
		t.Errorf("could not solve problem %s: %v", f, err)
	}
	if !sat {
		t.Errorf("problem was declared UNSAT")
	}
}

func TestString(t *testing.T) {
	f := And(Or(Var("a"), Not(Var("b"))), Not(Var("c")))
	const expected = "and(or(a, not(b)), not(c))"
	if f.String() != expected {
		t.Errorf("string representation of formula not as expected: wanted %q, got %q", expected, f.String())
	}
}

func ExampleSolve() {
	f := Not(Implies(
		And(Var("a"), Var("b")), And(Or(Var("c"), Not(Var("d"))),
			Not(And(Var("c"), Eq(Var("e"), Not(Var("c"))))), Not(Xor(Var("a"), Var("b"))))))
	sat, _, err := Solve(f)
	if err != nil {
		fmt.Printf("Error while solving problem: %v", err)
	} else if sat {
		fmt.Printf("Problem is satisfiable")
	} else {
		fmt.Printf("Problem is unsatisfiable")
	}
	// Output: Problem is satisfiable
}

func ExampleDimacs() {
	f := Eq(And(Or(Var("a"), Not(Var("b"))), Not(Var("a"))), Var("b"))
	if err := Dimacs(f, os.Stdout); err != nil {
		fmt.Printf("Could not generate DIMACS file: %v", err)
	}
	// Output:
	// p cnf 4 6
	// c a=2
	// c b=3
	// -2 -1 0
	// 3 -1 0
	// 1 2 3 0
	// 2 -3 -4 0
	// -2 -4 0
	// 4 -3 0
}

func ExampleSolve_sudoku() {
	const varFmt = "line-%d-col-%d:%d" // Scheme for variable naming
	f := And()
	// In each spot, exactly one number is written
	for line := 1; line <= 9; line++ {
		for col := 1; col <= 9; col++ {
			vars := make([]Formula, 9)
			for val := 1; val <= 9; val++ {
				vars[val-1] = Var(fmt.Sprintf(varFmt, line, col, val))
			}
			f = And(f, Or(vars...))
			for val1 := 1; val1 <= 8; val1++ {
				var1 := Var(fmt.Sprintf(varFmt, line, col, val1))
				for val2 := val1 + 1; val2 <= 9; val2++ {
					var2 := Var(fmt.Sprintf(varFmt, line, col, val2))
					f = And(f, Or(Not(var1), Not(var2)))
				}
			}
		}
	}
	// In each line, each number appears at least once.
	// Since there are 9 spots and 9 numbers, that means each number appears exactly once.
	for line := 1; line <= 9; line++ {
		for val := 1; val <= 9; val++ {
			var vars []Formula
			for col := 1; col <= 9; col++ {
				vars = append(vars, Var(fmt.Sprintf(varFmt, line, col, val)))
			}
			f = And(f, Or(vars...))
		}
	}
	// In each column, each number appears at least once.
	for col := 1; col <= 9; col++ {
		for val := 1; val <= 9; val++ {
			var vars []Formula
			for line := 1; line <= 9; line++ {
				vars = append(vars, Var(fmt.Sprintf(varFmt, line, col, val)))
			}
			f = And(f, Or(vars...))
		}
	}
	// In each 3x3 box, each number appears at least once.
	for lineB := 0; lineB < 3; lineB++ {
		for colB := 0; colB < 3; colB++ {
			for val := 1; val <= 9; val++ {
				var vars []Formula
				for lineOff := 1; lineOff <= 3; lineOff++ {
					line := lineB*3 + lineOff
					for colOff := 1; colOff <= 3; colOff++ {
						col := colB*3 + colOff
						vars = append(vars, Var(fmt.Sprintf(varFmt, line, col, val)))
					}
				}
				f = And(f, Or(vars...))
			}
		}
	}
	// Some spots already have a fixed value
	f = And(
		f,
		Var("line-1-col-1:5"),
		Var("line-1-col-2:3"),
		Var("line-1-col-5:7"),
		Var("line-2-col-1:6"),
		Var("line-2-col-4:1"),
		Var("line-2-col-5:9"),
		Var("line-2-col-6:5"),
		Var("line-3-col-2:9"),
		Var("line-3-col-3:8"),
		Var("line-3-col-8:6"),
		Var("line-4-col-1:8"),
		Var("line-4-col-5:6"),
		Var("line-4-col-9:3"),
		Var("line-5-col-1:4"),
		Var("line-5-col-4:8"),
		Var("line-5-col-6:3"),
		Var("line-5-col-9:1"),
		Var("line-6-col-1:7"),
		Var("line-6-col-5:2"),
		Var("line-6-col-9:6"),
		Var("line-7-col-2:6"),
		Var("line-7-col-7:2"),
		Var("line-7-col-8:8"),
		Var("line-8-col-4:4"),
		Var("line-8-col-5:1"),
		Var("line-8-col-6:9"),
		Var("line-8-col-9:5"),
		Var("line-9-col-5:8"),
		Var("line-9-col-8:7"),
		Var("line-9-col-9:9"),
	)
	sat, model, err := Solve(f)
	if err != nil {
		fmt.Printf("Error while solving grid: %v\n", err)
		return
	} else if !sat {
		fmt.Println("Error: solving grid was found unsat")
		return
	}
	fmt.Println("The grid has a solution")
	for line := 1; line <= 9; line++ {
		for col := 1; col <= 9; col++ {
			for val := 1; val <= 9; val++ {
				if model[fmt.Sprintf(varFmt, line, col, val)] {
					fmt.Printf("%d", val)
				}
			}
		}
		fmt.Println()
	}
	// Output:
	// The grid has a solution
	// 534678912
	// 672195348
	// 198342567
	// 859761423
	// 426853791
	// 713924856
	// 961537284
	// 287419635
	// 345286179
}