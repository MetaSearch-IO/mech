package bandcamp

import (
   "bytes"
   "encoding/json"
   "github.com/89z/mech"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

const Origin = "http://bandcamp.com"

var Verbose = mech.Verbose

type Track struct {
   Bandcamp_URL string
}

func (t *Track) Get(id int) error {
   req, err := http.NewRequest(
      "GET", Origin + "/api/mobile/24/tralbum_details", nil,
   )
   if err != nil {
      return err
   }
   val := url.Values{
      "band_id": {"1"},
      "tralbum_id": {
         strconv.Itoa(id),
      },
      "tralbum_type": {"t"},
   }
   req.URL.RawQuery = val.Encode()
   res, err := mech.RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   return json.NewDecoder(res.Body).Decode(t)
}

func (t *Track) Post(id int) error {
   body := map[string]string{
      "band_id": "1",
      "tralbum_id": strconv.Itoa(id),
      "tralbum_type": "t",
   }
   buf := new(bytes.Buffer)
   if err := json.NewEncoder(buf).Encode(body); err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", Origin + "/api/mobile/24/tralbum_details", buf,
   )
   if err != nil {
      return err
   }
   res, err := mech.RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   return json.NewDecoder(res.Body).Decode(t)
}

func (t *Track) PostForm(id int) error {
   val := url.Values{
      "band_id": {"1"},
      "tralbum_id": {
         strconv.Itoa(id),
      },
      "tralbum_type": {"t"},
   }
   req, err := http.NewRequest(
      "POST", Origin + "/api/mobile/24/tralbum_details",
      strings.NewReader(val.Encode()),
   )
   if err != nil {
      return err
   }
   res, err := mech.RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   return json.NewDecoder(res.Body).Decode(t)
}
