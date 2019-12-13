package handler

import (
	"imWebSocket/model"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"imchat/hello5.3/util"
)

func CreateChatGroup(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	//解析参数
	createChatGroups := model.ContactGroups{}
	if err := util.Bind(req, &createChatGroups); err != nil {
		lgorus.Println("create chat group binding request params failed.", err.Error())
		ResponseJson(rw, "", ErrBind)
		return
	}

	if len(createChatGroups.Name) == 0 {
		lgorus.Println("create chat group no set group name")
		ResponseJson(rw, "", CreateContactGroupsNoSetName)
		return
	}

	contactGroups := model.ContactGroups{
		Ownerid: createChatGroups.Ownerid,
	}

	//查找该用户作为群主的已经创建过群的数目
	if _, groupsNumber, err := contactGroups.GetAllByOwnerId(); err != nil {
		logrus.Println("get all contact groups failed.", err.Error())
		ResponseJson(rw, "", InternalServerError)
		return
	}

	//检查是否创建的额群聊数目超过 5个
	if groupsNumber > 5 {
		logrus.Println("owner contact groups number over 5。")
		ResponseJson(rw, "", CreateContatcGroupsTooMuch)
		return
	}

	if err := contactGroups.Create(); err != nil {
		logrus.Println("create contact groups failed.", err.Error())
		ResponseJson(rw, "", CreateContactGroupsFailed)
		return
	}

	ResponseJson(rw, "", OK)
}
