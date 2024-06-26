package dto

import "gitlab.strale.io/go-travel/internal/database"

func ConvertArray[
	V database.City | database.Airport | database.Comment | database.Route,
	T CityDto | AirportDto | CityAirportDto | CommentDto | RouteDto,
](list []V, convert func(V) *T) []T {
	if list == nil {
		return nil
	}
	result := make([]T, len(list))
	for i, v := range list {
		result[i] = *convert(v)
	}
	return result
}
