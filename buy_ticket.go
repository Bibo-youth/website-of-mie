package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// 初始化一个带 Cookie 管理的 HTTP 客户端
func createClient() (*http.Client, error) {
	// 创建 CookieJar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %v", err)
	}

	// 创建 HTTP 客户端并绑定 CookieJar
	client := &http.Client{
		Jar: jar,
	}
	return client, nil
}

// 登录 12306 并保存 Cookie
func login12306(client *http.Client, username, password string) error {
	loginURL := "https://kyfw.12306.cn/passport/web/login"
	// 登录表单数据
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("checkMode", "0")
	data.Set("randCode", "12345")
	data.Set("appid", "otn") // 12306 的登录 API 需要 appid 参数

	req, err := http.NewRequest("POST", loginURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发起请求
	resp, err := client.PostForm(loginURL, data)
	if err != nil {
		return fmt.Errorf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: %s", string(body))
	}
	fmt.Println("Login successful. Cookies stored in client.", string(body))
	fmt.Println("############################################")
	return nil
}

// 带 Cookie 的请求示例
func queryTickets(client *http.Client) error {
	queryURL := "https://kyfw.12306.cn/otn/leftTicket/queryO"
	cookieURL := "https://kyfw.12306.cn/otn/view/index.html"

	// 查询参数
	params := url.Values{}
	params.Set("leftTicketDTO.train_date", "2024-12-15") // 车票日期
	params.Set("leftTicketDTO.from_station", "SHH")      // 上海
	params.Set("leftTicketDTO.to_station", "BJP")        // 北京
	params.Set("purpose_codes", "ADULT")                 // 成人票

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", queryURL, params.Encode()), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	parsedURL, _ := url.Parse(cookieURL)
	cookies := client.Jar.Cookies(parsedURL)
	var builder strings.Builder
	last := len(cookies) - 1
	fmt.Println("Stored Cookies:")
	for i, cookie := range cookies {
		fmt.Println("cookie.Name: : cookie.Value:\n", cookie.Name, cookie.Value)
		str := cookie.Name + "=" + cookie.Value
		builder.WriteString(str)
		// 如果不是最后一个元素，则添加分号
		if i != last {
			builder.WriteByte(';')
		}
	}
	result := builder.String()
	fmt.Println("result:", result)
	req.Header.Set("Cookie", result)
	// 发起请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ticket query failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Query Response: %s\n", string(body))
	return nil
}

func main() {
	// 创建带 Cookie 管理的客户端
	client, err := createClient()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 登录
	username := "18510986700"
	password := "ysyyang0402yy"
	if err := login12306(client, username, password); err != nil {
		fmt.Println("Login Error:", err)
		return
	}

	// 查询车票
	if err := queryTickets(client); err != nil {
		fmt.Println("Query Error:", err)
		return
	}
}
