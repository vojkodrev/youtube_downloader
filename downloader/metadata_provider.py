from abc import ABC, abstractmethod


class MetadataProvider(ABC):
    @abstractmethod
    async def get_video_id(self, channel_id: str) -> str | None: ...

    @abstractmethod
    async def get_channel_title(self, channel_id: str) -> str: ...

    @abstractmethod
    def get_video_url(self, video_id: str) -> str: ...

    @abstractmethod
    async def get_playlist_videos(self, playlist_id: str) -> list[dict]: ...
