package model

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type BaseResponse struct {
	Status Status `json:"status"`
	Error  error  `json:"error"`
}

type ProjectsResponse struct {
	Response BaseResponse `json:"response"`
	Projects []*Project   `json:"projects"`
}

func PrepareResponse(code int, description string, err error) BaseResponse {
	baseResponse := BaseResponse{}
	baseResponse.Error = err
	baseResponse.Status.Code = code
	baseResponse.Status.Description = description
	return baseResponse
}

func SuccessResponse() BaseResponse {
	return PrepareResponse(200, "Success", nil)
}
