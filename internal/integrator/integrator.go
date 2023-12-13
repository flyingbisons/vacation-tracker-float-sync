package integrator

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/flyingbisons/vacation-tracker-float-sync/internal/float"
	"github.com/flyingbisons/vacation-tracker-float-sync/internal/vacation"
)

type Integration struct {
	store       RequestRepository
	floatClient *float.Client
	vtClient    *vacation.Client
	log         *slog.Logger

	vtUsers     map[string]string
	currentTime time.Time
}

func New(store RequestRepository, floatClient *float.Client, vtClient *vacation.Client, log *slog.Logger) (*Integration, error) {
	manager := &Integration{
		store:       store,
		floatClient: floatClient,
		vtClient:    vtClient,
		log:         log,
	}

	return manager, nil
}

func (i *Integration) getVTUsers(ctx context.Context) error {
	var err error
	i.vtUsers, err = i.vtClient.Users(ctx)
	if err != nil {
		return fmt.Errorf("unable to get VT users: %w", err)
	}
	return nil
}

func (i *Integration) Sync(ctx context.Context) error {
	//get VT users
	if i.vtUsers == nil {
		if err := i.getVTUsers(ctx); err != nil {
			return err
		}
	}

	vtRequests, err := i.vtClient.LeaveRequests(ctx)
	i.log.Debug("got leave requests", slog.Int("count", len(vtRequests)))
	if err != nil {
		return fmt.Errorf("unable to get VT leave requests: %w", err)
	}

	for _, vtRequest := range vtRequests {
		_, err := i.store.GetRequest(vtRequest.ID)
		if err != nil {
			if errors.Is(err, ErrorRequestNotFound) {
				err := i.createRequest(ctx, vtRequest)
				if err != nil {
					i.log.Error("unable to create request", slog.Any("error", err))
				}
				continue
			}
			i.log.Error("unable to get leave request from db", slog.Any("error", err))
		}
	}

	return nil
}

func (i *Integration) createRequest(ctx context.Context, vtRequest vacation.Leave) error {
	userEmail, ok := i.vtUsers[vtRequest.UserID]
	if !ok {
		return fmt.Errorf("unable to find user with id %s", vtRequest.UserID)
	}

	approverEmail, ok := i.vtUsers[vtRequest.ApproverID]
	if !ok {
		return fmt.Errorf("unable to find approver with id %s", vtRequest.ApproverID)
	}

	//get user id from float
	employee, err := i.floatClient.FindEmployeeByEmail(ctx, userEmail)
	if err != nil {
		return err
	}

	isFullDay := 1
	if vtRequest.IsPartDay {
		isFullDay = 0
	}

	var startHour string
	var hours float64
	if vtRequest.StartHour != nil && vtRequest.EndHour != nil {
		startHour = fmt.Sprintf("%02d:00", *vtRequest.StartHour)
		hours = float64(*vtRequest.EndHour - *vtRequest.StartHour)
	}

	typeID := float.TimeTypePaidID
	switch vtRequest.LeaveTypeID {
	case vacation.LeaveTypeFreeTime:
		typeID = float.TimeTypeFreeID
	}

	//create request in float
	timeOff := float.TimeOff{
		TypeID:    typeID,
		FullDay:   isFullDay,
		StartDate: vtRequest.StartDate,
		EndDate:   vtRequest.EndDate,
		StartTime: startHour,
		Hours:     hours,
		PeopleIDs: []string{strconv.Itoa(employee.ID)},
		Notes:     fmt.Sprintf("Created from vacation tracker, accepted by %s", approverEmail),
	}
	slog.Debug("creating time off", slog.Any("timeOff", timeOff.TypeID))
	floatTimeOff, err := i.floatClient.AddTimeOff(ctx, timeOff)

	if err != nil {
		return err
	}

	//create request in db
	err = i.store.CreateRequest(vtRequest.ID, int64(floatTimeOff.ID))
	if err != nil {
		return err
	}

	return err
}
