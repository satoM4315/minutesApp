package main

//goはスペースが意味を持っているっぽい（調べてない
//pythonのインデントの感じ
//違うのはエラーなのにエラーと出力されないことがあること

import (
	"strconv"

	//ginのインポート
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"sort"
)

func setupRouter() *gin.Engine {

	r := gin.Default()

	dbInit() //データベースマイグレート

	// セッションの設定
	store := cookie.NewStore([]byte("secret"))
	// 静的ファイルのディレクトリを指定
	r.Static("dist", "./dist")
	// HTML ファイルのディレクトリを指定
	r.LoadHTMLGlob("./dist/public/*.html")

	r.Use(sessions.Sessions("mysession", store))
	store.Options(sessions.Options{
		MaxAge:   60 * 60 * 24 * 10, //10日立つと無効になる。それまでは"session time out"が出せるようとっておく
		Secure:   false,
		HttpOnly: true,
	})

	// セッション管理のテーブルを更新
	r.Use(sessionStoreUpdate())
	//ミーティング一覧のページ
	r.GET("/", returnMeetingsPage)
	// ログインページを返す
	r.GET("/login", returnLoginPage)
	// ログイン動作を司る
	r.POST("/login", postLogin)
	//ユーザー登録ページを返す
	r.GET("/register", returnRegisterPage)
	//　ユーザー登録動作を司る
	r.POST("/register", postRegister)
	//ログイン済みを前提とした処理を行う。sessionIDのチェックsessionCheck()を行った上で実行される
	logedIn := r.Group("/", sessionCheck())
	{
		// ミーティング一覧を返す
		logedIn.GET("/meetings", handleGetMeetings)
		// ミーティングの追加
		logedIn.POST("/meetings", handlePostMeetings)

		//meetingIDで指定されたIDを持つ議事録があるかを調べるミドルウェア
		// :~　とするとこで gin.contextから呼び出せる
		mtgIn := logedIn.Group("/meetings/:meetingID", meetingExistCheck())
		{
			//議事録一覧ページを返す
			mtgIn.GET(".", returnMinutesPage)
			// /message に　GETリクエストが飛んできたらfetchMessage関数を実行
			mtgIn.GET("/message", fetchMessage)
			// /add_messageへのPOSTリクエストは、handleAddMessage関数でハンドル
			mtgIn.POST("/add_message", handleAddMessage)
			// 重要と考えられる単語を返す
			mtgIn.GET("/important_words", handleImportantWords)
			//　重要と考えられる文を返す
			mtgIn.GET("/important_sentences", handleImportantSentences)
		}
		// /update_messageへのPOSTリクエストは、handleUpdateMessage関数でハンドル
		logedIn.POST("/update_message", handleUpdateMessage)
		// /delete_messageへのPOSTリクエストは、handleDeleteMessage関数でハンドル
		logedIn.POST("/delete_message", handleDeleteMessage)
		// ユーザー情報を返す
		logedIn.GET("/user", fetchUserInfo)
	}
	//セッション情報の削除
	r.GET("/logout", postLogout)

	r.GET("/entrance", returnEntrancePage)

	return r
}

func main() {

	//Temp="test"
	router := setupRouter()
	// サーバーを起動しています
	router.Run(":10000")
}

// ResponseUserPublic は、公開ユーザー情報がクライアントへ返される時の形式です。
// JSON形式へマーシャルできます。
type ResponseUserPublic struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ResponseMessage は、メッセージがクライアントへ返される時の形式です。
// JSON形式へマーシャルできます。
type ResponseMessage struct {
	ID      uint               `json:"id"`
	AddedBy ResponseUserPublic `json:"addedBy"`
	Message string             `json:"message"`
}

// ResponseMeeting は、ミーティングがクライアントへ返される時の形式です。
type ResponseMeeting struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

//URLで指定されたIDを持つ議事録があるかを調べるミドルウェア
func meetingExistCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		m, _ := strconv.Atoi(ctx.Param("meetingID"))
		meetingID := uint(m)

		meeting := getMeetingByID(meetingID)

		if meeting.ID == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			ctx.Abort()
			return
		}
		return
	}
}

