package user

import (
	"strconv"

	"github.com/Sabbir185/gpi/entity"
)

type DesignationResponse struct {
	ID              *string `json:"id"`
	Title           string  `json:"title"`
	Status          string  `json:"status"`
	RejectionReason *string `json:"rejectionReason"`
}

type UserResponse struct {
	ID              string               `json:"id"`
	FullName        string               `json:"fullName"`
	Phone           *string              `json:"phone"`
	Email           *string              `json:"email"`
	AvatarURL       *string              `json:"avatarUrl"`
	Designation     *DesignationResponse `json:"designation"`
	Country         *string              `json:"country"`
	ProfileComplete bool                 `json:"profileComplete"`
}

type DesignationChoice struct {
	ID    *string `json:"id" validate:"omitempty,numeric"`
	Title string  `json:"title" validate:"required,min=2,max=60"`
}

type UpdateProfileRequest struct {
	FullName    string             `json:"fullName" validate:"required,min=2,max=80"`
	Designation *DesignationChoice `json:"designation"`
	Country     string             `json:"country" validate:"required,iso3166_1_alpha2"`
	AvatarURL   *string            `json:"avatarUrl" validate:"omitempty,https_url,max=2048"`
}

func ToUserResponse(user *entity.User) *UserResponse {
	response := &UserResponse{
		ID:              strconv.FormatInt(user.ID, 10),
		FullName:        user.Name,
		Phone:           user.Phone,
		Email:           user.Email,
		AvatarURL:       user.ProfileImg,
		Country:         user.CountryCode,
		ProfileComplete: user.ProfileCompletedAt != nil,
	}
	if user.DesignationTitle != nil && user.DesignationStatus != nil {
		var id *string
		if user.DesignationID != nil {
			value := strconv.FormatInt(*user.DesignationID, 10)
			id = &value
		}
		response.Designation = &DesignationResponse{
			ID:     id,
			Title:  *user.DesignationTitle,
			Status: string(*user.DesignationStatus),
		}
	}
	return response
}
