package telegram

import (
	"testing"
)

const (
	kTestToken = "5698315188:AAErx6jGOgcqlCZidS3bjmCtvPCEJmuVOo"
	kMyChatId = 110751423
)

func TestSendMessageWithBadHtmlTag(t *testing.T) {
	err := SendMessage(kTestToken, kMyChatId, "<200")
	if err == nil {
		t.Fatalf("Expected error, but succesfully sent")
	}
	t.Logf("Expected error: %v", err)
}

func TestSendMessageWithHiperlink(t *testing.T) { 
	messageWithLink := `                                          ?column?                                          
	--------------------------------------------------------------------------------------------
	 <a href="https://cmp.wildberries.ru/campaigns/list/active/edit/search/4122050">4122050</a> `
	err := SendMessage(kTestToken, kMyChatId, messageWithLink)
	if err != nil {
		t.Fatal(err)
	}
}
