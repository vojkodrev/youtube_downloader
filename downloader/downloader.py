from abc import ABC, abstractmethod


class Downloader(ABC):
    @abstractmethod
    async def download(self, url: str) -> None: ...
