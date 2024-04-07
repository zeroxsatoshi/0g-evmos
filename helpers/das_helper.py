#!/usr/bin/env python3

import asyncio
import json
import websockets
import pprint

async def on_register(event):
    pprint.pprint(event)

async def on_msg(msg):
    if 'result' in msg:
        result = msg['result']

        if not 'data' in result:
            return
        data = result['data']

        assert(data['type'] == 'tendermint/event/Tx')
        value = data['value']

        logs = json.loads(value['TxResult']['result']['log'])
        for log in logs:
            msg_index = log['msg_index']
            print('msg_index: ', msg_index)
            for event in log['events']:
                if event['type'] == 'register':
                    await on_register(event)
                else:
                    print('unkown event: ', event['type'])

async def main():
    uri = "ws://localhost:26657/websocket"
    async with websockets.connect(uri) as websocket:
        msg = {
            "jsonrpc": "2.0",
            "method": "subscribe",
            "id": 1,
            "params": {
                "query": "tm.event='Tx'"
            }
        }
        await websocket.send(json.dumps(msg))
        while True:
            try:
                raw_msg = await websocket.recv()
                msg = json.loads(raw_msg)
                await on_msg(msg)
            except websockets.exceptions.ConnectionClosedOK:
                print("WebSocket connection closed")
                break

asyncio.get_event_loop().run_until_complete(main())
