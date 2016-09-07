EXEC=kbdashboard

all:$(EXEC) run

$(EXEC):main.go ../kbd/config.go
	@go build -o $(EXEC)

run:$(EXEC)
	@./$(EXEC)

install:
	cp $(EXEC) ~/bin
