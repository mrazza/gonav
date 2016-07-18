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

	// Find the center of the bombsites
	start = time.Now()
	bombsiteA := mesh.GetPlaceByName("BombsiteA")
	aCenter, _ := bombsiteA.GetEstimatedCenter()
	aArea := mesh.GetNearestArea(aCenter, false)
	bombsiteB := mesh.GetPlaceByName("BombsiteB")
	bCenter, _ := bombsiteB.GetEstimatedCenter()
	bArea := mesh.GetNearestArea(bCenter, false)
	elapsed = time.Since(start)

	fmt.Printf("BombsiteA found at: %v\n", aArea)
	fmt.Printf("BombsiteB found at: %v\n", bArea)
	fmt.Printf("\tin: %v\n\n", elapsed)

	// Path find!
	start = time.Now()
	path, _ := gonav.SimpleBuildShortestPath(aArea, bArea)
	elapsed = time.Since(start)

	fmt.Printf("Path built in: %v (length: %v)\n", elapsed, path.GetCost())
	for _, currNode := range path.Nodes {
		fmt.Println(currNode.Area)
	}
}
