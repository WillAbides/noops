package main

//
//func exampleMaze(t *testing.T) *mazer.Maze {
//	t.Helper()
//	f, err := os.Open("exampleMaze.json")
//	require.NoError(t, err)
//	mz, err := mazer.ParseMaze(f)
//	require.NoError(t, err)
//	return mz
//}
//
//func Test_parseMaze(t *testing.T) {
//	f, err := os.Open("exampleMaze.json")
//	require.NoError(t, err)
//	mz, err := mazer.ParseMaze(f)
//	assert.NoError(t, err)
//	assert.Equal(t, "Maze #35 (10x10)", mz.Name)
//	fmt.Println(mz.Map)
//}
//
//func Test_printmaze(t *testing.T) {
//	mz := exampleMaze(t)
//	mp := mz.Map
//	out := make([]string, len(mp))
//	for i, ln := range mp {
//		for _, v := range ln {
//			out[i] += v
//		}
//	}
//	fmt.Println(strings.Join(out, "\n"))
//}
//
//
//
//func Test_positionMap_position(t *testing.T) {
//	mz := exampleMaze(t)
//	mp := mazer.MapPositions(mz.Map)
//	assert.Nil(t, mp.position(-1,0))
//	assert.Nil(t, mp.position(0,-1))
//	assert.Nil(t, mp.position(10,0))
//	assert.Nil(t, mp.position(0,10))
//	assert.NotNil(t, mp.position(0,9))
//	assert.NotNil(t, mp.position(9,0))
//}
//
//func Test_positionMap_updateAdjacentRoutes(t *testing.T) {
//	mz := exampleMaze(t)
//	mp := mazer.MapPositions(mz.Map)
//	base := mp.position(2,9)
//	updated := mp.updateAdjacentRoutes(base)
//	for _, p := range updated {
//		fmt.Println(p)
//	}
//}
//
//func TestSolve(t *testing.T) {
//	mz := exampleMaze(t)
//	mp := mazer.MapPositions(mz.Map)
//	base := mp.position(2,9)
//	answer := mp.solve(base)
//	assert.Equal(t, "NNNNEESEEENNNENNNWWWWWWWWSSSSE", answer)
//}
