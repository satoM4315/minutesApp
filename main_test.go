package main

import (

	//"fmt"

	"strings"

	//"reflect"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

//正規のsession情報を格納する
var mainCookie string = " "

//ログインごとにsession情報が入れ替わることをテスト
var oldCookie string = " "

//別アカウントのテスト用
var subCookie string = " "

//有効期限切れのcookieを扱う用
var tempCookie string = " "

//ダミーのsession情報
//jsonStr := `{"UserId":"dummy_user","Password":"dumdum00"}`
var dummyCookie string = `mysession=MTU5ODg0Mjg4NXxEdi1CQkFFQ180SUFBUkFCRUFBQU1fLUNBQUVHYzNSeWFXNW5EQXNBQ1ZObGMzTnBiMjVKUkFaemRISnBibWNNRWdBUU5UTkJTRmRwT1VoMllubFJlblF3ZEE9PXzqHFPM7diSqe2r0Kg2EzePlFc1iOf9Y2hBfzalMTSebA==; Path=/; Expires=Wed, 30 Sep 2020 03:01:25 GMT; Max-Age=2592000`

//サーバーのルーティング
var router = setupRouter()

//テストのためのmeetingインスタンス
var GetSelectMeetingPageRoute string = "/"
var GetMeetingRoute string = "/meetings"
var PostMeetingRoute string = "/meetings"

var GetMinutesPageRoute string = "/meetings/1"
var DummyMinutesPageRoute string = "/meetings/1234"

var EntranceRoute string = "/entrance"
var GetUserInfo string = "/user"

var GetMessageRoute string = GetMinutesPageRoute + "/message"
var DummyGetMessgaeRoute string = DummyMinutesPageRoute + "/message"
var PostMessageRoute string = GetMinutesPageRoute + "/add_message"
var DummyPostMessgaeRoute string = DummyMinutesPageRoute + "/add_message"
var UpdateMessageRoute string = "/update_message"
var DeleteMessageRoute string = "/delete_message"

var LoginRoute string = "/login"
var RegisterRoute string = "/register"
var LogoutRoute string = "/logout"

var GetImportantWordsRoute string = GetMinutesPageRoute + "/important_words"
var GetImportantSentencesRoute string = GetMinutesPageRoute + "/important_sentences"

//エントランスページはセッション情報がなくても取得できる
func Test_entrancePage(t *testing.T) {
	//testRequestの結果を保存するやつ
	resp := httptest.NewRecorder()
	//テストのためのhttp request
	req, _ := http.NewRequest("GET", EntranceRoute, nil)
	//requestをサーバーに流して結果をrespに記録
	router.ServeHTTP(resp, req)

	//bodyを取り出し
	body, _ := ioutil.ReadAll(resp.Body)
	//ステータスコードは200のはず
	assert.Equal(t, 200, resp.Code)
	//titleはEntrance
	assert.Contains(t, string(body), "<title>Entrance</title>")
}

//loginページはセッション情報がなくても取得できる
func Test_loginPage(t *testing.T) {
	//testRequestの結果を保存するやつ
	resp := httptest.NewRecorder()
	//テストのためのhttp request
	req, _ := http.NewRequest("GET", LoginRoute, nil)
	//requestをサーバーに流して結果をrespに記録
	router.ServeHTTP(resp, req)

	//bodyを取り出し
	body, _ := ioutil.ReadAll(resp.Body)
	//ステータスコードは200のはず
	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "<title>Login and Register</title>")
}

//エントランスページはセッション情報がなくても取得できる
func Test_registerPage(t *testing.T) {
	//testRequestの結果を保存するやつ
	resp := httptest.NewRecorder()
	//テストのためのhttp request
	req, _ := http.NewRequest("GET", RegisterRoute, nil)
	//requestをサーバーに流して結果をrespに記録
	router.ServeHTTP(resp, req)

	//bodyを取り出し
	body, _ := ioutil.ReadAll(resp.Body)
	//ステータスコードは200のはず
	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "<title>Login and Register</title>")
}

