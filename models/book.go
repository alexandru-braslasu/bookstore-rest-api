package models

type TitleStruct struct {
	Title string
}

type Book struct {
	Title string
	Author string
	Description string
	NrSamples int
}

var Books []Book