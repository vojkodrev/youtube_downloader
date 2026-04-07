import os
import tomllib

from dotenv import load_dotenv
from injector import Module, provider, singleton

from config import Config
from downloader_map import DownloaderMap
from metadata_provider_map import MetadataProviderMap
from twitch_downloader import TwitchDownloader
from twitch_metadata_provider import TwitchMetadataProvider
from youtube_live_downloader import YoutubeLiveDownloader
from youtube_metadata_provider import YoutubeMetadataProvider


class DownloaderModule(Module):
    def configure(self, binder):
        load_dotenv(os.path.join(os.path.dirname(__file__), ".env"))
        with open(os.path.join(os.path.dirname(__file__), "config.toml"), "rb") as f:
            binder.bind(Config, to=Config(tomllib.load(f)))

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
