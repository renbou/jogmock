package gpx

import "encoding/xml"

// GPX file representation of a track with only needed parts
// Latitude, Longitude and Elevation must be contained within the GPX file
type GPX struct {
	XMLName xml.Name `xml:"gpx"`
	Track   Track    `xml:"trk"`
}

// GPX trk
type Track struct {
	XMLName      xml.Name     `xml:"trk"`
	TrackSegment TrackSegment `xml:"trkseg"`
}

// GPX trkseg
type TrackSegment struct {
	XMLName    xml.Name    `xml:"trkseg"`
	TrackParts []TrackPart `xml:"trkpt"`
}

// GPX trkpt
type TrackPart struct {
	XMLName   xml.Name `xml:"trkpt"`
	Lat       string   `xml:"lat,attr"`
	Lon       string   `xml:"lon,attr"`
	Elevation float64  `xml:"ele"`
}

func UnmarshalGPX(b []byte) (*GPX, error) {
	gpx := new(GPX)
	if err := xml.Unmarshal(b, &gpx); err != nil {
		return nil, err
	}
	return gpx, nil
}
