package wttr

type Data struct {
	CurrentCondition []CurrentCondition `json:"current_condition"`
	Weather          []Weather          `json:"weather"`
}

type StringValue struct {
	Value string `json:"value"`
}

type CurrentCondition struct {
	FeelsLikeC    string        `json:"FeelsLikeC"`
	Humidity      string        `json:"humidity"`
	Pressure      string        `json:"pressure"`
	TempC         string        `json:"temp_C"`
	WeatherDesc   []StringValue `json:"weatherDesc"`
	WindspeedKmph string        `json:"windspeedKmph"`
}

type Weather struct {
	Astronomy []Astronomy `json:"astronomy"`
	Hourly    []Hourly    `json:"hourly"`
}

type Astronomy struct {
	SunRise string `json:"sunrise"`
	SunSet  string `json:"sunset"`
}

type Hourly struct {
	FeelsLikeC  string        `json:"FeelsLikeC"`
	TempC       string        `json:"tempC"`
	WeatherDesc []StringValue `json:"weatherDesc"`
}
