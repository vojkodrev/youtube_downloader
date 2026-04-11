import asyncio
import os

import yt_dlp

from typing import Callable
from injector import inject
from config import Config
from downloader import Downloader
from yt_dlp_logger import YtDlpLogger


class YoutubeLiveDownloader(Downloader):
    @inject
    def __init__(
        self, config: Config, yt_dlp_logger_factory: Callable[[], YtDlpLogger]
    ):
        self._config = config
        self._yt_dlp_logger_factory = yt_dlp_logger_factory

    async def download(self, url: str) -> None:
        if not url:
            raise ValueError("url is required")

        def _retry_sleep(n):  # noqa: ARG001
            return 15

        def sync():
            ydl_opts = {
                "logger": self._yt_dlp_logger_factory(),
                "format": "bestvideo+bestaudio/best",
                # CRITICAL: This flag tells yt-dlp to start from the beginning of the DVR
                "live_from_start": True,
                "merge_output_format": "mp4",
                "overwrites": False,
                "outtmpl": os.path.join(
                    self._config["output_folder"], "[%(uploader)s] %(title)s.%(ext)s"
                ),
                # --- The Retry Trio ---
                "retries": 10,  # Generic network retries
                "fragment_retries": 10,  # Video chunk/fragment retries
                "extractor_retries": 10,  # Website parsing/scraping retries
                "file_access_retries": 10,  # Local disk/NAS access retries
                # Sleep between each retry attempt (all retry types)
                "retry_sleep_functions": {
                    "http": _retry_sleep,
                    "fragment": _retry_sleep,
                    "file_access": _retry_sleep,
                    "extractor": _retry_sleep,
                },
                # --- Precise Timing Control ---
                "sleep_interval": 15,  # Seconds to wait between download tasks
                "max_sleep_interval": 15,  # Keep it strictly at 15s (no randomization)
                "sleep_requests": 5,  # Wait 5s between finding info for each video
                # --- Safety Buffers ---
                "socket_timeout": 30,  # Wait 30s before considering a socket "dead"
            }
            with yt_dlp.YoutubeDL(ydl_opts) as ydl:
                ydl.download([url])

        loop = asyncio.get_running_loop()
        await loop.run_in_executor(None, sync)
