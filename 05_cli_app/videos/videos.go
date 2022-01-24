package videos

import (
	"encoding/json"
	"io/ioutil"
)

type video struct {
	Id          string
	Title       string
	Description string
	ImageUrl    string
	Url         string
}

func GetVideos() []video {

	var videos []video

	fileBytes, error := ioutil.ReadFile("./videos.json")

	if error != nil {
		panic(error)
	}

	error = json.Unmarshal(fileBytes, &videos)

	if error != nil {
		panic(error)
	}

	return videos
}

func SaveVideos() {

}