//idとpasswordがそれぞれ８文字以上の英数字だと登録できる
func Test_canRegister_id_and_password_more8_and_alphanumeric(t *testing.T) {

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"test1234","Password":"qwer7890"}`

	req, _ := http.NewRequest(
		"POST",
		RegisterRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 200, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), "success")

}

//userIdの重複不可
func Test_cntRegister_same_id_and_password(t *testing.T) {

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"test1234","Password":"qwer7890"}`

	req, _ := http.NewRequest(
		"POST",
		RegisterRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 400, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), `error":"already use this id`)

}

//パスワードが同じ場合は許す
func Test_canRegister_same_password(t *testing.T) {

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"test5678","Password":"qwer7890"}`

	req, _ := http.NewRequest(
		"POST",
		RegisterRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 200, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), "success")

}

//提案
//７文字以下のuserIdは許さない
/*
func Test_cntRegister_id_less8(t *testing.T){

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"test567","Password":"qwer7890"}`

	req, _ := http.NewRequest(
			"POST",
			RegisterRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 400, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), `error":"いい感じのエラー文`)

}
*/

//提案
//７文字以下のpasswordは許さない
/*
func Test_cntRegister_password_less8(t *testing.T){

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"test5678","Password":"qwer789"}`

	req, _ := http.NewRequest(
			"POST",
			RegisterRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 400, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), `error":"いい感じのエラー文`)

}
*/

//提案
//英字のみのpasswordは許さない
/*
func Test_cntRegister_password_only_alphabet(t *testing.T){

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"test5678","Password":"qwertest"}`

	req, _ := http.NewRequest(
			"POST",
			RegisterRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 400, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), `error":"いい感じのエラー文`)

}
*/

//提案
//数字のみのpasswordは許さない
/*
func Test_cntRegister_password_only_num(t *testing.T){

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"test5678","Password":"11111111"}`

	req, _ := http.NewRequest(
			"POST",
			RegisterRoute,
			bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 400, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), `error":"いい感じのエラー文`)

}
*/

//登録済みのユーザーはログイン可能
func Test_canLogin_registered_user(t *testing.T) {

	resp := httptest.NewRecorder()

	jsonStr := `{"UserId":"test1234","Password":"qwer7890"}`
	req, _ := http.NewRequest(
		"POST",
		LoginRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")
	mainCookie = resp.Header().Get("Set-Cookie")
	oldCookie = resp.Header().Get("Set-Cookie")
}

//未登録のユーザーではログインできない
func Test_cntLogin_not_registered_user(t *testing.T) {

	resp := httptest.NewRecorder()

	jsonStr := `{"UserId":"te344567","Password":"3wer3333"}`
	req, _ := http.NewRequest(
		"POST",
		LoginRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `error":"user not exist`)
}

//ログインせずに議事録一覧ページにはいけない
//リダイレクト
func Test_redirect_meetingPage_not_logined(t *testing.T) {

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", GetSelectMeetingPageRoute, nil)
	router.ServeHTTP(resp, req)

	assert.Equal(t, 303, resp.Code)
}

//登録されていないユーザー情報を持ったsessionでは議事録一覧ページにアクセスできない
func Test_cntAccess_meetingPage_dummySession(t *testing.T) {

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", GetSelectMeetingPageRoute, nil)
	req.Header.Set("Cookie", dummyCookie)

	router.ServeHTTP(resp, req)

	//body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 303, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	//assert.Contains(t, string(body), "Invalid session ID")
}

//ログインしたなら議事録一覧ページに行ける
func Test_canAccess_meetingPage_logined(t *testing.T) {

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", GetSelectMeetingPageRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), "<title>Meetings</title>")
}

//登録していないユーザーは議事録一覧を取得不可
func Test_cntGetMeeting_not_logined_user(t *testing.T) {

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", GetMeetingRoute, nil)

	router.ServeHTTP(resp, req)

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `error":"Bad Request`)
}

//登録済みのユーザーは議事録を作成可能
func Test_canAddMeeting_logined_user(t *testing.T) {

	resp := httptest.NewRecorder()

	jsonStr := `{"meeting":"議事録テスト"}`
	req, _ := http.NewRequest(
		"POST",
		PostMeetingRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	req, _ = http.NewRequest("GET", GetMeetingRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), `[{"id":1,"name":"議事録テスト"}]`)
}

//同名の議事録は作成不可
func Test_cntAddMeeting_sameName(t *testing.T) {

	resp := httptest.NewRecorder()

	jsonStr := `{"meeting":"議事録テスト"}`
	req, _ := http.NewRequest(
		"POST",
		PostMeetingRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `error":"already use this name`)
}

