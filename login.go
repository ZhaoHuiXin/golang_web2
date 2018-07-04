package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func init() {
	app.HandleFunc("GET", "/login", loginHandler)
	app.HandleFunc("GET", "/logout", logoutHandler)
	app.HandleFunc("GET", "/weixin/auth", wait_weixin_auth_handler)
	app.HandleFunc("GET", "/weixin/waitbind", wait_weixin_bind_handler)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	eno := r.FormValue("eno")
	data := make(map[string]interface{})
	data["eno"] = eno

	email := r.FormValue("email")
	if len(email) > 0 {
		data["email"] = email
		id := 0
		s := `SELECT id FROM users WHERE email=?`
		err := app.db.QueryRow(s, email).Scan(&id)
		if err != nil {
			log.Println("QueryRow", err, s, email)
			data["eno"] = 2
		} else if id < 1 {
			data["eno"] = 3
			log.Println("login: no user with email", email)
		} else {
			expire_seconds := 120
			token := rand.Int63()
			data["token"] = token
			data["expire_seconds"] = expire_seconds
			data["qrcode"] = fmt.Sprintf("https://api.autoforce.net/wx/dms/qrcode?action_name=QR_STR_SCENE&scene=2:%d,%d&expire_seconds=%d", id, token, expire_seconds)
		}
	}
	render.Render(w, "login.html", data)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := session_store.Get(r, session_name)
	delete(session.Values, "uid")
	delete(session.Values, "role_id")
	session.Save(r, w)
	http.Redirect(w, r, `/`, http.StatusFound)
}

type WxUserScanEvent struct {
	Scene      string `json:"scene,omitempty"`
	Unionid    string `json:"unionid,omitempty"`
	Nickname   string `json:"nickname,omitempty"`
	Headimgurl string `json:"headimgurl,omitempty"`

	Sex json.Number `json:"sex,omitempty"`
}

func queryWxScanEvent(id, token int64) int {

	client := http.Client{
		Timeout: time.Duration(60 * time.Second),
	}
	event_point := fmt.Sprintf("https://api.autoforce.net/wx/dms/userscanevent?scene=2:%d,%d", id, token)
	resp, err := client.Get(event_point)
	if err != nil {
		log.Println("request event fail", err)
		return -1
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if body == nil || len(body) < 1 {
		return -2
	}

	res := WxUserScanEvent{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Println("request update, parse json fail", string(body), err)
		return -2
	}

	scenes := strings.Split(res.Scene, ",")
	if len(scenes) != 2 {
		log.Println("scene invalid", scenes)
		return -4
	}
	_id, _ := strconv.ParseInt(scenes[0], 10, 64)
	if id != _id {
		log.Println("id invalid", _id)
		return -5
	}
	tk, _ := strconv.ParseInt(scenes[1], 10, 64)
	if tk != token {
		log.Println("token invalid", tk)
		return -6
	}

	if len(res.Unionid) < 2 {
		log.Printf("len res.Unionid empty")
		return -7
	}
	_unionid := ""

	s := `SELECT unionid FROM cheshi.users WHERE id=?`
	err = app.db.QueryRow(s, id).Scan(&_unionid)
	if err != nil {
		log.Println("QueryRow", err, s, id)
		return -8
	}

	if strings.HasPrefix(_unionid, "waitbind:") {
		_, err := app.db.Exec(`UPDATE cheshi.users SET unionid=? where id=?`, res.Unionid, id)
		if err == nil {
			_unionid = res.Unionid
		} else {
			log.Println("bind user weixin", err)
			return -9
		}
	}
	if _unionid != res.Unionid {
		log.Println("user", id, "unionid err")
		return -10
	}

	_, err = app.db.Exec(
		`UPDATE cheshi.users SET sex=?,nickname=?,headimgurl=?,token=? where id=?`,
		res.Sex.String(), res.Nickname, res.Headimgurl, tk, id)
	if err != nil {
		log.Println("update user", err)
		return -11
	}
	return 0
}

func wait_weixin_auth_handler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	email := r.FormValue("email")
	tk, _ := strconv.ParseInt(r.FormValue("tk"), 10, 64)

	if email == "" {
		data["status"] = "fail"
		jsonp(w, r, data)
		return
	}
	if tk < 1 {
		data["status"] = "fail"
		jsonp(w, r, data)
		return
	}

	var id int64
	s := `SELECT id FROM cheshi.users WHERE email=?`
	err := app.db.QueryRow(s, email).Scan(&id)
	if err != nil {
		log.Println("QueryRow", err, s, email)
		data["status"] = "fail"
		jsonp(w, r, data)
		return
	}

	for i := 0; i < 120; i++ {
		time.Sleep(time.Second)
		if queryWxScanEvent(id, tk) != 0 {
			continue
		}
		var name sql.NullString
		var role_id sql.NullInt64
		var token int64
		s := `SELECT name,token,role_id FROM cheshi.users WHERE id=?`
		err := app.db.QueryRow(s, id).Scan(&name, &token, &role_id)
		if err != nil {
			log.Println("QueryRow", err, s, id)
			continue
		}

		if tk != token {
			log.Println("token invalid", tk, token)
			continue
		}
		session, _ := session_store.Get(r, session_name)
		session.Values["uid"] = int(id)
		session.Values["role_id"] = int(role_id.Int64)
		session.Values["name"], _ = name.Value()
		session.Values["token"] = token
		session.Save(r, w)
		data["status"] = "ok"
		jsonp(w, r, data)
		return
	}
	data["status"] = "fail"
	jsonp(w, r, data)
}

func wait_weixin_bind_handler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	for i := 0; i < 120; i++ {
		time.Sleep(time.Second)
		if queryWxScanEvent(id, 0) != 0 {
			continue
		}
		var unionid sql.NullString
		s := `SELECT unionid FROM cheshi.users WHERE id=?`
		err := app.db.QueryRow(s, id).Scan(&unionid)
		if err != nil {
			data["status"] = "fail"
			jsonp(w, r, data)
			return
		}

		if !isWbind(unionid.String) {
			data["status"] = "ok"
			jsonp(w, r, data)
			return
		}
	}
	data["status"] = "fail"
	jsonp(w, r, data)
}
