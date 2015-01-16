# Python 3.3.3 and 2.7.6
# python helloworld_python.py

from threading import Thread

i = 0

def thread1():
    global i
    for _ in range(1000000):
    	i += 1

def thread2():
    global i
    for _ in range(1000000):
    	i -= 1

def main():
    t1 = Thread(target=thread1, args=())
    t2 = Thread(target=thread2, args=())

    t1.start()
    t2.start()
    
    t1.join()
    t2.join()
    print("i = %i!" % i)


main()
