package ipfs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// 通过指定节点获取IPFS内容
func GetIPFSContentFromNode(hash string, nodeAddr string) (string, []byte, error) {
	// 修改URL构建方式，确保正确的格式
	url := fmt.Sprintf("%s/ipfs/%s", nodeAddr, hash)
	// 创建客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("请求IPFS节点失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("IPFS节点响应错误,状态码: %d", resp.StatusCode)
	}

	// 读取响应内容
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("读取内容失败: %v", err)
	}

	return resp.Request.URL.String(), content, nil
}

var NodeAddrs = []string{
	"http://127.0.0.1:8080",
	// "https://ipfs.io",
	// "https://dweb.link",
	// "https://gateway.pinata.cloud",
}

// 通过多个节点尝试获取IPFS内容
func GetIPFSContentWithFallback(hash string) (string, []byte, error) {
	var lastErr error

	// 依次尝试所有节点
	for _, nodeAddr := range NodeAddrs {
		url, content, err := GetIPFSContentFromNode(hash, nodeAddr)
		if err == nil {
			return url, content, nil
		}
		lastErr = err
	}

	return "", nil, fmt.Errorf("所有节点均访问失败，最后错误: %v", lastErr)
}
func GetIPFSServerHost() string {
	return NodeAddrs[0]
}
