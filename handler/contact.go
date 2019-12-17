package handler

import (
	"imWebSocket/model"
	"imWebSocket/pkg/request"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func AddFriend(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	//绑定参数 解析请求参数
	addFriendsRequest := AddFriendsRequest{}
	if err := request.Bind(req, &addFriendsRequest); err != nil {
		logrus.Println("add friends failed for binding failed.", err)
		ResponseJson(rw, "", ErrBind)
		return
	}

	//不能添加自己作为自己的好友
	if addFriendsRequest.OwnerId == addFriendsRequest.DstUserId {
		logrus.Println("can not add self as friend")
		ResponseJson(rw, "", AddFriendsCanNotSelf)
		return
	}

	//是否已经添加过该好友 重复添加
	contact := model.Contact{}
	contact.OwnerId = addFriendsRequest.OwnerId
	contact.DstUserId = addFriendsRequest.DstUserId
	contact.Cate = addFriendsRequest.Cate
	if err := contact.CheckIsRepeatedAddFriends(); err == nil {
		logrus.Println("add friends repeated. please confirm and retry.")
		ResponseJson(rw, "", AddFriendsRepeatedInSameCate)
		return
	}

	//插入好友关系
	if err := contact.Create(); err != nil {
		logrus.Println("add friends to db failed.", err)
		ResponseJson(rw, "", AddFriendsFailed)
		return
	}

	//添加好友成功
	ResponseJson(rw, contact, nil)
}

func JoinContactGroup(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	//绑定参数 解析请求参数
	joinContactGroupsRequest := JoinContactGroupRequest{}
	if err := request.Bind(req, &joinContactGroupsRequest); err != nil {
		logrus.Println("join contact groups request bind failed.", err.Error())
		ResponseJson(rw, "", ErrBind)
		return
	}

	//查看这个群是否存在
	contactGroups := model.ContactGroups{}
	contactGroups.ID = joinContactGroupsRequest.GroupsID
	if addFriendsRequest.OwnerId == addFriendsRequest.DstUserId {
		logrus.Println("can not add self as friend")
		ResponseJson(rw, "", AddFriendsCanNotSelf)
		return
	}

	//是否已经添加过该好友 重复添加
	contact := model.Contact{}
	contact.OwnerId = addFriendsRequest.OwnerId
	contact.DstUserId = addFriendsRequest.DstUserId
	contact.Cate = addFriendsRequest.Cate
	if err := contact.CheckIsRepeatedAddFriends(); err == nil {
		logrus.Println("add friends repeated. please confirm and retry.")
		ResponseJson(rw, "", AddFriendsRepeatedInSameCate)
		return
	}

	//插入好友关系
	if err := contact.Create(); err != nil {
		logrus.Println("add friends to db failed.", err)
		ResponseJson(rw, "", AddFriendsFailed)
		return
	}

	//添加好友成功
	ResponseJson(rw, contact, nil)
}

func LoadFriend(rw http.ResponseWriter, req *http.Request) {
	var loadFriendArgs ContactArg
	//绑定请求参数
	if err := request.Bind(req, &loadFriendArgs); err != nil {
		logrus.Println("load friend bind args failed.", err.Error())
		ResponseJson(rw, "", ErrBind)
		return
	}
	contact := model.Contact{}
	contact.UserId = loadFriendArgs.UserId
	if users, usersNumber, err := contact.SearchFriend(); err != nil {
		logrus.Println("serach my friend failed.", err.Error())
		ResponseJson(rw, "", SearchMyFriendFailed)
		return
	}
	response := map[string]interface{}{
		"userList": users,
		"count":    usersNumber,
	}

	ResponseJson(rw, response, nil)
}

func LoadCommunity(rw http.ResponseWriter, req *http.Request) {
	var loadCommunityArgs ContactArgs
	//绑定参数
	if err := request.Bind(req, &loadCommunityArgs); err != nil {
		logrus.Println("bind load community args failed.", err.Error())
		ResponseJson(rw, "", ErrBind)
		return
	}

	contact := model.Contact{}
	contact.UserId = loadCommunityArgs.User
	if communityies, count, err := contact.SearchCommunity(); err != nil {
		logrus.Println("search community failed.", err.Error())
		ResponseJson(rw, "", SearchMyGroupsFailed)
		return
	}

	response := map[string]interface{}{
		"List":  communityies,
		"count": count,
	}

	ResponseJson(rw, response, nil)
}
