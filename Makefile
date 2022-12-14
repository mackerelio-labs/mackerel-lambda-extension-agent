make build-extension:
	mkdir -p extensions
	GOOS=linux GOARCH=amd64 go build -o extensions/mackerel-lambda-extension-agent
	zip mackerel-lambda-extension-agent.zip extensions/*
	$(RM) -r extensions/
