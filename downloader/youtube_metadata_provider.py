from __future__ import annotations

import asyncio

from googleapiclient.discovery import build

from channel_title_not_found_error import ChannelTitleNotFoundError
from metadata_provider import MetadataProvider
from youtube_api_key_pool import YoutubeApiKeyPool


class YoutubeMetadataProvider(MetadataProvider):
    def __init__(self, api_keys: YoutubeApiKeyPool):
        self._api_keys = api_keys

    async def get_video_id(self, channel_id: str) -> str | None:
        if not channel_id:
            raise ValueError("channel_id is required")

        def get_video_id_sync():
            youtube = build("youtube", "v3", developerKey=self._api_keys.next())
            request = youtube.search().list(
                part="snippet",
                channelId=channel_id,
                eventType="live",
                type="video",
                maxResults=1,
            )
            response = request.execute()
            if response.get("items"):
                return response["items"][0]["id"].get("videoId")
            return None

        loop = asyncio.get_running_loop()
        return await loop.run_in_executor(None, get_video_id_sync)

    async def get_channel_title(self, channel_id: str) -> str:
        if not channel_id:
            raise ValueError("channel_id is required")

        def get_channel_title_sync():
            youtube = build("youtube", "v3", developerKey=self._api_keys.next())
            response = youtube.channels().list(part="snippet", id=channel_id).execute()
            if not response.get("items"):
                raise ChannelTitleNotFoundError(
                    f"Could not find channel with ID '{channel_id}'"
                )
            return response["items"][0]["snippet"]["title"]

        loop = asyncio.get_running_loop()
        return await loop.run_in_executor(None, get_channel_title_sync)

    def get_video_url(self, video_id: str) -> str:
        if not video_id:
            raise ValueError("video_id is required")
        return f"https://www.youtube.com/watch?v={video_id}"
