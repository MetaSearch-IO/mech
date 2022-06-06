package paramount
// github.com/89z

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
   wv "github.com/89z/mech/widevine"
)

func (s Session) Containers(mod *wv.Module) (wv.Containers, error) {
   in, err := mod.Marshal()
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", s.URL, bytes.NewReader(in),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Authorization", "Bearer " + s.LS_Session)
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   out, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   return mod.Unmarshal(out)
}

type Session struct {
   URL string
   LS_Session string
}

func NewSession(contentID string) (*Session, error) {
   token, err := newToken()
   if err != nil {
      return nil, err
   }
   var buf strings.Builder
   buf.WriteString("https://www.paramountplus.com/apps-api/v3.0/androidphone")
   buf.WriteString("/irdeto-control/anonymous-session-token.json")
   req, err := http.NewRequest("GET", buf.String(), nil)
   if err != nil {
      return nil, err
   }
   req.URL.RawQuery = "at=" + url.QueryEscape(token)
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   sess := new(Session)
   if err := json.NewDecoder(res.Body).Decode(sess); err != nil {
      return nil, err
   }
   sess.URL += contentID
   return sess, nil
}