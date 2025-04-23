package models

type SwiftCode struct {
	SwiftCode     string `bson:"swiftCode" json:"swiftCode"`
	SwiftPrefix   string `bson:"swiftPrefix" json:"-"`
	IsHeadquarter bool   `bson:"isHeadquarter" json:"isHeadquarter"`
	BankName      string `bson:"bankName" json:"bankName"`
	Address       string `bson:"address" json:"address"`
	CountryISO2   string `bson:"countryISO2" json:"countryISO2"`
	CountryName   string `bson:"countryName" json:"countryName"`
}
