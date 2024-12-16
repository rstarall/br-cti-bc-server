package fabric

import (
	"fmt"
	"github.com/righstar2020/br-cti-bc-server/global"
)

// 注册评论
func RegisterComment(txRawMsgData []byte) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码注册评论
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "RegisterComment", [][]byte{txRawMsgData})
	if err != nil {
		return "", fmt.Errorf("调用链码失败: %v", err)
	}

	return string(resp), nil
}

// 审核评论
func ApproveComment(txRawMsgData []byte) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码审核评论
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "ApproveComment", [][]byte{txRawMsgData})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 查询单个评论
func QueryComment(commentID string) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码查询评论
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryComment", [][]byte{[]byte(commentID)})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 查询指定文档的评论列表
func QueryAllCommentsByRefID(refID string) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 调用链码查询评论列表
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryAllCommentsByRefID", [][]byte{[]byte(refID)})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// 分页查询评论列表
func QueryCommentsByRefIDWithPagination(refID string, page int, pageSize int) (string, error) {
	// 创建通道客户端
	client, err := CreateChannelClient(global.FabricSDK)
	if err != nil {
		return "", err
	}

	// 构造查询参数
	args := [][]byte{
		[]byte(refID),
		[]byte(fmt.Sprintf("%d", page)),
		[]byte(fmt.Sprintf("%d", pageSize)),
	}

	// 调用链码分页查询评论
	resp, err := InvokeChaincode(client, global.MainChaincodeName, "QueryCommentsByRefIDWithPagination", args)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}
