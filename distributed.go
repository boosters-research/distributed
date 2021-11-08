// Package distributed is for building federated communities
package distributed

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	db "go.m3o.com/db"
	user "go.m3o.com/user"
)

var (
	// our global M3O API token
	token = os.Getenv("M3O_API_TOKEN")
	// client for the DB API
	dbService = db.NewDbService(token)
	// client for the User API
	userService = user.NewUserService(token)
	// csv of user ids
	mods = os.Getenv("DISTRIBUTED_MODS")
)

// Types

type Post struct {
	Id           string  `json:"id"`
	UserId       string  `json:"userId"`
	UserName     string  `json:"userName"`
	Content      string  `json:"content"`
	Created      string  `json:"created"`
	Upvotes      float32 `json:"upvotes"`
	Downvotes    float32 `json:"downvotes"`
	Score        float32 `json:"score"`
	Title        string  `json:"title"`
	Url          string  `json:"url"`
	Sub          string  `json:"sub"`
	CommentCount float32 `json:"commentCount"`
}

type Comment struct {
	Content   string  `json:"content"`
	Parent    string  `json:"sub"`
	Upvotes   float32 `json:"upvotes"`
	Downvotes float32 `json:"downvotes"`
	PostId    string  `json:"postId"`
	UserName  string  `json:"usernName"`
	UserId    string  `json:"userId"`
}

type PostRequest struct {
	Post      Post   `json:"post"`
	SessionID string `json:"sessionId"`
}

