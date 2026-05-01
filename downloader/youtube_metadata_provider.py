from __future__ import annotations

import asyncio

import isodate

from googleapiclient.discovery import build
from injector import inject

from channel_title_not_found_error import ChannelTitleNotFoundError
from video_metadata import VideoMetadata
from metadata_provider import MetadataProvider
from youtube_api_key_pool import YoutubeApiKeyPool


class YoutubeMetadataProvider(MetadataProvider):
    @inject
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

    async def get_playlist_videos(self, playlist_id: str) -> list[dict]:
        if not playlist_id:
            raise ValueError("playlist_id is required")

        def get_playlist_videos_sync():
            youtube = build("youtube", "v3", developerKey=self._api_keys.next())
            videos = []
            next_page_token = None

            while True:
                response = youtube.playlistItems().list(
                    part="snippet",
                    playlistId=playlist_id,
                    maxResults=50,
                    pageToken=next_page_token,
                ).execute()

                for item in response.get("items", []):
                    videos.append({
                        "video_id": item["snippet"]["resourceId"]["videoId"],
                        "title": item["snippet"]["title"],
                    })

                next_page_token = response.get("nextPageToken")
                if not next_page_token:
                    break

            return videos

        loop = asyncio.get_running_loop()
        return await loop.run_in_executor(None, get_playlist_videos_sync)

    async def get_video_metadata(self, video_id: str) -> VideoMetadata | None:
        if not video_id:
            raise ValueError("video_id is required")

        def get_video_metadata_sync():
            youtube = build("youtube", "v3", developerKey=self._api_keys.next())
            response = youtube.videos().list(
                part="liveStreamingDetails,snippet,contentDetails",
                id=video_id,
            ).execute()
            if not response.get("items"):
                return None
            item = response["items"][0]
            iso_duration = item.get("contentDetails", {}).get("duration")
            duration_seconds = isodate.parse_duration(iso_duration).total_seconds() if iso_duration else None
            return VideoMetadata(
                title=item.get("snippet", {}).get("title"),
                live_start_time=item.get("liveStreamingDetails", {}).get("actualStartTime"),
                duration_seconds=duration_seconds,
            )

        loop = asyncio.get_running_loop()
        return await loop.run_in_executor(None, get_video_metadata_sync)

    def get_video_url(self, video_id: str) -> str:
        if not video_id:
            raise ValueError("video_id is required")
        return f"https://www.youtube.com/watch?v={video_id}"
