// F1 24 Telemetry Packet (packetId = 6)
package telemetry

import (
	"time"
)

type PacketHeader struct {
	PacketFormat            uint16 // 2024
	GameYear                uint8  // 24
	GameMajorVersion        uint8
	GameMinorVersion        uint8
	PacketVersion           uint8
	PacketId                uint8
	SessionUID              uint64
	SessionTime             float32
	FrameIdentifier         uint32
	OverallFrameIdentifier  uint32
	PlayerCarIndex          uint8
	SecondaryPlayerCarIndex uint8
}

type CarTelemetryData struct {
	Speed                   uint16
	Throttle                float32
	Steer                   float32
	Brake                   float32
	Clutch                  uint8
	Gear                    int8
	EngineRPM               uint16
	DRS                     uint8
	RevLightsPercent        uint8
	RevLightsBitValue       uint16
	BrakesTemperature       [4]uint16
	TyresSurfaceTemperature [4]uint8
	TyresInnerTemperature   [4]uint8
	EngineTemperature       uint16
	TyresPressure           [4]float32
	SurfaceType             [4]uint8
}

type PacketCarTelemetryData struct {
	Header                 PacketHeader
	CarTelemetryData       [22]CarTelemetryData
	MFDPanelIndex          uint8
	MFDPanelIndexSecondary uint8
	SuggestedGear          int8
}

// ############################################// Session represents an F1 session from the OpenF1 API
type Session struct {
	MeetingKey       int       `json:"meeting_key"`
	SessionKey       int       `json:"session_key"`
	Location         string    `json:"location"`
	DateStart        time.Time `json:"date_start"`
	DateEnd          time.Time `json:"date_end"`
	SessionType      string    `json:"session_type"`
	SessionName      string    `json:"session_name"`
	CountryKey       int       `json:"country_key"`
	CountryCode      string    `json:"country_code"`
	CountryName      string    `json:"country_name"`
	CircuitKey       int       `json:"circuit_key"`
	CircuitShortName string    `json:"circuit_short_name"`
	GMTOffset        string    `json:"gmt_offset"`
	Year             int       `json:"year"`
}

type CarData struct {
	Date         time.Time `json:"date"`
	SessionKey   int       `json:"session_key"`
	MeetingKey   int       `json:"meeting_key"`
	DriverNumber int       `json:"driver_number"`
	Speed        int       `json:"speed"`
	Throttle     int       `json:"throttle"`
	Brake        int       `json:"brake"`
	RPM          int       `json:"rpm"`
	NGear        int       `json:"n_gear"`
	DRS          int       `json:"drs"`
}

type Driver struct {
	Broadcast_name string `json:"broadcast_name"`
	Country_code   string `json:"country_code"`
	Driver_number  int    `json:"driver_number"`
	First_name     string `json:"first_name"`
	Full_name      string `json:"full_name"`
	Headshot_url   string `json:"headshot_url"`
	Last_name      string `json:"last_name"`
	Meeting_key    int    `json:"meeting_key"`
	Name_acronym   string `json:"name_acronym"`
	Session_key    int    `json:"session_key"`
	Team_colour    string `json:"team_colour"`
	Team_name      string `json:"team_name"`
}

type Interval struct {
	Date         time.Time `json:"date"`
	DriverNumber int       `json:"driver_number"`
	GapToLeader  float64   `json:"gap_to_leader"`
	Interval     float64   `json:"interval"`
	MeetingKey   int       `json:"meeting_key"`
	SessionKey   int       `json:"session_key"`
}

type Lap struct {
	DateStart       time.Time `json:"date_start"`
	DriverNumber    int       `json:"driver_number"`
	DurationSector1 float64   `json:"duration_sector_1"`
	DurationSector2 float64   `json:"duration_sector_2"`
	DurationSector3 float64   `json:"duration_sector_3"`
	I1Speed         int       `json:"i1_speed"`
	I2Speed         int       `json:"i2_speed"`
	IsPitOutLap     bool      `json:"is_pit_out_lap"`
	LapDuration     float64   `json:"lap_duration"`
	LapNumber       int       `json:"lap_number"`
	MeetingKey      int       `json:"meeting_key"`
	SegmentsSector1 []int     `json:"segments_sector_1"`
	SegmentsSector2 []int     `json:"segments_sector_2"`
	SegmentsSector3 []int     `json:"segments_sector_3"`
	SessionKey      int       `json:"session_key"`
	STSpeed         int       `json:"st_speed"`
}

type Location struct {
	Date         time.Time `json:"date"`
	DriverNumber int       `json:"driver_number"`
	MeetingKey   int       `json:"meeting_key"`
	SessionKey   int       `json:"session_key"`
	X            int       `json:"x"`
	Y            int       `json:"y"`
	Z            int       `json:"z"`
}

