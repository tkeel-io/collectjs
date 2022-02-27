package collectjs

import (
	"bytes"
	"testing"
)

var planets = []*Collect{
	New(`["Venus", 0.815, 0.7]`),
	New(`["Earth", 1.0, 1.0]`),
	New(`["Mars", 0.107, 1.5]`),
}

func TestSort(t *testing.T) {

	// Closures that order the Planet structure.
	name := func(p1, p2 *Collect) bool {
		return bytes.Compare(p1.Get("[0]").Raw(), p2.Get("[0]").Raw()) > 0
	}
	mass := func(p1, p2 *Collect) bool {
		return bytes.Compare(p1.Get("[1]").Raw(), p2.Get("[1]").Raw()) > 0
	}
	distance := func(p1, p2 *Collect) bool {
		return bytes.Compare(p1.Get("[2]").Raw(), p2.Get("[2]").Raw()) > 0
	}

	// Sort the planets by the various criteria.
	By(name).Sort(planets)
	t.Log("By name:", string(planets[0].Raw()), string(planets[1].Raw()), string(planets[2].Raw()))

	By(mass).Sort(planets)
	t.Log("By mass:", string(planets[0].Raw()), string(planets[1].Raw()), string(planets[2].Raw()))

	By(distance).Sort(planets)
	t.Log("By distance:", string(planets[0].Raw()), string(planets[1].Raw()), string(planets[2].Raw()))

	// Output:
}
