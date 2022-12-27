package main

func Start(conf *Config, eventTimespan int, numEvents int, numIssues int) {
	eventsDriver := NewEventsHTTPDriver(config.WebHook, conf.Secret)
	eventsService := NewEventsService(eventsDriver)
	generatorHandler := NewGeneratorHandler(eventsService)
	generatorHandler.Run(eventTimespan, numEvents, numIssues)
}
