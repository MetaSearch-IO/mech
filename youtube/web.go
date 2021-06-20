package youtube
import "encoding/json"

type Microformat struct {
   PlayerMicroformatRenderer `json:"playerMicroformatRenderer"`
}

type PlayerMicroformatRenderer struct {
   AvailableCountries []string
   PublishDate string
}

type Web struct {
   Microformat `json:"microformat"`
   VideoDetails `json:"videoDetails"`
}

func NewWeb(id string) (Web, error) {
   res, err := post(id, "WEB", "1.19700101")
   if err != nil {
      return Web{}, err
   }
   defer res.Body.Close()
   var w Web
   if err := json.NewDecoder(res.Body).Decode(&w); err != nil {
      return Web{}, err
   }
   return w, nil
}
