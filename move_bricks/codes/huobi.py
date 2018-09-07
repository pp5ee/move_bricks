# -*- coding: utf-8 -*-
#author: 半熟的韭菜

from websocket import create_connection
import gzip
import time
import _thread
import json
import email_sender 
import threading
import sys
class Nas_USDT_Thread(threading.Thread):
    def __init__(self, threadID, name, price_high,price_low,mail):
        threading.Thread.__init__(self)
        self.threadID1 = threadID
        self.name1 = name
        self.price_high1 = price_high
        self.price_low1=price_low
        self.mail=mail
    def run(self):
        print ("开始线程：" + self.name1)
        threading.Lock()
        self.get_realtime_kline_nas_usdt1(self.price_high1,self.price_low1,self.mail)

        print ("退出线程：" + self.name1)
    
    def get_realtime_kline_nas_usdt1(self, price_high, price_low,mail):
        print("Nas-BTC")
        while(1):
            try:
                ws = create_connection("wss://api.huobipro.com/ws")
                break
            except:
                print('connect ws error,retry...')

            time.sleep(5)

        # 订阅 KLine 数据
        tradeStr = """{"sub":"market.nasusdt.trade.detail"}"""

        # 请求 KLine 数据
        # tradeStr="""{"req": "market.ethusdt.kline.1min","id": "id10", "from": 1513391453, "to": 1513392453}"""

        # 订阅 Market Depth 数据
        # tradeStr="""{"sub": "market.ethusdt.depth.step5", "id": "id10"}"""

        # 请求 Market Depth 数据
        # tradeStr="""{"req": "market.ethusdt.depth.step5", "id": "id10"}"""

        # 订阅 Trade Detail 数据
        # tradeStr="""{"sub": "market.ethusdt.trade.detail", "id": "id10"}"""

        # 请求 Trade Detail 数据
        # tradeStr="""{"req": "market.ethusdt.trade.detail", "id": "id10"}"""

        # 请求 Market Detail 数据
        # tradeStr="""{"req": "market.ethusdt.detail", "id": "id12"}"""

        ws.send(tradeStr)
        #_thread.start_new_thread(self.recieve(ws,tradeStr), ("Thread-1", 2, ))
        ret=None
        count=0
        while(count==0):
            # def recieve(self,ws,tradeStr):
            compressData = ws.recv()
            result = str(gzip.decompress(compressData).decode('utf-8'))
            if result[:7] == '{"ping"':
                ts = result[8:21]
                pong = '{"pong":'+ts+'}'
                ws.send(pong)
                ws.send(tradeStr)
            else:
                # judge=json.load(result)
                result_json = json.loads(result)
                # print(result_json)
                try:
                    sub_json = dict(result_json["tick"])
                    data = sub_json['data']
                    cur_price = float(data[0]['price'])
                    print(cur_price,price_high,price_low)
                    if cur_price >= price_high or cur_price <= price_low:
                        subject="NAS/USDT"
                        cont=cur_price
                        ms=email_sender.email_sender()
                        ret=ms.mail(cont,subject,mail)
                        print(ret)
                    if ret  :
                        #thread=threading.current_thread()
                        count=1
                        #my_sender.__init__(cont)
                        #my_sender.mail(None,cont)
                        #email_sender.mail()
                except Exception as e:
                    print(e)





