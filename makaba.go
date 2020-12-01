package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	//"mime/multipart"
	"strings"
	//"reflect"
	//"time"

	"github.com/tidwall/gjson"
)

const (
	makabaUrl  = "https://2ch.hk/makaba/makaba.fcgi"
	postingUrl = "https://2ch.hk/makaba/posting.fcgi?json=1"
)

type Passcode struct {
	Usercode string
	Error    bool
}

var CurrentUsercode Passcode = Passcode{
	Usercode: "",
	Error:    false,
}

func sendRepostMakaba(u User, p Payload) bool {
	// Skip if repost disabled
	if !u.Setting.Makaba {
		return false
	}
	log.Println("Sending payload to makaba")

	num := findThread(p)
	log.Printf("/%v/ : #%v", p.Board, num)

	valuesBase := prepareBase(p, p.Board, num)
	valuesFiles := prepareFiles(p.Files)

	client, ok := customClient()
	if ok == false {
		return false
	}

	if _, ok, number := makabaPost(client, postingUrl, valuesBase, valuesFiles); ok {
		log.Printf("https://2ch.hk/%v/res/%v.html#%v", p.Board, num, number)
		return ok
	}
	return false

}

func makabaPost(client *http.Client, url string, valuesBase map[string]io.Reader, valuesFiles map[string]io.Reader) (err error, success bool, num float64) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range valuesBase {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if fw, err = w.CreateFormField(key); err != nil {
			return
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err, false, num
		}

	}
	for key, r := range valuesFiles {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if fw, err = w.CreateFormFile(key, ""); err != nil {
			return
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err, false, num
		}

	}
	w.Close()

	// Prepare handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Высрать в тред
	res, err := client.Do(req)
	if err != nil {
		log.Println("client.Do(req) error:", err)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ioutil.ReadAll error:", err)
		return
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if result["Error"] != nil {
		log.Println("Makaba post error:", result)
	}
	log.Println(result)
	if result["Error"] == nil {
		log.Println("Successfully made post 👌🏻")
		success = true
		num = result["Num"].(float64)
		log.Printf("%v", result["Num"])
	}
	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return err, success, num
}

func getCatalog(board string) []byte {
	url := fmt.Sprintf("https://2ch.hk/%v/threads.json", board)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))
	return body
}

func findThread(p Payload) string {
	js := getCatalog(p.Board)
	threads := gjson.GetBytes(js, `threads.#.subject`)
	var ind int // Thread index
	for k, v := range threads.Array() {
		if strings.Contains(strings.ToLower(v.String()), p.Thread) == true {
			fmt.Println("Thread found; Index is:", k, "; subject is:", v)
			ind = k
		}
	}
	gjsonPath := fmt.Sprintf("threads.%v.num", ind)

	num := gjson.GetBytes(js, gjsonPath).String()
	//fmt.Println("Thread number is:", num)
	return num
}

func customClient() (*http.Client, bool) {
	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie
	auth := CurrentUsercode.PasscodeAuth()
	if auth == false {
		log.Println("Failed to authorize passcode. Skip.")
	}
	cookie := &http.Cookie{
		Name:   "passcode_auth",
		Value:  CurrentUsercode.Usercode,
		Path:   "/",
		Domain: "2ch.hk",
	}
	cookies = append(cookies, cookie)
	u, _ := url.Parse(postingUrl)
	jar.SetCookies(u, cookies)
	//log.Println(jar.Cookies(u))
	client := &http.Client{
		Jar: jar,
	}
	return client, auth
}

func prepareBase(p Payload, board string, num string) map[string]io.Reader {
	var baseReader map[string]io.Reader
	var comment string
	var name string
	comment = p.Caption
	//name = ""
	//comment = fmt.Sprintf("[sup]Стрим запустился! %v ⛓[/sup]\n\n", jsonPayload.Source)

	baseReader = map[string]io.Reader{
		"task":   strings.NewReader("post"),
		"board":  strings.NewReader(board),
		"thread": strings.NewReader(num),
		"name":   strings.NewReader(name), // Tripcode for attention whoring
		//"email": strings.NewReader(""), // U kid right?
		//"subject": strings.NewReader(jsonPayload.Person),
		"comment": strings.NewReader(comment), // Post text
	}
	return baseReader
}

