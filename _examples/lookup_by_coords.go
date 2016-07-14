// Copyright (C) 2016 Matt Razza
// Use of this source code is governed by
// the license found in the LICENSE file.

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
		fmt.Println("Looking up via the slice...")
		var area *gonav.NavArea
		start = time.Now()
		for _, currArea := range mesh.Areas {
			if currArea.ContainsPoint(point) {
				area = currArea
				break
			}
		}
		elapsed = time.Since(start)

		if area != nil {
			fmt.Printf("Found in %d...\n", elapsed.Nanoseconds())
			fmt.Println(area)
		} else {
			fmt.Printf("No area found containing the specified coords in %v.\n", elapsed)
		}

		fmt.Println()
		fmt.Println("Looking up via the quadtree...")
		start = time.Now()
		area = mesh.QuadTreeAreas.FindAreaByPoint(point)
		elapsed = time.Since(start)

		if area != nil {
			fmt.Printf("Found in %d...\n", elapsed.Nanoseconds())
			fmt.Println(area)
		} else {
			fmt.Printf("No area found containing the specified coords in %v.\n", elapsed)
		}

		fmt.Println()
	}
}