#NAS/BTC
class Nas_BTC_Thread(threading.Thread):
    def __init__(self, threadID, name, price_high,price_low,mail):
        threading.Thread.__init__(self)
        self.threadID2 = threadID
        self.name2 = name
        self.price_high2 = price_high
        self.price_low2=price_low
        self.mail=mail
    def run(self):
        print ("开始线程：" + self.name2)
        self.get_realtime_kline_nas_btc2(self.price_high2,self.price_low2,self.mail)
        print ("退出线程：" + self.name2)
    def get_realtime_kline_nas_btc2(self, price_high, price_low,mail):
        print("Nas-BTC")
        while(1):
            try:
                ws = create_connection("wss://api.huobipro.com/ws")
                break
            except:
                print('connect ws error,retry...')

            time.sleep(5)

        # 订阅 KLine 数据
        tradeStr = """{"sub":"market.nasbtc.trade.detail"}"""

        # 请求 KLine 数据
        # tradeStr="""{"req": "market.ethusdt.kline.1min","id": "id10", "from": 1513391453, "to": 1513392453}"""

        # 订阅 Market Depth 数据
        # tradeStr="""{"sub": "market.ethusdt.depth.step5", "id": "id10"}"""

        # 请求 Market Depth 数据
        # tradeStr="""{"req": "market.ethusdt.depth.step5", "id": "id10"}"""

        # 订阅 Trade Detail 数据
        # tradeStr="""{"sub": "market.ethusdt.trade.detail", "id": "id10"}"""

        # 请求 Trade Detail 数据
        # tradeStr="""{"req": "market.ethusdt.trade.detail", "id": "id10"}"""

        # 请求 Market Detail 数据
        # tradeStr="""{"req": "market.ethusdt.detail", "id": "id12"}"""

        ws.send(tradeStr)
        #_thread.start_new_thread(self.recieve(ws,tradeStr), ("Thread-1", 2, ))
        count=0
        while(count==0):

            # def recieve(self,ws,tradeStr):
            compressData = ws.recv()
            result = str(gzip.decompress(compressData).decode('utf-8'))
            if result[:7] == '{"ping"':
                ts = result[8:21]
                pong = '{"pong":'+ts+'}'
                ws.send(pong)
                ws.send(tradeStr)
            else:
                # judge=json.load(result)
                result_json = json.loads(result)
                # print(result_json)
                try:
                    sub_json = dict(result_json["tick"])
                    data = sub_json['data']
                    cur_price = float(data[0]['price'])
                    print(cur_price,price_high,price_low)
                    if cur_price >= price_high or cur_price <= price_low:
                        print("发送邮件")
                        subject="NAS/BTC"
                        content=cur_price
                        ms=email_sender.email_sender()
                        ret=ms.mail(content,subject,mail)
                        print(ret)
                        if ret :
                            thread=threading.current_thread()
                            count=1
                            #thread.join()
                            sys.exit() 
                        #email_sender.mail()
                except:
                    pass
    


    
#BTC/USDT
class BTC_USDT_Thread(threading.Thread):
    def __init__(self, threadID, name, price_high,price_low,mail):
        threading.Thread.__init__(self)
        self.threadID3 = threadID
        self.name3 = name
        self.price_high3 = price_high
        self.price_low3=price_low
        self.mail=mail
    def run(self):
        print ("开始线程：" + self.name3)
        self.get_realtime_kline_btc_usdt(self.price_high3,self.price_low3,self.mail)
        print ("退出线程：" + self.name3)


    def get_realtime_kline_btc_usdt(self, price_high, price_low,mail):
        print("BTC-USDT")
        while(1):
            try:
                ws = create_connection("wss://api.huobipro.com/ws")
                break
            except:
                print('connect ws error,retry...')

            time.sleep(5)

        # 订阅 KLine 数据
        tradeStr = """{"sub":"market.btcusdt.trade.detail"}"""

        # 请求 KLine 数据
        # tradeStr="""{"req": "market.ethusdt.kline.1min","id": "id10", "from": 1513391453, "to": 1513392453}"""

        # 订阅 Market Depth 数据
        # tradeStr="""{"sub": "market.ethusdt.depth.step5", "id": "id10"}"""

        # 请求 Market Depth 数据
        # tradeStr="""{"req": "market.ethusdt.depth.step5", "id": "id10"}"""

        # 订阅 Trade Detail 数据
        # tradeStr="""{"sub": "market.ethusdt.trade.detail", "id": "id10"}"""

        # 请求 Trade Detail 数据
        # tradeStr="""{"req": "market.ethusdt.trade.detail", "id": "id10"}"""

        # 请求 Market Detail 数据
        # tradeStr="""{"req": "market.ethusdt.detail", "id": "id12"}"""

        ws.send(tradeStr)
        #_thread.start_new_thread(self.recieve(ws,tradeStr), ("Thread-1", 2, ))
        count=0
        while(count==0):

            # def recieve(self,ws,tradeStr):
            compressData = ws.recv()
            result = str(gzip.decompress(compressData).decode('utf-8'))
            if result[:7] == '{"ping"':
                ts = result[8:21]
                pong = '{"pong":'+ts+'}'
                ws.send(pong)
                ws.send(tradeStr)
            else:
                # judge=json.load(result)
                result_json = json.loads(result)
                # print(result_json)
                try:
                    sub_json = dict(result_json["tick"])
                    data = sub_json['data']
                    cur_price = float(data[0]['price'])
                    print(cur_price,price_high,price_low)
                    if cur_price >= price_high or cur_price <= price_low:
                        print("发送邮件")
                        subject="BTC/USDT"
                        content=cur_price
                        ms=email_sender.email_sender()
                        ret=ms.mail(content,subject,mail)
                        print(ret)
                        if ret :
                            #thread=threading.current_thread()
                            count=1
                            #thread.join()
                            #sys.exit()
                        #email_sender.mail()
                except:
                    pass



