package ourwx

import "time"

const (
	// this url requires the addition of query parameters `lat` and `lon`
	ndfdURLFmt = "https://graphical.weather.gov/xml/sample_products/browser_interface/ndfdXMLclient.php"
)

// Forecast holds a forcast for a single point. Numeric fields are type `string`
// so that empty strings may represent missing data in  the struct. This package
// is built with an eye towards home automation applications and usefullness for
// most people. Some available fields are not included becuase they are outside
// the scope of this package.
type Forecast struct {
	Daily []struct {
		Day         time.Time
		Temperature struct {
			Units     string
			HighValue string
			LowValue  string
		}
		RelativeHumidity struct {
			Units     string
			HighValue string
			LowValue  string
		}
		PrecipitationProbability struct {
			Units string
			Day   string
			Night string
		}
	}
	Hourly []struct {
		Hour        time.Time
		Temperature struct {
			Units string
			Value string
		}
		ApparentTemperature struct {
			Units string
			Value string
		}
		RelativeHumidity struct {
			Units string
			Value string
		}
		WindSpeed struct {
			Units          string
			SustainedValue string
			GustValue      string
		}
		WindDirection struct {
			Units string
			Value string
		}
		CloudAmount struct {
			Units string
			Value string
		}
		PrecipitationAmount struct {
			Units       string
			LiquidValue string
			SnowValue   string
			IceValue    string
		}
	}
}

// forecast holds the unmarshalled useful parts of the raw XML
// https://graphical.weather.gov/xml/rest.php#XML_contents
// https://graphical.weather.gov/xml/DWMLgen/schema/DWML.xsd
type forecast struct {
	Data struct {
		TimeLayouts []struct {
			Key string `xml:"layout-kay"`
			// some layouts don't have end-valid-times. In those cases,
			// start-valid-time is being considered a point instead of a range
			StartValidTimes []string `xml:"start-valid-time"`
			EndValidTimes   []string `xml:"end-valid-time"`
		} `xml:"time-layout"`
		Parameters struct {
			Temperatures []struct {
				Name string `xml:"name"`
				// Type may be: maximum, minimum, hourly, apparent ...
				Type       string   `xml:"type,attr"`
				Units      string   `xml:"units,attr"`
				TimeLayout string   `xml:"time-layout,attr"`
				Values     []string `xml:"value"`
			} `xml:"temperature"`
			Precipitations []struct {
				Name string `xml:"name"`
				// Type may be: liquid, ice, snow ...
				Type       string   `xml:"type,attr"`
				Units      string   `xml:"units,attr"`
				TimeLayout string   `xml:"time-layout,attr"`
				Values     []string `xml:"value"`
			} `xml:"precipitation"`
			ProbabilityOfPrecipitations []struct {
				Name string `xml:"name"`
				// Type may be: 12 hour ...
				Type       string   `xml:"type,attr"`
				Units      string   `xml:"units,attr"`
				TimeLayout string   `xml:"time-layout,attr"`
				Values     []string `xml:"value"`
			} `xml:"probability-of-precipitation"`
			WindSpeeds []struct {
				Name string `xml:"name"`
				// Type may be: sustained, gust ...
				Type       string   `xml:"type,attr"`
				Units      string   `xml:"units,attr"`
				TimeLayout string   `xml:"time-layout,attr"`
				Values     []string `xml:"value"`
			} `xml:"wind-speed"`
			Direction []struct {
				Name string `xml:"name"`
				// Type may be: wind ...
				Type       string   `xml:"type,attr"`
				Units      string   `xml:"units,attr"`
				TimeLayout string   `xml:"time-layout,attr"`
				Values     []string `xml:"value"`
			} `xml:"wind-speed"`
			CloudAmounts []struct {
				Name string `xml:"name"`
				// Type may be: total ...
				Type       string   `xml:"type,attr"`
				Units      string   `xml:"units,attr"`
				TimeLayout string   `xml:"time-layout,attr"`
				Values     []string `xml:"value"`
			} `xml:"cloud-amount"`
			Humidities []struct {
				Name string `xml:"name"`
				// Type may be: relative, `maximum relative`, `minimum relative`
				Type       string   `xml:"type,attr"`
				Units      string   `xml:"units,attr"`
				TimeLayout string   `xml:"time-layout,attr"`
				Values     []string `xml:"value"`
			} `xml:"humidity"`
		} `xml:"parameters"`
	} `xml:"data"`
}
