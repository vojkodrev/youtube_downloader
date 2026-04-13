import asyncio

from injector import inject

from channel_poller import ChannelPoller
from config import Config


class MultiChannelPoller:
    @inject
    def __init__(self, poller: ChannelPoller, config: Config):
        self._poller = poller
        self._config = config

    async def poll_all(self) -> None:
        await asyncio.gather(
            *[
                self._poller.poll(ch["id"], mode=ch["mode"])
                for ch in self._config.get("channels") or []
            ]
        )
