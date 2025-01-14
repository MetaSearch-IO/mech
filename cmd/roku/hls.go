package main

import (
   "fmt"
   "github.com/89z/format"
   "github.com/89z/format/hls"
   "io"
   "net/http"
   "net/url"
   "os"
)

func (d downloader) HLS(bandwidth int64) error {
   video, err := d.Content.HLS()
   if err != nil {
      return err
   }
   fmt.Println("GET", video.URL)
   res, err := http.Get(video.URL)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   master, err := hls.New_Scanner(res.Body).Master()
   if err != nil {
      return err
   }
   stream := master.Streams.Get_Bandwidth(bandwidth)
   if !d.info {
      addr, err := res.Request.URL.Parse(stream.Raw_URI)
      if err != nil {
         return err
      }
      return download_HLS(addr, d.Base())
   }
   fmt.Println(d.Content)
   for _, each := range master.Streams {
      if each.Bandwidth == stream.Bandwidth {
         fmt.Print("!")
      }
      fmt.Println(each)
   }
   return nil
}

func download_HLS(addr *url.URL, base string) error {
   fmt.Println("GET", addr)
   res, err := http.Get(addr.String())
   if err != nil {
      return err
   }
   defer res.Body.Close()
   seg, err := hls.New_Scanner(res.Body).Segment()
   if err != nil {
      return err
   }
   file, err := os.Create(base + hls.TS)
   if err != nil {
      return err
   }
   defer file.Close()
   pro := format.Progress_Chunks(file, len(seg.Clear))
   for _, clear := range seg.Clear {
      addr, err := res.Request.URL.Parse(clear)
      if err != nil {
         return err
      }
      res, err := http.Get(addr.String())
      if err != nil {
         return err
      }
      pro.Add_Chunk(res.ContentLength)
      if _, err := io.Copy(pro, res.Body); err != nil {
         return err
      }
      if err := res.Body.Close(); err != nil {
         return err
      }
   }
   return nil
}
