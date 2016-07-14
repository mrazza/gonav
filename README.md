# gonav
A [Source Engine](https://en.wikipedia.org/wiki/Source_(game_engine)) bot Nav file parser written in Go. The specifics of the .nav format were reverse engineered using the information on [Valve's wiki](https://developer.valvesoftware.com/wiki/NAV) as a starting point. For more information on Source's Navigation Meshes see Valve's wiki: https://developer.valvesoftware.com/wiki/Navigation_Meshes

# Usage
See the `_examples` folder for examples of how to use this library. The basics, however, are pretty straightforward. Create a parser and pass in any `Reader` that contains the binary .nav data, then call the `Parse()` method. This will output a parsed `NavMesh` object which you can perform operations on. Here's a terse example:

```
f, _ := os.Open("de_dust2.nav") // Open the file
parser := gonav.Parser{Reader: f}
mesh, _ := parser.Parse() // Parse the file
area := mesh.QuadTreeAreas.FindAreaByPoint(gonav.Vector3{10, 10, 10}) // Find the nav area that contains the world point {10, 10, 10}
fmt.Println(area)
```

# License
This source is licensed under the GNU AFFERO GENERAL PUBLIC LICENSE. If this license is not acceptable for your project, let me know. Additionally, C# and C++ versions of this library exist under different licenses and can be licensed upon request.