//議事録本体ページ
func returnMinutesPage(ctx *gin.Context) {

	m, _ := strconv.Atoi(ctx.Param("meetingID"))
	meetingID := uint(m)

	meeting := getMeetingByID(meetingID)

	ctx.HTML(http.StatusOK, "template.html", gin.H{"title": meeting.Name, "header": "minuteHeader", "id": []string{"message"}})
}

//ログインページのhtmlを返す
func returnLoginPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "template.html", gin.H{"title": "Login and Register", "header": "loginHeader", "id": []string{"LoginAndRegister", "serverMessage"}})
}

//ユーザー登録ページのhtmlを返す
func returnRegisterPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "template.html", gin.H{"title": "Login and Register", "id": []string{"LoginAndRegister", "serverMessage"}})
}

func returnEntrancePage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "template.html", gin.H{"title": "Entrance", "header": "entranceHeader", "id": []string{"entrance", "serverMessage"}})
}

//議事録一覧ページ
//セッションのエラーが多いページなので、エラーは全てリダイレクトにしている
func returnMeetingsPage(ctx *gin.Context) {

	session := sessions.Default(ctx)
	user := session.Get("SessionID")

	if user == nil {
		ctx.Redirect(http.StatusSeeOther, "/entrance")
		ctx.Abort()
		return
	}

	if !(SessionExist(user.(string))) {
		session.Clear()
		session.Save()
		// 不当なセッション情報によるアクセス
		ctx.Redirect(http.StatusSeeOther, "/entrance")
		ctx.Abort()
		return
	}

	if SessionTimeOut(user.(string)) {
		//セッション有効時間が切れていた場合
		//セッションからデータを破棄する
		sessionDelete(user.(string))
		session.Clear()
		session.Save()

		sessionDelete(user.(string))
		ctx.Redirect(http.StatusSeeOther, "/entrance")
		ctx.Abort()
		return
	}

	ctx.HTML(http.StatusOK, "template.html", gin.H{"title": "Meetings", "header": "minuteHeader", "id": []string{"serverMessage", "meetings"}})
}

//messagesに含まれるものを jsonで返す
func fetchMessage(ctx *gin.Context) {
	m, _ := strconv.Atoi(ctx.Param("meetingID"))
	meetingID := uint(m)

	messagesInDB := MeetingMessageGetAll(meetingID)
	// データベースに保存されているメッセージの形式から、クライアントへ返す形式に変換する
	messages := make([]ResponseMessage, len(messagesInDB))
	for i, msg := range messagesInDB {
		// TODO データベースでJOIN？
		user := getUserByID(msg.UserID)
		messages[i] = ResponseMessage{
			ID: msg.ID,
			AddedBy: ResponseUserPublic{
				ID:   msg.UserID,
				Name: user.Username,
			},
			Message: msg.Message,
		}
	}
	ctx.JSON(http.StatusOK, messages)
}

// AddMessageRequest は、クライアントからのメッセージ追加要求のフォーマットです。
type AddMessageRequest struct {
	Message string `json:"message"`
}

func handleAddMessage(ctx *gin.Context) {
	// POST bodyからメッセージを獲得
	req := new(AddMessageRequest)
	err := ctx.BindJSON(req)

	m, _ := strconv.Atoi(ctx.Param("meetingID"))
	meetingID := uint(m)

	if err != nil {
		// メッセージがJSONではない、もしくは、content-typeがapplication/jsonになっていない
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request as JSON format is expected"})
		return
	}

	if req.Message == "" {
		// メッセージがない、無効なリクエスト
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request due to parameter 'message' being empty"})
		// 帰ることを忘れない
		return
	}

	user := getSessionUserID(ctx)

	//メッセージをデータベースへ追加
	messageInsert(req.Message, meetingID, user.ID)

	ctx.JSON(http.StatusOK, gin.H{"success": true})

	return
}

