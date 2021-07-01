package youtube_test

import (
   "github.com/89z/mech/youtube"
   "io"
   "testing"
)

const desc = "Provided to YouTube by Epitaph\n\nSnowflake · Kate Bush\n\n" +
"50 Words For Snow\n\n" +
"℗ Noble & Brite Ltd. trading as Fish People, under exclusive license to Anti Inc.\n\n" +
"Released on: 2011-11-22\n\nMusic  Publisher: Noble and Brite Ltd.\n" +
"Composer  Lyricist: Kate Bush\n\nAuto-generated by YouTube."

func TestMWeb(t *testing.T) {
   mw, err := youtube.NewMWeb("XeojXq6ySs4")
   if err != nil {
      t.Fatal(err)
   }
   if mw.PublishDate != "2020-11-05" {
      t.Fatalf("%+v\n", mw)
   }
   if mw.ShortDescription != desc {
      t.Fatalf("%+v\n", mw)
   }
   if mw.ViewCount == 0 {
      t.Fatalf("%+v\n", mw)
   }
}

func TestAndroid(t *testing.T) {
   a, err := youtube.NewAndroid("XeojXq6ySs4")
   if err != nil {
      t.Fatal(err)
   }
   if a.Title != "Snowflake" {
      t.Fatalf("%+v\n", a)
   }
   f := a.StreamingData.AdaptiveFormats.Filter(func(f youtube.Format) bool {
      return f.Height == 0
   })
   if err := f[0].Write(io.Discard); err != nil {
      t.Fatal(err)
   }
}
