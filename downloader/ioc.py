import os
import tomllib

from dependency_injector import containers, providers
from dotenv import load_dotenv

from channel_poller import ChannelPoller
from multi_channel_poller import MultiChannelPoller
from twitch_downloader import TwitchDownloader
from twitch_metadata_provider import TwitchMetadataProvider
from youtube_api_key_pool import YoutubeApiKeyPool
from youtube_live_downloader import YoutubeLiveDownloader
from youtube_metadata_provider import YoutubeMetadataProvider
from fibonacci_sleep_factory import FibonacciSleepFactory


class Ioc(containers.DeclarativeContainer):
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

    youtube_live_downloader = providers.Singleton(YoutubeLiveDownloader, config=config)
    twitch_downloader = providers.Singleton(TwitchDownloader, config=config)

    sleep_factory = providers.Singleton(FibonacciSleepFactory)

    meta_providers = providers.Dict(
        youtube_live=youtube_metadata_provider,
        twitch=twitch_metadata_provider,
    )
    downloaders = providers.Dict(
        youtube_live=youtube_live_downloader,
        twitch=twitch_downloader,
    )

    channel_poller = providers.Singleton(
        ChannelPoller,
        meta_providers=meta_providers,
        downloaders=downloaders,
        sleep_factory=sleep_factory,
    )

    multi_channel_poller = providers.Singleton(
        MultiChannelPoller,
        poller=channel_poller,
        config=config,
    )