// Meeting represents an F1 race weekend/event from the OpenF1 API
type Meeting struct {
	CircuitKey          int       `json:"circuit_key"`
	CircuitShortName    string    `json:"circuit_short_name"`
	CountryCode         string    `json:"country_code"`
	CountryKey          int       `json:"country_key"`
	CountryName         string    `json:"country_name"`
	DateStart           time.Time `json:"date_start"`
	GMTOffset           string    `json:"gmt_offset"`
	Location            string    `json:"location"`
	MeetingKey          int       `json:"meeting_key"`
	MeetingName         string    `json:"meeting_name"`
	MeetingOfficialName string    `json:"meeting_official_name"`
	Year                int       `json:"year"`
}

// Overtake represents an overtaking maneuver from the OpenF1 API
type Overtake struct {
	Date                   time.Time `json:"date"`
	MeetingKey             int       `json:"meeting_key"`
	OvertakenDriverNumber  int       `json:"overtaken_driver_number"`
	OvertakingDriverNumber int       `json:"overtaking_driver_number"`
	Position               int       `json:"position"`
	SessionKey             int       `json:"session_key"`
}

// Pit represents a pit stop from the OpenF1 API
type Pit struct {
	Date         time.Time `json:"date"`
	DriverNumber int       `json:"driver_number"`
	LapNumber    int       `json:"lap_number"`
	MeetingKey   int       `json:"meeting_key"`
	PitDuration  float64   `json:"pit_duration"`
	SessionKey   int       `json:"session_key"`
}

// Position represents driver position data from the OpenF1 API
type Position struct {
	Date         time.Time `json:"date"`
	DriverNumber int       `json:"driver_number"`
	MeetingKey   int       `json:"meeting_key"`
	Position     int       `json:"position"`
	SessionKey   int       `json:"session_key"`
}

// RaceControl represents race control messages from the OpenF1 API
type RaceControl struct {
	Category     string    `json:"category"`
	Date         time.Time `json:"date"`
	DriverNumber *int      `json:"driver_number"` // nullable
	Flag         *string   `json:"flag"`          // nullable
	LapNumber    *int      `json:"lap_number"`    // nullable
	MeetingKey   int       `json:"meeting_key"`
	Message      string    `json:"message"`
	Scope        string    `json:"scope"`
	Sector       *int      `json:"sector"` // nullable
	SessionKey   int       `json:"session_key"`
}

// SessionResult represents final session results from the OpenF1 API
type SessionResult struct {
	DNF          bool    `json:"dnf"`
	DNS          bool    `json:"dns"`
	DSQ          bool    `json:"dsq"`
	DriverNumber int     `json:"driver_number"`
	Duration     float64 `json:"duration"`
	GapToLeader  float64 `json:"gap_to_leader"`
	NumberOfLaps int     `json:"number_of_laps"`
	MeetingKey   int     `json:"meeting_key"`
	Position     int     `json:"position"`
	SessionKey   int     `json:"session_key"`
}

// StartingGrid represents starting grid positions from the OpenF1 API
type StartingGrid struct {
	Position     int     `json:"position"`
	DriverNumber int     `json:"driver_number"`
	LapDuration  float64 `json:"lap_duration"`
	MeetingKey   int     `json:"meeting_key"`
	SessionKey   int     `json:"session_key"`
}

// Stint represents tire stint data from the OpenF1 API
type Stint struct {
	Compound       string `json:"compound"`
	DriverNumber   int    `json:"driver_number"`
	LapEnd         int    `json:"lap_end"`
	LapStart       int    `json:"lap_start"`
	MeetingKey     int    `json:"meeting_key"`
	SessionKey     int    `json:"session_key"`
	StintNumber    int    `json:"stint_number"`
	TyreAgeAtStart int    `json:"tyre_age_at_start"`
}

// TeamRadio represents team radio communication from the OpenF1 API
type TeamRadio struct {
	Date         time.Time `json:"date"`
	DriverNumber int       `json:"driver_number"`
	MeetingKey   int       `json:"meeting_key"`
	RecordingURL string    `json:"recording_url"`
	SessionKey   int       `json:"session_key"`
}

// Weather represents weather conditions from the OpenF1 API
type Weather struct {
	AirTemperature   float64   `json:"air_temperature"`
	Date             time.Time `json:"date"`
	Humidity         int       `json:"humidity"`
	MeetingKey       int       `json:"meeting_key"`
	Pressure         float64   `json:"pressure"`
	Rainfall         int       `json:"rainfall"`
	SessionKey       int       `json:"session_key"`
	TrackTemperature float64   `json:"track_temperature"`
	WindDirection    int       `json:"wind_direction"`
	WindSpeed        float64   `json:"wind_speed"`
}
