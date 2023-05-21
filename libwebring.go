package libwebring

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Data contains the webring data.
type Data struct {
	Version int    `json:"version"` // 1
	Name    string `json:"name,omitempty"`
	Root    string `json:"root,omitempty"`
	Ring    Ring   `json:"ring"`
}

func (w Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Version int `json:"version"`
		Data
	}{
		Version: 1,
		Data:    w,
	})
}

func (w *Data) UnmarshalJSON(b []byte) error {
	type raw Data
	if err := json.Unmarshal(b, (*raw)(w)); err != nil {
		return err
	}
	if w.Version != 1 {
		return errors.New("only version 1 is supported")
	}
	return nil
}

// Ring is a ring of links.
type Ring []Link

// Filter filters the ring with the given filter function.
func (r Ring) Filter(f func(Link) bool) Ring {
	filtered := make(Ring, 0, len(r))
	for _, link := range r {
		if f(link) {
			filtered = append(filtered, link)
		}
	}
	return filtered
}

// ExcludeAnomalies excludes the anomalies from the ring.
func (r Ring) ExcludeAnomalies(anomalies Anomalies) Ring {
	return r.Filter(func(link Link) bool {
		_, ok := anomalies[link.Link]
		return !ok
	})
}

// SurroundingLinks returns the links surrounding the given link.
// It returns empty links if the link is not found.
func (r Ring) SurroundingLinks(link Link) (Link, Link) {
	for i, l := range r {
		if l == link {
			return r.SurroundingIndex(i)
		}
	}
	return Link{}, Link{}
}

// SurroundingIndex returns the index of the links surrounding the given index.
// If the index is out of bounds, it returns empty links.
func (r Ring) SurroundingIndex(i int) (Link, Link) {
	if len(r) == 0 || i < 0 || i >= len(r) {
		return Link{}, Link{}
	}

	var prev, next Link

	if i > 0 {
		prev = r[i-1]
	} else {
		prev = r[len(r)-1]
	}

	if i < len(r)-1 {
		next = r[i+1]
	} else {
		next = r[0]
	}

	return prev, next
}

// Link is a link in the webring.
type Link struct {
	Name string `json:"name"`
	Link string `json:"link"` // may not have a scheme
}

// Anomalies is a map of anomalies. It only contains the links that are not
// working.
type Anomalies map[string]LinkStatus

// StatusData contains the status data of the webring. It only contains the
// anomalies.
type StatusData struct {
	Version   int       `json:"version"` // 1
	Anomalies Anomalies `json:"anomalies"`
}

func (w StatusData) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Version int `json:"version"`
		StatusData
	}{
		Version:    1,
		StatusData: w,
	})
}

func (w *StatusData) UnmarshalJSON(b []byte) error {
	type raw StatusData
	if err := json.Unmarshal(b, (*raw)(w)); err != nil {
		return err
	}
	if w.Version != 1 {
		return errors.New("only version 1 is supported")
	}
	return nil
}

// LinkStatus contains the status of a link.
type LinkStatus struct {
	Dead           bool `json:"dead,omitempty"`
	MissingWebring bool `json:"missingWebring,omitempty"`
}

// FetchData fetches the webring data from the given URL.
func FetchData(ctx context.Context, webringURL string) (*Data, error) {
	return fetchJSON[Data](ctx, webringURL)
}

// FetchStatusForWebring fetches the status data for the given webring URL.
// It tries to guess the status URL from the webring URL.
func FetchStatusForWebring(ctx context.Context, webringURL string) (*StatusData, error) {
	return fetchJSON[StatusData](ctx, GuessStatusURL(webringURL))
}

// FetchStatus fetches the status data from the given URL.
func FetchStatus(ctx context.Context, statusURL string) (*StatusData, error) {
	return fetchJSON[StatusData](ctx, statusURL)
}

func fetchJSON[T any](ctx context.Context, url string) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("failed to fetch webring data: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch webring data: " + resp.Status)
	}

	var data T
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("failed to decode webring data: " + err.Error())
	}

	return &data, nil
}

// GuessStatusURL returns the URL of the status file for the given webring URL.
func GuessStatusURL(webringURL string) string {
	return strings.Replace(webringURL, ".json", ".status.json", 1)
}
