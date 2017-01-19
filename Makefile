EXEC=kbdashboard
VER=`grep "const VERSION" cmd_version.go  | cut -d "=" -f 2 | cut -d '"' -f 2`
TAR=$(EXEC)-$(VER).tar.gz

all:
	@go build -o $(EXEC)

run:$(EXEC)
	@./$(EXEC)

install:
	cp $(EXEC) /usr/local/bin
	sudo cp ./kbdashboard-completion /etc/bash_completion.d/kbdashboard

archive:
	@echo "archive to $(TAR)"
	@git archive master --prefix="$(EXEC)-$(VER)/" --format tar.gz -o $(TAR)
