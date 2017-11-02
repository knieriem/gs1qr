// This example program reads GS1 elements from stdin and writes
// an GS1 QR code to ./qr.png. Some information about the process is
// printed to stdout. Run
//	go build
//	./gs1qr < data
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"rsc.io/qr/coding"

	"github.com/knieriem/gs1qr"
	"github.com/knieriem/gs1qr/ai"
)

func main() {
	list, err := ParseElements(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(list)

	el := gs1qr.ConvertElements(list)
	fmt.Println("\n* symbol string: ", strings.Join(el.Strings(), ""))
	clist := el.Compile()
	c, p, err := gs1qr.Encode(clist, coding.M, 2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n* Data/Check/Version", p.DataBytes, p.CheckBytes, p.Version)
	b, err := gs1qr.Bytes(clist, p.Level, p.Version)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n* bytes")
	fmt.Printf("\t% x\n", b)

	png := c.PNG()
	ioutil.WriteFile("./qr.png", png, 0666)
	fmt.Println("\n* compiled elements")
	for i := range clist {
		fmt.Printf("\t%+v\n", clist[i])
	}
}

var appIDs = map[string]*ai.AppID{
	"GTIN": ai.GTIN,

	"BATCH/LOT":   ai.BatchLot,
	"PROD DATE":   ai.ProdDate,
	"BEST BEFORE": ai.BestBefore,

	"VARIANT": ai.Variant,
	"SERIAL":  ai.Serial,

	"VAR. COUNT": ai.VarCount,

	"SHIP TO LOC": ai.ShipToLoc,

	"CPID":        ai.CPID,
	"PRODUCT URL": ai.ProductURL,
}

func ParseElements(r io.Reader) ([]ai.Elem, error) {
	var list []ai.Elem
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}
		f := strings.SplitN(line, "\t", 2)
		if len(f) != 2 {
			continue
		}
		var e ai.Elem
		name := strings.TrimSpace(f[0])
		ai, ok := appIDs[name]
		if !ok {
			log.Println("AI not implemented:", name)
			continue
		}
		e.AI = ai
		e.Data = strings.Join(strings.Fields(f[1]), "")
		list = append(list, e)
	}
	err := s.Err()
	if err != nil {
		return nil, err
	}
	return list, nil
}
