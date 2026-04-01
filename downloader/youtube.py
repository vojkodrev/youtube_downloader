from __future__ import annotations

import asyncio
import itertools
import os
import sys
import tomllib
from abc import ABC, abstractmethod
from dependency_injector import containers, providers
from dotenv import load_dotenv
from googleapiclient.discovery import build
from loguru import logger
import yt_dlp


class ChannelTitleNotFoundError(Exception):
    pass


class MetadataProvider(ABC):
    @abstractmethod
    async def get_video_id(self, channel_id: str) -> str | None: ...

    @abstractmethod
    async def get_channel_title(self, channel_id: str) -> str: ...

    @abstractmethod
    def get_video_url(self, video_id: str) -> str: ...


class YoutubeMetadataProvider(MetadataProvider):
    def __init__(self, api_keys: YoutubeApiKeyPool):
        self._api_keys = api_keys

    async def get_video_id(self, channel_id: str) -> str | None:
        # DEBUG
        # if channel_id == "UCUyeluBRhGPCW4rPe_UvBZQ":
        #     return "smHwBiT1Zrk"

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
        # DEBUG
        # if channel_id == "UCUyeluBRhGPCW4rPe_UvBZQ":
        #     return "primetime"

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


class TwitchMetadataProvider(MetadataProvider):
    async def get_video_id(self, channel_id: str) -> str | None:
        return channel_id

    async def get_channel_title(self, channel_id: str) -> str:
        return channel_id

    def get_video_url(self, video_id: str) -> str:
        return f"https://www.twitch.tv/{video_id}"


class MetadataProviderFactory:
    def __init__(
        self, youtube: YoutubeMetadataProvider, twitch: TwitchMetadataProvider
    ):
        self._youtube = youtube
        self._twitch = twitch

    def create(self, mode: str) -> MetadataProvider:
        if mode == "youtube_live":
            return self._youtube
        if mode == "twitch":
            return self._twitch
        raise ValueError(f"Unsupported mode: {mode}")


class Downloader(ABC):
    @abstractmethod
    async def download(self, url: str) -> None: ...


class YoutubeLiveDownloader(Downloader):
    def __init__(self, config):
        self._config = config

    async def download(self, url: str) -> None:
        if not url:
            raise ValueError("url is required")

        def sync():
            ydl_opts = {
                "format": "bestvideo+bestaudio/best",
                # CRITICAL: This flag tells yt-dlp to start from the beginning of the DVR
                "live_from_start": True,
                "merge_output_format": "mp4",
                "overwrites": True,
                "outtmpl": os.path.join(
                    self._config["output_folder"], "%(title)s.%(ext)s"
                ),
                # "http_headers": {"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36"},
                # Optional: retry if the stream connection drops
                # "ignoreerrors": True,
                # "concurrent_fragment_downloads": 10,  # Download 10 chunks at once
            }
            with yt_dlp.YoutubeDL(ydl_opts) as ydl:
                ydl.download([url])

        loop = asyncio.get_running_loop()
        await loop.run_in_executor(None, sync)


class TwitchDownloader(Downloader):
    def __init__(self, config):
        self._config = config

    async def download(self, url: str) -> None:
        if not url:
            raise ValueError("url is required")

        def sync():
            ydl_opts = {
                "format": "best",
                "merge_output_format": "mp4",
                "overwrites": True,
                "outtmpl": os.path.join(
                    self._config["output_folder"],
                    "%(uploader)s %(upload_date>%Y-%m-%d)s.%(ext)s",
                ),
            }
            with yt_dlp.YoutubeDL(ydl_opts) as ydl:
                ydl.download([url])

        loop = asyncio.get_running_loop()
        await loop.run_in_executor(None, sync)


class DownloaderFactory:
    def __init__(self, youtube_live: YoutubeLiveDownloader, twitch: TwitchDownloader):
        self._youtube_live = youtube_live
        self._twitch = twitch

    def create(self, mode: str) -> Downloader:
        if mode == "youtube_live":
            return self._youtube_live
        if mode == "twitch":
            return self._twitch
        raise ValueError(f"Unsupported mode: {mode}")


class FibonacciSleepFactory:
    def __init__(self, short: FibonacciSleep, long: FibonacciSleep):
        self._short = short
        self._long = long

    def create(self, mode: str) -> FibonacciSleep:
        if mode == "short":
            return self._short
        elif mode == "long":
            return self._long
        else:
            raise ValueError(f"Unsupported mode: {mode}")


class FibonacciSleep:
    def __init__(self, intervals):
        self._intervals = intervals
        self._index = 0

    async def sleep(self):
        interval = self._intervals[self._index]
        self._index = min(self._index + 1, len(self._intervals) - 1)
        await asyncio.sleep(interval * 60)
        return interval

    def peek(self):
        return self._intervals[self._index]

    def reset(self):
        self._index = 0


