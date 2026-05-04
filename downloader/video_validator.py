from abc import ABC, abstractmethod


class VideoValidator(ABC):
    @abstractmethod
    async def validate(self, video_id: str, mode: str, streamer: str = "") -> list[str]: ...
