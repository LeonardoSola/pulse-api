package device

type CreateDeviceDTO struct {
	Name string `json:"name"`
}

type UpdateDeviceDTO struct {
	Name *string `json:"name"`
}
