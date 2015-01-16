// gcc 4.7.2 +
// gcc -std=gnu99 -Wall -g -o helloworld_c helloworld_c.c -lpthread

#include <pthread.h>
#include <stdio.h>
#include <stdint.h>


int32_t i = 0;

// Note the return type: void*
void* thread1() {
    for (uint32_t j = 0; j < 1000000; j++) {
    	i++;
    }
    return NULL;
}

void* thread2() {
	for (uint32_t j = 0; j < 1000000; j++) {
    	i--;
    }
    return NULL;
}



int main(){
    pthread_t t1, t2;
    pthread_create(&t1, NULL, thread1, NULL);
    pthread_create(&t2, NULL, thread2, NULL);
    // Arguments to a thread would be passed here ---------^
    
    pthread_join(t1, NULL);
    pthread_join(t2, NULL);
    printf("i = %i!\n", i);
    return 0;
    
}
