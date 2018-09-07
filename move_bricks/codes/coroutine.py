#-*- coding=utf-8 -*-
import requests
from multiprocessing import Process
import gevent
from gevent import monkey; monkey.patch_all()
import time
import sys
import ccxt
sys.setrecursionlimit(1000000)
def fetch(url):
    try:
        print(url)
        s = requests.Session()
        r = s.get(url,timeout=1)#在这里抓取页面
        print(r.json())
    except Exception as e:
        print(e) 
    return ''
 
def process_start(url_list):
    #print(url_list)
    tasks = []
    for url in url_list:
        tasks.append(gevent.spawn(fetch,url))
    gevent.joinall(tasks)#协程执行

def task_start(filepath,flag = 4000):#每4000条url启动一个进程
    with open(filepath,'r') as reader:#从给定的文件中读取url
        url = reader.readline().strip()
        url_list = []#这个list用于存放协程任务
        i = 0 
        while url!='':
            i += 1
            url_list.append(url)
            if i == flag:
                p = Process(target=process_start,args=(url_list,))
                p.start()
                url_list = [] #重置url队列
                i = 0 #重置计数器
            url = reader.readline().strip()
            #print(url)
        if url_list is not []:
            p = Process(target=process_start,args=(url_list,))
            p.start()
  
if __name__ == '__main__':
    tiem1=int(round(time.time() * 1000))
    task_start('./url.txt')
    total=int(round(time.time() * 1000))-tiem1
    print(total)
