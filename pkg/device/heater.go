package device

import (
	"fmt"
	"io"
	"net/http"
)

type HeaterService struct {
	BaseDeviceService
	HEATER_TOGGLE int
}

func NewHeaterService(HOST string, PORT int, API string, HEATER_TOGGLE int) HeaterService {
	return HeaterService{
		BaseDeviceService{
			HOST: HOST,
			PORT: PORT,
			API:  API,
		},
		HEATER_TOGGLE,
	}
}

func (s *HeaterService) PowerOn() (any, error) {
	return s.changePower(true)
}

func (s *HeaterService) PowerOff() (any, error) {
	return s.changePower(false)
}

func (s *HeaterService) changePower(on bool) (any, error) {
	deviceUrl, _ := s.GetDeviceUrl()

	uri := fmt.Sprintf("%s/Switch.Set?id=0&on=%t", deviceUrl, on)
	if on {
		uri = fmt.Sprintf("%s&toggle_after=%d", uri, s.HEATER_TOGGLE)
	}
	resp, err := http.Get(uri)

	if err != nil {
		return "", err
	}

	// var x any = interface{}
	// err = json.NewDecoder(resp.Body).Decode(x)
	bytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
