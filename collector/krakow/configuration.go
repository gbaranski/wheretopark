package krakow

import "wheretopark/collector/krakow/meters"

var METER_TOTAL_SPOTS map[meters.Code]map[string]uint = map[meters.Code]map[string]uint{}
