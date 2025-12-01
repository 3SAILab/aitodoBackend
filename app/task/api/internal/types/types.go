package types

type CreateTaskReq struct {
	TypeId        string `json:"typeId"`                 // 任务类型ID
	Title         string `json:"title"`                  // 标题
	Description   string `json:"description,optional"`   // 描述 (可选)
	AssigneeId    string `json:"assigneeId,optional"`    // 负责人ID (可选)
	SalesPersonId string `json:"salesPersonId,optional"` // 销售ID (可选)
	DueDate       int64  `json:"dueDate,optional"`       // 截止时间(时间戳, 可选)
	Priority      int    `json:"priority,optional"`      // 排序权重 (可选)
}

type CreateTaskResp struct {
	Id string `json:"id"`
}

type DeleteTaskReq struct {
	Id string `path:"id"`
}

type CreateSalesReq struct {
	Name  string `json:"name"`
	Phone string `json:"phone,optional"`
}

type UpdateSalesReq struct {
	Id    string `path:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone,optional"`
}

type DeleteSalesReq struct {
	Id string `path:"id"`
}

type SalesResp struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type ListSalesResp struct {
	List []SalesResp `json:"list"`
}

type CreateTaskTypeReq struct {
	Name      string `json:"name"`
	ColorCode string `json:"colorCode"` // 例如 "#FF5733"
}

type UpdateTaskTypeReq struct {
	Id        string `path:"id"`
	Name      string `json:"name"`
	ColorCode string `json:"colorCode"`
}

type DeleteTaskTypeReq struct {
	Id string `path:"id"`
}

type TaskTypeResp struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	ColorCode string `json:"colorCode"`
}

type ListTaskTypeResp struct {
	List []TaskTypeResp `json:"list"`
}

type UpdateTaskReq struct {
	Id            string `path:"id"`
	TypeId        string `json:"typeId"`
	Title         string `json:"title"`
	Description   string `json:"description,optional"`
	AssigneeId    string `json:"assigneeId,optional"`
	SalesPersonId string `json:"salesPersonId,optional"`
	DueDate       int64  `json:"dueDate,optional"`
	Status        string `json:"status"` // 例如 TODO, DOING, DONE
	Priority      int    `json:"priority,optional"`
}

type TaskResp struct {
	Id            string `json:"id"`
	TypeId        string `json:"typeId"`
	CreatorId     string `json:"creatorId"`
	AssigneeId    string `json:"assigneeId"`
	SalesPersonId string `json:"salesPersonId"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	Priority      int    `json:"priority"`
	DueDate       int64  `json:"dueDate"`
	CreatedAt     int64  `json:"createdAt"`
	CompletedAt   int64  `json:"complete_at"`
}

type ListTaskResp struct {
	List []TaskResp `json:"list"`
}

type CreateTaskProgressReq struct {
	TaskId  string `path:"taskId"`
	Content string `json:"content"`
}

type TaskProgressResp struct {
	Id        string `json:"id"`
	TaskId    string `json:"taskId"`
	Content   string `json:"content"`
	CreatedBy string `json:"createdBy"`
	CreatedAt int64  `json:"createdAt"`
}

type ListTaskProgressReq struct {
	TaskId string `path:"taskId"`
}

type ListTaskProgressResp struct {
	List []TaskProgressResp `json:""list`
}
