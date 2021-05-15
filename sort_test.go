package collectjs

import (
	"bytes"
	"fmt"
)

var planets = []*Collect{
	New(`["Venus", 0.815, 0.7]`),
	New(`["Earth", 1.0, 1.0]`),
	New(`["Mars", 0.107, 1.5]`),
}

func Example_a() {

	// Closures that order the Planet structure.
	name := func(p1, p2 *Collect) bool {
		return bytes.Compare(p1.Get("[0]").raw, p2.Get("[0]").raw) > 0
	}
	mass := func(p1, p2 *Collect) bool {
		return bytes.Compare(p1.Get("[1]").raw, p2.Get("[1]").raw) > 0
	}
	distance := func(p1, p2 *Collect) bool {
		return bytes.Compare(p1.Get("[2]").raw, p2.Get("[2]").raw) > 0
	}

	// Sort the planets by the various criteria.
	By(name).Sort(planets)
	fmt.Println("By name:", string(planets[0].raw), string(planets[1].raw), string(planets[2].raw))

	By(mass).Sort(planets)
	fmt.Println("By mass:", string(planets[0].raw), string(planets[1].raw), string(planets[2].raw))

	By(distance).Sort(planets)
	fmt.Println("By distance:", string(planets[0].raw), string(planets[1].raw), string(planets[2].raw))

	// Output:
}