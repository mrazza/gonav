/*
	gonav - A Source Engine navigation mesh file parser written in Go.
	Copyright (C) 2016  Matt Razza

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published
	by the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mrazza/gonav"
)

func main() {
	fmt.Print("Enter file name (ex: de_dust2.nav): ")
	var file string
	fmt.Scanf("%s\n", &file)
	f, ok := os.Open(file) // Open the file

	if ok != nil {
		fmt.Printf("Failed to open file: %v\n", ok)
		return
	}

	defer f.Close()
	start := time.Now()
	parser := gonav.Parser{Reader: f}
	mesh, nerr := parser.Parse() // Parse the file
	elapsed := time.Since(start)

	if nerr != nil {
		fmt.Printf("Failed to parse: %v\n", nerr)
		return
	}

	fmt.Printf("Parse OK in %v\n\n", elapsed)

	// Loop forever and ask for coords.
	for {
		var x, y, z float32
		fmt.Print("Enter X coord: ")
		fmt.Scanf("%f\n", &x)

		fmt.Print("Enter Y coord: ")
		fmt.Scanf("%f\n", &y)

		fmt.Print("Enter Z coord: ")
		fmt.Scanf("%f\n", &z)

		point := gonav.Vector3{X: x, Y: y, Z: z}

		fmt.Println()
		fmt.Println("Looking up via the quadtree...")
		start = time.Now()
		area := mesh.QuadTreeAreas.FindAreaByPoint(point, false)
		elapsed = time.Since(start)

		if area != nil {
			fmt.Printf("Found in %fus...\n", float64(elapsed.Nanoseconds())/10000.0)
			fmt.Println(area)
		} else {
			fmt.Printf("No area found containing the specified coords in %v.\n", elapsed)
		}

		fmt.Println()
	}
}
