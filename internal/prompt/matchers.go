package prompt

func Always(prompt *Prompt) bool {
	return true
}

func AnsweredValueIs(prompt *Prompt, key string, value string) bool {
	return prompt.GetAnswer(key) == value
}
