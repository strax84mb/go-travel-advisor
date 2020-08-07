package database

import (
	"fmt"
	"math"
)

// SaveRoute - save new route
func SaveRoute(sourceID int64, destinationID int64, price float32) error {
	source, err := loadAirportByAirportID(sourceID)
	if err != nil {
		return err
	}
	destination, err := loadAirportByAirportID(destinationID)
	if err != nil {
		return err
	}
	route := Route{
		SourceID:      source.ID,
		DestinationID: destination.ID,
		Price:         price,
	}
	if gdb.Create(&route).Error() != nil {
		return err
	}
	return nil
}

type routeTreeNode struct {
	parent        *routeTreeNode
	children      []routeTreeNode
	routeID       int64
	startID       int64
	destinationID int64
	sumPrice      float32
}

func (r routeTreeNode) initFirst(startingPoint int64) {
	r.destinationID = startingPoint
}

func (r routeTreeNode) init(parent *routeTreeNode, route Route) {
	r.parent = parent
	r.routeID = route.ID
	r.startID = route.SourceID
	r.destinationID = route.DestinationID
	r.sumPrice = route.Price + parent.sumPrice
}

func (r routeTreeNode) destroy() {
	for _, c := range r.children {
		c.destroy()
	}
	r.children = nil
	r.parent = nil
}

func (r routeTreeNode) getAllFlights() []int64 {
	curr := &r
	count := 0
	for curr != nil && curr.routeID != 0 {
		count++
		curr = curr.parent
	}
	result := make([]int64, count)
	curr = &r
	count = 0
	for curr != nil && curr.routeID != 0 {
		result[count] = curr.routeID
		curr = curr.parent
	}
	return result
}

func (r routeTreeNode) getAllStops() []int64 {
	curr := &r
	count := 0
	for curr != nil {
		count++
		curr = curr.parent
	}
	result := make([]int64, count)
	curr = &r
	count = 0
	for curr != nil {
		result[count] = curr.destinationID
		curr = curr.parent
	}
	return result
}

type routeTree struct {
	sourceID         int64
	destinationID    int64
	root             *routeTreeNode
	cheapestPrice    float32 // Needs to be set to max possible value
	cheapestEndpoint *routeTreeNode
}

func (r routeTree) searchCheapestPath() (PathDto, error) {
	err := r.searchNode(r.root)
	if err != nil {
		return PathDto{}, err
	}
	if r.cheapestEndpoint == nil {
		return PathDto{
			flights:  nil,
			start:    PathRouteDto{},
			sumPrice: 0,
		}, nil
	}
	flightIDs := r.cheapestEndpoint.getAllFlights()
	flights, err := loadFullRoutes(flightIDs)
	if err != nil {
		return PathDto{}, err
	}
	pathDto := PathDto{
		sumPrice: r.cheapestPrice,
		start: PathRouteDto{
			RouteID: flights[0].ID,
			Airport: flights[0].Source.Name,
			City:    flights[0].Source.City.Name,
			Country: flights[0].Source.City.Country,
		},
	}
	flightDtos := make([]PathRouteDto, len(flights))
	for i, route := range flights {
		flightDtos[i] = PathRouteDto{
			RouteID: route.ID,
			Airport: route.Destination.Name,
			City:    route.Destination.City.Name,
			Country: route.Destination.City.Country,
		}
	}
	pathDto.flights = flightDtos
	return pathDto, nil
}

func (r routeTree) searchNode(n *routeTreeNode) error { // return error
	stops := n.getAllStops()
	possibleRoutes, err := findBySourceIDAndDestinationIDNotIn(n.destinationID, stops)
	if err != nil {
		return err
	}
	for _, route := range possibleRoutes {
		newNode := routeTreeNode{}
		newNode.init(n, route)
		if r.destinationID == route.DestinationID {
			if r.cheapestPrice > newNode.sumPrice {
				r.cheapestPrice = newNode.sumPrice
				r.cheapestEndpoint = &newNode
			}
		} else {
			if r.cheapestPrice > newNode.sumPrice {
				r.searchNode(&newNode)
			}
		}
	}
	return nil
}

func (r routeTree) destroy() {
	r.cheapestEndpoint = nil
	r.root.destroy()
	r.root = nil
}

func findBySourceIDAndDestinationIDNotIn(start int64, excludedDestinations []int64) ([]Route, error) {
	var count int
	if err := gdb.Model(&Route{}).Where("source_id = ? AND destination_id NOT IN (?)", &start, &excludedDestinations).
		Count(&count).Error(); err != nil {
		return nil, fmt.Errorf("Error while counting for start %d. Error: %s", start, err.Error())
	}
	result := make([]Route, count)
	if err := gdb.Where("source_id = ? AND destination_id NOT IN (?)", &start, &excludedDestinations).
		Find(&result).Error(); err != nil {
		return nil, fmt.Errorf("Error while loading destinations for start %d. Error: %s", start, err.Error())
	}
	return result, nil
}

func loadFullRoutes(ids []int64) ([]Route, error) {
	routes := make([]Route, len(ids))
	if err := gdb.Where("id IN(?)", &ids).
		Preload("Source").Preload("Source.City").
		Preload("Destination").Preload("Destination.City").
		Find(&routes).Error(); err != nil {
		return nil, err
	}
	result := make([]Route, len(ids))
	for i, id := range ids {
		result[i] = fetchRoute(routes, id)
	}
	return result, nil
}

func fetchRoute(routes []Route, ID int64) Route {
	for _, r := range routes {
		if r.ID == ID {
			return r
		}
	}
	panic("ID not found")
}

// FindCheapesRoute - find cheapes route
func FindCheapesRoute(start int64, end int64) (PathDto, error) {
	startAirport, err := loadAirportByAirportID(start)
	if err != nil {
		return PathDto{}, err
	}
	endAirport, err := loadAirportByAirportID(end)
	if err != nil {
		return PathDto{}, err
	}
	root := routeTreeNode{}
	root.initFirst(startAirport.ID)

	rt := routeTree{
		destinationID: endAirport.ID,
		root:          &root,
		cheapestPrice: math.MaxFloat32,
	}
	pathDto, err := rt.searchCheapestPath()
	if err != nil {
		return PathDto{}, err
	}
	return pathDto, nil
}
