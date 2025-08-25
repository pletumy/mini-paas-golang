package api

type CreateAppRequest struct {
	Name        string `json:"name" binding:"required"`
	GitURL      string `json:"git_url" binding:"required,url"`
	Description string `json:"description"`
}

type CreateAppResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type AppItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desciption"`
	GitURL      string `json:"git_url"`
	ImageURL    string `json:"image_url"`
	DeployURL   string `json:"deploy_url"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
