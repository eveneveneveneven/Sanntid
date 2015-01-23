# Python 3.3.3 and 2.7.6
# python oving2.py

from threading import Thread, Lock
from time import time

x, y = 0, 0
lock = Lock()

def fast_thread1():
    global x
    with lock:
        for _ in range(1000000):
            x += 1

def fast_thread2():
    global x
    with lock:
        for _ in range(1000000 - 1):
            x -= 1  
    
def slow_thread1():
    global y
    for _ in range(1000000):
        with lock:
            y += 1
        
def slow_thread2():
    global y
    for _ in range(1000000 - 1):
        with lock:
            y -= 1

def main():
    ft1 = Thread(target=fast_thread1, args=())
    ft2 = Thread(target=fast_thread2, args=())
    st1 = Thread(target=slow_thread1, args=())
    st2 = Thread(target=slow_thread2, args=())
    
    start = time()
    ft1.start()
    ft2.start()
    
    ft1.join()
    ft2.join()
    diff = time() - start
    print("fast x = %i, took %0.3f ms" % (x, diff*1000.0))
    
    start = time()
    st1.start()
    st2.start()
    
    st1.join()
    st2.join()
    diff = time() - start
    print("slow x = %i, took %0.3f ms" % (y, diff*1000.0))


main()
