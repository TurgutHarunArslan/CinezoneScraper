package api

import (
	"cinezonescraper/classes"
	"cinezonescraper/scraper"
	"cinezonescraper/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)
var ctx = context.Background()
var client = utils.RedisClient()
const portNum string = "0.0.0.0:3000"

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func Search(w http.ResponseWriter, r *http.Request){
    params := r.URL.Query()
    query := params.Get("q")
    if query == ""{
        http.Error(w,"Search Query Does Not exist",400)
        return
    }
    
    movies,err := scraper.Search(query)
    for _,x := range movies {
        client.Set(ctx,x.Id,x.Title,0)
    }
    if err != nil{
        http.Error(w,"No ResponseWriter",400)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(movies)
}

func Show(w http.ResponseWriter , r *http.Request){
    params := r.URL.Query()
    query := params.Get("TrackId")
    if query == ""{
        http.Error(w,"Track Query Does Not exist",400)
        return
    }
    eps,err := scraper.GetEps(query)
    if err != nil {
        http.Error(w,"This Show Does Not Exist",400)
        return
    }
    title,_ := client.Get(ctx,query).Result()
    data := classes.VideoData{
        Title: title,
        EpList: eps,
    }
    json.NewEncoder(w).Encode(data)
}

func Episode(w http.ResponseWriter, r *http.Request){
    params := r.URL.Query()
    query := params.Get("id")
    if query == ""{
        http.Error(w,"Id Query Does Not exist",400)
        return
    }
    servers, err := scraper.GetServers(query)
    if err != nil       {
        http.Error(w,"Error servers not found",404)
        return
    }
    serverUrl,err := scraper.GetServerLink(servers[0])

    if err != nil{
        http.Error(w,"Error server link not found",404)
        return
    }

    fmt.Fprint(w,serverUrl)
}

func Home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w,"Hello World")
}

func StartServer(){
	log.Println("Starting The Http server.")
    mux := http.NewServeMux()
	mux.HandleFunc("/",Home)
    mux.HandleFunc("/search",Search)
    mux.HandleFunc("/movie",Show)
    mux.HandleFunc("/episode",Episode)
	log.Println("Server Has Started at port",portNum)
    handlerWithCORS := corsMiddleware(mux)
    err := http.ListenAndServe(portNum, handlerWithCORS)
    if err != nil {
        log.Fatal(err)
    }
}