package main

import (
	"fmt"
	"github.com/godcong/go-ipfs-restapi"
)

func main() {
	s := shell.NewShell("localhost:5001")
	objects, e := s.AddDir("D:\\workspace\\goproject\\go-ipfs-restapi\\issues\\issues16\\")
	if e != nil {
		panic(e)
	}
	for _, object := range objects {
		fmt.Printf("%+v", object)
	}
	//Output:
	//&{Hash:QmatoEW5Xu8ndYKDeW9h3DcA1oeU7iUHWS6sZ1hYwztf4K Name:D:\workspace\goproject\go-ipfs-restapi\issues\issues16\/main.go Size:321}&{Hash:QmW8hBLAUpvfhFad9L8faSqz6MVyEcz9WTDiDyANLCgQnG Name:D:\workspace\goproject\go-ipfs-restapi\issues\issues16\/test.jpg Size:57804}&{Hash:QmVUo6UeEWpcQCroqRY1uw47eAmZKEFdPJA
	//r8dh1kpJpx9 Name:D:\workspace\goproject\go-ipfs-restapi\issues\issues16\ Size:58231}
}
