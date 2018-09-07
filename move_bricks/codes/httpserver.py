from aiohttp import web
from  huobi import Nas_USDT_Thread
from  huobi import BTC_USDT_Thread
from  huobi import Nas_BTC_Thread
from  huobi import Buy1_Thread
from  huobi import Period_Float_thread
from  huobi import EOS_USDT_Thread
from  huobi import EOS_BTC_Thread
from  huobi import date_to_timestamp
from  move_bricks import move_bricks_thread
import json
import _thread
#BTC/USDT
async def handle1(request):
    high = request.match_info.get('high')
    low =request.match_info.get('low')
    mail=request.match_info.get('mail')
    price_high = float(high)
    price_low= float(low)
    thread=BTC_USDT_Thread(1,"btc/usdt",price_high,price_low,mail)
    thread.start()
    text1='设定的最高价:'+str(price_high)+'设定的最低价'+str(price_low)
    return web.Response(text=text1)
#NAS/USDT
async def handle2(request): 
    high = request.match_info.get('high')
    low =request.match_info.get('low')
    price_high = float(high)
    price_low= float(low)
    mail=request.match_info.get('mail')
    thread=Nas_USDT_Thread(2,"NAS/USDT",price_high,price_low,mail)
    thread.start()
    text1='设定的最高价:'+str(price_high)+'设定的最低价'+str(price_low)
    return web.Response(text=text1)
#NAS/BTC
async def handle3(request):
    high = request.match_info.get('high')
    low =request.match_info.get('low')
    price_high = float(high)
    mail=request.match_info.get('mail')
    price_low= float(low)
    thread=Nas_BTC_Thread(3,"NAS/BTC",price_high,price_low,mail)
    thread.start()
    text1='设定的最高价:'+str(price_high)+'设定的最低价'+str(price_low)

    return web.Response(text=text1,headers=["Access-Control-Allow-Origin:*"])

async def handle4(request):
    #thread=Buy1_Thread(4,"Buy1")
    #thread.start()
    return web.Response()
    #价格区间波动
async def handle5(request):
    start=request.match_info.get('start')
    end=request.match_info.get('end')
    coin=request.match_info.get('name')
    
    start_stamp=date_to_timestamp(start)
    end_stamp=date_to_timestamp(end)
    print(start,end,start_stamp,end_stamp)
    if(start_stamp>end_stamp):
            return web.Response(text="请设置正确的时间！")
    threadN=Period_Float_thread(5,"PFT",start_stamp,end_stamp,coin)
    threadN.start()
    text1='设定的最高价:'+'设定的最低价'
    return web.Response(text=text1)
#EOS/USDT
async def handle6(request):
    high = request.match_info.get('high')
    low =request.match_info.get('low')
    price_high = float(high)
    mail=request.match_info.get('mail')
    price_low= float(low)
    threadN=EOS_USDT_Thread(6,"EU",price_high,price_low,mail)
    threadN.start()
    text1='设定的最高价:'+str(price_high)+'设定的最低价'+str(price_low)
    return web.Response(text=text1)
#EOS/BTC
async def handle7(request):
    high = request.match_info.get('high')
    low =request.match_info.get('low')
    price_high = float(high)
    mail=request.match_info.get('mail')
    price_low= float(low)

    threadN=EOS_BTC_Thread(7,"EB",price_high,price_low,mail)
    threadN.start()
    text1='设定的最高价:'+str(price_high)+'设定的最低价'+str(price_low)
    return web.Response(text=text1)
async def handle8(request):
    rate = float(request.match_info.get('rate'))
    #mail=request.match_info.get('mail')
    thread_move=move_bricks_thread(8,"move",rate)
    thread_move.start()
    text1='启动成功'
    return web.Response(text=text1)
async def handle9( request):
    print(request)
    return "TEST TEST"
async def pshandle(request):
    price=request.match_info.post("price")
    return price
async def wshandle(request):
    ws = web.WebSocketResponse()
    await ws.prepare(request)

    async for msg in ws:
        if msg.type == web.WSMsgType.text:
            await ws.send_str("Hello, {}".format(msg.data))
        elif msg.type == web.WSMsgType.binary:
            await ws.send_bytes(msg.data)
        elif msg.type == web.WSMsgType.close:
            break

    return ws


app = web.Application()
""" app.router.add_static('/page',
                       path='/static',
                       name='index') """
app.add_routes([
                web.get('/BU/{high}/{low}/{mail}', handle1),
                web.get('/NU/{high}/{low}/{mail}',handle2),
                web.get('/NB/{high}/{low}/{mail}',handle3),
                web.get('/',handle4),
                web.get('/CAL/{start}/{end}/{name}',handle5),
                web.get('/EU/{high}/{low}/{mail}',handle6),
                web.get('/EB/{high}/{low}/{mail}',handle7),
                web.get('/move/{rate}',handle8),
                web.post('/test',handle9)
                ])
web.run_app(app, host='0.0.0.0', port=3999)