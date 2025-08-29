package location

import "fmt"

// LocationArea represents a location domain entity
type LocationArea struct {
	Name    string
	Pokemon []string
}

// LocationService handles location exploration
type LocationService struct {
	currentAreas []LocationArea
	currentPage  int
}

// NewLocationService creates a new location service
func NewLocationService() *LocationService {
	return &LocationService{
		currentAreas: make([]LocationArea, 0),
		currentPage:  0,
	}
}

// DisplayAreas shows the current location areas
func (ls *LocationService) DisplayAreas(areas []LocationArea) {
	ls.currentAreas = areas
	fmt.Println("\n--- Location Areas ---")
	for i, area := range areas {
		fmt.Printf("%d. %s\n", i+1, area.Name)
	}
}

// ExploreArea shows Pokemon in a specific location
func (ls *LocationService) ExploreArea(area LocationArea) {
	if len(area.Pokemon) == 0 {
		fmt.Printf("Found no Pokemon in %s...\n", area.Name)
		return
	}

	fmt.Printf("Exploring %s...\n", area.Name)
	fmt.Println("Found Pokemon:")
	for _, pokemon := range area.Pokemon {
		fmt.Printf("  - %s\n", pokemon)
	}
}

// GetAreaByName finds an area by name from current areas
func (ls *LocationService) GetAreaByName(name string) (LocationArea, bool) {
	for _, area := range ls.currentAreas {
		if area.Name == name {
			return area, true
		}
	}
	return LocationArea{}, false
}

// HasAreas checks if there are any loaded areas
func (ls *LocationService) HasAreas() bool {
	return len(ls.currentAreas) > 0
}

// GetAreaCount returns the number of current areas
func (ls *LocationService) GetAreaCount() int {
	return len(ls.currentAreas)
}
