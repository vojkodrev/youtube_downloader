import os
import tomllib

from dotenv import load_dotenv
from typing import Callable

from injector import Module, provider, singleton, SingletonScope

from channel_poller import ChannelPoller
from download_request_poller import DownloadRequestPoller
from config import Config
from downloader_map import DownloaderMap
from fibonacci_sleep_factory import FibonacciSleepFactory
from metadata_provider_map import MetadataProviderMap
from multi_channel_poller import MultiChannelPoller
from twitch_downloader import TwitchDownloader
from twitch_metadata_provider import TwitchMetadataProvider
from youtube_api_key_pool import YoutubeApiKeyPool
from youtube_live_downloader import YoutubeLiveDownloader
from youtube_video_downloader import YoutubeVideoDownloader
from youtube_metadata_provider import YoutubeMetadataProvider
from yt_dlp_logger import YtDlpLogger


class DownloaderModule(Module):
    def configure(self, binder):
        load_dotenv(os.path.join(os.path.dirname(__file__), ".env"))
        with open(os.path.join(os.path.dirname(__file__), "config.toml"), "rb") as f:
            binder.bind(Config, to=Config(tomllib.load(f)))
        binder.bind(YoutubeMetadataProvider, to=YoutubeMetadataProvider, scope=SingletonScope)
        binder.bind(TwitchMetadataProvider, to=TwitchMetadataProvider, scope=SingletonScope)
        binder.bind(YoutubeLiveDownloader, to=YoutubeLiveDownloader, scope=SingletonScope)
        binder.bind(YoutubeVideoDownloader, to=YoutubeVideoDownloader, scope=SingletonScope)
        binder.bind(TwitchDownloader, to=TwitchDownloader, scope=SingletonScope)
        binder.bind(YoutubeApiKeyPool, to=YoutubeApiKeyPool, scope=SingletonScope)
        binder.bind(FibonacciSleepFactory, to=FibonacciSleepFactory, scope=SingletonScope)
        binder.bind(ChannelPoller, to=ChannelPoller, scope=SingletonScope)
        binder.bind(MultiChannelPoller, to=MultiChannelPoller, scope=SingletonScope)
        binder.bind(DownloadRequestPoller, to=DownloadRequestPoller, scope=SingletonScope)

    @provider
    @singleton
    def yt_dlp_logger_factory(self) -> Callable[[], YtDlpLogger]:
        return YtDlpLogger

    @provider
    @singleton
    def meta_providers(
        self,
        youtube: YoutubeMetadataProvider,
        twitch: TwitchMetadataProvider,
    ) -> MetadataProviderMap:
        return MetadataProviderMap({
            "youtube_live": youtube,
            "twitch": twitch,
        })

    @provider
    @singleton
    def downloaders(
        self,
        youtube_live: YoutubeLiveDownloader,
        youtube_video: YoutubeVideoDownloader,
        twitch: TwitchDownloader,
    ) -> DownloaderMap:
        return DownloaderMap({
            "youtube_live": youtube_live,
            "youtube_video": youtube_video,
            "twitch": twitch,
        })
