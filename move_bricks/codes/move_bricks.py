import ccxt
import time
import _thread
import json
import email_sender
import threading
import sys
import datetime
gateio=ccxt.gateio()
huobi  = ccxt.huobipro()

#hitbtc_markets = hitbtc.load_markets()

#print(hitbtc.id, hitbtc_markets)
#print(bitmex.id, bitmex.load_markets())
#print(huobi.id, huobi.load_markets())

#print(hitbtc.fetch_order_book(hitbtc.symbols[0]))
#print(bitmex.fetch_ticker('BTC/USD'))
#print(huobi.fetch_trades('LTC/CNY'))
#print(huobi.fetch_l2_order_book("NAS/USDT"))
#print(hu
# obi.fetch_ticker("NAS/USDT"))

class move_bricks_thread(threading.Thread):

    def __init__(self, threadID, name,rate):
        threading.Thread.__init__(self)
        self.threadID1 = threadID
        self.name1 = name
        self.rate=rate
        self.counter=0
        self.times=0
        self.huobi_usdt_balance=0
        self.huobi_nas_balance=0
        self.gate_usdt_balance=0
        self.gate_nas_balance=0
    def run(self):
        print ("开始线程：" + self.name1)
        threading.Lock()
        self.move_nas_usdt_huobi_gateio(self.rate)
        print ("退出线程：" + self.name1)
