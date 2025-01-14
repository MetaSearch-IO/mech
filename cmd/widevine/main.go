package main

import (
   "flag"
   "fmt"
   "os"
   "path/filepath"
)

func main() {
   home, err := os.UserHomeDir()
   if err != nil {
      panic(err)
   }
   var f flags
   // a
   flag.StringVar(&f.address, "a", "", "license address")
   // b
   flag.StringVar(&f.key_id, "b", "", "key ID")
   // c
   f.client = filepath.Join(home, "mech/client_id.bin")
   flag.StringVar(&f.client, "c", f.client, "client ID")
   // h
   flag.StringVar(&f.header, "h", "", "header")
   // k
   f.private_key = filepath.Join(home, "mech/private_key.pem")
   flag.StringVar(&f.private_key, "k", f.private_key, "private key")
   // v
   var verbose bool
   flag.BoolVar(&verbose, "v", false, "verbose")
   flag.Parse()
   if verbose {
      level.Level = 1
   }
   if f.key_id != "" {
      contents, err := f.contents()
      if err != nil {
         panic(err)
      }
      fmt.Println(contents.Content())
   } else {
      flag.Usage()
   }
}