#一段时间的涨跌                  
class Period_Float_thread(threading.Thread):
    def __init__(self, threadID, name,start,end,type1):
        threading.Thread.__init__(self)
        self.threadID3 = threadID
        self.name3 = name
        self.c=None
        self.a=None
        self.b=None
        self.type=type1
        self.start_t=start
        self.end_t=end
    def run(self):
        print ("开始线程：" + self.name3)
        self.get_realtime_trade()
        print ("退出线程：" + self.name3)


    def get_realtime_trade(self):
        while(1):
            try:
                ws = create_connection("wss://api.huobipro.com/ws")
                break
            except:
                print('connect ws error,retry...')
            time.sleep(5)
        type_str=str(self.type)
        #print(type_str)
        sub_trade=None
        if type_str=="NAS-BTC":
            sub_trade="nasbtc"
        elif  type_str=="NAS-USDT":
            sub_trade="nasusdt"
        elif type_str=="BTC-USDT":
            sub_trade="btcusdt"

        # 订阅 KLine 历史数据
        #tradeStr = """{"req":"market."++".trade.detail"}"""
        # 请求 KLine 数据
        {"id":"1524645873543.1.0","req":"market.eosbtc.kline.1min","from":1524549213,"to":1524554635}
        tradeStr="""{"req": "market."""+sub_trade+""".kline.1min","id": "id10", "from":"""+str(self.start_t)+""", "to":"""+str(self.end_t)+"""}"""
        print(tradeStr)
        # 订阅 Market Depth 数据
        # tradeStr="""{"sub": "market.ethusdt.depth.step5", "id": "id10"}"""

        # 请求 Market Depth 数据
        # tradeStr="""{"req": "market.ethusdt.depth.step5", "id": "id10"}"""

        # 订阅 Trade Detail 数据
        # tradeStr="""{"sub": "market.ethusdt.trade.detail", "id": "id10"}"""

        # 请求 Trade Detail 数据
        # tradeStr="""{"req": "market.ethusdt.trade.detail", "id": "id10"}"""

        # 请求 Market Detail 数据
        # tradeStr="""{"req": "market.ethusdt.detail", "id": "id12"}"""

        ws.send(tradeStr)
        #_thread.start_new_thread(self.recieve(ws,tradeStr), ("Thread-1", 2, ))
        
        while(1):

            # def recieve(self,ws,tradeStr):
            compressData = ws.recv()
            result = str(gzip.decompress(compressData).decode('utf-8'))
            if result[:7] == '{"ping"':
                ts = result[8:21]
                pong = '{"pong":'+ts+'}'
                ws.send(pong)
                ws.send(tradeStr)
            else:
                # judge=json.load(result)
                result_json = json.loads(result)
                print(result_json)
                try:
                    sub = result_json['tick']
                    ch =result_json['ch']
                    # print(ch)
                    #print(sub['bids'])
                    bids=list(sub['bids'])
                    asks=list(sub['asks'])
                    bids_one=list(bids[0])#买一
                    asks_one=list(asks[0])#卖一
                    if 'nasbtc' in ch:
                        self.a=float(bids_one[0])
                    elif 'btcusdt' in ch:
                        self.c=float(bids_one[0])
                    elif 'nasusdt' in ch:
                        self.b=float(asks_one[0])
                    #print("买1:"+str(self.c))
                    asw=(self.b-self.a*self.c)/self.b
                    print(asw)
                except:
                    pass
