package history

import (
	"testing"
	"time"
)

var forecastBase = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func forecastEntry(host string, port int, event EventType, t time.Time) Entry {
	return Entry{
		Host:      host,
		Ports:     []int{port},
		Event:     event,
		Timestamp: t,
	}
}

func TestForecast_PredictsFutureEvent(t *testing.T) {
	entries := []Entry{
		forecastEntry("host-a", 80, EventOpened, forecastBase),
		forecastEntry("host-a", 80, EventOpened, forecastBase.Add(24*time.Hour)),
		forecastEntry("host-a", 80, EventOpened, forecastBase.Add(48*time.Hour)),
	}

	now := forecastBase.Add(48 * time.Hour)
	forecasts := Forecast(entries, now)

	if len(forecasts) == 0 {
		t.Fatal("expected at least one forecast, got none")
	}

	f := forecasts[0]
	if f.Host != "host-a" {
		t.Errorf("expected host-a, got %s", f.Host)
	}
	if f.Port != 80 {
		t.Errorf("expected port 80, got %d", f.Port)
	}
	expected := forecastBase.Add(72 * time.Hour)
	if !f.LikelyAt.Equal(expected) {
		t.Errorf("expected LikelyAt %v, got %v", expected, f.LikelyAt)
	}
}

func TestForecast_ConfidenceIncreasesWithCount(t *testing.T) {
	var entries []Entry
	for i := 0; i < 8; i++ {
		entries = append(entries, forecastEntry("host-b", 443, EventOpened, forecastBase.Add(time.Duration(i)*24*time.Hour)))
	}

	now := forecastBase.Add(8 * 24 * time.Hour)
	forecasts := Forecast(entries, now)

	if len(forecasts) == 0 {
		t.Fatal("expected forecast")
	}
	if forecasts[0].Confidence <= 0.5 {
		t.Errorf("expected high confidence, got %.2f", forecasts[0].Confidence)
	}
	if forecasts[0].Confidence > 0.95 {
		t.Errorf("confidence should be capped at 0.95, got %.2f", forecasts[0].Confidence)
	}
}

func TestForecast_Empty(t *testing.T) {
	forecasts := Forecast(nil, forecastBase)
	if len(forecasts) != 0 {
		t.Errorf("expected empty forecasts, got %d", len(forecasts))
	}
}

func TestForecastByHost_Filters(t *testing.T) {
	entries := []Entry{
		forecastEntry("host-a", 80, EventOpened, forecastBase),
		forecastEntry("host-a", 80, EventOpened, forecastBase.Add(24*time.Hour)),
		forecastEntry("host-a", 80, EventOpened, forecastBase.Add(48*time.Hour)),
		forecastEntry("host-b", 22, EventOpened, forecastBase),
		forecastEntry("host-b", 22, EventOpened, forecastBase.Add(24*time.Hour)),
		forecastEntry("host-b", 22, EventOpened, forecastBase.Add(48*time.Hour)),
	}

	now := forecastBase.Add(48 * time.Hour)
	forecasts := ForecastByHost(entries, "host-a", now)

	for _, f := range forecasts {
		if f.Host != "host-a" {
			t.Errorf("expected only host-a, got %s", f.Host)
		}
	}
}
