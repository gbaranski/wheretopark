import * as S from 'wheretopark-shared'

export import Coordinate = S.app.wheretopark.shared.Coordinate
export type ParkingLotID = string
export import ParkingLot = S.app.wheretopark.shared.ParkingLot
export import ParkingLotMetadata = S.app.wheretopark.shared.ParkingLotMetadata
export import ParkingLotState = S.app.wheretopark.shared.ParkingLotState
export import ParkingLotSpotType = S.app.wheretopark.shared.ParkingSpotType
export import ParkingLotRule = S.app.wheretopark.shared.ParkingLotRule
export import ParkingLotPricingRule = S.app.wheretopark.shared.ParkingLotPricingRule
export import ParkingLotResource = S.app.wheretopark.shared.ParkingLotResource
export import parseParkingLots = S.app.wheretopark.shared.parseParkingLots
export import StorekeeperClient = S.app.wheretopark.shared.JSStorekeeperClient
export import AuthorizationClient = S.app.wheretopark.shared.JSAuthorizationClient
export import fetchParkingLots = S.app.wheretopark.shared.fetchParkingLots
export import toRecord = S.app.wheretopark.shared.toRecord
export import toArray = S.app.wheretopark.shared.toArray
export import instantToJSDate = S.app.wheretopark.shared.instantToJSDate
export import durationToISO = S.app.wheretopark.shared.durationToISO
