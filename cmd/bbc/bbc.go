package main

import (
   "flag"
   "fmt"
   "github.com/89z/mech"
   "github.com/89z/mech/bbc"
   "github.com/89z/parse/m3u"
   "net/http"
   "os"
   "path"
   "strconv"
   "strings"
)


type choice struct {
   format bool
   ids map[string]bool
}

func main() {
   cHLS := choice{
      ids: make(map[string]bool),
   }
   flag.BoolVar(&cHLS.format, "hf", false, "HLS formats")
   flag.Func("h", "HLS IDs", func(id string) error {
      cHLS.ids[id] = true
      return nil
   })
   flag.Parse()
   if len(os.Args) == 1 {
      fmt.Println("bbc [flags] [URL]")
      flag.PrintDefaults()
      return
   }
   if !cHLS.format && len(cHLS.ids) == 0 {
      return
   }
   addr := flag.Arg(0)
   if err := cHLS.HLS(addr); err != nil {
      panic(err)
   }
}

func (c choice) HLS(addr string) error {
   news, err := bbc.NewNewsVideo(addr)
   if err != nil {
      return err
   }
   media, err := news.Media()
   if err != nil {
      return err
   }
   ext, err := mech.ExtensionByType(media.Type)
   if err != nil {
      return err
   }
   video, err := media.Video()
   if err != nil {
      return err
   }
   forms, err := video.HLS()
   if err != nil {
      return err
   }
   base := path.Base(addr)
   for _, form := range forms {
      if c.format {
         if !strings.HasPrefix(form.URI.File, "keyframes/") {
            fmt.Printf("%+v\n", form)
         }
      } else if c.ids[strconv.Itoa(form.ID)] {
         fmt.Println("GET", form.URI)
         res, err := http.Get(form.URI.String())
         if err != nil {
            return err
         }
         defer res.Body.Close()
         forms, err := m3u.Decode(res.Body, form.URI.Dir)
         if err != nil {
            return err
         }
         file, err := os.Create(news.Caption + "-" + base + ext)
         if err != nil {
            return err
         }
         defer file.Close()
         for _, form := range forms {
            fmt.Println("GET", form.URI)
            res, err := http.Get(form.URI.String())
            if err != nil {
               return err
            }
            defer res.Body.Close()
            if _, err := file.ReadFrom(res.Body); err != nil {
               return err
            }
         }
      }
   }
   return nil
}
