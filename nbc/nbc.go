package nbc

import (
   "crypto/hmac"
   "crypto/sha256"
   "encoding/hex"
   "encoding/json"
   "errors"
   "github.com/89z/format"
   "github.com/89z/mech"
   "io"
   "net/http"
   "strconv"
   "strings"
   "time"
)

const persisted_query = "872a3dffc3ae6cdb3dc69fe3d9a949b539de7b579e95b2942e68d827b1a6ec62"

var (
   Log format.Log
   secret_key = []byte("2b84a073ede61c766e4c0b3f1e656f7f")
)

func authorization() string {
   now := strconv.FormatInt(time.Now().UnixMilli(), 10)
   buf := new(strings.Builder)
   buf.WriteString("NBC-Security key=android_nbcuniversal,version=2.4")
   buf.WriteString(",time=")
   buf.WriteString(now)
   buf.WriteString(",hash=")
   mac := hmac.New(sha256.New, secret_key)
   io.WriteString(mac, now)
   hex.NewEncoder(buf).Write(mac.Sum(nil))
   return buf.String()
}

func New_Bonanza_Page(guid int64) (*Bonanza_Page, error) {
   var p page_request
   p.Extensions.Persisted_Query.SHA_256_Hash = persisted_query
   p.Variables.App = "nbc"
   p.Variables.Name = strconv.FormatInt(guid, 10)
   p.Variables.One_App = true
   p.Variables.Platform = "android"
   p.Variables.Type = "VIDEO"
   buf, err := mech.Encode(p)
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://friendship.nbc.co/v2/graphql", buf,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/json")
   Log.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   var page struct {
      Data struct {
         BonanzaPage Bonanza_Page
      }
   }
   if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
      return nil, err
   }
   return &page.Data.BonanzaPage, nil
}

func (b Bonanza_Page) Base() string {
   return mech.Clean(b.Analytics.ConvivaAssetName)
}

type Bonanza_Page struct {
   Analytics struct {
      ConvivaAssetName string
   }
   Metadata struct {
      MpxAccountId string
   }
   Name string
}

func (b Bonanza_Page) Video() (*Video, error) {
   var in video_request
   in.Device = "android"
   in.Device_ID = "android"
   in.External_Advertiser_ID = "NBC"
   in.Mpx.Account_ID = b.Metadata.MpxAccountId
   body, err := mech.Encode(in)
   if err != nil {
      return nil, err
   }
   var addr strings.Builder
   addr.WriteString("http://access-cloudpath.media.nbcuni.com")
   addr.WriteString("/access/vod/nbcuniversal/")
   addr.WriteString(b.Name)
   req, err := http.NewRequest("POST", addr.String(), body)
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Authorization": {authorization()},
      "Content-Type": {"application/json"},
   }
   Log.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   out := new(Video)
   if err := json.NewDecoder(res.Body).Decode(out); err != nil {
      return nil, err
   }
   return out, nil
}

type Video struct {
   ManifestPath string // this is only valid for one minute
}

type page_request struct {
   Extensions struct {
      Persisted_Query struct {
         SHA_256_Hash string `json:"sha256Hash"`
      } `json:"persistedQuery"`
   } `json:"extensions"`
   Variables struct {
      App string `json:"app"`
      Name string `json:"name"` // String cannot represent a non string value
      One_App bool `json:"oneApp"`
      Platform string `json:"platform"`
      Type string `json:"type"`
      User_ID string `json:"userId"` // can be empty
   } `json:"variables"`
}

type video_request struct {
   Device string `json:"device"`
   Device_ID string `json:"deviceId"`
   External_Advertiser_ID string `json:"externalAdvertiserId"`
   Mpx struct {
      Account_ID string `json:"accountId"`
   } `json:"mpx"`
}