func prepareFiles(files []string) map[string]io.Reader {
	var filesReader map[string]io.Reader
	//files := p.Files
	url := []string{}

	for k, v := range files {
		fmt.Println(k, "is:", v)
		url = append(url, v)
	}
	if len(url) == 0 {
		return filesReader
	}
	count := len(url)
	// I know, I know. But it works...
	switch count {
	case 1:
		fmt.Println("One file")
		for k, v := range files {
			fmt.Println(k, "is:", v)
			resp1, e := http.Get(v)
			if e != nil {
				fmt.Println("http.Get error:", e)
				log.Printf("%v", e)
			}
			//defer resp.Body.Close()
			filesReader = map[string]io.Reader{
				`files1`: resp1.Body,
			}
		}
	case 2:
		fmt.Println("Two files")
		resp1, e := http.Get(url[0])
		if e != nil {
			fmt.Println("http.Get error:", e)
			log.Printf("%v", e)
		}
		resp2, e := http.Get(url[1])
		if e != nil {
			fmt.Println("http.Get error:", e)
			log.Printf("%v", e)
		}
		//defer resp.Body.Close()
		filesReader = map[string]io.Reader{
			`files1`: resp1.Body,
			`files2`: resp2.Body,
		}
	case 3:
		fmt.Println("Three files")
		resp1, e := http.Get(url[0])
		if e != nil {
			fmt.Println("http.Get error:", e)
			log.Printf("%v", e)
		}
		resp2, e := http.Get(url[1])
		if e != nil {
			fmt.Println("http.Get error:", e)
			log.Printf("%v", e)
		}
		resp3, e := http.Get(url[2])
		if e != nil {
			fmt.Println("http.Get error:", e)
			log.Printf("%v", e)
		}
		//defer resp.Body.Close()
		filesReader = map[string]io.Reader{
			`files1`: resp1.Body,
			`files2`: resp2.Body,
			`files3`: resp3.Body,
		}
	default:
		fmt.Println(len(url))
		resp1, e := http.Get(url[0])
		if e != nil {
			fmt.Println("http.Get error:", e)
			log.Printf("%v", e)
		}
		resp2, e := http.Get(url[1])
		if e != nil {
			fmt.Println("http.Get error:", e)
			log.Printf("%v", e)
		}
		resp3, e := http.Get(url[2])
		if e != nil {
			fmt.Println("http.Get error:", e)
			log.Printf("%v", e)
		}
		resp4, e := http.Get(url[3])
		if e != nil {
			fmt.Println("http.Get error:", e)
			log.Printf("%v", e)
		}
		//defer resp.Body.Close()
		filesReader = map[string]io.Reader{
			`files1`: resp1.Body,
			`files2`: resp2.Body,
			`files3`: resp3.Body,
			`files4`: resp4.Body,
		}
	}
	return filesReader
}

// PasscodeAuth is used to authorize your passcode to get usercode. Used to bypass captcha
func (c *Passcode) PasscodeAuth() bool {
	formData := url.Values{
		"json":     {"1"},
		"task":     {"auth"},
		"usercode": {cfg.Passcode}}
	resp, err := http.PostForm(makabaUrl, formData)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false
	}
	//log.Println(string(body))

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(err)
		return false
	}
	//log.Println(result)
	if result["result"].(float64) == 0 {
		log.Println(result["description"])
		return false
	}
	if result["result"].(float64) == 1 {
		hash := fmt.Sprint(result["hash"])
		log.Println("✅ Got passcode_auth:", result["hash"])
		c.Usercode = hash
		c.Error = false
		return true
	}

	return false
}
