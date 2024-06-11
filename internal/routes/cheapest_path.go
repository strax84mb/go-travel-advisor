package routes

import "gitlab.strale.io/go-travel/internal/database"

type CRNode struct {
	Destinations     []*CRNode
	Parent           *CRNode
	Route            database.Route
	AccumulatedPrice int64
}

func listPassedCities(crn *CRNode, ids []int64) []int64 {
	if ids == nil {
		ids = []int64{crn.Route.Destination.CityID}
	} else {
		ids = append(ids, crn.Route.Destination.CityID)
	}
	if crn.Parent == nil {
		return ids
	} else {
		return listPassedCities(crn.Parent, ids)
	}
}

func (crn *CRNode) ListPassedCities() []int64 {
	return listPassedCities(crn, nil)
}

type CheapestRouteWalk struct {
	CheapestPrice       int64
	CheapestPathEndNode *CRNode
	Expanded            bool
	FinalCityID         int64
	Root                CRNode
}

func newCheapestRouteWalk(start int64) CheapestRouteWalk {
	return CheapestRouteWalk{
		CheapestPrice:       0,
		CheapestPathEndNode: nil,
		Expanded:            false,
		Root: CRNode{
			Route: database.Route{
				Destination: &database.Airport{
					CityID: start,
				},
			},
		},
	}
}

func (crw *CheapestRouteWalk) destroy(node *CRNode) {
	if node.Destinations != nil {
		for i, n := range node.Destinations {
			crw.destroy(n)
			node.Destinations[i] = nil
		}
	}
	node.Parent = nil
}

func (crw *CheapestRouteWalk) Reset(startCityID int64) {
	crw.destroy(&crw.Root)
	crw.Root = CRNode{
		Parent:           nil,
		Destinations:     nil,
		AccumulatedPrice: 0,
		Route: database.Route{
			DestinationID: startCityID,
		},
	}
}

func (crw *CheapestRouteWalk) CollectCheapestPath() []*database.Route {
	if crw.CheapestPathEndNode == nil {
		return nil
	}
	var (
		node  *CRNode = crw.CheapestPathEndNode
		count int
	)
	for {
		count++
		node = node.Parent
		if node == nil {
			break
		}
	}
	node = crw.CheapestPathEndNode
	list := make([]*database.Route, count)
	for {
		count--
		list[count] = &node.Route
		node = node.Parent
		if node == nil {
			break
		}
	}
	return list
}
