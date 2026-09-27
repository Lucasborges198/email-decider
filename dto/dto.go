package dto

// clientMessage represents the request body that Flask Laya service ll receive
type ClientMessageRequest struct {
	DecisionType string `json:"decisionType"`
	Message string `json:"message"`
}

type LayaResponse struct {
  Answer DataResponse `json:"answer"`
}

type DataResponse struct {
	Choice string `json:"choice"`
    Confidence int `json:"confidence"`
    Status string `json:"status"`
    Probabilities Probabilities `json:"probabilities"`
}

type Probabilities struct {
	Cancellment int `json:"cancellment"`
    Support int `json:"support"`
    FollowUpMeeting int `json:"follow_up_metting"`
    NewSoftware int `json:"new_software"`
    NoneOfTheOptions int `json:"none_of_the_options"`
}

func (d ClientMessageRequest) ValidateRequest() bool {
    return d.DecisionType != "" && d.Message != ""
}

