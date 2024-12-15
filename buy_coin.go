package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
)

// 登录请求数据
type LoginRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	VerifyCode string `json:"verifyCode"`
}

// 预约请求数据
type ReserveRequest struct {
	NetPointID string `json:"netPointId"`
	CoinType   string `json:"coinType"`
	Amount     int    `json:"amount"`
}

// 初始化 HTTP 客户端
func createClientCoin() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %v", err)
	}
	client := &http.Client{
		Jar: jar,
	}
	return client, nil
}

// 登录
func login(client *http.Client, username, password, verifyCode string) error {
	loginURL := "https://eapply.abchina.com/coin/login"

	// 准备登录请求数据
	data := LoginRequest{
		Username:   username,
		Password:   password,
		VerifyCode: verifyCode,
	}
	jsonData, _ := json.Marshal(data)

	// 发送登录请求
	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create login request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	// 检查登录结果
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: %s", string(body))
	}

	fmt.Println("Login successful!")
	return nil
}

// 获取网点信息
func queryNetPoints(client *http.Client) ([]map[string]interface{}, error) {
	queryURL := "https://eapply.abchina.com/coin/queryNetPoints"
	// 发送 GET 请求
	resp, err := client.Get(queryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to query net points: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Query Response:", string(body))
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// 提取网点数据
	points, ok := result["data"].([]map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}
	return points, nil
}

// 提交预约请求
func reserve(client *http.Client, reserveData ReserveRequest) error {
	reserveURL := "https://eapply.abchina.com/coin/reserve"

	// 准备预约数据
	jsonData, _ := json.Marshal(reserveData)
	req, err := http.NewRequest("POST", reserveURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create reserve request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("reservation request failed: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Reserve Response:", string(body))
	return nil
}

// 下载验证码图片
func downloadCaptcha(client *http.Client) (string, error) {
	captchaURL := "https://eapply.abchina.com/coin/captcha"
	resp, err := client.Get(captchaURL)
	if err != nil {
		return "", fmt.Errorf("failed to download captcha: %v", err)
	}
	defer resp.Body.Close()

	// 保存图片到本地
	fileName := "/Users/bytedance/Downloads/captcha.jpg"
	file, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to save captcha: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write captcha to file: %v", err)
	}

	fmt.Println("Captcha saved to", fileName)
	return fileName, nil
}

func main() {
	// 初始化客户端
	client, err := createClientCoin()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//// 下载验证码
	//captchaFile, err := downloadCaptcha(client)
	//if err != nil {
	//	fmt.Println("Captcha Error:", err)
	//	return
	//}
	//
	//// 用户输入验证码
	//var verifyCode string
	//fmt.Printf("Enter Captcha (see %s): ", captchaFile)
	//fmt.Scan(&verifyCode)
	//
	//// 登录账号
	//username := "your_username"
	//password := "your_password"
	//if err := login(client, username, password, verifyCode); err != nil {
	//	fmt.Println("Login Error:", err)
	//	return
	//}

	// 查询网点信息
	points, err := queryNetPoints(client)
	if err != nil {
		fmt.Println("Query Error:", err)
		return
	}

	// 打印可预约的网点
	for _, point := range points {
		fmt.Printf("NetPoint: %v\n", point)
	}

	// 提交预约
	reserveData := ReserveRequest{
		NetPointID: "net_point_id", // 替换为实际网点 ID
		CoinType:   "coin_type",    // 替换为实际纪念币类型
		Amount:     1,              // 替换为实际预约数量
	}
	if err := reserve(client, reserveData); err != nil {
		fmt.Println("Reserve Error:", err)
		return
	}

	fmt.Println("Reservation completed successfully!")
}
