package serializer

type SendHeartbeatRequest struct {
	Username  string `json:"username"`
	SessionID string `json:"sessionId"`
}

type SendHeartbeatResponse struct {
	Data           string  `json:"data"`
	Message        *string `json:"message"`
	SessionTimeOut bool    `json:"sessionTimeOut"`
	Success        bool    `json:"success"`
}
