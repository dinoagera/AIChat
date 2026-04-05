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
	Text        string `json:"text" binding:"required"`
	PhoneNumber string `json:"phone_number,omitempty"`
	UserID      string `json:"user_id,omitempty"`
}
type EmergencyResponse struct {
	BrigadeID   int64   `json:"brigade_id"`
	BrigadeName string  `json:"brigade_name"`
	ETAMinutes  int     `json:"eta_minutes"`
	DistanceKm  float64 `json:"distance_km"`
	Address     string  `json:"address"`
	Priority    string  `json:"priority"`
	Message     string  `json:"message"`
	RequestID   int64   `json:"request_id"`
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

// // Response 200 OK
// {
//   "brigade_id": 1,
//   "brigade_name": "Бригада #1",
//   "eta_minutes": 12,
//   "distance_km": 2.3,
//   "weather_delay_minutes": 3,
//   "message": "Бригада выехала"
// }
