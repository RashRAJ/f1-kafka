package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"f1sim/config"
	"f1sim/internal/telemetry"
)

var appConfig = config.LoadAppConfig()
var baseURL = appConfig.BaseURL
var httpClient = newHttpClient()

func newHttpClient() *http.Client {
	trp := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   2 * time.Minute,
			KeepAlive: 3 * time.Minute,
		}).DialContext,
		TLSHandshakeTimeout:   2 * time.Minute,
		ResponseHeaderTimeout: 3 * time.Minute,
		ExpectContinueTimeout: 1 * time.Minute,
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
	}
	client := &http.Client{
		Timeout:   10 * time.Minute,
		Transport: trp,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	return client
}

// FetchSessions retrieves F1 sessions for a given year
func FetchSessions(year int) ([]telemetry.Session, error) {
	url := fmt.Sprintf("%s/sessions?year=%d", baseURL, year)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sessions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(resp.Body, 1024)
		body, readErr := io.ReadAll(limitedReader)
		if readErr != nil {
			return nil, fmt.Errorf("API returned status %d (failed to read error body: %w)", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReaderSize(resp.Body, 8*1024)
	var sessions []telemetry.Session
	if err := json.NewDecoder(reader).Decode(&sessions); err != nil {
		return nil, fmt.Errorf("failed to decode sessions: %w", err)
	}

	return sessions, nil
}

// FetchCarData retrieves car telemetry data for a session
// If driverNumber is nil, fetches data for all drivers
func FetchCarData(sessionKey int, driverNumber *int) ([]telemetry.CarData, error) {
	url := fmt.Sprintf("%s/car_data?session_key=%d", baseURL, sessionKey)
	if driverNumber != nil {
		url = fmt.Sprintf("%s&driver_number=%d", url, *driverNumber)
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch car data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(resp.Body, 1024)
		body, readErr := io.ReadAll(limitedReader)
		if readErr != nil {
			return nil, fmt.Errorf("API returned status %d (failed to read error body: %w)", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReaderSize(resp.Body, 32*1024)
	var carData []telemetry.CarData
	if err := json.NewDecoder(reader).Decode(&carData); err != nil {
		return nil, fmt.Errorf("failed to decode car data: %w", err)
	}

	return carData, nil
}

// FetchDriverData retrieves driver information for a session
// If driverNumber is nil, fetches data for all drivers
func FetchDriverData(sessionKey int, driverNumber *int) ([]telemetry.Driver, error) {
	url := fmt.Sprintf("%s/drivers?session_key=%d", baseURL, sessionKey)
	if driverNumber != nil {
		url = fmt.Sprintf("%s&driver_number=%d", url, *driverNumber)
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch driver data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(resp.Body, 1024)
		body, readErr := io.ReadAll(limitedReader)
		if readErr != nil {
			return nil, fmt.Errorf("API returned status %d (failed to read error body: %w)", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReaderSize(resp.Body, 8*1024)
	var drivers []telemetry.Driver
	if err := json.NewDecoder(reader).Decode(&drivers); err != nil {
		return nil, fmt.Errorf("failed to decode driver data: %w", err)
	}

	return drivers, nil
}

func FetchIntervals(sessionKey int, driverNumber *int) ([]telemetry.Interval, error) {
	url := fmt.Sprintf("%s/intervals?session_key=%d", baseURL, sessionKey)
	if driverNumber != nil {
		url = fmt.Sprintf("%s&driver_number=%d", url, *driverNumber)
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch intervals: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(resp.Body, 1024)
		body, readErr := io.ReadAll(limitedReader)
		if readErr != nil {
			return nil, fmt.Errorf("API returned status %d (failed to read error body: %w)", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReaderSize(resp.Body, 16*1024)
	var intervals []telemetry.Interval
	if err := json.NewDecoder(reader).Decode(&intervals); err != nil {
		return nil, fmt.Errorf("failed to decode intervals: %w", err)
	}

	return intervals, nil
}

func FetchLaps(sessionKey int, driverNumber *int) ([]telemetry.Lap, error) {
	url := fmt.Sprintf("%s/laps?session_key=%d", baseURL, sessionKey)
	if driverNumber != nil {
		url = fmt.Sprintf("%s&driver_number=%d", url, *driverNumber)
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch laps: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(resp.Body, 1024)
		body, readErr := io.ReadAll(limitedReader)
		if readErr != nil {
			return nil, fmt.Errorf("API returned status %d (failed to read error body: %w)", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReaderSize(resp.Body, 32*1024)
	var laps []telemetry.Lap
	if err := json.NewDecoder(reader).Decode(&laps); err != nil {
		return nil, fmt.Errorf("failed to decode laps: %w", err)
	}

	return laps, nil
}

// FetchLocation retrieves GPS location data for a specific session and driver
func FetchLocation(sessionKey int, driverNumber *int) ([]telemetry.Location, error) {
	url := fmt.Sprintf("%s/location?session_key=%d", baseURL, sessionKey)
	if driverNumber != nil {
		url = fmt.Sprintf("%s&driver_number=%d", url, *driverNumber)
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch location data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(resp.Body, 1024)
		body, readErr := io.ReadAll(limitedReader)
		if readErr != nil {
			return nil, fmt.Errorf("API returned status %d (failed to read error body: %w)", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReaderSize(resp.Body, 64*1024)
	var locations []telemetry.Location
	if err := json.NewDecoder(reader).Decode(&locations); err != nil {
		return nil, fmt.Errorf("failed to decode location data: %w", err)
	}

	return locations, nil
}

func FetchMeetings(year int, country_name *string) ([]telemetry.Meeting, error) {
	url := fmt.Sprintf("%s/meetings?year=%d", baseURL, year)
	if country_name != nil {
		url = fmt.Sprintf("%s&country_name=%s", url, *country_name)
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch meetings: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(resp.Body, 1024)
		body, readErr := io.ReadAll(limitedReader)
		if readErr != nil {
			return nil, fmt.Errorf("API returned status %d (failed to read error body: %w)", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReaderSize(resp.Body, 8*1024)
	var meetings []telemetry.Meeting
	if err := json.NewDecoder(reader).Decode(&meetings); err != nil {
		return nil, fmt.Errorf("failed to decode meetings: %w", err)
	}

	return meetings, nil
}

// currently in beta and may complicate things

// func FetchOvertakes(sessionKey int) ([]telemetry.Overtake, error) {
// 	url := fmt.Sprintf("%s/overtakes?session_key=%d", baseURL, sessionKey)

// 	resp, err := httpClient(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch overtakes: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		body, _ := io.ReadAll(resp.Body)
// 		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read response body: %w", err)
// 	}

// 	var overtakes []telemetry.Overtake
// 	if err := json.Unmarshal(body, &overtakes); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal overtakes: %w", err)
// 	}

// 	return overtakes, nil
// }

func FetchPits(sessionKey int) ([]telemetry.Pit, error) {
	url := fmt.Sprintf("%s/pits?session_key=%d", baseURL, sessionKey)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pits: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(resp.Body, 1024)
		body, readErr := io.ReadAll(limitedReader)
		if readErr != nil {
			return nil, fmt.Errorf("API returned status %d (failed to read error body: %w)", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReaderSize(resp.Body, 8*1024)
	var pits []telemetry.Pit
	if err := json.NewDecoder(reader).Decode(&pits); err != nil {
		return nil, fmt.Errorf("failed to decode pits: %w", err)
	}

	return pits, nil
}

func FetchPositions(sessionKey int, driverNumber *int) ([]telemetry.Position, error) {
	url := fmt.Sprintf("%s/positions?session_key=%d", baseURL, sessionKey)
	if driverNumber != nil {
		url = fmt.Sprintf("%s&driver_number=%d", url, *driverNumber)
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch positions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(resp.Body, 1024)
		body, readErr := io.ReadAll(limitedReader)
		if readErr != nil {
			return nil, fmt.Errorf("API returned status %d (failed to read error body: %w)", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReaderSize(resp.Body, 32*1024)
	var positions []telemetry.Position
	if err := json.NewDecoder(reader).Decode(&positions); err != nil {
		return nil, fmt.Errorf("failed to decode positions: %w", err)
	}

	return positions, nil
}

func FetchRaceControl(sessionKey int) ([]telemetry.RaceControl, error) {
	url := fmt.Sprintf("%s/race_control?session_key=%d", baseURL, sessionKey)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch race control messages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(resp.Body, 1024)
		body, readErr := io.ReadAll(limitedReader)
		if readErr != nil {
			return nil, fmt.Errorf("API returned status %d (failed to read error body: %w)", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReaderSize(resp.Body, 8*1024)
	var raceControls []telemetry.RaceControl
	if err := json.NewDecoder(reader).Decode(&raceControls); err != nil {
		return nil, fmt.Errorf("failed to decode race control messages: %w", err)
	}

	return raceControls, nil
}

func FetchSessionResults(sessionKey int) ([]telemetry.SessionResult, error) {
	url := fmt.Sprintf("%s/session_results?session_key=%d", baseURL, sessionKey)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch session results: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(resp.Body, 1024)
		body, readErr := io.ReadAll(limitedReader)
		if readErr != nil {
			return nil, fmt.Errorf("API returned status %d (failed to read error body: %w)", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReaderSize(resp.Body, 8*1024)
	var sessionResults []telemetry.SessionResult
	if err := json.NewDecoder(reader).Decode(&sessionResults); err != nil {
		return nil, fmt.Errorf("failed to decode session results: %w", err)
	}

	return sessionResults, nil
}

func FetchStartingGrid(sessionKey int) ([]telemetry.StartingGrid, error) {
	url := fmt.Sprintf("%s/starting_grid?session_key=%d", baseURL, sessionKey)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch starting grid: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(resp.Body, 1024)
		body, readErr := io.ReadAll(limitedReader)
		if readErr != nil {
			return nil, fmt.Errorf("API returned status %d (failed to read error body: %w)", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReaderSize(resp.Body, 8*1024)
	var startingGrid []telemetry.StartingGrid
	if err := json.NewDecoder(reader).Decode(&startingGrid); err != nil {
		return nil, fmt.Errorf("failed to decode starting grid: %w", err)
	}

	return startingGrid, nil
}

func FetchStints(sessionKey int, driverNumber *int) ([]telemetry.Stint, error) {
	url := fmt.Sprintf("%s/stints?session_key=%d", baseURL, sessionKey)
	if driverNumber != nil {
		url = fmt.Sprintf("%s&driver_number=%d", url, *driverNumber)
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stints: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(resp.Body, 1024)
		body, readErr := io.ReadAll(limitedReader)
		if readErr != nil {
			return nil, fmt.Errorf("API returned status %d (failed to read error body: %w)", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReaderSize(resp.Body, 8*1024)
	var stints []telemetry.Stint
	if err := json.NewDecoder(reader).Decode(&stints); err != nil {
		return nil, fmt.Errorf("failed to decode stints: %w", err)
	}

	return stints, nil
}

func FetchTeamRadio(sessionKey int, driverNumber *int) ([]telemetry.TeamRadio, error) {
	url := fmt.Sprintf("%s/team_radio?session_key=%d", baseURL, sessionKey)
	if driverNumber != nil {
		url = fmt.Sprintf("%s&driver_number=%d", url, *driverNumber)
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch team radio messages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(resp.Body, 1024)
		body, readErr := io.ReadAll(limitedReader)
		if readErr != nil {
			return nil, fmt.Errorf("API returned status %d (failed to read error body: %w)", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReaderSize(resp.Body, 8*1024)
	var teamRadios []telemetry.TeamRadio
	if err := json.NewDecoder(reader).Decode(&teamRadios); err != nil {
		return nil, fmt.Errorf("failed to decode team radio messages: %w", err)
	}

	return teamRadios, nil
}

func FetchWeather(sessionKey int) ([]telemetry.Weather, error) {
	url := fmt.Sprintf("%s/weather?session_key=%d", baseURL, sessionKey)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(resp.Body, 1024)
		body, readErr := io.ReadAll(limitedReader)
		if readErr != nil {
			return nil, fmt.Errorf("API returned status %d (failed to read error body: %w)", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReaderSize(resp.Body, 8*1024)
	var weathers []telemetry.Weather
	if err := json.NewDecoder(reader).Decode(&weathers); err != nil {
		return nil, fmt.Errorf("failed to decode weather data: %w", err)
	}

	return weathers, nil
}
