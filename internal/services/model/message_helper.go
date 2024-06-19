package model

import "fmt"

func CreateStartGameMsg(contentURL string) string {
	return fmt.Sprintf(`
	{
		"type": "room",
		"data": {
			"action": "start",
			"data": {
				"contentLink": "%s"
			}			
		}
	}
	`, contentURL)
}
