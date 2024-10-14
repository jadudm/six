package vcap

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

//var vcap *gjson.Result = nil

// Alias
// type A = B
// New type
// type A B
type Credentials = gjson.Result

type VcapServices struct {
	Source string
	VCAP   gjson.Result
}

var VCS *VcapServices = nil

func (vcs *VcapServices) GetCredentials(service string, name string) (Credentials, error) {
	query_string := fmt.Sprintf("%s.#(name==%s).credentials", service, name)

	r := vcs.VCAP.Get(query_string)
	if r.Exists() {
		return Credentials(r), nil
	} else {
		return Credentials{}, errors.Errorf("No <%s> credentials found for '%s'", service, name)
	}
}

func ReadVCAPConfig() *VcapServices {
	// Remotely, read it in from the VCAP_SERVICES env var, which will
	// provide a large JSON structure.
	vcap_string := os.Getenv("VCAP_SERVICES")
	json := gjson.Parse(vcap_string)
	vcs := VcapServices{
		Source: "env",
		VCAP:   json,
	}
	VCS = &vcs
	return VCS
}

func ReadVCAPConfigFile(filename string) *VcapServices {
	fp := ""
	_, err := os.Open(filename)
	if errors.Is(err, os.ErrNotExist) {
		_, err = os.Open(filepath.Join(os.Getenv("HOME"), ".fac", filename))
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("CGOVUTIL Cannot find config. Exiting.")
		} else {
			fp = filepath.Join(os.Getenv("HOME"), ".fac", filename)
		}
	} else {
		fp = filename
	}

	bytes, err := os.ReadFile(fp)
	if err != nil {
		log.Printf("CGOVUTIL could not load config from file.")
		os.Exit(-1)
	}
	json := gjson.ParseBytes(bytes)
	VCS = &VcapServices{
		Source: filename,
		VCAP:   json,
	}
	return VCS
}
