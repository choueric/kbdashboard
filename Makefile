EXEC=kbdashboard

all:
	@go build -o $(EXEC)

run:$(EXEC)
	@./$(EXEC)

install:
	cp $(EXEC) /usr/local/bin
	sudo cp ./kbdashboard-completion /etc/bash_completion.d/kbdashboard
