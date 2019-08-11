package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	swagger "github.com/jedruniu/plotted/swagger-generated"

	"github.com/antihax/optional"
	gopoly "github.com/twpayne/go-polyline"

	"golang.org/x/oauth2"

	"net/http"
)

var (
	stravaClientID = flag.String("strava_clientID", "", "Strava client ID")
	stravaSecret   = flag.String("strava_secret", "", "Strava Secret")
	mapBoxToken    = flag.String("mapbox", "", "Mapbox API Access token")

	layout = "02/01/2006"
)

func init() {
	flag.Parse()
}

var code string
var token string
var state string

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	ctx := context.Background()

	conf := &oauth2.Config{
		ClientID:     *stravaClientID,
		ClientSecret: *stravaSecret,
		Scopes:       []string{"activity:write,activity:read_all,profile:read_all"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.strava.com/oauth/authorize",
			TokenURL: "https://www.strava.com/oauth/token",
		},
		RedirectURL: "http://localhost:8888/auth_callback",
	}

	http.HandleFunc("/auth_callback", func(w http.ResponseWriter, r *http.Request) {
		code = r.URL.Query().Get("code")
		callbackState := r.URL.Query().Get("state")
		if callbackState != state {
			http.Error(w, fmt.Sprintf("state verification failed"), http.StatusBadRequest)
			return
		}

		tok, err := conf.Exchange(ctx, code)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not exchange ouath2 token, err: %v", err), http.StatusInternalServerError)
			return
		}
		token = tok.AccessToken

		http.Redirect(w, r, "http://localhost:8888/map?after=30/05/2019&before=30/09/2019", 302)
	})

	http.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {
		cfg := swagger.NewConfiguration()
		client := swagger.NewAPIClient(cfg)

		ctx = context.WithValue(ctx, swagger.ContextAccessToken, token)

		opts := swagger.GetLoggedInAthleteActivitiesOpts{}

		unparsedAfter := r.URL.Query().Get("after")
		unparsedBefore := r.URL.Query().Get("before")

		after, _ := time.Parse(layout, unparsedAfter)
		after = after.AddDate(0, 0, -1)
		before, _ := time.Parse(layout, unparsedBefore)
		before = before.AddDate(0, 0, 1)

		var activities []swagger.SummaryActivity

		for i := 1; i < 3; i++ {
			opts.After = optional.NewInt32(int32(after.Unix()))
			opts.Before = optional.NewInt32(int32(before.Unix()))
			opts.Page = optional.NewInt32(int32(i))
			opts.PerPage = optional.NewInt32(200)

			summary, resp, err := client.ActivitiesApi.GetLoggedInAthleteActivities(ctx, &opts)
			if err != nil {
				http.Error(w, err.Error(), resp.StatusCode)
				return
			}
			if len(summary) == 0 {
				break
			}
			activities = append(activities, summary...)
		}

		var polylines [][][]float64

		for _, activity := range activities {
			var polyline []byte

			cachedFileName := fmt.Sprintf("cache/%d.cache", activity.Id)
			cacheContent, err := ioutil.ReadFile(cachedFileName)

			if err != nil {
				detailed, _, err := client.ActivitiesApi.GetActivityById(ctx, activity.Id, nil)
				if err != nil {
					log.Printf("err for activity %d, err: %v", activity.Id, err)
					continue
				}
				if detailed.Map_.Polyline == "" {
					continue
				}
				polyline = []byte(detailed.Map_.Polyline)

				file, err := os.Create(cachedFileName)
				if err != nil {
					log.Printf("error when creating %s, err: %v", cachedFileName, err)
					continue
				}
				defer file.Close()
				_, err = file.Write(polyline)
				if err != nil {
					log.Printf("error when writing to %s, err: %v", cachedFileName, err)
				}
			} else {
				log.Printf("cache hit for file %s\n", cachedFileName)
				polyline = cacheContent
			}

			var polylineDecoded [][]float64

			polylineDecoded, _, err = gopoly.DecodeCoords(polyline)
			if err != nil {
				log.Printf("could not decode polyline from file %d, err: %v", activity.Id, err)
			} else {
				polylines = append(polylines, polylineDecoded)
			}

		}

		templ, _ := template.ParseFiles("index_tmpl.html")

		data := struct {
			EncodedRoutes [][][]float64
			MapboxToken   string
		}{
			polylines,
			*mapBoxToken,
		}
		templ.Execute(w, data)

	})
	state = uuid.New().String()
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	log.Println(url)
	templ, err := template.ParseFiles("static/index_tmpl.html")
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)

	data := struct{ Auth string }{url}
	_ = templ.Execute(buf, data)
	os.Remove("static/index.html")
	file, _ := os.Create("static/index.html")
	defer file.Close()
	file.Write(buf.Bytes())

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

	log.Fatal(http.ListenAndServe(":8888", nil))
}

type Storage interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Exists(string) (bool, error)
}

type FilesStorage struct {
	cache sync.Map
	prefix string
}

func (s *FilesStorage) Exists(key string) (bool, error) {
	cachedFileName := fmt.Sprintf("%s/%s", s.prefix, key)
	cacheContent, err := ioutil.ReadFile(cachedFileName)
	if err != nil {
		return false, err
	}
	s.cache.Store(key, cacheContent)
	return true, nil
}

func (s *FilesStorage) Get(key string) ([]byte, error) {
	cachedFileName := fmt.Sprintf("%s/%s", s.prefix, key)
	v, ok := s.cache.Load(cachedFileName)

	if !ok {
		cacheContent, err := ioutil.ReadFile(cachedFileName)
		if err != nil {
			return []byte{}, err
		} else {
			return cacheContent, nil
		}
	}

	content, assertOk := v.([]byte)
	if assertOk {
		return content, nil
	}
	return []byte{}, fmt.Errorf("🤷")
}

func (s *FilesStorage) Set(key string, value []byte) error {
	s.cache.Store(key, value)
	cachedFileName := fmt.Sprintf("%s/%s", s.prefix, key)
	file, err := os.Create(cachedFileName)
	if err != nil {
		return fmt.Errorf("error when creating %s, err: %v", cachedFileName, err)
	}
	defer file.Close()
	_, err = file.Write(value)
	if err != nil {
		return fmt.Errorf("error when writing to %s, err: %v", cachedFileName, err)
	}
	return nil
}

func NewFileStorage(cacheDir string) (*FilesStorage, error) {
	_, err := os.Create(cacheDir)
	if err != nil {
		if err != os.ErrExist {
			return nil, err
		}
	}
	return &FilesStorage{prefix:cacheDir}, nil
}
// try to move to app engine
// code clean up
//package main
//
//import (
//	"fmt"
//	"io/ioutil"
//	"os"
//)
//
//func main() {
//	cachedFileName := fmt.Sprintf("cache/%d.cache", 9686632701111)
//	_, err := ioutil.ReadFile(cachedFileName)
//	if  os.IsNotExist(err) {
//		fmt.Println(err)
//	} else {
//		fmt.Println("coś jebło")
//	}
//}