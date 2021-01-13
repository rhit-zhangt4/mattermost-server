// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package extchat

import (
	"fmt"
	"os"
	"time"

	"github.com/Arman92/go-tdlib"
	"github.com/mattermost/mattermost-server/v5/app"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/otiai10/copy"
)

var client *tdlib.Client

type TelegramAdapter struct {
}

func (adapter *TelegramAdapter) StartAuthentication(a app.AppIface, username string) *model.AppError {
	fmt.Println("Start Auth" + username)
	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./logs/tdlib.log")
	var err *model.AppError

	client, err = adapter.getTdClient(a, username)
	if err != nil {
		return err
	}
	fmt.Println("Client OK")
	// //_, _ = client.LogOut()
	go adapter.sendPhoneNumberToAuthorize(client, username)

	// else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
	// 	fmt.Print("Enter code: ")
	// 	var code string
	// 	fmt.Scanln(&code)
	// 	_, err := client.SendAuthCode(code)
	// 	if err != nil {
	// 		fmt.Printf("Error sending auth code : %v", err)
	// 	}
	// } else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPasswordType {
	// 	fmt.Print("Enter Password: ")
	// 	var password string
	// 	fmt.Scanln(&password)
	// 	_, err := client.SendAuthPassword(password)
	// 	if err != nil {
	// 		fmt.Printf("Error sending auth password: %v", err)
	// 	}
	// } else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
	// 	fmt.Println("Authorization Ready! Let's rock")
	// 	break
	// }

	// contact, _ := client.GetContacts()
	// fmt.Println("Here are your contacts")
	// fmt.Println(contact.TotalCount)
	// fmt.Println(contact.UserIDs)

	// for _, s := range contact.UserIDs {
	// 	user, _ := client.GetUser(s)
	// 	fmt.Print(user.FirstName + " " + user.LastName + " -- ID: ")
	// 	fmt.Println(user.ID)

	// }

	// chat, _ := client.CreatePrivateChat(1359977993, false)

	// inputMsgTxt := tdlib.NewInputMessageText(tdlib.NewFormattedText("/start", nil), true, true)
	// client.SendMessage(chat.ID, 0, false, true, nil, inputMsgTxt)
	// fmt.Print("Test chat ID changes ")
	// fmt.Println(chat.ID)

	// // Main loop
	// for update := range client.GetRawUpdatesChannel(10) {
	// 	// Show all updates
	// 	data := update.Data
	// 	if data["@type"].(string) == "updateNewMessage" {
	// 		chatId := data["message"].(map[string]interface{})["chat_id"]
	// 		chat, _ := client.GetChat(int64(chatId.(float64)))
	// 		chatType := chat.Type
	// 		fmt.Println("--------------- TYPE: " + chatType.GetChatTypeEnum() + "---------------")
	// 		fmt.Println("--------------- TITLE: " + chat.Title + "---------------")

	// 		sender := data["message"].(map[string]interface{})["sender"].(map[string]interface{})["user_id"]
	// 		senderUser, _ := client.GetUser(int32(sender.(float64)))
	// 		fmt.Print(senderUser.FirstName + " " + senderUser.LastName + ":  ")

	// 		text := data["message"].(map[string]interface{})["content"].(map[string]interface{})["text"].(map[string]interface{})["text"]
	// 		fmt.Println(text)
	// 	}

	// 	// fmt.Print("\n\n")
	// }
	return nil
}

func (adapter *TelegramAdapter) sendPhoneNumberToAuthorize(client *tdlib.Client, username string) {
	for {
		currentState, _ := client.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			fmt.Println("Sending Phone")
			//fmt.Print("Enter phone: ")
			//var number string
			//fmt.Scanln(&number)
			_, err := client.SendPhoneNumber(username)
			if err != nil {

				//error
				fmt.Printf("Error sending phone number: %v", err)
			}
			// defer client.DestroyInstance()
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (adapter *TelegramAdapter) verifyPasscode(client *tdlib.Client, username string, code string) {
	for {
		currentState, err := client.Authorize()
		if err != nil {
			fmt.Printf("Error verifying code: %v", err)
			continue
		}
		fmt.Print("currentState is : ")
		fmt.Println(currentState)
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			fmt.Println("Sending Phone")
			//fmt.Print("Enter phone: ")
			//var number string
			//fmt.Scanln(&number)
			_, err := client.SendPhoneNumber(username)
			if err != nil {

				//error
				fmt.Printf("Error sending phone number: %v", err)
			}
			break
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
			fmt.Println("Sending Code")
			//fmt.Print("Enter phone: ")
			//var number string
			//fmt.Scanln(&number)
			_, err := client.SendAuthCode(code)
			if err != nil {

				//error
				fmt.Printf("Error sending code number: %v", err)
			}
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (adapter *TelegramAdapter) VerifyPasscode(a app.AppIface, username string, code string) (*model.ExtRef, *model.AppError) {
	fmt.Print("Start verifying: ")
	// client, err := adapter.getTdClient(a, username)
	// if err != nil {
	// 	return nil, err
	// }
	fmt.Print("client is: ")
	fmt.Println(client)
	fmt.Print("Code is: ")
	fmt.Println(code)
	go adapter.verifyPasscode(client, username, code)
	return nil, nil
}

func (adapter *TelegramAdapter) getTdClient(a app.AppIface, username string) (*tdlib.Client, *model.AppError) {
	apiId, err := a.Srv().Store.Secret().GetBySecretName("TelegramAPIId")
	if err != nil {
		//error
		fmt.Printf("ApiId Err %v", err)
	}
	apiHash, err := a.Srv().Store.Secret().GetBySecretName("TelegramAPIHash")
	if err != nil {
		//error
		fmt.Printf("ApiHash Err %v", err)
	}

	var baseDirectory string

	if _, err := os.Stat("./data/extchat/telegram/" + username); os.IsNotExist(err) {
		// use temp
		baseDirectory = "./temp/extchat/telegram/" + username
	} else {
		// use data
		baseDirectory = "./data/extchat/telegram/" + username
	}

	client := tdlib.NewClient(tdlib.Config{
		APIID:               apiId.SecretValue,
		APIHash:             apiHash.SecretValue,
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
		UseMessageDatabase:  true,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseTestDataCenter:   false,
		DatabaseDirectory:   baseDirectory + "/tdlib-db",
		FileDirectory:       baseDirectory + "/tdlib-files",
		IgnoreFileNames:     false,
	})
	return client, nil
}

func (adapter *TelegramAdapter) copyDataFromTemp(username string) *model.AppError {
	err := copy.Copy("./temp/extchat/telegram/"+username, "./data/extchat/telegram/"+username)
	if err != nil {
		//error
		return nil
	}
	return nil
}
