
#include <stdio.h>

#define INITIAL_DIVIDER 2

int main() {
  int target = 0;

  if (scanf("%d", &target) < 1 || target <= 1) {
    fprintf(stderr, "Invalid number, should greater 1\n");
    return 1;
  }

  printf("%d:", target);

  int divider = INITIAL_DIVIDER;
  while(target > 1) {

    while (target % divider == 0) {
      printf(" %d", divider);
      target /= divider;
    }

    divider++;
  }
  printf("\n");
  return 0;
}