// UpdateMessageRequest は、クライアントからのメッセージ追加要求のフォーマットです。
type UpdateMessageRequest struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func handleUpdateMessage(ctx *gin.Context) {
	// POST bodyからメッセージを獲得
	req := new(UpdateMessageRequest)
	err := ctx.BindJSON(req)
	if err != nil {
		// メッセージがJSONではない、もしくは、content-typeがapplication/jsonになっていない
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request as JSON format is expected"})
		return
	}

	if req.Message == "" {
		// メッセージがない、無効なリクエスト
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request due to parameter 'message' being empty"})
		// 帰ることを忘れない
		return
	}

	id, _ := strconv.Atoi(req.ID)

	session := sessions.Default(ctx)

	if session.Get("SessionID") == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	user := getSessionUserID(ctx)
	msg := dbGetOne(id)

	if user.ID != msg.UserID {
		// 権限がない
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request due to privileges"})
		// 帰ることを忘れない
		return
	}

	//データベースにある指定されたメッセージを更新
	dbUpdate(id, req.Message)

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteMessageRequest は、クライアントからのメッセージ追加要求のフォーマットです。
type DeleteMessageRequest struct {
	ID string `json:"id"`
}

func handleDeleteMessage(ctx *gin.Context) {
	// POST bodyからメッセージを獲得
	req := new(DeleteMessageRequest)
	err := ctx.BindJSON(req)
	if err != nil {
		// メッセージがJSONではない、もしくは、content-typeがapplication/jsonになっていない
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request as JSON format is expected"})
		return
	}

	if req.ID == "" {
		// IDがない、無効なリクエスト
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request due to parameter 'id' being empty"})
		// 帰ることを忘れない
		return
	}

	id, _ := strconv.Atoi(req.ID)

	session := sessions.Default(ctx)

	if session.Get("SessionID") == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	user := getSessionUserID(ctx)
	msg := dbGetOne(id)

	if user.ID != msg.UserID {
		// 権限がない
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request due to privileges"})
		// 帰ることを忘れない
		return
	}

	//データベースにある指定されたメッセージを更新
	dbDelete(id)

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// ユーザー自身の情報を返す
func fetchUserInfo(ctx *gin.Context) {

	user := getSessionUserID(ctx)
	userInfo := ResponseUserPublic{
		ID:   user.ID,
		Name: user.Username,
	}
	ctx.JSON(http.StatusOK, userInfo)
}

//ログイン試行時にクライアントから送られてくるフォーマット
type userInfo struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

//登録動作
func postRegister(ctx *gin.Context) {
	// POST bodyからメッセージを獲得
	req := new(userInfo)
	err := ctx.BindJSON(req)

	if err != nil {
		// メッセージがJSONではない、もしくは、content-typeがapplication/jsonになっていない
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request as JSON format is expected"})
		return
	}

	if req.UserId == "" || req.Password == "" {
		// メッセージがない、無効なリクエスト
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request due to parameter 'userId' or 'password' being empty"})
		// 帰ることを忘れない
		return
	}

	// DBにユーザーの情報を登録
	if err := createUser(req.UserId, req.Password); err != nil {
		// ログインIDがすでに使用されている
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "already use this id"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
	return
}

//ログイン処理
func postLogin(ctx *gin.Context) {
	// POST bodyからメッセージを獲得
	req := new(userInfo)
	err := ctx.BindJSON(req)

	if err != nil {
		// メッセージがJSONではない、もしくは、content-typeがapplication/jsonになっていない
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request as JSON format is expected"})
		return
	}

	if req.UserId == "" || req.Password == "" {
		// メッセージがない、無効なリクエスト
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request due to parameter 'userId' or 'password' being empty"})
		// 帰ることを忘れない
		return
	}

	// 入力されたIDをもとにDBからレコードを取得
	user := getUser(req.UserId)

	if user.ID == 0 {
		// DBにユーザーの情報がない
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user not exist"})
		return
	}

	if err := comparePassword(user.Password, req.Password); err != nil {
		// パスワードが間違っている
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "wrong password"})
		return
	}

	//セッション管理
	sessionID := createSession(user.Username)
	if sessionID == "" {
		ctx.Redirect(http.StatusSeeOther, "/login")
		ctx.Abort()
		return
	}

	//セッションにデータを格納する
	session := sessions.Default(ctx)
	session.Set("SessionID", sessionID)
	session.Save()

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

//ログアウト処理
func postLogout(ctx *gin.Context) {

	//セッションからデータを破棄する
	session := sessions.Default(ctx)

	sessionID := session.Get("SessionID").(string)
	sessionDelete(sessionID)

	session.Clear()
	session.Save()

	ctx.Redirect(http.StatusSeeOther, "/entrance")

}

//tfidfを元に重要度を計算し、重要と考えられる単語を返す
func handleImportantWords(ctx *gin.Context) {
	m, _ := strconv.Atoi(ctx.Param("meetingID"))
	meetingID := uint(m)

	messagesInDB := MeetingMessageGetAll(meetingID)

	messages := make([]string, len(messagesInDB))
	for i, msg := range messagesInDB {
		messages[i] = msg.Message
	}
	allTfIdf := allTfIdf(messages)
	bestTfIdf := map[string]float64{}

	//複数文書に現れる単語のtfidfは最も大きい値を採用
	for _, tfidfs := range allTfIdf {
		for term := range tfidfs {
			if _, ok := bestTfIdf[term]; ok {
				if tfidfs[term] > bestTfIdf[term] {
					bestTfIdf[term] = tfidfs[term]
				}
			} else {
				bestTfIdf[term] = tfidfs[term]
			}
		}
	}
	//tfidfが大きい順にソート
	sortedTfIdf := List{}
	for k, v := range bestTfIdf {
		e := Items{k, v}
		sortedTfIdf = append(sortedTfIdf, e)
	}
	sort.Sort(sortedTfIdf)
	//上位10個(1０未満だったらその数だけ)を返す
	n := 10
	if n > len(sortedTfIdf) {
		n = len(sortedTfIdf)
	}
	result := make([]string, n)
	for i, item := range sortedTfIdf {
		if i == n {
			break
		}
		result[i] = item.name
	}
	ctx.JSON(http.StatusOK, result)
}

//以下mapのvalueを基準としてソートするために必要なもの
type Items struct {
	name  string
	value float64
}
type List []Items

func (l List) Len() int {
	return len(l)
}

func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l List) Less(i, j int) bool {
	if l[i].value == l[j].value {
		return (l[i].name < l[j].name)
	} else {
		return (l[i].value > l[j].value)
	}
}

//ここまで

func handleImportantSentences(ctx *gin.Context) {
	m, _ := strconv.Atoi(ctx.Param("meetingID"))
	meetingID := uint(m)

	messagesInDB := MeetingMessageGetAll(meetingID)
	messages := make([]string, len(messagesInDB))
	for i, msg := range messagesInDB {
		messages[i] = msg.Message
	}
	ranking := getImportantSentence(messages)
	n := 5
	if n > len(messages) {
		n = len(messages)
	}
	result := make([]string, n)
	for i, rank := range ranking {
		if i == n {
			break
		}
		result[i] = messages[rank]
	}
	ctx.JSON(http.StatusOK, result)
}

//議事録一覧を取得
func handleGetMeetings(ctx *gin.Context) {
	ms := getAllMeeting()
	ret := make([]ResponseMeeting, len(ms))
	for i, meeting := range ms {
		ret[i] = ResponseMeeting{
			ID:   meeting.ID,
			Name: meeting.Name,
		}
	}
	ctx.JSON(http.StatusOK, ret)
}

//議事録追加の際にクライアントから送られるフォーマット
type AddMeetingRequest struct {
	Meeting string `json:"meeting"`
}

//議事録を追加
func handlePostMeetings(ctx *gin.Context) {
	// POST bodyからメッセージを獲得
	req := new(AddMeetingRequest)
	err := ctx.BindJSON(req)

	if err != nil {
		// メッセージがJSONではない、もしくは、content-typeがapplication/jsonになっていない
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request as JSON format is expected"})
		return
	}

	if req.Meeting == "" {
		// メッセージがない、無効なリクエスト
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Malformed request due to parameter 'message' being empty"})
		// 帰ることを忘れない
		return
	}

	user := getSessionUserID(ctx)

	//議事録を作成
	if err := createMeeting(req.Meeting, user.ID); err != nil {
		// ログインIDがすでに使用されている
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "already use this name"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})

	return
}
