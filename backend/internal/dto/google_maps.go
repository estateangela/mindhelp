package dto

import "github.com/go-playground/validator/v10"

// GeocodeRequest 地理編碼請求
type GeocodeRequest struct {
	Address string `json:"address" binding:"required,min=1,max=500" validate:"required,min=1,max=500"`
	Language string `json:"language" binding:"omitempty,len=2" validate:"omitempty,len=2"` // zh-TW, en
	Region   string `json:"region" binding:"omitempty,len=2" validate:"omitempty,len=2"`   // tw, us
}

// GeocodeResponse 地理編碼回應
type GeocodeResponse struct {
	Results []GeocodeResult `json:"results"`
	Status  string          `json:"status"`
}

// GeocodeResult 地理編碼結果
type GeocodeResult struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress  string            `json:"formatted_address"`
	Geometry          Geometry          `json:"geometry"`
	PlaceID           string            `json:"place_id"`
	Types             []string          `json:"types"`
}

// AddressComponent 地址組件
type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

// Geometry 幾何資訊
type Geometry struct {
	Location     LocationPoint `json:"location"`
	LocationType string        `json:"location_type"`
	Viewport     Bounds        `json:"viewport"`
	Bounds       Bounds        `json:"bounds"`
}

// LocationPoint 位置點
type LocationPoint struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Bounds 邊界
type Bounds struct {
	Northeast LocationPoint `json:"northeast"`
	Southwest LocationPoint `json:"southwest"`
}

// ReverseGeocodeRequest 反向地理編碼請求
type ReverseGeocodeRequest struct {
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180" validate:"required,min=-180,max=180"`
	Language  string  `json:"language" binding:"omitempty,len=2" validate:"omitempty,len=2"`
	ResultType string `json:"result_type" binding:"omitempty" validate:"omitempty"` // street_address, route, intersection, political, country, administrative_area_level_1, etc.
	LocationType string `json:"location_type" binding:"omitempty" validate:"omitempty"` // ROOFTOP, RANGE_INTERPOLATED, GEOMETRIC_CENTER, APPROXIMATE
}

// PlacesSearchRequest 地點搜尋請求
type PlacesSearchRequest struct {
	Query     string  `json:"query" binding:"omitempty,min=1,max=500" validate:"omitempty,min=1,max=500"`
	Location  string  `json:"location" binding:"omitempty" validate:"omitempty"` // "lat,lng"
	Radius    int     `json:"radius" binding:"omitempty,min=1,max=50000" validate:"omitempty,min=1,max=50000"`
	Type      string  `json:"type" binding:"omitempty" validate:"omitempty"`     // hospital, health, establishment
	Language  string  `json:"language" binding:"omitempty,len=2" validate:"omitempty,len=2"`
	Region    string  `json:"region" binding:"omitempty,len=2" validate:"omitempty,len=2"`
}

// PlacesSearchResponse 地點搜尋回應
type PlacesSearchResponse struct {
	Results []PlaceResult `json:"results"`
	Status  string        `json:"status"`
	NextPageToken string  `json:"next_page_token,omitempty"`
}

// PlaceResult 地點結果
type PlaceResult struct {
	PlaceID       string          `json:"place_id"`
	Name          string          `json:"name"`
	FormattedAddress string       `json:"formatted_address"`
	Geometry      Geometry        `json:"geometry"`
	Rating        float64         `json:"rating,omitempty"`
	PriceLevel    int             `json:"price_level,omitempty"`
	Types         []string        `json:"types"`
	Photos        []PlacePhoto    `json:"photos,omitempty"`
	OpeningHours  OpeningHours    `json:"opening_hours,omitempty"`
	Vicinity      string          `json:"vicinity,omitempty"`
}

// PlacePhoto 地點照片
type PlacePhoto struct {
	PhotoReference string `json:"photo_reference"`
	Height         int    `json:"height"`
	Width          int    `json:"width"`
	HTMLAttributions []string `json:"html_attributions"`
}

// OpeningHours 營業時間
type OpeningHours struct {
	OpenNow     bool     `json:"open_now"`
	Periods     []Period `json:"periods,omitempty"`
	WeekdayText []string `json:"weekday_text,omitempty"`
}

// Period 時間段
type Period struct {
	Close DayTime `json:"close,omitempty"`
	Open  DayTime `json:"open,omitempty"`
}

// DayTime 日期時間
type DayTime struct {
	Day  int    `json:"day"`  // 0-6 (Sunday-Saturday)
	Time string `json:"time"` // HHMM format
}

// DirectionsRequest 路線規劃請求
type DirectionsRequest struct {
	Origin      string `json:"origin" binding:"required,min=1,max=500" validate:"required,min=1,max=500"`
	Destination string `json:"destination" binding:"required,min=1,max=500" validate:"required,min=1,max=500"`
	Mode        string `json:"mode" binding:"omitempty" validate:"omitempty"` // driving, walking, bicycling, transit
	Language    string `json:"language" binding:"omitempty,len=2" validate:"omitempty,len=2"`
	Region      string `json:"region" binding:"omitempty,len=2" validate:"omitempty,len=2"`
	Alternatives bool  `json:"alternatives"` // 是否提供替代路線
	Avoid       string `json:"avoid" binding:"omitempty" validate:"omitempty"` // tolls, highways, ferries, indoor
	Units       string `json:"units" binding:"omitempty" validate:"omitempty"` // metric, imperial
}