class Buy1_Thread(threading.Thread):
    def __init__(self, threadID, name):
        threading.Thread.__init__(self)
        self.threadID3 = threadID
        self.name3 = name
        self.c=None
        self.a=None
        self.b=None
    def run(self):
        print ("开始线程：" + self.name3)
        self.get_realtime_kline_btc_usdt()
        print ("退出线程：" + self.name3)


    def get_realtime_kline_btc_usdt(self):
        while(1):
            try:
                ws = create_connection("wss://api.huobipro.com/ws")
                break
            except:
                print('connect ws error,retry...')

            time.sleep(5)

        # 订阅 KLine 数据
        tradeStr = """{"sub":"market.btcusdt.depth.step0","pick":["bids.7","asks.7"]}"""
        tradeStr2 = """{"sub":"market.nasusdt.depth.step0","pick":["bids.7","asks.7"]}"""
        tradeStr3="""{"sub":"market.nasbtc.depth.step0","pick":["bids.7","asks.7"]}"""
        # 请求 KLine 数据
        # tradeStr="""{"req": "market.ethusdt.kline.1min","id": "id10", "from": 1513391453, "to": 1513392453}"""

        # 订阅 Market Depth 数据
        # tradeStr="""{"sub": "market.ethusdt.depth.step5", "id": "id10"}"""

        # 请求 Market Depth 数据
        # tradeStr="""{"req": "market.ethusdt.depth.step5", "id": "id10"}"""

        # 订阅 Trade Detail 数据
        # tradeStr="""{"sub": "market.ethusdt.trade.detail", "id": "id10"}"""

        # 请求 Trade Detail 数据
        # tradeStr="""{"req": "market.ethusdt.trade.detail", "id": "id10"}"""

        # 请求 Market Detail 数据
        # tradeStr="""{"req": "market.ethusdt.detail", "id": "id12"}"""

        ws.send(tradeStr)
        ws.send(tradeStr2)
        ws.send(tradeStr3)
        #_thread.start_new_thread(self.recieve(ws,tradeStr), ("Thread-1", 2, ))
        
        while(1):

            # def recieve(self,ws,tradeStr):
            compressData = ws.recv()
            result = str(gzip.decompress(compressData).decode('utf-8'))
            if result[:7] == '{"ping"':
                ts = result[8:21]
                pong = '{"pong":'+ts+'}'
                ws.send(pong)
                ws.send(tradeStr)
            else:
                # judge=json.load(result)
                result_json = json.loads(result)
                #print(result_json)
                try:
                    sub = result_json['tick']
                    ch =result_json['ch']
                    # print(ch)
                    #print(sub['bids'])
                    bids=list(sub['bids'])
                    asks=list(sub['asks'])
                    bids_one=list(bids[0])#买一
                    asks_one=list(asks[0])#卖一
                    if 'nasbtc' in ch:
                        self.a=float(bids_one[0])
                    elif 'btcusdt' in ch:
                        self.c=float(bids_one[0])
                    elif 'nasusdt' in ch:
                        self.b=float(asks_one[0])
                    #print("买1:"+str(self.c))
                    asw=(self.b-self.a*self.c)/self.b
                    print(asw)
                except:
                    pass
#eos/usdt
class EOS_USDT_Thread(threading.Thread):
    def __init__(self, threadID, name, price_high,price_low,mail):
        threading.Thread.__init__(self)
        self.threadID3 = threadID
        self.name3 = name
        self.price_high3 = price_high
        self.price_low3=price_low
        self.mail=mail
    def run(self):
        print ("开始线程：" + self.name3)
        self.get_realtime_kline_eos_usdt(self.price_high3,self.price_low3,self.mail)
        print ("退出线程：" + self.name3)


    def get_realtime_kline_eos_usdt(self, price_high, price_low,mail):
        print("eos/usdt")
        while(1):
            try:
                ws = create_connection("wss://api.huobipro.com/ws")
                break
            except:
                print('connect ws error,retry...')

            time.sleep(5)

        # 订阅 KLine 数据
        tradeStr = """{"sub":"market.eosusdt.trade.detail"}"""

        # 请求 KLine 数据
        # tradeStr="""{"req": "market.ethusdt.kline.1min","id": "id10", "from": 1513391453, "to": 1513392453}"""

        # 订阅 Market Depth 数据
        # tradeStr="""{"sub": "market.ethusdt.depth.step5", "id": "id10"}"""

        # 请求 Market Depth 数据
        # tradeStr="""{"req": "market.ethusdt.depth.step5", "id": "id10"}"""

        # 订阅 Trade Detail 数据
        # tradeStr="""{"sub": "market.ethusdt.trade.detail", "id": "id10"}"""

        # 请求 Trade Detail 数据
        # tradeStr="""{"req": "market.ethusdt.trade.detail", "id": "id10"}"""

        # 请求 Market Detail 数据
        # tradeStr="""{"req": "market.ethusdt.detail", "id": "id12"}"""

        ws.send(tradeStr)
        #_thread.start_new_thread(self.recieve(ws,tradeStr), ("Thread-1", 2, ))
        count=0
        while(count==0):

            # def recieve(self,ws,tradeStr):
            compressData = ws.recv()
            result = str(gzip.decompress(compressData).decode('utf-8'))
            print("EU:",result)
            if result[:7] == '{"ping"':
                ts = result[8:21]
                pong = '{"pong":'+ts+'}'
                ws.send(pong)
                ws.send(tradeStr)
            else:
                # judge=json.load(result)
                result_json = json.loads(result)
                # print(result_json)
                try:
                    sub_json = dict(result_json["tick"])
                    data = sub_json['data']
                    cur_price = float(data[0]['price'])
                    print(cur_price,price_high,price_low)
                    if cur_price >= price_high or cur_price <= price_low:
                        print("发送邮件")
                        subject="EOS/USDT"
                        content=cur_price
                        ms=email_sender.email_sender()
                        ret=ms.mail(content,subject,mail)
                        print(ret)
                        if ret :
                            #thread=threading.current_thread()
                            count=1
                            #thread.join()
                            #sys.exit()
                        #email_sender.mail()
                except:
                    pass

