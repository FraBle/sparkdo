package main

type DigitalOceanResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
	Uid          int    `json:"uid,omitempty"`
	Info         Info   `json:"info,omitempty"`
}

type Info struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type Droplets struct {
	Droplets []Droplet `json:"droplets,omitempty"`
}

type Droplet struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
	Id     int    `json:"id,omitempty"`
}

type SingleDroplet struct {
	Droplet Droplet `json:"droplet,omitempty"`
}
