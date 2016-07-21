package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
	"strings"
	"io/ioutil"
	"net/http"
	"encoding/base64"
	"encoding/json"
	"flag"

)

type Emailhash struct {
	Get string `json:"GET"`
}
const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

)

var coder = base64.NewEncoding(base64Table)
var e = flag.String("e", "email@email.email", "email-address")
func base64Encode(src []byte) []byte {
	dst := coder.EncodeToString(src)
	dst = strings.Replace(dst, "+", "-", -1)
	dst = strings.Replace(dst, "/", "_", -1)
	return []byte(dst)
}

func base64Decode(src []byte) ([]byte, error) {
	dst := string(src)
	dst = strings.Replace(dst, "-", "+", -1)
	dst = strings.Replace(dst, "_", "/", -1)
	return coder.DecodeString(dst)
}

func httpGet(rurl string) string {
	resp, err := http.Get(rurl)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	return  string(body)
}




func runY(ii string,file1 *os.File,proxy string) {
	sibada:
	for {
		time.Sleep(5 * time.Second)
		timeChan := time.NewTimer(120 * time.Second)
		fmt.Println("timer go")
		chg := make(chan int,1)
		cmd := exec.Command("./qq","-proxy="+proxy)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println("stdout pipe err")
			chg <- 0
		}
		cmd.Start()
		done := make(chan error,1)
		go func() {
			done <- cmd.Wait()
		}()
		mi := log.New(file1, "YY" + ii + ": ", log.Ldate | log.Ltime)
		buff := bufio.NewScanner(stdout)
		go func() {
			for buff.Scan() {
				timeChan.Reset(120 * time.Second)
				mi.Printf(buff.Text())
				fmt.Println(buff.Text())
				if strings.Contains(buff.Text(), "fail") {
					chg <- 1
					break
				}
			}
		}()

		select {
		case <-chg:
			fmt.Println("(Child said fail Err. Restart.)")
			if err = cmd.Process.Signal(os.Interrupt); err != nil {
				log.Fatalln("Failed to open log file", "os", ":", err)
			}
			continue  sibada

		case <-timeChan.C:
			fmt.Println("Time Cycle Out. Restart.")
			if err = cmd.Process.Signal(os.Interrupt); err != nil {
				log.Fatalln("Failed to open log file", "os", ":", err)
			}
			continue sibada

		case err = <-done:
			fmt.Println("CMD Restart.")
			continue sibada
		}
	}

}

func main() {
	flag.Parse()
	file2, err := os.OpenFile("/weblog.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", "os", ":", err)
	}

	ss := httpGet("http://httpmq-server:1280/GET/" + *e + "?type=json")
	jsonSrc := []byte(ss)
	var eh Emailhash
	json.Unmarshal(jsonSrc, &eh)
	ssk := strings.Split(eh.Get,"|")
	b64d := strings.Replace(ssk[2], "_", "/", -1)
	ddd, _ := base64.StdEncoding.DecodeString(b64d)
	_ = ioutil.WriteFile("/lc.gob", ddd, 0666)
	de64,_ := base64Decode([]byte(ssk[0]))
	fmt.Println(string(de64))
	var i = 0
	for  {
		i++
		var uu string
		uu = strconv.Itoa(i)
		runY(uu,file2,ssk[1])
		fmt.Println("for restart")
	}
}

