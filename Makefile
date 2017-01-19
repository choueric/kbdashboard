EXEC = kbdashboard
SRCS = *.go
COMP = $(EXEC).bash-completion

VER = `grep "const VERSION" cmd_version.go  | cut -d "=" -f 2 | cut -d '"' -f 2`
TAR = $(EXEC)-$(VER).tar.gz

BUILD_TIME = `date +%Y-%m-%d:%H:%M:%S`

# Comment this if do not want to include build-time-string in the executable file.
X_ARGS += -X main.BUILD_TIME=$(BUILD_TIME)
X_ARGS += -X main.COMP_FILENAME=$(COMP)

BIN = $(DESTDIR)/usr/bin
COMP_DIR = $(DESTDIR)/etc/bash_completion.d

all:bin $(COMP)

bin:
	@echo "Build Version: $(VER)"
	@go build -ldflags "$(X_ARGS)" -o $(EXEC) -v

$(COMP): $(EXEC)
	@./$(EXEC) completion

install:$(EXEC) $(COMP)
	install -d $(BIN) $(COMP_DIR)
	install $(EXEC) $(BIN)
	install $(COMP) $(COMP_DIR)/$(EXEC)

clean:
	@rm -rfv $(EXEC) $(COMP)

archive:
	@echo "archive to $(TAR)"
	@git archive master --prefix="$(EXEC)-$(VER)/" --format tar.gz -o $(TAR)
