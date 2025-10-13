package models

import "time"

type APIInfoResponse struct {
	Name        string       `json:"name" example:"Base Golang RESTful API"`
	Version     string       `json:"version" example:"v1"`
	Description string       `json:"description" example:"A comprehensive RESTful API"`
	Status      string       `json:"status" example:"operational"`
	Timestamp   time.Time    `json:"timestamp"`
	Features    []string     `json:"features"`
	Endpoints   APIEndpoints `json:"endpoints"`
}

type APIEndpoints struct {
	Auth          string `json:"auth" example:"/api/v1/auth"`
	Users         string `json:"users" example:"/api/v1/users"`
	Products      string `json:"products" example:"/api/v1/products"`
	Email         string `json:"email" example:"/api/v1/email"`
	Notifications string `json:"notifications" example:"/api/v1/notifications"`
	Documentation string `json:"documentation" example:"/swagger/index.html"`
}

type APIVersionResponse struct {
	Version            string   `json:"version" example:"v1"`
	ExtractedFrom      string   `json:"extracted_from" example:"header"`
	SupportedVersions  []string `json:"supported_versions"`
	DeprecatedVersions []string `json:"deprecated_versions"`
}

type HealthDetailedResponse struct {
	Status    string        `json:"status" example:"healthy"`
	Timestamp time.Time     `json:"timestamp"`
	Version   string        `json:"version" example:"v1"`
	Uptime    string        `json:"uptime" example:"2h30m15s"`
	System    SystemMetrics `json:"system"`
	Services  ServiceStatus `json:"services"`
}

type SystemMetrics struct {
	GoVersion    string `json:"go_version" example:"go1.21.0"`
	NumGoroutine int    `json:"num_goroutine" example:"25"`
	NumCPU       int    `json:"num_cpu" example:"8"`
	MemoryAlloc  uint64 `json:"memory_alloc" example:"5242880"`
	MemoryTotal  uint64 `json:"memory_total" example:"10485760"`
	MemorySys    uint64 `json:"memory_sys" example:"20971520"`
	NumGC        uint32 `json:"num_gc" example:"10"`
}

type ServiceStatus struct {
	Database      bool `json:"database" example:"true"`
	Redis         bool `json:"redis" example:"true"`
	Email         bool `json:"email" example:"true"`
	Notifications bool `json:"notifications" example:"true"`
}

type ServerTimeResponse struct {
	UTC       time.Time `json:"utc"`
	Local     time.Time `json:"local"`
	Unix      int64     `json:"unix" example:"1697184000"`
	UnixMilli int64     `json:"unix_milli" example:"1697184000000"`
	Timezone  string    `json:"timezone" example:"UTC"`
	Formatted string    `json:"formatted" example:"2023-10-13T10:00:00Z"`
}
