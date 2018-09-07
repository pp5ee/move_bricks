# -*- coding: utf-8 -*-
#author: 半熟的韭菜

from websocket import create_connection
import gzip
import time
import pandas as pd
import json
import talib as tb
import email_sender

kl1m_dict = {}
kl30m_dict={}

if __name__ == '__main__':


    
    while(1):
        try:
            ws = create_connection("wss://api.huobipro.com/ws")
            break
        except:
            print('connect ws error,retry...')
            time.sleep(5)

    # 订阅 KLine 数据
    tradeStr = """{"sub": "market.nasusdt.kline.1min","id": "id10"}"""
    tradeStr_1day="""{"sub":"market.nasbtc.kline.1day"}"""
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
    ws.send(tradeStr_1day)
    while(1):
        compressData = ws.recv()
        result = gzip.decompress(compressData).decode('utf-8')
        if result[:7] == '{"ping"':
            ts = result[8:21]
            pong = '{"pong":'+ts+'}'
            ws.send(pong)
            ws.send(tradeStr)
        else:
            dct = json.loads(result)
            #print(dct)
            #print(dct)

            if 'ch' in dct and 'market.nasusdt.kline.1min' in dct['ch']:
                #print("1分钟筛选",dct)
               
                id = dct['tick']['id']
                kl1m_dict[id] = dct['tick']
                df = pd.DataFrame.from_records(list(kl1m_dict.values()),'id')
                #print(df['close'])
                #MACD

                macd, macdsignal, macdhist = tb.MACD(df['close'], fastperiod=12, slowperiod=26, signalperiod=9)
                macd_arr=macd.tail(1).values
                macdsignal_arr=macdsignal.tail(1).values
                macdhist_arr=macdhist.tail(1).values

                macd_last=macd_arr[0]
                macdsignal_last=macdsignal_arr[0]
                macdhist_last=macdhist_arr[0]
                slowk, slowd = tb.STOCH(df['high'], df['low'], df['close'], fastk_period=5, slowk_period=3, slowk_matype=0, slowd_period=3, slowd_matype=0)
                slowk_float_value=(slowk.tail(1).values)[0]
                slowd_float_value=(slowd.tail(1).values)[0]
                print(slowk_float_value,slowd_float_value)
                print("K值，D值,J值:",slowk_float_value,slowd_float_value,(3*slowd_float_value-2*slowk_float_value))
                #如果参数加载完成执行以下操作
                if not 'nan' in str(macd_last) and not 'nan' in str(macdsignal_last) and not 'nan' in str(macdhist_last):
                    macd_float_value=float(macd_last)
                    macdhist_float_value=float(macdhist_last)
                    macdsignal_float_value=float(macdsignal_last)
                    #出现金叉和死叉时
                    macd_is_dead_or_live=(macd_float_value==macdhist_float_value==macdsignal_float_value)
                    slowk, slowd = tb.STOCH(df['high'], df['low'], df['close'], fastk_period=5, slowk_period=3, slowk_matype=0, slowd_period=3, slowd_matype=0)
                    slowk_float_value=(slowk.tail(1).values)[0]
                    slowd_float_value=(slowd.tail(1).values)[0]
                    kdj_is_dead_or_alive=slowd_float_value==slowk
                    print("K值，D值,J值:",slowk_float_value,slowd_float_value,(3*slowd_float_value-2*slowk_float_value))
                    if macd_is_dead_or_live and kdj_is_dead_or_alive:
                        subject="金叉死叉提醒"
                        kdj_str="slowd:"+str(slowd_float_value)+",slowk:"+str(slowk_float_value)
                        macd_str=",macd:"+str(macd_float_value)+",macdhist:"+str(macdhist_float_value)+",macdsignal:"+str(macdsignal_float_value)
                        content=kdj_str+macd_str
                        ms=email_sender.email_sender()
                        ret=ms.mail(content,subject)
                #print("MACD:",macd)
                #print("macd:",macd_arr[0])
                #print("macdsignal:",macdsignal.tail(1).values)
                #print("macdhist:",macdhist.tail(1).values)
                print('---------------------------------------------------------------------------------------------------------------------------------------------------------------------------')
                #print()
                #KDJ
               
        
                # clear dict
                if len(kl1m_dict) >= 200:
                    l = len(kl1m_dict)
                    keys = list(kl1m_dict.keys())[0:int(l/2)]
                    for key in keys:
                        del kl1m_dict[key]


