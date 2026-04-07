import os
import tomllib

from dotenv import load_dotenv
from injector import Module, provider, singleton

from channel_poller import ChannelPoller
from downloader_map import DownloaderMap
from fibonacci_sleep_factory import FibonacciSleepFactory
from metadata_provider_map import MetadataProviderMap
from multi_channel_poller import MultiChannelPoller
from twitch_downloader import TwitchDownloader
from twitch_metadata_provider import TwitchMetadataProvider
from youtube_api_key_pool import YoutubeApiKeyPool
from youtube_live_downloader import YoutubeLiveDownloader
from youtube_metadata_provider import YoutubeMetadataProvider



class Ioc(Module):
    def __init__(self):
        load_dotenv(os.path.join(os.path.dirname(__file__), ".env"))
        with open(os.path.join(os.path.dirname(__file__), "config.toml"), "rb") as f:
            self._config = tomllib.load(f)

    @provider
    @singleton
    def api_keys(self) -> YoutubeApiKeyPool:
        return YoutubeApiKeyPool(os.getenv("API_KEYS", ""))

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
        youtube: YoutubeLiveDownloader,
        twitch: TwitchDownloader,
    ) -> DownloaderMap:
        return DownloaderMap({
            "youtube_live": youtube,
            "twitch": twitch,
        })

    @provider
    @singleton
    def youtube_live_downloader(self) -> YoutubeLiveDownloader:
        return YoutubeLiveDownloader(self._config)

    @provider
    @singleton
    def twitch_downloader(self) -> TwitchDownloader:
        return TwitchDownloader(self._config)

    @provider
    @singleton
    def multi_channel_poller(self, poller: ChannelPoller) -> MultiChannelPoller:
        return MultiChannelPoller(poller, self._config)

    @provider
    @singleton
    def sleep_factory(self) -> FibonacciSleepFactory:
        return FibonacciSleepFactory()
