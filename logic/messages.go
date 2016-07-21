package logic

const (
	AskForLocation = "Please, provide your location.\n" +
		"Share your location via Telegram app for iOS or Android" +
		"and then write you query like `coffee`\n" +
		"or just write a search query in a format `something in place` like `coffee in New York`"
	LocationAccepted = "Roger that. Location accepted. I'll keep it for a while."
	UnknownMessage = "Hmmm... I don't understand you..." +
		"You can share your location via Telegram app for iOS or Android" +
		"and then write you query like `chinese`\n" +
		"or just write a search query in a format `something in place` like `chinese in Berlin Mitte`"
	WelcomeMessage = "Hi *%s*, it nice that you decided to join us.\n" +
		"How to use it: share your location via Telegram app for iOS or Android," +
		"and then write you query like `park`\n" +
		"or just write a search query in a format `something in place` like `park in London`"
	ErrorHappened = "Sorry, error happened"
	NothingFound = "Sorry, but I found nothing. Try to search something else."
)
