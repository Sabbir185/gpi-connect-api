package user

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Sabbir185/gpi/entity"
	"github.com/Sabbir185/gpi/internal/designation"
	"github.com/Sabbir185/gpi/internal/device"
	"github.com/Sabbir185/gpi/pkg/httpx"
)

type UserService interface {
	Me(ctx context.Context, userID int64) (*UserResponse, error)
	UpdateProfile(ctx context.Context, userID int64, payload *UpdateProfileRequest) (*UserResponse, error)
	Devices(ctx context.Context, userID int64, sessionID string) ([]device.DeviceResponse, error)
	RegisterDevice(ctx context.Context, userID int64, payload *device.DeviceRequest) error
	SignOutDevice(ctx context.Context, userID int64, deviceID string) error
}

type userService struct {
	userRepo        UserRepository
	deviceRepo      device.DeviceRepository
	designationRepo designation.DesignationRepository
}

func NewUserService(userRepo UserRepository, deviceRepo device.DeviceRepository, designationRepo designation.DesignationRepository) UserService {
	return &userService{userRepo: userRepo, deviceRepo: deviceRepo, designationRepo: designationRepo}
}

func (s *userService) Me(ctx context.Context, userID int64) (*UserResponse, error) {
	account, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, httpx.NewAppError(http.StatusNotFound, httpx.CodeNotFound, "User not found", nil)
	}
	return ToUserResponse(account), nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID int64, payload *UpdateProfileRequest) (*UserResponse, error) {
	account, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, httpx.NewAppError(http.StatusNotFound, httpx.CodeNotFound, "User not found", nil)
	}

	profile := &ProfileUpdate{
		Name:        strings.TrimSpace(payload.FullName),
		CountryCode: payload.Country,
		ProfileImg:  payload.AvatarURL,
	}
	if payload.Designation != nil {
		if err := s.applyDesignation(ctx, account, payload.Designation, profile); err != nil {
			return nil, err
		}
	}

	if err := s.userRepo.UpdateProfile(ctx, userID, profile); err != nil {
		return nil, err
	}
	return s.Me(ctx, userID)
}

// applyDesignation trusts the catalogue for picked titles. A typed-in title waits for
// review, unless it is the one the user already has, which keeps its current status.
func (s *userService) applyDesignation(ctx context.Context, account *entity.User, choice *DesignationChoice, profile *ProfileUpdate) error {
	if choice.ID != nil {
		id, err := strconv.ParseInt(*choice.ID, 10, 64)
		if err != nil {
			return errUnknownDesignation()
		}
		found, err := s.designationRepo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if found == nil {
			return errUnknownDesignation()
		}
		status := entity.DesignationStatusApproved
		profile.DesignationID = &found.ID
		profile.DesignationTitle = &found.Title
		profile.DesignationStatus = &status
		return nil
	}

	title := strings.TrimSpace(choice.Title)
	status := entity.DesignationStatusPending
	if account.DesignationID == nil && account.DesignationTitle != nil && account.DesignationStatus != nil &&
		strings.EqualFold(*account.DesignationTitle, title) {
		status = *account.DesignationStatus
	}
	profile.DesignationTitle = &title
	profile.DesignationStatus = &status
	return nil
}

func (s *userService) Devices(ctx context.Context, userID int64, sessionID string) ([]device.DeviceResponse, error) {
	currentDeviceID, err := s.deviceRepo.SessionDeviceID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	devices, err := s.deviceRepo.ListActiveForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]device.DeviceResponse, 0, len(devices))
	for _, d := range devices {
		name := d.Platform
		if d.OSVersion != nil {
			name += " " + *d.OSVersion
		}
		result = append(result, device.DeviceResponse{
			ID:           strconv.FormatInt(d.ID, 10),
			Name:         name,
			Platform:     d.Platform,
			LastActiveAt: d.LastActiveAt.UTC().Format(time.RFC3339),
			Current:      currentDeviceID != nil && *currentDeviceID == d.ID,
		})
	}
	return result, nil
}

func (s *userService) RegisterDevice(ctx context.Context, userID int64, payload *device.DeviceRequest) error {
	_, err := s.deviceRepo.Upsert(ctx, userID, payload, false)
	if errors.Is(err, device.ErrOwnedByAnotherUser) {
		return httpx.NewAppError(http.StatusForbidden, httpx.CodeForbidden, "This device is registered to another account", nil)
	}
	return err
}

func (s *userService) SignOutDevice(ctx context.Context, userID int64, deviceID string) error {
	id, err := strconv.ParseInt(deviceID, 10, 64)
	if err != nil {
		return errDeviceNotFound()
	}
	revoked, err := s.deviceRepo.RevokeSessions(ctx, userID, id)
	if err != nil {
		return err
	}
	if revoked == 0 {
		return errDeviceNotFound()
	}
	return s.deviceRepo.ClearPushToken(ctx, id)
}

func errUnknownDesignation() error {
	return httpx.NewAppError(http.StatusBadRequest, httpx.CodeValidationError, "Unknown designation", nil)
}

func errDeviceNotFound() error {
	return httpx.NewAppError(http.StatusNotFound, httpx.CodeNotFound, "Device not found", nil)
}