//登録していないユーザーは議事録を作成不可
func Test_cntAddMeeting_not_logined_user(t *testing.T) {

	resp := httptest.NewRecorder()

	jsonStr := `{"meeting":"議事録ダミー"}`
	req, _ := http.NewRequest(
		"POST",
		PostMeetingRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `error":"Bad Request`)

}

//ログインせずに議事録ページにはいけない
//リダイレクト
func Test_redirect_minutesPage_not_logined(t *testing.T) {

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", GetMinutesPageRoute, nil)
	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `error":"Bad Request`)
}

//登録されていないユーザー情報を持ったsessionではアクセスできない
func Test_cntAccess_minutesPage_dummySession(t *testing.T) {

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", GetMinutesPageRoute, nil)
	req.Header.Set("Cookie", dummyCookie)

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), "Invalid session ID")
}

//ログインしたなら議事録ページに行ける
func Test_canAccess_minutesPage_logined(t *testing.T) {

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", GetMinutesPageRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), "<title>議事録テスト</title>")
}

//存在しない議事録ページにはいけない
func Test_cntAccess_minutesPage_notExist(t *testing.T) {

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", DummyMinutesPageRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 404, resp.Code)
	assert.Contains(t, string(body), `error":"Not Found`)
}

//ログアウト後に議事録ページにいけない
func Test_logout(t *testing.T) {

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", LogoutRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)

	assert.Equal(t, 303, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	mainCookie = resp.Header().Get("Set-Cookie")
	//ちょうど mysession = valueの部分が取り出せる
	assert.NotEqual(t, strings.Split(mainCookie, " ")[0], strings.Split(oldCookie, " ")[0])

	req, _ = http.NewRequest("GET", GetMinutesPageRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)

	assert.Equal(t, 303, resp.Code)
}

//登録していないユーザーはメッセージを取得不可
func Test_cntGetMessge_not_logined_user(t *testing.T) {

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", GetMessageRoute, nil)

	router.ServeHTTP(resp, req)
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `error":"Bad Request`)

}