type VoteRequest struct {
	Id        string `json:"id"`
	SessionID string `json:"sessionId"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type CommentRequest struct {
	Comment   Comment `json:"comment"`
	SessionID string  `json:"sessionId"`
}

type CommentsRequest struct {
	PostId string `json:"postId"`
}

type PostsRequest struct {
	Min   int32  `json:"min"`
	Max   int32  `json:"max"`
	Limit int32  `json:"limit"`
	Sub   string `json:"sub"`
}

// Endpoints

// upvote or downvote a post or a comment
func vote(w http.ResponseWriter, req *http.Request, upvote bool, isComment bool, t VoteRequest) error {
	if t.Id == "" {
		return fmt.Errorf("missing post id")
	}
	table := "posts"
	if isComment {
		table = "comments"
	}
	rsp, err := dbService.Read(&db.ReadRequest{
		Table: table,
		Id:    t.Id,
	})
	if err != nil {
		return err
	}
	if len(rsp.Records) == 0 {
		return fmt.Errorf("post or comment not found")
	}

	// auth
	sessionRsp, err := userService.ReadSession(&user.ReadSessionRequest{
		SessionId: t.SessionID,
	})
	if err != nil {
		return err
	}
	if sessionRsp.Session.UserId == "" {
		return fmt.Errorf("user id not found")
	}

	// prevent double votes
	checkTable := table + "votecheck"
	checkId := t.Id + sessionRsp.Session.UserId
	checkRsp, err := dbService.Read(&db.ReadRequest{
		Table: checkTable,
		Id:    checkId,
	})
	mod := isMod(sessionRsp.Session.UserId, mods)
	if err == nil && (checkRsp != nil && len(checkRsp.Records) > 0) {
		if !mod {
			return fmt.Errorf("already voted")
		}
	}
	val := float64(1)
	if mod {
		rand.Seed(time.Now().UnixNano())
		val = float64(rand.Intn(17-4) + 4)
	}

	if !mod {
		_, err = dbService.Create(&db.CreateRequest{
			Table: checkTable,
			Record: map[string]interface{}{
				"id": checkId,
			},
		})
		if err != nil {
			return err
		}
	}

	obj := rsp.Records[0]
	key := "upvotes"
	if !upvote {
		key = "downvotes"
	}

	if _, ok := obj["upvotes"].(float64); !ok {
		obj["upvotes"] = float64(0)
	}
	if _, ok := obj["downvotes"].(float64); !ok {
		obj["downvotes"] = float64(0)
	}

	obj[key] = obj[key].(float64) + val
	obj["score"] = obj["upvotes"].(float64) - obj["downvotes"].(float64)

	_, err = dbService.Update(&db.UpdateRequest{
		Table:  table,
		Id:     t.Id,
		Record: obj,
	})
	return err
}

func isMod(userId, s string) bool {
	arr := strings.Split(s, ",")
	for _, v := range arr {
		if v == userId {
			return true
		}
	}
	return false
}

func VoteWrapper(upvote bool, isComment bool) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if Cors(w, req) {
			return
		}

		decoder := json.NewDecoder(req.Body)
		var t VoteRequest
		err := decoder.Decode(&t)
		if err != nil {
			respond(w, nil, err)
			return
		}
		err = vote(w, req, upvote, isComment, t)
		respond(w, nil, err)
	}
}

func Signup(w http.ResponseWriter, req *http.Request) {
	if Cors(w, req) {
		return
	}

	decoder := json.NewDecoder(req.Body)
	var t SignupRequest
	err := decoder.Decode(&t)
	if err != nil {
		respond(w, err, err)
		return
	}
	_, err = userService.Read(&user.ReadRequest{
		Username: t.Username,
	})
	if err != nil {
		createRsp, err := userService.Create(&user.CreateRequest{
			Username: t.Username,
			Email:    t.Email,
			Password: t.Password,
		})
		if err != nil {
			respond(w, createRsp, err)
			return
		}
	}
	logRsp, err := userService.Login(&user.LoginRequest{
		Username: t.Username,
		Password: t.Password,
	})
	respond(w, logRsp, err)
}

func Login(w http.ResponseWriter, req *http.Request) {
	if Cors(w, req) {
		return
	}
	decoder := json.NewDecoder(req.Body)
	var t LoginRequest
	err := decoder.Decode(&t)
	if err != nil {
		respond(w, err, err)
		return
	}
	logRsp, err := userService.Login(&user.LoginRequest{
		Username: t.Username,
		Password: t.Password,
	})
	respond(w, logRsp, err)
}

func ReadSession(w http.ResponseWriter, req *http.Request) {
	if Cors(w, req) {
		return
	}

	decoder := json.NewDecoder(req.Body)
	var t user.ReadSessionRequest
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("%v", err.Error()))
	}
	rsp, err := userService.ReadSession(&t)
	if err != nil {
		respond(w, rsp, err)
		return
	}
	readRsp, err := userService.Read(&user.ReadRequest{
		Id: rsp.Session.UserId,
	})
	respond(w, map[string]interface{}{
		"session": rsp.Session,
		"account": readRsp.Account,
	}, err)
}

func NewPost(w http.ResponseWriter, req *http.Request) {
	if Cors(w, req) {
		return
	}

	decoder := json.NewDecoder(req.Body)
	var t PostRequest
	err := decoder.Decode(&t)
	if err != nil {
		respond(w, nil, err)
		return
	}
	if t.Post.Sub == "" || t.Post.Title == "" {
		respond(w, nil, fmt.Errorf("both title and sub are required"))
		return
	}
	if t.Post.Url == "" && t.Post.Content == "" {
		respond(w, nil, fmt.Errorf("url or content required"))
		return
	}
	if len(t.Post.Title) > 200 || len(t.Post.Url) > 200 {
		respond(w, nil, fmt.Errorf("post url or title too long"))
		return
	}
	if len(t.Post.Sub) > 50 {
		respond(w, nil, fmt.Errorf("post sub too long"))
		return
	}
	if len(t.Post.Content) > 3000 {
		respond(w, nil, fmt.Errorf("post content too long"))
		return
	}
	userID := ""
	userName := ""
	if t.SessionID != "" {
		rsp, err := userService.ReadSession(&user.ReadSessionRequest{
			SessionId: t.SessionID,
		})
		if err != nil {
			respond(w, rsp, err)
			return
		}
		userID = rsp.Session.UserId
		readRsp, err := userService.Read(&user.ReadRequest{
			Id: userID,
		})
		if err != nil {
			respond(w, rsp, err)
			return
		}
		userName = readRsp.Account.Username
	}
	dbService.Create(&db.CreateRequest{
		Table: "posts",
		Record: map[string]interface{}{
			"id":        uuid.NewV4(),
			"userId":    userID,
			"userName":  userName,
			"content":   t.Post.Content,
			"url":       t.Post.Url,
			"upvotes":   float64(0),
			"downvotes": float64(0),
			"score":     float64(0),
			"sub":       t.Post.Sub,
			"title":     t.Post.Title,
			"created":   time.Now(),
		},
	})
}

func NewComment(w http.ResponseWriter, req *http.Request) {
	if Cors(w, req) {
		return
	}

	decoder := json.NewDecoder(req.Body)
	var t CommentRequest
	err := decoder.Decode(&t)
	if err != nil {
		respond(w, nil, err)
		return
	}
	userID := ""
	userName := ""
	// get user if available
	if t.SessionID != "" {
		rsp, err := userService.ReadSession(&user.ReadSessionRequest{
			SessionId: t.SessionID,
		})
		if err != nil {
			respond(w, rsp, err)
			return
		}
		userID = rsp.Session.UserId
		readRsp, err := userService.Read(&user.ReadRequest{
			Id: userID,
		})
		if err != nil {
			respond(w, rsp, err)
			return
		}
		userName = readRsp.Account.Username
	}
	if t.Comment.PostId == "" {
		respond(w, nil, fmt.Errorf("no post id"))
		return
	}

	// get post to update comment counter
	readRsp, err := dbService.Read(&db.ReadRequest{
		Table: "posts",
		Id:    t.Comment.PostId,
	})
	if err != nil {
		respond(w, nil, err)
		return
	}
	if readRsp == nil || len(readRsp.Records) == 0 {
		respond(w, nil, fmt.Errorf("post not found"))
		return
	}
	if len(readRsp.Records) > 1 {
		respond(w, nil, fmt.Errorf("multiple posts found"))
		return
	}

	// create comment
	_, err = dbService.Create(&db.CreateRequest{
		Table: "comments",
		Record: map[string]interface{}{
			"id":        uuid.NewV4(),
			"userId":    userID,
			"userName":  userName,
			"content":   t.Comment.Content,
			"parent":    t.Comment.Parent,
			"postId":    t.Comment.PostId,
			"upvotes":   float64(0),
			"downvotes": float64(0),
			"score":     float64(0),
			"created":   time.Now(),
		},
	})
	if err != nil {
		respond(w, nil, err)
		return
	}

	// update counter
	oldCount, ok := readRsp.Records[0]["commentCount"].(float64)
	if !ok {
		oldCount = 0
	}
	oldCount++
	readRsp.Records[0]["commentCount"] = oldCount
	_, err = dbService.Update(&db.UpdateRequest{
		Table:  "posts",
		Id:     t.Comment.PostId,
		Record: readRsp.Records[0],
	})
	respond(w, nil, err)
}

func score(m map[string]interface{}) float64 {
	score, ok := m["score"].(float64)
	if !ok {
		return -10000
	}
	sign := float64(1)
	if score == 0 {
		sign = 0
	}
	if score < 0 {
		sign = -1
	}
	order := math.Log10(math.Max(math.Abs(score), 1))
	var created int64
	switch v := m["created"].(type) {
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			fmt.Println(err)
		}
		created = t.Unix()
	case float64:
		created = int64(v)
	case int64:
		created = v
	}

	seconds := created - 1134028003
	return sign*order + float64(seconds)/45000
}

func Posts(w http.ResponseWriter, req *http.Request) {
	if Cors(w, req) {
		return
	}

	var t PostsRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&t)
	r := &db.ReadRequest{
		Table:   "posts",
		Order:   "desc",
		OrderBy: "created",
		Limit:   1000,
	}
	query := ""

	// @TODO this should be != 0 but that causes an empty new page
	if t.Min > 0 {
		query += "score >= " + fmt.Sprintf("%v", t.Min)
	}
	if t.Max > 0 {
		if query != "" {
			query += " and "
		}
		query += "score <= " + fmt.Sprintf("%v", t.Max)
	}
	if t.Sub != "all" && t.Sub != "" {
		if query != "" {
			query += " and "
		}
		query += fmt.Sprintf("sub == '%v'", t.Sub)
	}
	if query != "" {
		r.Query = query
	}

	rsp, err := dbService.Read(r)
	sort.Slice(rsp.Records, func(i, j int) bool {
		return score(rsp.Records[i]) > score(rsp.Records[j])
	})
	respond(w, rsp, err)
}

func Comments(w http.ResponseWriter, req *http.Request) {
	if Cors(w, req) {
		return
	}

	var t CommentsRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("%v", err.Error()))
	}
	rsp, err := dbService.Read(&db.ReadRequest{
		Table:   "comments",
		Order:   "desc",
		Query:   "postId == '" + t.PostId + "'",
		OrderBy: "created",
	})
	sort.Slice(rsp.Records, func(i, j int) bool {
		return score(rsp.Records[i]) > score(rsp.Records[j])
	})
	respond(w, rsp, err)
}

// Utils

func Cors(w http.ResponseWriter, req *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

func respond(w http.ResponseWriter, i interface{}, err error) {
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	}
	if i == nil {
		i = map[string]interface{}{}
	}
	if err != nil {
		i = map[string]interface{}{
			"error": err.Error(),
		}
	}
	bs, _ := json.Marshal(i)
	fmt.Fprintf(w, fmt.Sprintf("%v", string(bs)))
}

// Run the server on a given address
func Run(address string) {
	http.HandleFunc("/upvotePost", VoteWrapper(true, false))
	http.HandleFunc("/downvotePost", VoteWrapper(false, false))
	http.HandleFunc("/upvoteComment", VoteWrapper(true, true))
	http.HandleFunc("/downvoteComment", VoteWrapper(false, true))
	http.HandleFunc("/posts", Posts)
	http.HandleFunc("/post", NewPost)
	http.HandleFunc("/comment", NewComment)
	http.HandleFunc("/comments", Comments)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/readSession", ReadSession)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
  <head>
    <title>Go Distributed</title>
    <meta charset="UTF-8">
    <meta name="description" content="A federated community API">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <body>
    <h1>Go Distributed</h1>
    <p>A federated community API<p>
    <p>Visit <a href="https://github.com/m3o/go-distributed">Github</a> for more info</p>
  </body>
</html>
`))
	})

	http.ListenAndServe(address, nil)
}
