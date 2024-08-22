package classes

import "fmt"

type Movie struct{
	Id string
	Title string
	Img string
	Link string
}

type Episodes struct{
	Id string
	Season string
	EpNumber string
}

func (m *Movie) Print() {
	fmt.Printf("Id : %s\nTitle : %s \nImg : %s   \nLink : %s\n",m.Id,m.Title,m.Img,m.Link)
}

type ServerJson struct{
	Status int `json:"status"`
	Result interface{}`json:"result"`
}

type VideoData struct {
	Title string
	EpList []Episodes
}