#NAS和USDT的搬砖
    def move_nas_usdt_huobi_gateio(self,rate):
        gateio=ccxt.gateio()
        huobi  = ccxt.huobipro()
        huobi.apiKey='7548cd5e-b95cad3b-33ebac0a-f809d'
        huobi.secret='63a58654-0cf1b5fa-9ded4823-a9165'
        gateio.apiKey="1C465DD1-3657-4AEB-B6AA-62F558F397ED"
        gateio.secret=""
        while True:
            try:
                depth_huobi=huobi.fetch_order_book("NAS/USDT",limit=7)

                depth_gateio=gateio.fetch_order_book("NAS/USDT",limit=7)

                huobi_bids=list(depth_huobi['bids'])
                hb_bids_1=list(huobi_bids[0])
        #买一价
                hb_bids_1_price=float(hb_bids_1[0])
        #买一量
                hb_bids_1_amount=float(hb_bids_1[1])
        #获取卖一
                hb_huobi_asks=list(depth_huobi['asks'])
                hb_asks_1=list(hb_huobi_asks[0])
        #卖一价
                hb_aks_1_price=float(hb_asks_1[0])
        #卖一量
                hb_asks_1_amount=float(hb_asks_1[1])


                gateio_bids=list(depth_gateio['bids'])
                gate_bids_1=list(gateio_bids[0])
                gate_bids_7=list(gateio_bids[-1])
                gate_bids_7_price=float(gate_bids_7[0])
                gate_bids_1_price=float(gate_bids_1[0])
                gate_bids_1_amout=float(gate_bids_1[1])
                gateio_asks=list(depth_gateio['asks'])
                gate_asks_7=list(gateio_asks[-1])
                gate_asks_1=list(gateio_asks[0])
                gate_asks_1_price=float(gate_asks_1[0])
                gate_asks_7_price=float(gate_asks_7[0])
                gate_asks_1_amout=float(gate_asks_1[1])


                condition_1=float((hb_bids_1_price-gate_asks_1_price)/gate_asks_1_price)
                condition_2=float((gate_bids_1_price-hb_aks_1_price)/hb_aks_1_price)
                print("CONDITION_1：",condition_1,"CONDITION_2:",condition_2)


                if condition_1>=rate or condition_2>=rate:
                #获取火币余额
                    self. times+=1
                    huobi_all_balance=dict(huobi.fetch_balance()["info"])
                    huobi_data=dict(huobi_all_balance["data"])
                    huobi_list=list(huobi_data["list"])
                    for i,val in enumerate(huobi_list):
                        val=dict(val)
                    #print(val)
                        if val["type"]=="trade" and val["currency"]=="usdt":
                            self.huobi_usdt_balance=float(val["balance"])
                        if val["type"]=="trade" and val["currency"]=="nas":
                            self.huobi_nas_balance=float(val["balance"])
                #获取gate余额
                    gateio_all_balance=dict(dict(gateio.fetch_balance()["info"])["available"])
                    self.gate_nas_balance=float(gateio_all_balance["NAS"])
                    self.gate_usdt_balance=float(gateio_all_balance["USDT"])
                    print(datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'),"-----交易前----","gate_NAS余额:",self.gate_nas_balance,"gate_USDT 余额:",self.gate_usdt_balance,"火币_NAS余额:",self.huobi_nas_balance,"火币USDT余额:",self.huobi_usdt_balance,"交易前NAS总量:",(self.gate_nas_balance+self.huobi_nas_balance),"交易前USDT总和:",(self.huobi_usdt_balance+self.gate_usdt_balance))

                    #if condition_1>=rate and hb_asks_1_amount
                    #条件2
                    if condition_1>=rate and hb_asks_1_amount>=2 and gate_bids_1_amout>=2 and self.gate_usdt_balance>=gate_asks_1_price and self.huobi_nas_balance>=1:
                        print("火币卖,gate买")
                        huobi_sell_order=huobi.create_market_sell_order("NAS/USDT",1)
                        gate_buy_order=gateio.create_limit_buy_order("NAS/USDT",1/0.998,gate_asks_1_price*1.05)
                        print(datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'),"--交易订单:[","火币卖单-",huobi_sell_order,"],[gate买单-",gate_buy_order,"]")
                        #print(datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'),"-----交易后----","gateNAS余额:",self.gate_nas_balance,"gate USDT 余额:",self.gate_usdt_balance,"火币NAS余额:",self.huobi_nas_balance,"火币USDT余额:",self.huobi_usdt_balance,"交易后NAS总量:",self.gate_nas_balance+self.huobi_nas_balance,"交易后USDT总和:",self.huobi_usdt_balance+self.gate_usdt_balance)
                        self.counter+=1

                    else:
                        if self.gate_usdt_balance<gate_asks_1_price:
                            print("gate_usdt余额不足！")
                            subject="gate_usdt余额不足！"
                            content="gate_usdt余额不足！"
                            ms=email_sender.email_sender()
                            ret=ms.mail(content,subject,"1017498887@qq.com")
                        elif self.huobi_nas_balance<1:
                            print("火币_NAS余额不足！")
                            subject="火币_NAS余额不足！"
                            content="火币_NAS余额不足！"
                            ms=email_sender.email_sender()
                            ret=ms.mail(content,subject,"1017498887@qq.com")
                        elif hb_asks_1_amount<2 and gate_bids_1_amout<2:
                            print("火币卖gate买交易量不满足")

                    if condition_2 >=rate and hb_asks_1_amount>=2 and gate_bids_1_amout>=2 and self.huobi_usdt_balance>=hb_bids_1_price and self.gate_nas_balance>=1 :
                        print("火币买,gate卖")
                        huobi_buy_order=huobi.create_market_buy_order("NAS/USDT",(hb_bids_1_price)*(1.005/0.998))
                        gate_sell_order=gateio.create_limit_sell_order("NAS/USDT",1,gate_bids_7_price*0.95)
                        print(datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'),"--交易订单:[","火币买单-",huobi_buy_order,"],[gate卖单-",gate_sell_order,"]")
                        self.counter+=1
                        #print(datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'),"-----交易后----","gateNAS余额:",self.gate_nas_balance,"gate USDT 余额:",self.gate_usdt_balance,"火币NAS余额:",self.huobi_nas_balance,"火币USDT余额:",self.huobi_usdt_balance,"交易后NAS总量:",self.gate_nas_balance+self.huobi_nas_balance,"交易后USDT总和:",self.huobi_usdt_balance+self.gate_usdt_balance)
                    else:
                        if self.huobi_usdt_balance<hb_bids_1_price:
                            print("火币_usdt余额不足！")
                            subject="火币_usdt余额不足！"
                            content="火币_usdt余额不足！"
                            ms=email_sender.email_sender()
                            ret=ms.mail(content,subject,"1017498887@qq.com")
                        elif self.gate_nas_balance<1:
                            print("GATE_NAS余额不足！")
                            subject="GATE_NAS余额不足！"
                            content="GATE_NAS余额不足！"
                            ms=email_sender.email_sender()
                            ret=ms.mail(content,subject,"1017498887@qq.com")
                        elif hb_asks_1_amount<2 and gate_bids_1_amout<2:
                            print("火币买gate卖交易量不满足")
                #time.sleep(0.0005)
                print(datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'),"--已搬砖次数:",self.counter,"--检测到可搬砖次数:",self.times)
            except:
                pass
