import asyncio
import os

import yt_dlp

from typing import Callable
from injector import inject
from config import Config
from downloader import Downloader
from yt_dlp_logger import YtDlpLogger
from yt_dlp_settings import SHARED_YT_DLP_SETTINGS


class YoutubeLiveDownloader(Downloader):
    @inject
    def __init__(
        self, config: Config, yt_dlp_logger_factory: Callable[[], YtDlpLogger]
    ):
        self._config = config
        self._yt_dlp_logger_factory = yt_dlp_logger_factory

    async def download(self, url: str, streamer: str = "-") -> None:
        if not url:
            raise ValueError("url is required")

        def sync():
            yt_dlp_logger = self._yt_dlp_logger_factory()
            yt_dlp_logger.set_streamer(streamer)
            ydl_opts = {
                **SHARED_YT_DLP_SETTINGS,
                "logger": yt_dlp_logger,
                "js_runtimes": {"node": {}},
                "remote_components": ["ejs:github"],
                "format": "bestvideo+bestaudio/best",
                # CRITICAL: This flag tells yt-dlp to start from the beginning of the DVR
                "live_from_start": True,
                "merge_output_format": "mp4",
                "overwrites": True,
                "outtmpl": os.path.join(
                    self._config["output_folder"], "[%(uploader)s] %(title)s.%(ext)s"
                ),
            }
            with yt_dlp.YoutubeDL(ydl_opts) as ydl:
                ydl.download([url])

        loop = asyncio.get_running_loop()
        await loop.run_in_executor(None, sync)
