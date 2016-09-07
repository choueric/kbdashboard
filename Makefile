EXEC=kbdashboard

all:$(EXEC) run

$(EXEC):main.go
	@go build -o $(EXEC)

run:$(EXEC)
	@./$(EXEC)
