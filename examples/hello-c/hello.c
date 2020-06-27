#include <stdio.h>

int main(int argc, char *argv[]) {
    if (argc >= 2) {
        printf("Hello from C, %s!", argv[1]);
    } else {
        printf("Hello from C, please pass an argument to this command!");
    }
}