//再度ログインが必要
//登録済みのユーザーはログイン可能
//session情報は毎回変わる
func Test_canLogin_registered_user_useDifferentSessionInfo(t *testing.T) {

	resp := httptest.NewRecorder()

	jsonStr := `{"UserId":"test1234","Password":"qwer7890"}`
	req, _ := http.NewRequest(
		"POST",
		LoginRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	//同じアカウントでも毎回session情報が変わることのテスト
	mainCookie = resp.Header().Get("Set-Cookie")
	//ちょうど mysession = valueの部分が取り出せる
	assert.NotEqual(t, strings.Split(mainCookie, " ")[0], strings.Split(oldCookie, " ")[0])
}

//存在しない議事録のメッセージは取得できない
func Test_cntGetMessage_notExistMinutes(t *testing.T) {

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", DummyGetMessgaeRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 404, resp.Code)
	assert.Contains(t, string(body), `error":"Not Found`)
}

//存在しない議事録へはメッセージを送信できない
func Test_cntPostMessage_notExistMinutes(t *testing.T) {

	resp := httptest.NewRecorder()

	jsonStr := `{"message":"カシスオレンジ"}`
	req, _ := http.NewRequest(
		"POST",
		DummyPostMessgaeRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 404, resp.Code)
	assert.Contains(t, string(body), `error":"Not Found`)
}

//登録済みのユーザーはメッセージを送信可能
func Test_canAddMessge_logined_user(t *testing.T) {

	resp := httptest.NewRecorder()

	jsonStr := `{"message":"カシスオレンジ"}`
	req, _ := http.NewRequest(
		"POST",
		PostMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	req, _ = http.NewRequest("GET", GetMessageRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), `[{"id":1,"addedBy":{"id":1,"name":"test1234"},"message":"カシスオレンジ"}]`)
}

//登録していないユーザーはメッセージを送信不可
func Test_cntAddMessge_not_logined_user(t *testing.T) {

	resp := httptest.NewRecorder()

	jsonStr := `{"message":"ジン"}`
	req, _ := http.NewRequest(
		"POST",
		PostMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `error":"Bad Request`)

}

//テストのために他ユーザーセッションを取得
func Test_getSubCookie(t *testing.T) {

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"subA2222","Password":"qwegds890"}`

	req, _ := http.NewRequest(
		"POST",
		RegisterRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	req, _ = http.NewRequest(
		"POST",
		LoginRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	subCookie = resp.Header().Get("Set-Cookie")

}

//異なるユーザーはメッセージを更新不可
func Test_cntUpdateMessge_different_user(t *testing.T) {

	resp := httptest.NewRecorder()

	jsonStr := `{"id":"1","message":"ストロングゼロ"}`
	req, _ := http.NewRequest(
		"POST",
		UpdateMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", subCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `"error":"Malformed request due to privileges"`)

	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", GetMessageRoute, nil)
	req.Header.Set("Cookie", subCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), `[{"id":1,"addedBy":{"id":1,"name":"test1234"},"message":"カシスオレンジ"}]`)
	assert.NotContains(t, string(body), `[{"id":1,"addedBy":{"id":1,"name":"test1234"},"message":"ストロングゼロ"}]`)
}

//登録していないユーザーはメッセージを更新不可
func Test_cntUpdateMessge_not_logined_user(t *testing.T) {

	resp := httptest.NewRecorder()

	jsonStr := `{"id":"1","message":"ストロングゼロ"}`
	req, _ := http.NewRequest(
		"POST",
		UpdateMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `"error":"Bad Request"`)

}

//同じユーザーはメッセージを更新可能
func Test_canUpdateMessge_same_user(t *testing.T) {

	resp := httptest.NewRecorder()

	jsonStr := `{"id":"1","message":"ストロングゼロ"}`
	req, _ := http.NewRequest(
		"POST",
		UpdateMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", GetMessageRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), `[{"id":1,"addedBy":{"id":1,"name":"test1234"},"message":"ストロングゼロ"}]`)
	assert.NotContains(t, string(body), `[{"id":1,"addedBy":{"id":1,"name":"test1234"},"message":"カシスオレンジ"}]`)
}

//異なるユーザーはメッセージを削除不可
func Test_cntDeleteMessge_different_user(t *testing.T) {

	resp := httptest.NewRecorder()

	jsonStr := `{"id":"1"}`
	req, _ := http.NewRequest(
		"POST",
		DeleteMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", subCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `"error":"Malformed request due to privileges"`)

	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", GetMessageRoute, nil)
	req.Header.Set("Cookie", subCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), `[{"id":1,"addedBy":{"id":1,"name":"test1234"},"message":"ストロングゼロ"}]`)
}

//登録していないユーザーはメッセージを削除不可
func Test_cntDeleteMessge_not_logined_user(t *testing.T) {

	resp := httptest.NewRecorder()

	jsonStr := `{"id":"1"}`
	req, _ := http.NewRequest(
		"POST",
		DeleteMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `"error":"Bad Request"`)

}

//同じユーザーはメッセージを削除可能
func Test_canDeleteMessge_same_user(t *testing.T) {

	resp := httptest.NewRecorder()

	jsonStr := `{"id":"1"}`
	req, _ := http.NewRequest(
		"POST",
		DeleteMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", GetMessageRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.NotContains(t, string(body), `[{"id":1,"addedBy":{"id":1,"name":"test1234"},"message":"ストロングゼロ"}]`)
}

//ログインしているユーザは重要単語を取得できる
func Test_canGetImportantWords_logined_user(t *testing.T) {
	resp := httptest.NewRecorder()

	jsonStr := `{"message":"寿司が食べたい。"}`
	req, _ := http.NewRequest(
		"POST",
		PostMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	req, _ = http.NewRequest("GET", GetImportantWordsRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), `["。","が","たい","寿司","食べ"]`)
}

//登録していないユーザーは重要単語の取得不可
func Test_cntGetImportantWords_not_logined_user(t *testing.T) {

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest(
		"GET",
		GetImportantWordsRoute,
		nil,
	)

	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `"error":"Bad Request"`)
	//念のため、後のテストのために"寿司が食べたい。"は消しておく
	resp = httptest.NewRecorder()

	jsonStr := `{"id":"2"}`
	req, _ = http.NewRequest(
		"POST",
		DeleteMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", GetMessageRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.NotContains(t, string(body), `[{"id":2,"addedBy":{"id":1,"name":"test1234"},"message":"寿司が食べたい。"}]`)
}

//ログインしている人は重要な文を取得できる
func Test_canGetImportantSentences_logined_user(t *testing.T) {
	resp := httptest.NewRecorder()

	jsonStr := `{"message":"支払い方法変更"}`
	req, _ := http.NewRequest(
		"POST",
		PostMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	jsonStr = `{"message":"別の支払い方法を希望"}`
	req, _ = http.NewRequest(
		"POST",
		PostMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	jsonStr = `{"message":"支払い方法変更お願いいたします"}`
	req, _ = http.NewRequest(
		"POST",
		PostMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	jsonStr = `{"message":"デビッドカード支払いはできないのでしょうか 別の支払い方方法はないのでしょうか？"}`
	req, _ = http.NewRequest(
		"POST",
		PostMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	jsonStr = `{"message":"クレジットカードでの支払いでしたが、フィッシング詐欺にあってクレジットカードを停止しました。支払い方法を変更してコンビニ支払いにしたいので支払い方法をお知らせください。"}`
	req, _ = http.NewRequest(
		"POST",
		PostMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	jsonStr = `{"message":"お世話になっております。次回の[製品名]のお支払い方法をクレジットカードにしたいのですが、[不満内容]その先がわかりません。お返事お願いします。[氏名]"}`
	req, _ = http.NewRequest(
		"POST",
		PostMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	req, _ = http.NewRequest("GET", GetImportantSentencesRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), `["別の支払い方法を希望","クレジットカードでの支払いでしたが、フィッシング詐欺にあってクレジットカードを停止しました。支払い方法を変更してコンビニ支払いにしたいので支払い方法をお知らせください。","支払い方法変更お願いいたします","お世話になっております。次回の[製品名]のお支払い方法をクレジットカードにしたいのですが、[不満内容]その先がわかりません。お返事お願いします。[氏名]","支払い方法変更"]`)
}

//登録していないユーザーは重要文の取得不可
func Test_cntGetImportantSentences_not_logined_user(t *testing.T) {

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest(
		"GET",
		GetImportantSentencesRoute,
		nil,
	)

	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `"error":"Bad Request"`)
	//念のため、後のテストのために前に入力した文章は消しておく
	resp = httptest.NewRecorder()

	jsonStr := `{"id":"3"}`
	req, _ = http.NewRequest(
		"POST",
		DeleteMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", GetMessageRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.NotContains(t, string(body), `[{"id":3,"addedBy":{"id":1,"name":"test1234"},"message":"支払い方法変更"}]`)

	jsonStr = `{"id":"4"}`
	req, _ = http.NewRequest(
		"POST",
		DeleteMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", GetMessageRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.NotContains(t, string(body), `[{"id":4,"addedBy":{"id":1,"name":"test1234"},"message":"別の支払い方法を希望"}]`)

	jsonStr = `{"id":"5"}`
	req, _ = http.NewRequest(
		"POST",
		DeleteMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", GetMessageRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.NotContains(t, string(body), `[{"id":5,"addedBy":{"id":1,"name":"test1234"},"message":"支払い方法変更お願いいたします"}]`)

	jsonStr = `{"id":"6"}`
	req, _ = http.NewRequest(
		"POST",
		DeleteMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", GetMessageRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.NotContains(t, string(body), `[{"id":6,"addedBy":{"id":1,"name":"test1234"},"message":"デビッドカード支払いはできないのでしょうか 別の支払い方方法はないのでしょうか？"}]`)

	jsonStr = `{"id":"7"}`
	req, _ = http.NewRequest(
		"POST",
		DeleteMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", GetMessageRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.NotContains(t, string(body), `[{"id":7,"addedBy":{"id":1,"name":"test1234"},"message":"クレジットカードでの支払いでしたが、フィッシング詐欺にあってクレジットカードを停止しました。支払い方法を変更してコンビニ支払いにしたいので支払い方法をお知らせください。"}]`)

	jsonStr = `{"id":"8"}`
	req, _ = http.NewRequest(
		"POST",
		DeleteMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", mainCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), "success")

	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", GetMessageRoute, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.NotContains(t, string(body), `[{"id":8,"addedBy":{"id":1,"name":"test1234"},"message":"お世話になっております。次回の[製品名]のお支払い方法をクレジットカードにしたいのですが、[不満内容]その先がわかりません。お返事お願いします。[氏名]"}]`)
}

//ログインしていないユーザーはユーザー情報は帰らない
func Test_cntGetUserInfo_not_logined_user(t *testing.T) {

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", GetUserInfo, nil)

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `"error":"Bad Request"`)

}

//ログインしているユーザーはユーザー情報が帰る
//セキュリティ的に大丈夫か？
func Test_canGetUserInfo_logined_user(t *testing.T) {

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", GetUserInfo, nil)
	req.Header.Set("Cookie", mainCookie)

	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Header().Get("Set-Cookie"))
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, string(body), `{"id":1,"name":"test1234"}`)

}

//有効期限切れのsessionIDの利用をテストするユーザを作成
func Test_Login_registered_user_tempUse00(t *testing.T) {

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"temp00","Password":"temptemp00"}`

	req, _ := http.NewRequest(
		"POST",
		RegisterRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	req, _ = http.NewRequest(
		"POST",
		LoginRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	tempCookie = resp.Header().Get("Set-Cookie")
}

//時間切れのセッションは適宜データベースから破棄されている
func Test_session_database_update(t *testing.T) {

	//DBのsession情報を時間切れになるように無理やり設定
	sessionTimeSet10daysLater("temp00")

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", GetMinutesPageRoute, nil)
	req.Header.Set("Cookie", tempCookie)

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), "Invalid session ID")

}

