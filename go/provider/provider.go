package provider

type Common interface {
	Name() string
}

// func RunProvidersHTTP(
// 	port int,
// 	providers []Common,
// 	cache wheretopark.CacheProvider,
// 	getParkingLotsFn func([]Common, wheretopark.CacheProvider, string) error
// ) error {

// }
