from __future__ import annotations

from metadata_provider import MetadataProvider


class TwitchMetadataProvider(MetadataProvider):
    async def get_video_id(self, channel_id: str) -> str | None:
        return channel_id

    async def get_channel_title(self, channel_id: str) -> str:
        return channel_id

    def get_video_url(self, video_id: str) -> str:
        return f"https://www.twitch.tv/{video_id}"
