#import ccxt

import time
import random
import requests
from requests import Request, Session

if __name__ == '__main__':
    
    """ huobi  = ccxt.huobipro()
    huobi.apiKey='7548cd5e-b95cad3b-33ebac0a-f809d111'
    huobi.secret='63a58654-0cf1b5fa-9ded4823-a9165111'
        
    while 1==1:
        tiem1=int(round(time.time() * 1000))
        counter=0
        depth_huobi=huobi.fetch_order_book("NAS/USDT",limit=7)

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
        #print("火币买一量:",hb_bids_1_amount,"火币卖一量;",hb_asks_1_amount)
        #print("-----------------------------------------------------------------------------------------")
        #print("火币买1:",hb_bids_1_price,"火币卖1:",hb_aks_1_price)
        #huobi_buy_order=huobi.create_market_buy_order("NAS/USDT",hb_aks_1_price*(1/0.998))
        # data_list=huobi.fetch_orders("NAS/USDT",1525500000000,100,None)
        # for item in data_list:
        #     print(item)
        #print(huobi.fetch_orders("NAS/USDT",1525500000000,100,None))
        total=int(round(time.time() * 1000))-tiem1
        print(total)
        #print(depth_huobi
        #print(huobi.fetch_free_balance()) """



    counter=0
    while True:
        
        rand=random.random()*10
        s=requests.Session()
        url="https://testnet.nebulas.io/claim/api/claim/"+str(rand)+"@qq.com/n1cS26LrgGEUazvaGHGW2Nr6AMPBQ52NoC3/"
        r1= s.get(url, verify=False)
        status=int(r1.status_code)
        if status==200:
            counter+=1
            print("已获取"+str(counter)+"个NAS")
            print(url+"-------"+str(r1._content))
        

