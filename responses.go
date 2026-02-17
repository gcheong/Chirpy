package main

type chirpError struct {
	Err string `json:"error"`
}

type chirpValid struct {
	CleanedBody string `json:"cleaned_body"`
}