// DirectionsResponse 路線規劃回應
type DirectionsResponse struct {
	Routes []Route `json:"routes"`
	Status string  `json:"status"`
}

// Route 路線
type Route struct {
	Summary         string            `json:"summary"`
	Legs            []Leg             `json:"legs"`
	OverviewPolyline Polyline         `json:"overview_polyline"`
	Bounds          Bounds            `json:"bounds"`
	Copyrights      string            `json:"copyrights"`
	Warnings        []string          `json:"warnings,omitempty"`
	WaypointOrder   []int             `json:"waypoint_order,omitempty"`
}

// Leg 路段
type Leg struct {
	Distance      Distance     `json:"distance"`
	Duration      Duration     `json:"duration"`
	DurationInTraffic Duration `json:"duration_in_traffic,omitempty"`
	StartAddress  string       `json:"start_address"`
	EndAddress    string       `json:"end_address"`
	StartLocation LocationPoint `json:"start_location"`
	EndLocation   LocationPoint `json:"end_location"`
	Steps         []Step       `json:"steps"`
	TrafficSpeedEntry []TrafficSpeedEntry `json:"traffic_speed_entry,omitempty"`
}

// Distance 距離
type Distance struct {
	Text  string `json:"text"`
	Value int    `json:"value"` // 公尺
}

// Duration 時間
type Duration struct {
	Text  string `json:"text"`
	Value int    `json:"value"` // 秒
}

// Step 步驟
type Step struct {
	Distance         Distance     `json:"distance"`
	Duration         Duration     `json:"duration"`
	EndLocation      LocationPoint `json:"end_location"`
	HTMLInstructions string       `json:"html_instructions"`
	Maneuver         string       `json:"maneuver,omitempty"`
	Polyline         Polyline     `json:"polyline"`
	StartLocation    LocationPoint `json:"start_location"`
	TravelMode       string       `json:"travel_mode"`
}

// Polyline 多邊線
type Polyline struct {
	Points string `json:"points"`
}

// TrafficSpeedEntry 交通速度條目
type TrafficSpeedEntry struct {
	OffsetMeters int    `json:"offset_meters"`
	SpeedCategory string `json:"speed_category"`
}

// DistanceMatrixRequest 距離矩陣請求
type DistanceMatrixRequest struct {
	Origins      []string `json:"origins" binding:"required,min=1,max=25" validate:"required,min=1,max=25"`
	Destinations []string `json:"destinations" binding:"required,min=1,max=25" validate:"required,min=1,max=25"`
	Mode         string   `json:"mode" binding:"omitempty" validate:"omitempty"`
	Language     string   `json:"language" binding:"omitempty,len=2" validate:"omitempty,len=2"`
	Region       string   `json:"region" binding:"omitempty,len=2" validate:"omitempty,len=2"`
	Units        string   `json:"units" binding:"omitempty" validate:"omitempty"`
	TrafficModel string   `json:"traffic_model" binding:"omitempty" validate:"omitempty"` // best_guess, pessimistic, optimistic
	DepartureTime string  `json:"departure_time" binding:"omitempty" validate:"omitempty"` // Unix timestamp
	ArrivalTime   string  `json:"arrival_time" binding:"omitempty" validate:"omitempty"`   // Unix timestamp
}

// DistanceMatrixResponse 距離矩陣回應
type DistanceMatrixResponse struct {
	DestinationAddresses []string                `json:"destination_addresses"`
	OriginAddresses      []string                `json:"origin_addresses"`
	Rows                 []DistanceMatrixRow     `json:"rows"`
	Status               string                  `json:"status"`
}

// DistanceMatrixRow 距離矩陣行
type DistanceMatrixRow struct {
	Elements []DistanceMatrixElement `json:"elements"`
}

// DistanceMatrixElement 距離矩陣元素
type DistanceMatrixElement struct {
	Distance Distance `json:"distance"`
	Duration Duration `json:"duration"`
	Status   string   `json:"status"`
}

// GoogleMapsError Google Maps API 錯誤回應
type GoogleMapsError struct {
	ErrorMessage string `json:"error_message"`
	Status       string `json:"status"`
	Results      []interface{} `json:"results"`
}

// Validate 驗證地理編碼請求
func (r *GeocodeRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate 驗證反向地理編碼請求
func (r *ReverseGeocodeRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate 驗證地點搜尋請求
func (r *PlacesSearchRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate 驗證路線規劃請求
func (r *DirectionsRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate 驗證距離矩陣請求
func (r *DistanceMatrixRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// BatchGeocodeRequest 批次地理編碼請求
type BatchGeocodeRequest struct {
	Addresses []string `json:"addresses" binding:"required,min=1,max=100" validate:"required,min=1,max=100"`
	Language  string   `json:"language" binding:"omitempty,len=2" validate:"omitempty,len=2"`
	Region    string   `json:"region" binding:"omitempty,len=2" validate:"omitempty,len=2"`
}

// BatchGeocodeResponse 批次地理編碼回應
type BatchGeocodeResponse struct {
	Results []GeocodeResponse `json:"results"`
	Total   int               `json:"total"`
	Status  string            `json:"status"`
}

// Validate 驗證批次地理編碼請求
func (r *BatchGeocodeRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
