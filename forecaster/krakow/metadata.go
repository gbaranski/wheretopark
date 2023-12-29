package krakow

import (
	"encoding/xml"
	"io"
	"os"

	"golang.org/x/text/encoding/charmap"
)

type Placemark struct {
	Zone        string `xml:"name"`
	Card        string `xml:"card"`
	Model       string `xml:"model"`
	Code        string `xml:"parkingmeter"`
	Address     string `xml:"address"`
	Coordinates struct {
		Latitude  float64 `xml:"latitude"`
		Longitude float64 `xml:"longitude"`
	} `xml:"coordinates"`
}

type Folder struct {
	Placemarks []Placemark `xml:"placemark"`
}

func LoadPlacemarks(path string) ([]Placemark, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	byteValue, err := io.ReadAll(charmap.ISO8859_2.NewDecoder().Reader(file))
	if err != nil {
		return nil, err
	}
	var folder Folder
	err = xml.Unmarshal(byteValue, &folder)
	if err != nil {
		return nil, err
	}
	return folder.Placemarks, nil
}
