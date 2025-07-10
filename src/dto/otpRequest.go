package dto

type GetOtpRequest struct {
	MobileNumber string `json:"mobileNumber" binding:"required"`
}
