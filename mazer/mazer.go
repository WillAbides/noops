package mazer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const APIBase = `https://api.noopschallenge.com`

type position struct {
	x              int
	y              int
	value          string
	routeFromStart string
}

type positionMap map[int]map[int]*position

func mapPositions(input [][]string) positionMap {
	ySize := len(input)
	xSize := 0
	if ySize > 0 {
		xSize = len(input[0])
	}
	output := make(positionMap, xSize)
	for x := 0; x < xSize; x++ {
		output[x] = make(map[int]*position, ySize)
		for y := 0; y < ySize; y++ {
			output[x][y] = &position{
				value: input[y][x],
				x:     x,
				y:     y,
			}
		}
	}
	return output
}

func (m positionMap) validX(x int) bool {
	return x >= 0 && x <= len(m)-1
}

func (m positionMap) validY(y int) bool {
	maxY := 0
	if m.validX(0) {
		maxY = len(m[0]) - 1
	}
	return y >= 0 && y <= maxY
}

func (m positionMap) position(x, y int) *position {
	if !m.validX(x) || !m.validY(y) {
		return nil
	}
	return m[x][y]
}

func (m positionMap) nextPosition(base *position, dir string) *position {
	switch dir {
	case "N":
		return m.position(base.x, base.y-1)
	case "E":
		return m.position(base.x+1, base.y)
	case "W":
		return m.position(base.x-1, base.y)
	case "S":
		return m.position(base.x, base.y+1)
	}
	return nil
}

func (m positionMap) solve(base *position) string {
	updated := []*position{base}
	for {
		if len(updated) == 0 {
			return ""
		}
		for _, p := range updated {
			if p.value == "B" {
				return p.routeFromStart
			}
		}
		updated = m.multiUpdate(updated...)
	}
}

func (m positionMap) multiUpdate(bases ...*position) []*position {
	var updates []*position
	for _, base := range bases {
		updates = append(updates, m.updateAdjacentRoutes(base)...)
	}
	return updates
}

func (m positionMap) updateAdjacentRoutes(base *position) []*position {
	dirs := []string{"N", "E", "S", "W"}
	updates := make([]*position, 0, 4)
	routeFromStart := base.routeFromStart
	for _, dir := range dirs {
		next := m.nextPosition(base, dir)
		if next == nil || len(next.routeFromStart) > 0 || next.value == "X" {
			continue
		}

		if m.validX(next.x) && m.validY(next.y) {
			m[next.x][next.y].routeFromStart = routeFromStart + dir
			updates = append(updates, next)
		}
	}
	return updates
}

func (m positionMap) updatePositionRoute(p *position, newRoute string) {
	if p == nil || !m.validX(p.x) || !m.validY(p.y) {
		return
	}
	m[p.x][p.y].routeFromStart = newRoute
}

type maze struct {
	Name             string
	MazePath         string
	StartingPosition [2]int
	EndingPosition   [2]int
	Map              [][]string
}

func (m maze) postSolution(directions string) (result *Result, err error) {
	u := fmt.Sprintf("%s%s", APIBase, m.MazePath)
	vals := map[string]string{
		"directions": directions,
	}
	jsonValue, err := json.Marshal(vals)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(u, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	result = &Result{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func parseMaze(r io.Reader) (*maze, error) {
	mz := &maze{}
	err := json.NewDecoder(r).Decode(&mz)
	return mz, err
}

func doMaze(url string) (result *Result, err error) {
	if !strings.HasPrefix(url, APIBase) {
		url = fmt.Sprintf("%s%s", APIBase, url)
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	mz, err := parseMaze(resp.Body)
	if err != nil {
		return nil, err
	}
	startPos := mz.StartingPosition
	mp := mapPositions(mz.Map)
	answer := mp.solve(mp.position(startPos[0], startPos[1]))
	return mz.postSolution(answer)
}

func DoRandomMaze(minSize, maxSize int) (result *Result, err error) {
	if maxSize == 0 {
		maxSize = 200
	}
	u := fmt.Sprintf("%s/mazebot/random?minSize=%d&maxSize=%d", APIBase, minSize, maxSize)
	return doMaze(u)
}

type Result struct {
	Result                 string
	Elapsed                float64
	ShortestSolutionLength int
	YourSolutionLength     int
	NextMaze               string
}

func (r *Result) String() string {
	return fmt.Sprintf(`Result                 %s
Elapsed                %f
ShortestSolutionLength %d
YourSolutionLength     %d
NextMaze               %s
`, r.Result, r.Elapsed, r.ShortestSolutionLength, r.YourSolutionLength, r.NextMaze)
}

func startRace(login string) (string, error) {
	nextMaze := ""
	raceForm, err := json.Marshal(map[string]string{
		"login": login,
	})
	if err != nil {
		return nextMaze, err
	}
	u := fmt.Sprintf("%s/mazebot/race/start", APIBase)
	resp, err := http.Post(u, "application/json", bytes.NewBuffer(raceForm))
	if err != nil {
		return nextMaze, err
	}
	got := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		return nextMaze, err
	}
	var ok bool
	nextMaze, ok = got["nextMaze"].(string)
	if !ok {
		return nextMaze, errors.New("unexpected nextmaze")
	}
	return nextMaze, nil
}

func RunRace(login string) error {
	nextMaze, err := startRace(login)
	if err != nil {
		return err
	}
	for i := 0; nextMaze != ""; i++ {
		result, err := doMaze(nextMaze)
		if err != nil {
			return err
		}
		fmt.Println("\n", result)
		nextMaze = result.NextMaze
	}
	return nil
}