//有効期限切れのsessionIDの利用をテストするユーザを作成
func Test_Login_registered_user_tempUse01(t *testing.T) {

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"temp00","Password":"temptemp00"}`

	req, _ := http.NewRequest(
		"POST",
		RegisterRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	req, _ = http.NewRequest(
		"POST",
		LoginRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	tempCookie = resp.Header().Get("Set-Cookie")
}

//時間切れのセッションはセッションタイムアウトになる
func Test_session_timeout(t *testing.T) {

	//DBのsession情報を時間切れになるように無理やり設定
	sessionTimeSetNow("temp00")

	resp := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", GetMinutesPageRoute, nil)
	req.Header.Set("Cookie", tempCookie)

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	//順序注意　assert.Contains 第二引数に第三引数の要素が含まれているか
	assert.Contains(t, string(body), "Session time out")

}

//middleware内全ての処理は登録していないユーザは受け付けない**必ずlogedIn内のリクエスでテストすること**
//登録していないユーザーはメッセージを送信不可
func Test_middeleware_logedIn_not_logined_user(t *testing.T) {

	resp := httptest.NewRecorder()

	jsonStr := `{"message":"ジン"}`
	req, _ := http.NewRequest(
		"POST",
		PostMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), `error":"Bad Request`)

}

//middleware内全ての処理は不当なセッションIDは受け付けない**必ずlogedIn内のリクエスでテストすること**
//不当なセッションIDのユーザーはメッセージを送信不可
func Test_middeleware_logedIn_invalid_user(t *testing.T) {

	resp := httptest.NewRecorder()
	jsonStr := `{"message":"ジン"}`
	req, _ := http.NewRequest(
		"POST",
		PostMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", dummyCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), "Invalid session ID")
}

//有効期限切れのsessionIDの利用をテストするユーザを作成
func Test_Login_registered_user_tempUse02(t *testing.T) {

	resp := httptest.NewRecorder()
	//送信するjson
	jsonStr := `{"UserId":"temp00","Password":"temptemp00"}`

	req, _ := http.NewRequest(
		"POST",
		RegisterRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	req, _ = http.NewRequest(
		"POST",
		LoginRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	tempCookie = resp.Header().Get("Set-Cookie")
}

//middleware内全ての処理は有効期限の切れたセッションIDは受け付けない**必ずlogedIn内のリクエスでテストすること**
//有効期限切れのユーザーはメッセージを送信不可
func Test_middeleware_logedIn_session_timeout(t *testing.T) {

	sessionTimeSetNow("temp00")
	resp := httptest.NewRecorder()
	jsonStr := `{"message":"ジン"}`
	req, _ := http.NewRequest(
		"POST",
		PostMessageRoute,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Cookie", tempCookie)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, string(body), "Session time out")
}
