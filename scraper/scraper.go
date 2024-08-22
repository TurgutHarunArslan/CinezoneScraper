package scraper

import (
	"cinezonescraper/classes"
	"cinezonescraper/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"regexp"
	"github.com/PuerkitoBio/goquery"
)

var BaseUrl string = "https://cinezone.to"

func Search(query string) ([]classes.Movie,error) {
	res,err := http.Get(BaseUrl + "/filter?keyword=" + query + "&sort=trending")
	if err != nil {
		return nil,errors.New("No response")
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil,errors.New("Invalid Html")
	}

	
	movies := []classes.Movie{}
	doc.Find("div.item").Each(func(i int, s *goquery.Selection) {
		a := s.Find("a.title")
		title := a.Text()
		link,_ := a.Attr("href")
		img,_ := s.Find("img").Attr("data-src")
		id,_ := s.Find("[data-tip]").Attr("data-tip")
		id = strings.Split(id,"?")[0]
		if err != nil{
			return
		}
		movie := classes.Movie{
			Id: id,
			Title:title,
			Img:img,
			Link:BaseUrl + link,
		}
		movies = append(movies, movie)
	})
	return movies,nil
}


func GetServers(id string) ([]string,error){
	res,err := http.Get(BaseUrl + "/ajax/server/list/" + id + "?vrf=" + utils.Vrf_encrypt(id))		
	if err != nil {
		return nil,errors.New("Request Error")	
	}

	defer res.Body.Close()

	var jsonData classes.ServerJson
    if err := json.NewDecoder(res.Body).Decode(&jsonData); err != nil {
        return nil,errors.New("Failed to parse JSON")
    }

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(jsonData.Result.(string)))
	if err != nil{
		return nil,errors.New("Document Parse Error")
	}
	
	var server_ids []string
	doc.Find("span").Each(func(i int, s *goquery.Selection) {
		id,exist := s.Attr("data-link-id")
		if exist == true{
			server_ids = append(server_ids,id)
		}
	})
	return server_ids,nil
}

func extractURL(s string) (string, error) {
    re := regexp.MustCompile(`(?i)\{.*"status":\s*200,\s*"result":\s*\{\s*"url":\s*"([^"]*)".*\}`)
    matches := re.FindStringSubmatch(s)
    if len(matches) < 2 {
        return "", errors.New("Url Not Found")
    }
    return matches[1], nil
}

func GetEps(Id string) ([]classes.Episodes,error){
	res,err := http.Get(BaseUrl + "/ajax/episode/list/" + Id + "?vrf=" + utils.Vrf_encrypt(Id))		
	if err != nil {
		return nil,errors.New("Request Error")	
	}

	defer res.Body.Close()

	var jsonData classes.ServerJson
    if err := json.NewDecoder(res.Body).Decode(&jsonData); err != nil {
        return nil,errors.New("Failed to parse JSON")
    }

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(jsonData.Result.(string)))
	if err != nil{
		return nil,errors.New("Document Parse Error")
	}
	
	var Eps []classes.Episodes
	doc.Find("ul").Each(func(i int, s *goquery.Selection) {
		seasonNumber,_ := s.Attr("data-season")
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			id,has := s.Attr("data-id")
			if has == false {
				return
			}
			epnumber,_ := s.Attr("data-num")
			ep := classes.Episodes{
				Id: id,
				EpNumber: epnumber,
				Season: seasonNumber,
			}
			Eps = append(Eps, ep)
		})
	})
	return Eps,nil
}

func GetServerLink(id string) (string,error){
	res,err := http.Get(BaseUrl + "/ajax/server/" + id + "?vrf=" + utils.Vrf_encrypt(id))	
	if err != nil{
		return "",errors.New("Request Error")
	}
	defer res.Body.Close()

	servertext,_ := io.ReadAll(res.Body)
	ServerTextStr := string(servertext)
	ServerTextStr,_ = extractURL(ServerTextStr)
	url := utils.Vrf_decrypt(ServerTextStr)
	//url = strings.Split(url,"&")[0]
	return url,nil
}
