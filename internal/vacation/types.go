package vacation

import "time"

type Leavs struct {
	Status    string  `json:"status"`
	NextToken *string `json:"nextToken"`
	Data      []Leave `json:"data"`
}

type Leave struct {
	ID                   string    `json:"id"`
	UserID               string    `json:"userId"`
	ApproverID           string    `json:"approverId"`
	AutoApproved         bool      `json:"autoApproved"`
	DurationCalendarDays float64   `json:"calendarDays"`
	DurationWorkingDays  float64   `json:"workingDays"`
	StartDate            string    `json:"startDate"`
	EndDate              string    `json:"endDate"`
	IsPartDay            bool      `json:"isPartDay"`
	LeaveTypeID          string    `json:"leaveTypeID"`
	LocationID           string    `json:"locationID"`
	DepartmentID         string    `json:"departmentID"`
	Status               string    `json:"status"`
	CreatedAt            time.Time `json:"createdAt"`
	StartHour            *int64    `json:"partDayStartHour"`
	EndHour              *int64    `json:"partDayEndHour"`
}

type GetLeaveRequestByDate struct {
	LeaveRequests []Leave `json:"leaveRequests"`
}

type UsersResponse struct {
	Status    string  `json:"status"`
	Users     []User  `json:"data"`
	NextToken *string `json:"nextToken"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