class EOS_BTC_Thread(threading.Thread):
    def __init__(self, threadID, name, price_high,price_low,mail):
        threading.Thread.__init__(self)
        self.threadID3 = threadID
        self.name3 = name
        self.price_high3 = price_high
        self.price_low3=price_low
        self.mail=mail
    def run(self):
        print ("开始线程：" + self.name3)
        self.get_realtime_kline_eos_btc(self.price_high3,self.price_low3,self.mail)
        print ("退出线程：" + self.name3)


    def get_realtime_kline_eos_btc(self, price_high, price_low,mail):
        print("EOS/BTC")
        while(1):
            try:
                ws = create_connection("wss://api.huobipro.com/ws")
                break
            except:
                print('connect ws error,retry...')

            time.sleep(5)

        # 订阅 KLine 数据
        tradeStr = """{"sub":"market.eosbtc.trade.detail"}"""

        # 请求 KLine 数据
        # tradeStr="""{"req": "market.ethusdt.kline.1min","id": "id10", "from": 1513391453, "to": 1513392453}"""

        # 订阅 Market Depth 数据
        # tradeStr="""{"sub": "market.ethusdt.depth.step5", "id": "id10"}"""

        # 请求 Market Depth 数据
        # tradeStr="""{"req": "market.ethusdt.depth.step5", "id": "id10"}"""

        # 订阅 Trade Detail 数据
        # tradeStr="""{"sub": "market.ethusdt.trade.detail", "id": "id10"}"""

        # 请求 Trade Detail 数据
        # tradeStr="""{"req": "market.ethusdt.trade.detail", "id": "id10"}"""

        # 请求 Market Detail 数据
        # tradeStr="""{"req": "market.ethusdt.detail", "id": "id12"}"""

        ws.send(tradeStr)
        #_thread.start_new_thread(self.recieve(ws,tradeStr), ("Thread-1", 2, ))
        count=0
        while(count==0):

            # def recieve(self,ws,tradeStr):
            compressData = ws.recv()
            result = str(gzip.decompress(compressData).decode('utf-8'))
            print("EB:",result)
            if result[:7] == '{"ping"':
                ts = result[8:21]
                pong = '{"pong":'+ts+'}'
                ws.send(pong)
                ws.send(tradeStr)
            else:
                # judge=json.load(result)
                result_json = json.loads(result)
                # print(result_json)
                try:
                    sub_json = dict(result_json["tick"])
                    data = sub_json['data']
                    cur_price = float(data[0]['price'])
                    print(cur_price,price_high,price_low)
                    if cur_price >= price_high or cur_price <= price_low:
                        print("发送邮件")
                        subject="EOS/BTC"
                        content=cur_price
                        ms=email_sender.email_sender()
                        ret=ms.mail(content,subject,mail)
                        print(ret)
                        if ret :
                            #thread=threading.current_thread()
                            count=1
                            #thread.join()
                            #sys.exit()
                        #email_sender.mail()
                except:
                    pass
def date_to_timestamp(date, format_string="%Y-%m-%dT%H:%M:%S"):
    time_array = time.strptime(date, format_string)
    time_stamp = int(time.mktime(time_array))
    return time_stamp