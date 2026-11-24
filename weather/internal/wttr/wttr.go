package wttr

type Data struct {
	CurrentCondition []CurrentCondition `json:"current_condition"`
	NearestArea      []NearestArea      `json:"nearest_area"`
	Request          []Request          `json:"request"`
	Weather          []Weather          `json:"weather"`
}

type StringValue struct {
	Value string `json:"value"`
}

type CurrentCondition struct {
	FeelsLikeC       string        `json:"FeelsLikeC"`
	FeelsLikeF       string        `json:"FeelsLikeF"`
	CloudCover       string        `json:"cloudcover"`
	Humidity         string        `json:"humidity"`
	LocalObsDateTime string        `json:"localObsDateTime"`
	ObservationTime  string        `json:"observation_time"`
	PrecipInches     string        `json:"precipInches"`
	PrecipMM         string        `json:"precipMM"`
	Pressure         string        `json:"pressure"`
	PressureInches   string        `json:"pressureInches"`
	TempC            string        `json:"temp_C"`
	TempF            string        `json:"temp_F"`
	UvIndex          string        `json:"uvIndex"`
	Visibility       string        `json:"visibility"`
	VisibilityMiles  string        `json:"visibilityMiles"`
	WeatherCode      string        `json:"weatherCode"`
	WeatherDesc      []StringValue `json:"weatherDesc"`
	WeatherIconUrl   []StringValue `json:"weatherIconUrl"`
	WindDir16Point   string        `json:"winddir16Point"`
	WindDirDegree    string        `json:"winddirDegree"`
	WindspeedKmph    string        `json:"windspeedKmph"`
	WindspeedMiles   string        `json:"windspeedMiles"`
}

type NearestArea struct {
	AreaName   []StringValue `json:"areaName"`
	Country    []StringValue `json:"country"`
	Latitude   string        `json:"latitude"`
	Longitude  string        `json:"longitude"`
	Population string        `json:"population"`
	Region     []StringValue `json:"region"`
	WeatherURL []StringValue `json:"weatherUrl"`
}

type Request struct {
	Query string `json:"query"`
	Type  string `json:"type"`
}

type Weather struct {
	Astronomy   []Astronomy `json:"astronomy"`
	AVGTempC    string      `json:"avgtempC"`
	AVGTempF    string      `json:"avgtempF"`
	Date        string      `json:"date"`
	Hourly      []Hourly    `json:"hourly"`
	MaxTempC    string      `json:"maxtempC"`
	MaxTempF    string      `json:"maxtempF"`
	MinTempC    string      `json:"mintempC"`
	MinTempF    string      `json:"mintempF"`
	SunHour     string      `json:"sunHour"`
	TotalSnowCM string      `json:"totalSnow_cm"`
	UVIndex     string      `json:"uvIndex"`
}

type Astronomy struct {
	MoonIllumination string `json:"moon_illumination"`
	MoonPhase        string `json:"moon_phase"`
	MoonRise         string `json:"moonrise"`
	MoonSet          string `json:"moonset"`
	SunRise          string `json:"sunrise"`
	SunSet           string `json:"sunset"`
}

type Hourly struct {
	DewPointC        string        `json:"DewPointC"`
	DewPointF        string        `json:"DewPointF"`
	FeelsLikeC       string        `json:"FeelsLikeC"`
	FeelsLikeF       string        `json:"FeelsLikeF"`
	HeatIndexC       string        `json:"HeatIndexC"`
	HeatIndexF       string        `json:"HeatIndexF"`
	WindChillC       string        `json:"WindChillC"`
	WindChillF       string        `json:"WindChillF"`
	WindGustKmph     string        `json:"WindGustKmph"`
	WindGustMiles    string        `json:"WindGustMiles"`
	ChanceOfFog      string        `json:"chanceoffog"`
	ChanceOfFrost    string        `json:"chanceoffrost"`
	ChanceOfHighTemp string        `json:"chanceofhightemp"`
	ChanceOfOverCast string        `json:"chanceofovercast"`
	ChanceOfRain     string        `json:"chanceofrain"`
	ChanceOfRemDry   string        `json:"chanceofremdry"`
	ChanceOfSnow     string        `json:"chanceofsnow"`
	ChanceOfSunshine string        `json:"chanceofsunshine"`
	ChanceOfThunder  string        `json:"chanceofthunder"`
	ChanceOfWindy    string        `json:"chanceofwindy"`
	CloudCover       string        `json:"cloudcover"`
	Humidity         string        `json:"humidity"`
	PrecipInches     string        `json:"precipInches"`
	PrecipMM         string        `json:"precipMM"`
	Pressure         string        `json:"pressure"`
	PressureInches   string        `json:"pressureInches"`
	TempC            string        `json:"tempC"`
	TempF            string        `json:"tempF"`
	Time             string        `json:"time"`
	UVIndex          string        `json:"uvIndex"`
	Visibility       string        `json:"visibility"`
	VisibilityMiles  string        `json:"visibilityMiles"`
	WeatherCode      string        `json:"weatherCode"`
	WeatherDesc      []StringValue `json:"weatherDesc"`
	WeatherIconUrl   []StringValue `json:"weatherIconUrl"`
	WindDir16Point   string        `json:"winddir16Point"`
	WindDirDegree    string        `json:"winddirDegree"`
	WindSpeedKmph    string        `json:"windspeedKmph"`
	WindSpeedMiles   string        `json:"windspeedMiles"`
}
