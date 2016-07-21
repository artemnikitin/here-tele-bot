package logic

const (
	askForLocation = "Please, provide your location.\n" +
		"Share your location via Telegram app for iOS or Android" +
		"and then write you query like `coffee`\n" +
		"or just write a search query in a format `something in place` like `coffee in New York`"
	locationAccepted = "Roger that. Location accepted. I'll keep it for a while."
	unknownMessage   = "Hmmm... I don't understand you..." +
		"You can share your location via Telegram app for iOS or Android" +
		"and then write you query like `chinese`\n" +
		"or just write a search query in a format `something in place` like `chinese in Berlin Mitte`"
	welcomeMessage = "Hi *%s*, it nice that you decided to join us.\n" +
		"How to use it: share your location via Telegram app for iOS or Android," +
		"and then write you query like `park`\n" +
		"or just write a search query in a format `something in place` like `park in London`"
	errorHappened = "Sorry, error happened"
	nothingFound  = "Sorry, but I found nothing. Try to search something else."
)
