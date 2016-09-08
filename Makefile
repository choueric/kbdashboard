EXEC=kbdashboard

all:
	@go build -o $(EXEC)

run:$(EXEC)
	@./$(EXEC)

install:
	cp $(EXEC) ~/bin
