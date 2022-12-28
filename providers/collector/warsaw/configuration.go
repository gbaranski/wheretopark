package warsaw

import (
	wheretopark "wheretopark/go"

	_ "embed"
)

var configuration = &Configuration{
	ParkingLots: make(map[string]wheretopark.Metadata),
}

func init() {
	for k, v := range ztpParkingLots {
		v.PaymentMethods = append(v.PaymentMethods, ztpBasePaymentMethods...)
		v.Features = append(v.Features, ztpBaseFeatures...)
		v.Resources = append(v.Resources, ztpBasicResources...)
		configuration.ParkingLots[k] = v
	}
	for k, v := range prParkingLots {
		v.Resources = append(v.Resources, prBaseResources...)
		v.Features = append(v.Features, prBaseFeatures...)
		v.Comment = prComment
		v.Rules = prRules
		configuration.ParkingLots[k] = v
	}
}

type Configuration struct {
	ParkingLots map[wheretopark.ID]wheretopark.Metadata
}
