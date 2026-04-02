package handler

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}
type ResponseWithTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
type Response struct {
	Message string `json:"message"`
}
type EmergencyRequest struct {
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Address  string  `json:"address,omitempty"`
	Priority string  `json:"priority,omitempty"`
}

type EmergencyResponse struct {
	BrigadeID       int64   `json:"brigade_id"`
	BrigadeName     string  `json:"brigade_name"`
	ETAMinutes      int     `json:"eta_minutes"`
	DistanceKm      float64 `json:"distance_km"`
	WeatherDelayMin int     `json:"weather_delay_minutes,omitempty"`
	Message         string  `json:"message"`
}
type AddBrigadeRequest struct {
	Name   string  `json:"name"`
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	Status string  `json:"status,omitempty"`
}

type AddBrigadeResponse struct {
	ID      int64   `json:"brigade_id"`
	Name    string  `json:"brigade_name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Status  string  `json:"status"`
	Message string  `json:"message"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=free busy offline"`
}

func (r AddBrigadeRequest) Validate() bool {
	if r.Name == "" {
		return false
	}
	if r.Lat < -90 || r.Lat > 90 {
		return false
	}
	if r.Lon < -180 || r.Lon > 180 {
		return false
	}
	if r.Status != "" && r.Status != "free" && r.Status != "busy" {
		return false
	}
	return true
}

// {
//   "lat": 55.788,
//   "lon": 49.122,
//   "address": "ул. Баумана, 15",
//   "priority": "critical"
// }

// // Response 200 OK
// {
//   "brigade_id": 1,
//   "brigade_name": "Бригада #1",
//   "eta_minutes": 12,
//   "distance_km": 2.3,
//   "weather_delay_minutes": 3,
//   "message": "Бригада выехала"
// }