class ChannelPoller:
    def __init__(
        self,
        meta_provider_factory: MetadataProviderFactory,
        downloader_factory: DownloaderFactory,
        sleep_factory: FibonacciSleepFactory,
    ):
        self._meta_provider_factory = meta_provider_factory
        self._downloader_factory = downloader_factory
        self._sleep_factory = sleep_factory

    async def poll(self, channel_id: str, mode: str) -> None:
        meta = self._meta_provider_factory.create(mode)
        downloader = self._downloader_factory.create(mode)

        log = logger.bind(streamer=channel_id)
        channel_title = None

        sleep_offline = self._sleep_factory.create("long")
        sleep_err = self._sleep_factory.create("short")

        while True:
            try:
                if channel_title is None:
                    channel_title = await meta.get_channel_title(channel_id)
                    log = logger.bind(streamer=channel_title)
                    log.info(f"Resolved channel ID '{channel_id}'. Polling started...")

                video_id = await meta.get_video_id(channel_id)

                if video_id:
                    sleep_offline.reset()
                    url = meta.get_video_url(video_id)
                    log.info(f"Streamer is LIVE! Downloading from: {url}")
                    await downloader.download(url)
                    sleep_err.reset()
                    log.info("Download finished. Resuming poll...")
                    await sleep_offline.sleep()
                else:
                    log.info(
                        f"Streamer is offline. Checking again in {sleep_offline.peek()} minutes..."
                    )
                    await sleep_offline.sleep()
            except ChannelTitleNotFoundError:
                raise
            except Exception as e:
                log.error(f"Error: {e}. Retrying in {sleep_err.peek()} minutes...")
                await sleep_err.sleep()


class MultiChannelPoller:
    def __init__(self, poller: ChannelPoller, config: dict):
        self._poller = poller
        self._config = config

    async def poll_all(self) -> None:
        await asyncio.gather(
            *[
                self._poller.poll(ch["id"], mode=ch["mode"])
                for ch in self._config["channels"]
            ]
        )


class YoutubeApiKeyPool:
    def __init__(self, api_keys_str: str):
        keys = [k.strip() for k in api_keys_str.split(",") if k.strip()]
        if not keys:
            raise RuntimeError("No API keys found in API_KEYS env var")
        self._cycle = itertools.cycle(keys)

    def next(self):
        return next(self._cycle)


class Container(containers.DeclarativeContainer):
    load_dotenv(os.path.join(os.path.dirname(__file__), ".env"))

    with open(os.path.join(os.path.dirname(__file__), "config.toml"), "rb") as _f:
        _config = tomllib.load(_f)

    config = providers.Object(_config)

    api_keys = providers.Singleton(
        YoutubeApiKeyPool, api_keys_str=os.getenv("API_KEYS", "")
    )

    youtube_metadata_provider = providers.Singleton(
        YoutubeMetadataProvider, api_keys=api_keys
    )
    twitch_metadata_provider = providers.Singleton(TwitchMetadataProvider)
    metadata_provider_factory = providers.Singleton(
        MetadataProviderFactory,
        youtube=youtube_metadata_provider,
        twitch=twitch_metadata_provider,
    )

    youtube_live_downloader = providers.Singleton(YoutubeLiveDownloader, config=config)
    twitch_downloader = providers.Singleton(TwitchDownloader, config=config)
    downloader_factory = providers.Singleton(
        DownloaderFactory,
        youtube_live=youtube_live_downloader,
        twitch=twitch_downloader,
    )

    sleep_factory = providers.Singleton(
        FibonacciSleepFactory,
        short=providers.Factory(FibonacciSleep, intervals=[5, 8, 13, 21, 34]),
        long=providers.Factory(FibonacciSleep, intervals=[13, 21, 34, 55]),
    )

    channel_poller = providers.Singleton(
        ChannelPoller,
        meta_provider_factory=metadata_provider_factory,
        downloader_factory=downloader_factory,
        sleep_factory=sleep_factory,
    )

    multi_channel_poller = providers.Singleton(
        MultiChannelPoller,
        poller=channel_poller,
        config=config,
    )


def main():
    container = Container()
    config = container.config()

    logger.remove()
    logger.add(sys.stderr, format=config["log_format"])

    if not os.path.isdir(config["output_folder"]):
        logger.bind(streamer="-").error(
            f"Output folder does not exist: {config['output_folder']}"
        )
        exit(1)

    multi_poller = container.multi_channel_poller()
    asyncio.run(multi_poller.poll_all())


if __name__ == "__main__":
    main()
