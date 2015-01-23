// gcc 4.7.2 +
// gcc -std=gnu99 -Wall -g -o oving2 oving2.c -lpthread

#include <pthread.h>
#include <stdio.h>
#include <stdint.h>
#include <time.h>

pthread_mutex_t lock;
int32_t x = 0, y = 0;

void* fast_thread1()
{
    pthread_mutex_lock(&lock);

    for (uint32_t j = 0; j < 1000000; j++) {
    	x++;
    }
    
    pthread_mutex_unlock(&lock);

    return NULL;
}

void* fast_thread2()
{
    pthread_mutex_lock(&lock);

    for (uint32_t j = 0; j < 1000000 - 1; j++) {
    	x--;
    }
    
    pthread_mutex_unlock(&lock);

    return NULL;
}

void* slow_thread1()
{
    for (uint32_t j = 0; j < 1000000; j++) {
        pthread_mutex_lock(&lock);
    	y++;
        pthread_mutex_unlock(&lock);
    }

    return NULL;
}

void* slow_thread2()
{
    for (uint32_t j = 0; j < 1000000 - 1; j++) {
        pthread_mutex_lock(&lock);
    	y--;
        pthread_mutex_unlock(&lock);
    }

    return NULL;
}



int main()
{
    clock_t start, diff;
    int msec;
    pthread_t ft1, ft2, st1, st2;
    pthread_mutex_init(&lock, NULL);
    
    // Fast locking
    start = clock();
    pthread_create(&ft1, NULL, fast_thread1, NULL);
    pthread_create(&ft2, NULL, fast_thread2, NULL);
    
    pthread_join(ft1, NULL);
    pthread_join(ft2, NULL);
    diff = clock() - start;
    msec = diff * 1000 / CLOCKS_PER_SEC;
    printf("fast x = %d, took %d milliseconds!\n", x, msec);
    
    // Slow locking
    start = clock();
    pthread_create(&st1, NULL, slow_thread1, NULL);
    pthread_create(&st2, NULL, slow_thread2, NULL);
    
    pthread_join(st1, NULL);
    pthread_join(st2, NULL);
    diff = clock() - start;
    msec = diff * 1000 / CLOCKS_PER_SEC;
    printf("slow y = %d, took %d milliseconds!\n", y, msec);
    
    pthread_mutex_destroy(&lock);
    return 0;
    
}
