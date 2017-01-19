EXEC=kbdashboard
SRCS=*.go
COMP=$(EXEC)-completion

VER=`grep "const VERSION" cmd_version.go  | cut -d "=" -f 2 | cut -d '"' -f 2`
TAR=$(EXEC)-$(VER).tar.gz

GIT_COMMIT=`git log --pretty=format:"%h" -1`
GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`
BUILD_TIME=`date +%Y-%m-%d:%H:%M:%S`
X_ARGS="-X main.GIT_COMMIT=$(GIT_COMMIT) -X main.GIT_BRANCH=$(GIT_BRANCH) -X main.BUILD_TIME=$(BUILD_TIME)"

all:$(EXEC) $(COMP)

$(EXEC): $(SRCS)
	@echo "Build Version: $(VER)-$(GIT_COMMIT)-$(GIT_BRANCH)"
	@go build -ldflags $(X_ARGS) -o $(EXEC) -v

$(COMP): $(EXEC)
	@./$(EXEC) completion

install:$(EXEC) $(COMP)
	@cp -v $(EXEC) /usr/local/bin
	@cp -v $(COMP) /etc/bash_completion.d/kbdashboard

clean:
	@rm -rfv $(EXEC) $(COMP)

archive:
	@echo "archive to $(TAR)"
	@git archive master --prefix="$(EXEC)-$(VER)/" --format tar.gz -o $(TAR)
