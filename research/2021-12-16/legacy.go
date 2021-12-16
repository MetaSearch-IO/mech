package main

import (
   "encoding/hex"
   "fmt"
   "golang.org/x/crypto/blowfish"
)

const data = 
   "0b5b3f806abef32879a802a0749e65e9bea1623d9ff53d4c47e1db0a11135f61e8de2089919ef34facde3af32214db0ec6a249a41f1dee680de4a53dd649b4e447abe2f430167fbef7d0c40952bf86b" +
   "4f89edfd49c846b348affb88374876158e61e7a1897570b3c5d649871a9f7a843eaea7ad91b6568cce8a4825e79fe600a497a074213cba320157cb173b2f7160be19ebcd105667b7cd8e3c16f3cb658" +
   "8ec972859029ba17df0d643f758c918b8021c242b6fc63cb77698d9f20cd456757ae960d62daf16609f129dd0a75014a89649c7cb29abc3154a41b24843187b22c416e43635bc12a2aac171ac10dcb0" +
   "ae32a8c36352df5b63177dfa7ddc0fd577b94ab4855ef1746ccf72a766bd4bbb1a6599fc44629c0c24b90e7ae1a12a4a47f698ad33d515193ab4f3eed1001037610909df3e58804038634eeec596eed" +
   "c58f1d067b2d4a7b6e05136801425788d872e512be24936fa65e2a4ef9b796b76f3ea86242c72346da40f27a39516e3e887ca3ffbd0003be0e09735ef77a8ced7125b443bdb0ca3470fe098755d5d7d" +
   "559ccf552d87546587a1ecf9cb940b2a3417f2f080c3e3cd5dda618392fa73d82b896e12f222c09c70ed7e67df541a956a0c86f570dfed8d921a32f18e82f705579b232dc6ba4debf9d52d0bf280d5d" +
   "4b5b3f07a0b5c1965561912c2a452e255b799ba187c2c6eab13480a41275f6bddd777fddc2380d380bab700f77a58ee7b37328b66d7a7f3fbfd266"

const blockSize = 8

func main() {
   src, err := hex.DecodeString(data)
   if err != nil {
      panic(err)
   }
   dst := make([]byte, len(src))
   c, err := blowfish.NewCipher([]byte(`6#26FRL$ZWD`))
   if err != nil {
      panic(err)
   }
   for bs, be := 0, blockSize; bs < len(src); bs, be = bs+blockSize, be+blockSize {
      c.Decrypt(dst[bs:be], src[bs:be])
   }
   fmt.Printf("%q\n", dst)
}
