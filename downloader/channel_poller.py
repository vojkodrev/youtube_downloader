from __future__ import annotations

from loguru import logger

from channel_title_not_found_error import ChannelTitleNotFoundError
from downloader import Downloader
from fibonacci_sleep_factory import FibonacciSleepFactory
from metadata_provider import MetadataProvider


class ChannelPoller:
    def __init__(
        self,
        meta_providers: dict[str, MetadataProvider],
        downloaders: dict[str, Downloader],
        sleep_factory: FibonacciSleepFactory,
    ):
        self._meta_providers = meta_providers
        self._downloaders = downloaders
        self._sleep_factory = sleep_factory

    async def poll(self, channel_id: str, mode: str) -> None:
        meta = self._meta_providers[mode]
        downloader = self._downloaders[mode]

        log = logger.bind(streamer=channel_id)
        channel_title = None

        sleep_offline = self._sleep_factory.create("long", mode)
        sleep_err = self._sleep_factory.create("short", mode)

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
                    log.info(f"Downloading from: {url}")
                    await downloader.download(url)
                    sleep_err.reset()
                    log.info(
                        f"Download finished. Resuming poll in {sleep_offline.peek()} minutes..."
                    )
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
