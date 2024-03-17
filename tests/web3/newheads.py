import asyncio
from web3 import AsyncWeb3
from web3.providers import WebsocketProviderV2

async def ws_v2_subscription():
    async with AsyncWeb3.persistent_websocket(WebsocketProviderV2(f"ws://localhost:8546")) as w3:
        # subscribe to new block headers
        subscription_id = await w3.eth.subscribe("newHeads")

        async for response in w3.ws.process_subscriptions():
            print(f"{response}\n")
            # handle responses here

        # still an open connection, make any other requests and get
        # responses via send / receive
        latest_block = await w3.eth.get_block("latest")
        print(f"Latest block: {latest_block}")

asyncio.run(ws_v2_subscription())
