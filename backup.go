package main
import (
	"github.com/tedroden/goflickr"
	"fmt"
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
	"log"
)


var (
	API_KEY string = os.Getenv("FLICKR_KEY")
	API_SECRET string = os.Getenv("FLICKR_SECRET")
)


func getFrob() goflickr.Frob {
	r := &goflickr.Request{
		ApiKey: API_KEY,
		ApiSecret: API_SECRET,		
	}	
	frob := r.FrobGet()
	fmt.Println(frob)
	fmt.Println("Leaving getFrob")

	auth_url := r.AuthUrl(frob.Payload, "write")

	fmt.Println(auth_url)
	cmd := exec.Command("xdg-open", auth_url)

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Go to your browser and authenticate.")
	fmt.Println("Then press `Enter` to continue: ")
	_, _ = reader.ReadString('\n')
	return frob
}

func getToken() goflickr.Auth {
	r := &goflickr.Request{
		ApiKey: API_KEY,
		ApiSecret: API_SECRET,		
	}	
	frob := getFrob()
	fmt.Println("Frob: " + frob.Payload)
	return r.GetToken(frob)
}

func main() {
	
	auth := getToken()

	
	path := "/home/troden/Dropbox/Photos/Headshots/ted_15apr2011"
	files, _ := ioutil.ReadDir(path)
	fmt.Printf("Backing up %d files.\n", len(files))
	
	r := &goflickr.Request{
		ApiKey: API_KEY,
		ApiSecret: API_SECRET,
		AuthToken: auth.Token,
	}

	var set goflickr.Photoset
	total := len(files)
	i := 0
	for _, f := range files {
		i += 1
		full_path := fmt.Sprintf("%s/%s", path, f.Name())
		resp, err := r.Upload(full_path, "image/jpeg")
		if err != nil {
			fmt.Println("everything is wrong")
		} else {
			fmt.Printf("%d of %d: Uploaded %s\n", i, total, full_path)
			if set.Id == "" {
				set = r.PhotosetCreate("working backup test", resp.Id)
			} else {
				fmt.Printf("%d of %d: Added to set\n", i, total)
				r.PhotosetAddPhoto(set, resp.Id)
			}
		}
	}
}
