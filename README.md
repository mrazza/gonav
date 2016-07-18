# gonav
A [Source Engine](https://en.wikipedia.org/wiki/Source_(game_engine)) bot Nav file parser written in Go. This was written for CS:GO but will likely work with little or no modification for other Source titles. The specifics of the .nav format were reverse engineered using the information on [Valve's wiki](https://developer.valvesoftware.com/wiki/NAV) as a starting point. For more information on Source's Navigation Meshes see Valve's wiki: https://developer.valvesoftware.com/wiki/Navigation_Meshes

# Usage
See the `_examples` folder for examples of how to use this library. The basics, however, are pretty straightforward. Create a parser and pass in any `Reader` that contains the binary .nav data, then call the `Parse()` method. This will output a parsed `NavMesh` object which you can perform operations on. Here's a terse example:

```
f, _ := os.Open("de_dust2.nav") // Open the file
parser := gonav.Parser{Reader: f}
mesh, _ := parser.Parse() // Parse the file
area := mesh.QuadTreeAreas.FindAreaByPoint(gonav.Vector3{10, 10, 10}, true) // Find the nav area that contains the world point {10, 10} with a Z-value closest to 10
fmt.Println(area)
```

# Path Finding
This library also supports path finding across nav meshes via an A* implementation.

```
bombsiteA := mesh.GetPlaceByName("BombsiteA")
aCenter, _ := bombsiteA.GetEstimatedCenter()
aArea := mesh.GetNearestArea(aCenter, false)
bombsiteB := mesh.GetPlaceByName("BombsiteB")
bCenter, _ := bombsiteB.GetEstimatedCenter()
bArea := mesh.GetNearestArea(bCenter, false)
path, _ := gonav.SimpleBuildShortestPath(aArea, bArea)

for _, currNode := range path.Nodes {
	fmt.Println(currNode.Area)
}
```

Example output of path finding from the center of A-site on de_nuke to the center of B-site:
```
Path length: 2381.591
AreaID: 1419 [BombsiteA] @ {{675 -925 -414.96875}, {700 -900 -415.35532}}
AreaID: 7 [BombsiteA] @ {{575 -1325 -415.96875}, {950 -925 -415.96875}}
AreaID: 2285 [BombsiteA] @ {{550 -1300 -415.96875}, {575 -1225 -415.96875}}
AreaID: 3107 [BombsiteA] @ {{500 -1400 -415.96875}, {550 -1225 -415.96875}}
AreaID: 3368 [BombsiteA] @ {{475 -1375 -415.96875}, {500 -1225 -415.96875}}
AreaID: 3076 [BombsiteA] @ {{425 -1375 -415.96875}, {475 -1325 -415.96875}}
AreaID: 3075 [BombsiteA] @ {{425 -1400 -415.96875}, {475 -1375 -415.96875}}
AreaID: 1673 [Vents] @ {{425 -1450 -607.96875}, {625 -1425 -607.96875}}
AreaID: 700 [Vents] @ {{625 -1525 -607.96875}, {650 -1425 -607.96875}}
AreaID: 2478 [Vents] @ {{650 -1450 -607.96875}, {700 -1425 -607.96875}}
AreaID: 2479 [Vents] @ {{700 -1450 -607.96875}, {750 -1425 -607.96875}}
AreaID: 2477 [Vents] @ {{750 -1450 -607.96875}, {825 -1425 -607.96875}}
AreaID: 2300 [Vents] @ {{825 -1450 -607.96875}, {850 -1400 -607.96875}}
AreaID: 2301 [Vents] @ {{825 -1400 -607.96875}, {850 -1350 -607.96875}}
AreaID: 1690 [Vents] @ {{825 -1350 -607.96875}, {850 -1325 -639.96875}}
AreaID: 2265 [BombsiteB] @ {{775 -1325 -639.96875}, {900 -1275 -639.96875}}
AreaID: 3589 [BombsiteB] @ {{900 -1325 -639.96875}, {950 -1275 -639.96875}}
AreaID: 3591 [BombsiteB] @ {{900 -1275 -639.96875}, {950 -1050 -639.96875}}
AreaID: 3593 [BombsiteB] @ {{900 -1050 -639.96875}, {950 -900 -639.96875}}
AreaID: 2490 [BombsiteB] @ {{878 -1050 -597.96875}, {888 -900 -597.96875}}
AreaID: 2690 [BombsiteB] @ {{850 -1050 -767.96875}, {888 -900 -768.9288}}
AreaID: 3500 [BombsiteB] @ {{550 -1075 -770.25446}, {850 -900 -767.96875}}
AreaID: 3502 [BombsiteB] @ {{550 -900 -769.49255}, {850 -725 -767.96875}}
```

# License
This source is licensed under the GNU AFFERO GENERAL PUBLIC LICENSE. If this license is not acceptable for your project, let me know. Additionally, more feature-rich versions of this library exist, written in C# and C++, and can be licensed upon request.
