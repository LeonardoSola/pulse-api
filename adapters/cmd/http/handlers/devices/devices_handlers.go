package device_handler

import (
	"net/http"

	"github.com/LeonardoSola/wagon-api/adapters/bootstrap"
	"github.com/LeonardoSola/wagon-api/adapters/cmd/http/utils"
	device_service "github.com/LeonardoSola/wagon-api/internal/services/device"
	chi "github.com/go-chi/chi/v5"
)

func InitDevices(router *chi.Mux) {
	router.Get("/devices", getDevices)
	router.Get("/devices/{id}", getDevice)
	router.Post("/devices", createDevice)
	router.Put("/devices/{id}", updateDevice)
	router.Delete("/devices/{id}", deleteDevice)

	router.Get("/devices/{id}/qr", getDeviceQR)
}

func getDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := bootstrap.Services.DeviceService.GetDevices(r.Context())
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	utils.SendResponse(devices, http.StatusOK, w)
}

func getDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")

	device, err := bootstrap.Services.DeviceService.GetDevice(r.Context(), deviceID)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	utils.SendResponse(device, http.StatusOK, w)
}

func createDevice(w http.ResponseWriter, r *http.Request) {
	var payload device_service.CreateDeviceDTO

	err := utils.DecodeBody(r, &payload)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	device, err := bootstrap.Services.DeviceService.CreateDevice(r.Context(), payload)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	utils.SendResponse(device, http.StatusServiceUnavailable, w)
}

func updateDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")

	var payload device_service.UpdateDeviceDTO

	err := utils.DecodeBody(r, &payload)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	err = bootstrap.Services.DeviceService.UpdateDevice(r.Context(), deviceID, payload)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	utils.SendResponse("success", http.StatusServiceUnavailable, w)
}

func deleteDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")

	err := bootstrap.Services.DeviceService.DeleteDevice(r.Context(), deviceID)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	utils.SendResponse("success", http.StatusOK, w)
}

func getDeviceQR(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")

	qrChan, err := bootstrap.Services.DeviceService.ConnectDevice(r.Context(), deviceID)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	utils.StreamSSE(w, r, qrChan)
}
