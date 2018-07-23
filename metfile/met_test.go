package metfile_test

import (
	"bytes"
	"encoding/json"
	"github.com/wenerme/goform/metfile"
	"io/ioutil"
	"os"
	"testing"
)

func TestMetFileWriter(t *testing.T) {
	var met *metfile.MetForm
	var org []byte
	var err error
	{
		org, err = ioutil.ReadFile("testdata/server.met")
		if err != nil {
			t.Fatal(err)
		}
		reader := metfile.NewMetReader(bytes.NewBuffer(org))
		if met, err = reader.ReadForm(); err != nil {
			panic(err)
		}
	}
	out := &bytes.Buffer{}

	writer := metfile.NewMetWriter(met, out)
	if err := writer.WriteForm(); err != nil {
		panic(err)
	}

	if !bytes.HasPrefix(org, out.Bytes()) {
		t.Fatal("incorrect write")
	}

	// Debug only
	//ioutil.WriteFile("testdata/server-out.met",out.Bytes(),os.ModePerm)

	/* File created by ed2k.has.it will append
	00006610  03 01 00 91 0f 00 11 00  03 01 00 92 0b 00 00 00  |................|
	00006620  1d e0 4d 01 61 73 73 65  6d 62 6c 65 64 20 62 79  |..M.assembled by|
	00006630  20 68 74 74 70 3a 2f 2f  65 64 32 6b 2e 68 61 73  | http://ed2k.has|
	00006640  2e 69 74 2f 00                                    |.it/.|
	*/

}

func TestMetFileReader(t *testing.T) {
	file, err := os.Open("testdata/server.met")
	if err != nil {
		t.Fatal(err)
	}
	reader := metfile.NewMetReader(file)
	// Debug only
	//reader.Logger = func(f string, args ...interface{}) {
	//	fmt.Printf(f, args...)
	//	fmt.Print("\n")
	//}
	var met *metfile.MetForm
	if met, err = reader.ReadForm(); err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(met, "", "  ")
	if err != nil {
		panic(err)
	}
	println(string(b))
}
