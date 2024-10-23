package diff

import (
	"strings"

	"github.com/jolovicdev/nora/internal/types"
)


func SimpleMyers(oldText, newText []string) []types.DiffStep {
	n := len(oldText)
	m := len(newText)
	max := n + m

	v := make(map[int]int)
	v[1] = 0
	trace := make([]map[int]int, 0)

	for d := 0; d <= max; d++ {
		vCopy := make(map[int]int)
		for k, v := range v {
			vCopy[k] = v
		}
		trace = append(trace, vCopy)

		for k := -d; k <= d; k += 2 {
			var x int
			if k == -d || (k != d && v[k-1] < v[k+1]) {
				x = v[k+1]
			} else {
				x = v[k-1] + 1
			}

			y := x - k

			for x < n && y < m && oldText[x] == newText[y] {
				x++
				y++
			}

			v[k] = x

			if x >= n && y >= m {
				return Backtrack(trace, oldText, newText)
			}
		}
	}

	return nil
}

func Backtrack(trace []map[int]int, oldText, newText []string) []types.DiffStep {
	steps := make([]types.DiffStep, 0)
	x := len(oldText)
	y := len(newText)
	
	for d := len(trace) - 1; d >= 0; d-- {
		v := trace[d]
		k := x - y
		
		var prevK int
		if k == -d || (k != d && v[k-1] < v[k+1]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}
		
		prevX := v[prevK]
		prevY := prevX - prevK
		
		for x > prevX && y > prevY {
			steps = append([]types.DiffStep{{
				Type:     "keep",
				Content:  oldText[x-1],
				Position: x - 1,
			}}, steps...)
			x--
			y--
		}
		
		if d > 0 {
			if x > prevX {
				steps = append([]types.DiffStep{{
					Type:     "delete",
					Content:  oldText[x-1],
					Position: x - 1,
				}}, steps...)
				x--
			} else if y > prevY {
				steps = append([]types.DiffStep{{
					Type:     "add",
					Content:  newText[y-1],
					Position: y - 1,
				}}, steps...)
				y--
			}
		}
	}
	
	return steps
}
func CalculateDiff(oldContent, newContent string) []types.DiffStep {
	oldLines := strings.Split(oldContent, "\n")
	newLines := strings.Split(newContent, "\n")
	
	return SimpleMyers(oldLines, newLines)
}