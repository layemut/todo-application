package model

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type BaseResponse struct {
	Status Status `json:"status"`
	Error  string `json:"error"`
}

type ProjectsResponse struct {
	Response BaseResponse `json:"response"`
	Projects []*Project   `json:"projects"`
}

type TasksResponse struct {
	Response BaseResponse `json:"response"`
	Tasks    []*Task      `json:"tasks"`
}

func PrepareResponse(code int, description string, err string) BaseResponse {
	baseResponse := BaseResponse{}
	baseResponse.Error = err
	baseResponse.Status.Code = code
	baseResponse.Status.Description = description
	return baseResponse
}

func SuccessResponse() BaseResponse {
	return PrepareResponse(200, "Success", "")
}

func BadRequestResponse(err error) BaseResponse {
	return PrepareResponse(400, "Bad Request", err.Error())
}
