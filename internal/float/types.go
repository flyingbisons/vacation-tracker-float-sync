package float

type TimeOff struct {
	ID          int      `json:"timeoff_id,omitempty"`
	TypeID      int      `json:"timeoff_type_id" validate:"required"`
	TypeName    string   `json:"timeoff_type_name,omitempty"`
	StartDate   string   `json:"start_date" validate:"required"`
	EndDate     string   `json:"end_date" validate:"required"`
	StartTime   string   `json:"start_time,omitempty"`
	Hours       float64  `json:"hours,omitempty" validate:"gte=0,lt=24"`
	Notes       string   `json:"timeoff_notes"`
	ModifiedBy  int      `json:"modified_by,omitempty"`
	CreatedBy   int      `json:"created_by,omitempty"`
	Created     string   `json:"created,omitempty"`
	Modified    string   `json:"modified,omitempty"`
	RepeatState int      `json:"repeat_state,omitempty"`
	FullDay     int      `json:"full_day,omitempty" validate:"oneof=1 0"`
	PeopleIDs   []string `json:"people_ids" validate:"required,gt=0,dive,required"`
	Status      int      `json:"status,omitempty"`
}

type TimeOffType struct {
	ID   int    `json:"timeoff_type_id"`   //The ID of this time off type ,
	Name string `json:"timeoff_type_name"` //The name of this time off type. Note that the name "*" is reserved ,
}

type Employee struct {
	ID         int    `json:"people_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	JobTitle   string `json:"job_title"`
	Department struct {
		DepartmentId int    `json:"department_id"`
		ParentId     int    `json:"parent_id"`
		Name         string `json:"name"`
	} `json:"department"`
	Notes string `json:"notes"`
}
