package bleep

import (
   "bytes"
   "encoding/json"
   "github.com/89z/mech"
   "github.com/89z/parse/net"
   "io"
   "net/http"
   "strconv"
   "strings"
   "time"
)

// 8728-1-1
func Parse(track string) (*Track, error) {
   split := strings.SplitN(track, "-", 3)
   err := mech.Strings(split).Has(2)
   if err != nil {
      return nil, err
   }
   rel, err := strconv.ParseInt(split[0], 10, 64)
   if err != nil {
      return nil, err
   }
   dis, err := strconv.ParseInt(split[1], 10, 64)
   if err != nil {
      return nil, err
   }
   num, err := strconv.ParseInt(split[2], 10, 64)
   if err != nil {
      return nil, err
   }
   return &Track{ReleaseID: rel, Disc: dis, Number: num}, nil
}

const origin = "https://bleep.com"

type Meta []net.Node

func NewMeta(releaseID int64) (Meta, error) {
   addr := []byte(origin)
   addr = append(addr, "/release/"...)
   addr = strconv.AppendInt(addr, releaseID, 10)
   req, err := http.NewRequest(
      "GET", string(addr), nil,
   )
   if err != nil {
      return nil, err
   }
   // this redirects, so we cannot use RoundTrip
   res, err := new(http.Client).Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return net.ReadHTML(res.Body, "meta"), nil
}

func (m Meta) Image() string {
   for _, node := range m {
      if node.Attr["property"] == "og:image" {
         return node.Attr["content"]
      }
   }
   return ""
}

// can be either one of these:
//  2001-05-01 00:00:00.0
//  Tue May 01 00:00:00 UTC 2001
func (m Meta) ReleaseDate() (time.Time, error) {
   for _, node := range m {
      if node.Attr["property"] == "music:release_date" {
         value := node.Attr["content"]
         date, err := time.Parse(time.UnixDate, value)
         if err != nil {
            return time.Parse("2006-01-02 15:04:05.9", value)
         }
         return date, nil
      }
   }
   return time.Time{}, mech.NotFound{"music:release_date"}
}

type Track struct {
   Artist string
   Title string
   ReleaseID int64
   Disc int64
   Number int64
}

func Release(releaseID int64) ([]Track, error) {
   body := []byte("type=ReleaseProduct&id=")
   body = strconv.AppendInt(body, releaseID, 10)
   req, err := http.NewRequest(
      "POST", origin + "/player/addToPlaylist", bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   var rel []Track
   if err := json.NewDecoder(res.Body).Decode(&rel); err != nil {
      return nil, err
   }
   return rel, nil
}

func (t Track) Resolve() (string, error) {
   req, err := http.NewRequest(
      "GET", origin + "/player/resolve/" + t.String(), nil,
   )
   if err != nil {
      return "", err
   }
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return "", err
   }
   defer res.Body.Close()
   dst, err := io.ReadAll(res.Body)
   if err != nil {
      return "", err
   }
   return string(dst), nil
}

func (t Track) String() string {
   track := strconv.AppendInt(nil, t.ReleaseID, 10)
   track = append(track, '-')
   track = strconv.AppendInt(track, t.Disc, 10)
   track = append(track, '-')
   track = strconv.AppendInt(track, t.Number, 10)
   return string(track)
}