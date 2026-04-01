package adb

import (
	"fmt"
	"strings"
	"time"
)

// PairResult holds the outcome of a wireless ADB pairing attempt.
type PairResult struct {
	Success bool
	Message string
}

// Pair initiates wireless debugging pairing with an Android device.
// The device must have its pairing server running (e.g. after scanning a QR code
// or selecting "Pair device with pairing code" in Developer Options).
//
// Parameters:
//   - host: IP address of the device's pairing server
//   - port: TCP port of the device's pairing server
//   - password: pairing password (6-digit code or 10-digit QR password)
//
// This sends the host:pair command to the ADB server, which handles the
// SPAKE2+TLS handshake with the device.
func (c *Adb) Pair(host string, port int, password string) (*PairResult, error) {
	req := fmt.Sprintf("host:pair:%s:%s:%d", password, host, port)

	resp, err := roundTripSingleResponseTimeout(c.server, req, 30*time.Second)
	if err != nil {
		// Check if the error message indicates success (some ADB versions
		// return the success message as part of the response).
		errStr := err.Error()
		if strings.Contains(errStr, "Successfully paired") {
			return &PairResult{
				Success: true,
				Message: errStr,
			}, nil
		}
		return &PairResult{
			Success: false,
			Message: fmt.Sprintf("pair failed: %v", err),
		}, nil
	}

	respStr := string(resp)
	if strings.Contains(respStr, "Successfully paired") {
		return &PairResult{
			Success: true,
			Message: respStr,
		}, nil
	}

	return &PairResult{
		Success: false,
		Message: respStr,
	}, nil
}
