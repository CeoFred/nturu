#!/bin/bash


BINARY_NAME="nturu"
INSTALL_PATH="/usr/local/bin"

build() {
    go build -o "$BINARY_NAME"
}

install() {
    build
    sudo cp "$BINARY_NAME" "$INSTALL_PATH"
}

test_generate() {
    install
    echo "Running tests..."
    ./"$BINARY_NAME" generate

    echo "Running tests... 2"
      ./"$BINARY_NAME" generate -framework fiber
    echo "Tests passed successfully."
}

uninstall() {
    sudo rm -f "$INSTALL_PATH/$BINARY_NAME"
}

clean() {
    rm -f "$BINARY_NAME"
}

case "$1" in
    build)
        build
        ;;
    install)
        install
        ;;
    test_generate)
        test_generate
        ;;
    uninstall)
        uninstall
        ;;
    clean)
        clean
        ;;
    *)
        echo "Usage: $0 {build|install|test_generate|uninstall|clean}"
        exit 1
esac

exit 0
